package model

import (
	"time"
)

type Task struct {
	ID          int        `json:"id" db:"id"`
	UserID      int        `json:"user_id" db:"user_id"`
	Description string     `json:"description" db:"description"`
	Duration    int        `json:"duration" db:"duration"` // in minutes
	Completed   bool       `json:"completed" db:"completed"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
}
