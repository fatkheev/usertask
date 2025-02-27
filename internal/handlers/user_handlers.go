package handlers

import (
	"net/http"
	"strconv"
	"usertask/internal/auth"
	"usertask/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUserGin создает нового пользователя.
// @Summary Создать пользователя
// @Description Регистрирует нового пользователя и выдаёт токен
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.RequestCreateUser true "Данные пользователя"
// @Success 201 {object} models.ResponseCreateUser "Пользователь успешно создан"
// @Router /users/create [post]
func (h *UserHandler) CreateUserGin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, token, err := h.userService.CreateUser(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

// SetReferrerGin устанавливает реферала для пользователя.
// @Summary Установить реферала
// @Description Позволяет пользователю указать, кто его пригласил. Если успешен, рефереру начисляется бонус и создается запись в tasks.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param request body models.RequestSetReferrer true "ID реферера"
// @Success 200 {object} models.ResponseSetReferrer "Реферальный код успешно установлен"
// @Failure 500 {object} models.ErrorSetReferrerConflict "Реферальный код уже установлен"
// @Router /users/{id}/referrer [post]
func (h *UserHandler) SetReferrerGin(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		ReferrerID int `json:"referrer_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.ReferrerID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if err := h.userService.SetReferrer(userID, req.ReferrerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "referrer set successfully"})
}

// GetUserStatusGin возвращает информацию о пользователе.
// @Summary Получить статус пользователя
// @Description Возвращает детали пользователя по его ID
// @Tags Users
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User
// @Router /users/{id}/status [get]
func (h *UserHandler) GetUserStatusGin(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userService.GetUserStatus(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetLeaderboardGin возвращает список лидеров по очкам.
// @Summary Получить лидерборд
// @Description Возвращает список пользователей с наибольшим количеством очков.
// @Tags Leaderboard
// @Security BearerAuth
// @Param limit query int false "Количество пользователей в списке (по умолчанию 10)"
// @Success 200 {array} models.User "Список лидеров"
// @Router /users/leaderboard [get]
func (h *UserHandler) GetLeaderboardGin(c *gin.Context) {
	limit := 10

	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	users, err := h.userService.GetLeaderboard(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get leaderboard"})
		return
	}

	c.JSON(http.StatusOK, users)
}


// CompleteTaskGin выполняет задание.
// @Summary Завершить задание
// @Description Добавляет пользователю очки за выполнение задания
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param request body models.RequestCompleteTask true "Детали задания"
// @Success 200 {object} models.ResponseCompleteTask "Задание успешно завершено"
// @Router /users/{id}/task/complete [post]
func (h *UserHandler) CompleteTaskGin(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		TaskType string `json:"task_type"`
		Points   int    `json:"points"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.TaskType == "" || req.Points <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	err = h.userService.CompleteTask(userID, req.TaskType, req.Points)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task completed",
		"points_awarded": req.Points,
	})
}

// RefreshTokenGin обновляет токен пользователя.
// @Summary Обновить токен
// @Description Генерирует новый JWT-токен для пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RequestRefreshToken true "ID пользователя"
// @Success 200 {object} models.ResponseRefreshToken "Новый токен успешно сгенерирован"
// @Failure 500 {object} models.ErrorRefreshTokenUserNotFound "Пользователь не найден"
// @Router /users/token/refresh [post]
func (h *UserHandler) RefreshTokenGin(c *gin.Context) {
	var req struct {
		UserID int `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.UserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, err := h.userService.GetUserStatus(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	token, err := auth.GenerateToken(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "new token generated",
		"token":   token,
	})
}
