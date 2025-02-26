package main

import (
	"fmt"
	"log"
	"os"
	"usertask/internal/database"
	"usertask/internal/handlers"
	"usertask/internal/repository"
	"usertask/internal/service"
	"usertask/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить .env, используются переменные окружения")
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer database.CloseDB()

	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.POST("/users/create", userHandler.CreateUserGin)

	// Закрываем эндпоинты JWT-аутентификацией
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users/:id/status", userHandler.GetUserStatusGin)
		protected.POST("/users/:id/task/complete", userHandler.CompleteTaskGin)
		protected.POST("/users/:id/referrer", userHandler.SetReferrerGin)
	}

	fmt.Printf("Starting server on :%s\n", port)
	log.Fatal(router.Run(":" + port))
}
