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

type AuthenticatedUser struct {
	ID    uint   `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
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

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"gte=6,lte=32"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (req CreateUserRequest) Validate() error {
	return validator.New().Struct(&req)
}
