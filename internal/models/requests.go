package models

// для создания пользователя.
type RequestCreateUser struct {
	Username string `json:"username" example:"user123"`
}

// JSON-запрос для выполнения задания.
type RequestCompleteTask struct {
	TaskType string `json:"task_type" example:"math_problem"`
	Points   int    `json:"points" example:"50"`
}

// JSON-запрос для установки реферера.
type RequestSetReferrer struct {
	ReferrerID int `json:"referrer_id" example:"1"`
}

// JSON-запрос для обновления токена.
type RequestRefreshToken struct {
	UserID int `json:"user_id" example:"1"`
}

// JSON-запрос для решения математической задачи.
type RequestSolveMathProblem struct {
	Answer int `json:"answer" example:"42"`
}
