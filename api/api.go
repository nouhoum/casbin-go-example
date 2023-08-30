package api

import (
	"github.com/go-playground/validator/v10"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gte=6,lte=32"`
}

func (auth AuthRequest) Validate() error {
	return validator.New().Struct(&auth)
}

type CreateOrUpdateTodoItemRequest struct {
	Title       string `json:"title" validate:"gte=6,lte=32"`
	Description string `json:"description"`
}

func (req CreateOrUpdateTodoItemRequest) Validate() error {
	return validator.New().Struct(&req)
}

type TodoItemCompleteRequest struct {
	IsComplete bool `json:"description"`
}
