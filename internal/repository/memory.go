package repository

import (
	"github.com/learning/go-todo-clean/internal/domain"
)

type InMemoryRepo struct {
	tasks  []domain.Task
	nextID int
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		tasks:  []domain.Task{},
		nextID: 1,
	}
}

func (r *InMemoryRepo) Save(task *domain.Task) error {
	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, *task)
	return nil
}

func (r *InMemoryRepo) GetAll() ([]domain.Task, error) {
	return r.tasks, nil
}
