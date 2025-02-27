package models

// ответ при регистрации пользователя.
type ResponseCreateUser struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}


// ответ при установке реферала.
type ResponseSetReferrer struct {
	Message string `json:"message" example:"referrer set successfully"`
}

// реферальный код уже установлен
type ErrorSetReferrerConflict struct {
	Error string `json:"error" example:"реферальный код уже установлен"`
}


// ответ при обновлении токена.
type ResponseRefreshToken struct {
	Message string `json:"message" example:"new token generated"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// пользователь не найден
type ErrorRefreshTokenUserNotFound struct {
	Error string `json:"error" example:"пользователь не найден"`
}


// ответ при завершении задания
type ResponseCompleteTask struct {
	Message       string `json:"message" example:"task completed"`
	PointsAwarded int    `json:"points_awarded" example:"50"`
}


// математическая задача
type ResponseMathProblem struct {
	Operand1  int    `json:"operand1" example:"5"`
	Operand2  int    `json:"operand2" example:"3"`
	Operation string `json:"operation" example:"+"`
}

// ответ после успешного решения задачи.
type ResponseSolveMathProblem struct {
	Message       string `json:"message" example:"correct answer!"`
	PointsAwarded int    `json:"points_awarded" example:"50"`
}

// неверный ответ к математической задаче
type ErrorSolveMathIncorrectAnswer struct {
	Error string `json:"error" example:"неверный ответ"`
}