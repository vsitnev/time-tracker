package model

import (
	"time"
)

type User struct {
	ID             int    `json:"id" db:"id"`
	Name           string `json:"name" db:"username"`
	Surname        string `json:"surname" db:"surname"`
	Patronymic     string `json:"patronymic" db:"patronymic"`
	PassportNumber string `json:"passport_number" db:"passport_number"`
	Address        string `json:"address" db:"address"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
