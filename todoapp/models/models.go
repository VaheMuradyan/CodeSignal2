package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	DueDate   time.Time `json:"due_date"`
	Priority  string    `json:"priority"`
}

type Library struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Type     string `json:"type"`
	IsOpen   bool   `json:"is_open"`
}

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
}

type Credentials struct {
	Username string `json:"username" validate:"min=5"`
	Password string `json:"password" validate:"min=8"`
}
