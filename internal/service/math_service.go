package service

import (
	"errors"
	"math/rand"
	"time"
)

type MathProblem struct {
	Operand1  int    `json:"operand1"`
	Operand2  int    `json:"operand2"`
	Operation string `json:"operation"`
	Answer    int    `json:"-"`
}

func GenerateMathProblem() MathProblem {
	rand.Seed(time.Now().UnixNano())

	operand1 := rand.Intn(10) + 1
	operand2 := rand.Intn(10) + 1
	operation := "+"
	answer := operand1 + operand2

	return MathProblem{
		Operand1:  operand1,
		Operand2:  operand2,
		Operation: operation,
		Answer:    answer,
	}
}

func CheckMathAnswer(problem MathProblem, userAnswer int) (bool, error) {
	if userAnswer == problem.Answer {
		return true, nil
	}
	return false, errors.New("неверный ответ")
}
