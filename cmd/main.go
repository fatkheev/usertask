package main

import (
	"fmt"
	"log"
	"os"
	"usertask/docs"
	"usertask/internal/database"
	"usertask/internal/handlers"
	"usertask/internal/middleware"
	"usertask/internal/repository"
	"usertask/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title UserTask API
// @version 1.0
// @description API для управления пользователями и заданиями.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		protected.GET("/users/leaderboard", userHandler.GetLeaderboardGin)
	}

	fmt.Printf("Starting server on :%s\n", port)
	log.Fatal(router.Run(":" + port))
}
