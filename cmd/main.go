package main

import (
	"fmt"
	"log"
	"usertask/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/users/:id/status", handlers.GetUserStatusGin)
	router.GET("/users/leaderboard", handlers.GetLeaderboardGin)
    router.POST("/users/:id/task/complete", handlers.CompleteTaskGin)
    router.POST("/users/:id/referrer", handlers.SetReferrerGin)

	fmt.Println("Starting server on :8080")
	log.Fatal(router.Run(":8080"))
}