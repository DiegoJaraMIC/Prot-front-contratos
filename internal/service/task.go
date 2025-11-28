package service

import (
	"errors"

	"github.com/learning/go-todo-clean/internal/domain"
)

type TaskService struct {
	repo domain.TaskRepository
}

func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(title string) (*domain.Task, error) {
	if title == "" {
		return nil, errors.New("el título no puede estar vacío")
	}

	newTask := &domain.Task{
		Title:       title,
		IsCompleted: false,
	}

	err := s.repo.Save(newTask)
	if err != nil {
		return nil, err
	}

	return newTask, nil
}

func (s *TaskService) GetTasks() ([]domain.Task, error) {
	return s.repo.GetAll()
}
