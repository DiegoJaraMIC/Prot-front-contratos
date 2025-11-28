package domain

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskRepository interface {
	Save(task *Task) error
	GetAll() ([]Task, error)
}
