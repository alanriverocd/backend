package handlers

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/finatiol/backend/config"
	"github.com/finatiol/backend/middleware"
	"github.com/finatiol/backend/models"
	"github.com/gin-gonic/gin"
)

// GetCurrentUser retorna los datos del usuario autenticado.
func GetCurrentUser(c *gin.Context) {
	uid, ok := middleware.GetUserUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	doc, err := config.FirestoreClient.Collection("users").Doc(uid).Get(context.Background())
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

// UpdateCurrentUser actualiza los datos del usuario autenticado.
func UpdateCurrentUser(c *gin.Context) {
	uid, ok := middleware.GetUserUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := []firestore.Update{
		{Path: "updatedAt", Value: time.Now()},
	}
	if req.DisplayName != "" {
		updates = append(updates, firestore.Update{Path: "displayName", Value: req.DisplayName})
	}
	if req.PhotoURL != "" {
		updates = append(updates, firestore.Update{Path: "photoUrl", Value: req.PhotoURL})
	}

	_, err := config.FirestoreClient.Collection("users").Doc(uid).Update(context.Background(), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado correctamente"})
}
