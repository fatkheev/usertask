package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TaskType  string    `json:"task_type"`
	Points    int       `json:"points"`
	CompletedAt time.Time `json:"completed_at"`
}
