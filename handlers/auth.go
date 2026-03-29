package handlers

import (
	"context"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/finatiol/backend/config"
	"github.com/finatiol/backend/models"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"displayName" binding:"required"`
}

type LoginRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}

// Register crea un nuevo usuario en Firebase Auth y Firestore.
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password).
		DisplayName(req.DisplayName).
		EmailVerified(false)

	userRecord, err := config.AuthClient.CreateUser(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario: " + err.Error()})
		return
	}

	now := time.Now()
	user := models.User{
		UID:         userRecord.UID,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = config.FirestoreClient.Collection("users").Doc(userRecord.UID).Set(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar usuario en Firestore"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado exitosamente", "uid": userRecord.UID})
}

// Login verifica el ID token de Firebase y retorna los datos del usuario.
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := config.AuthClient.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalido"})
		return
	}

	doc, err := config.FirestoreClient.Collection("users").Doc(token.UID).Get(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer datos del usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
