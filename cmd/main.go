package main

import (
	"fmt"
	"log"
	"os"
	"usertask/internal/database"
	"usertask/internal/handlers"
	"usertask/internal/middleware"
	"usertask/internal/repository"
	"usertask/internal/service"

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
	mathHandler := handlers.NewMathHandler(userService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.POST("/users/create", userHandler.CreateUserGin)
	router.POST("/users/token/refresh", userHandler.RefreshTokenGin)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users/:id/status", userHandler.GetUserStatusGin)
		protected.POST("/users/:id/task/complete", userHandler.CompleteTaskGin)
		protected.POST("/users/:id/referrer", userHandler.SetReferrerGin)
		protected.GET("/users/:id/task/math", mathHandler.GetMathProblem)
		protected.POST("/users/:id/task/math/solve", mathHandler.SolveMathProblem)
	}

	fmt.Printf("Starting server on :%s\n", port)
	log.Fatal(router.Run(":" + port))
}
