package main

import (
	"fmt"
	"os"

	"my-personal-web/api/config"
	"my-personal-web/api/routes"
	"my-personal-web/api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	utils.MigrateSeeders()

	r := gin.Default()

	// Setup CORS
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// register routes
	routes.RegisterRoutes(r)

	r.Static("/public", "./public")

	// start server
	fmt.Println("🚀 Server running at http://localhost:8080")
	r.Run(":8080")
}
