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
