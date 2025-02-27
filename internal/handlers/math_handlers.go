package handlers

import (
	"net/http"
	"strconv"
	"usertask/internal/service"

	"github.com/gin-gonic/gin"
)

type MathHandler struct {
	userService *service.UserService
}

func NewMathHandler(userService *service.UserService) *MathHandler {
	return &MathHandler{userService: userService}
}

var MathProblemStorage = make(map[int]service.MathProblem)

// GetMathProblem выдаёт пользователю случайную математическую задачу.
// @Summary Получить математическую задачу
// @Description Генерирует и возвращает пользователю случайную математическую задачу.
// @Tags MathTasks
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.ResponseMathProblem "Математическая задача успешно сгенерирована"
// @Router /users/{id}/task/math [get]
func (h *MathHandler) GetMathProblem(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	problem := service.GenerateMathProblem()
	MathProblemStorage[userID] = problem

	c.JSON(http.StatusOK, gin.H{
		"operand1":  problem.Operand1,
		"operand2":  problem.Operand2,
		"operation": problem.Operation,
	})
}

// SolveMathProblem проверяет ответ на математическую задачу.
// @Summary Решить математическую задачу
// @Description Проверяет правильность ответа на задачу, начисляет очки за верный ответ.
// @Tags MathTasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param request body models.RequestSolveMathProblem true "Ответ пользователя на задачу"
// @Success 200 {object} models.ResponseSolveMathProblem "Ответ правильный, начислены очки"
// @Failure 400 {object} models.ErrorSolveMathIncorrectAnswer "Неверный ответ"
// @Router /users/{id}/task/math/solve [post]
func (h *MathHandler) SolveMathProblem(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Answer int `json:"answer"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	problem, exists := MathProblemStorage[userID]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "math problem not found"})
		return
	}

	correct, err := service.CheckMathAnswer(problem, req.Answer)
	if !correct {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	points := 50
	err = h.userService.CompleteTask(userID, "math_problem", points)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	delete(MathProblemStorage, userID)

	c.JSON(http.StatusOK, gin.H{
		"message":       "correct answer!",
		"points_awarded": points,
	})
}
