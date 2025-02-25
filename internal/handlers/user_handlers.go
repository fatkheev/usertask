package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Возвращаем информацию о пользователе
func GetUserStatusGin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User status endpoint"})
}

// Возвращаем список лидеров по баллам
func GetLeaderboardGin(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Leaderboard endpoint"})
}

// Отмечаем выполнение задания
func CompleteTaskGin(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Complete task endpoint"})
}

// Рефералка
func SetReferrerGin(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Set referrer endpoint"})
}