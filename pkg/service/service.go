package service

import "github.com/azonnix/todo-app/pkg/repository"

type Authorization struct {
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
	return &Service{}
}
