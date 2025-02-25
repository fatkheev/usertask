package main

import (
	"fmt"
	"log"
	"os"
	"usertask/internal/database"
	"usertask/internal/handlers"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.GET("/users/:id/status", handlers.GetUserStatusGin)
	router.GET("/users/leaderboard", handlers.GetLeaderboardGin)
	router.POST("/users/:id/task/complete", handlers.CompleteTaskGin)
	router.POST("/users/:id/referrer", handlers.SetReferrerGin)

	fmt.Printf("Starting server on :%s\n", port)
	log.Fatal(router.Run(":" + port))
}
