package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/Syzua/degree-auth/backend/handlers"
	"github.com/Syzua/degree-auth/backend/middleware"
	"github.com/Syzua/degree-auth/backend/services"
)

func main() {
	if err := services.InitSDK(); err != nil {
		log.Printf("Warning: SDK init failed (running in mock mode): %v", err)
	}
	defer services.CloseSDK()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Static("/frontend", "./frontend")

	api := r.Group("/api/v1")

	api.POST("/login", handlers.Login)

	auth := api.Group("")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/education", middleware.RequireRole("university"), handlers.AddEducation)
		auth.PUT("/education", middleware.RequireRole("university"), handlers.UpdateEducation)
		auth.GET("/education/:certNo", handlers.QueryEducationByID)
		auth.GET("/education/verify", handlers.VerifyEducation)
		auth.GET("/education/:certNo/history", handlers.GetHistory)
		auth.POST("/education/:certNo/authorize", middleware.RequireRole("student"), handlers.AuthorizeViewer)
	}

	fmt.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}
