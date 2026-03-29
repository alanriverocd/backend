package main

import (
	"log"
	"os"

	"github.com/finatiol/backend/config"
	"github.com/finatiol/backend/handlers"
	"github.com/finatiol/backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	config.InitFirebase()

	r := gin.Default()

	// Configuracion CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "app": "FINATIOL"})
	})

	api := r.Group("/api/v1")
	{
		// Rutas publicas
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Rutas protegidas
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/users/me", handlers.GetCurrentUser)
			protected.PUT("/users/me", handlers.UpdateCurrentUser)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("FINATIOL backend corriendo en el puerto %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
