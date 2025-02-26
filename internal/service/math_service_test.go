package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMathProblem(t *testing.T) {
	problem := GenerateMathProblem()

	assert.NotZero(t, problem.Operand1)
	assert.NotZero(t, problem.Operand2)
	assert.Equal(t, "+", problem.Operation)
}

func TestCheckMathAnswer(t *testing.T) {
	problem := MathProblem{Operand1: 2, Operand2: 3, Operation: "+", Answer: 5}

	correct, err := CheckMathAnswer(problem, 5)
	assert.True(t, correct)
	assert.NoError(t, err)

	incorrect, err := CheckMathAnswer(problem, 6)
	assert.False(t, incorrect)
	assert.Error(t, err)
}
