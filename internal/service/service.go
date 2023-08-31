package service

import "errors"

var (
	ErrEmailTaken             = errors.New("email already taken")
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrNoSuchUser             = errors.New("no such user")
	ErrTodoItemNotFound       = errors.New("todo item not found")
)
