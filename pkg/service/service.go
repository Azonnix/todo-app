package service

import (
	"github.com/azonnix/todo-app"
	"github.com/azonnix/todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateTocken(username, password string) (string, error)
}

type TodoList struct {
}

type TodoItem struct {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
