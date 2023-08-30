package model

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Email     string     `gorm:"index" json:"email"`
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"lastname"`
	Password  string     `json:"-"`
	CreatedAt *time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

type TodoItem struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CompletedAt *time.Time `gorm:"index" json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

func (item TodoItem) IsComplete() bool {
	return item.CompletedAt != nil
}
