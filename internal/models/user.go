package models

import "time"

type User struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	Points     int        `json:"points"`
	ReferrerID *int       `json:"referrer_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}