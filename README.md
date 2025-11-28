```go
package service

import (
	"errors"

	"github.com/learning/go-todo-clean/internal/domain"
)

// TaskService define la interfaz del servicio de tareas.
// Encapsula la lógica de negocio asociada a la gestión de tareas.
type TaskService struct {
	repo domain.TaskRepository
}

// NewTaskService crea una nueva instancia de TaskService.
//
// Parametros:
//   - repo: implementación de domain.TaskRepository que provee acceso a almacenamiento.
//
// Retorna:
//   - *TaskService: instancia inicializada del servicio de tareas.
func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

// CreateTask crea una nueva tarea con el título proporcionado.
//
// Reglas de negocio:
//   - El título no puede ser vacío.
//
// Parametros:
//   - title: título de la tarea a crear.
//
// Retorna:
//   - *domain.Task: tarea creada.
//   - error: error en caso de validación o fallo al persistir.
func (s *TaskService) CreateTask(title string) (*domain.Task, error) {
	if title == "" {
		return nil, errors.New("el título no puede estar vacío")
	}

	newTask := &domain.Task{
		Title:       title,
		IsCompleted: false,
	}

	if err := s.repo.Save(newTask); err != nil {
		return nil, err
	}

	return newTask, nil
}

// GetTasks obtiene todas las tareas registradas.
//
// Retorna:
//   - []domain.Task: listado de tareas.
//   - error: error al acceder al repositorio.
func (s *TaskService) GetTasks() ([]domain.Task, error) {
	return s.repo.GetAll()
}
```