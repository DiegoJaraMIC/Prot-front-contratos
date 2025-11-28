```go
// Package repository proporciona una implementación en memoria para almacenar
// y gestionar tareas (domain.Task) durante la ejecución de la aplicación.
//
// Esta implementación está pensada principalmente para:
//   - Entornos de desarrollo y pruebas.
//   - Prototipado rápido sin dependencia de una base de datos persistente.
//   - Casos donde se requiera un repositorio simple sin I/O externo.
//
// NOTA: Esta implementación **no es segura para concurrencia**. En caso de uso
// en entornos concurrentes, se recomienda envolverla con mecanismos de
// sincronización (por ejemplo, sync.Mutex) o implementar un repositorio
// específico para producción.
package repository

import (
	"github.com/learning/go-todo-clean/internal/domain"
)

// InMemoryRepo representa un repositorio en memoria de tareas.
type InMemoryRepo struct {
	tasks  []domain.Task
	nextID int
}

// NewInMemoryRepo crea y retorna una nueva instancia de InMemoryRepo
// inicializada sin tareas y con el ID inicial configurado en 1.
func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		tasks:  []domain.Task{},
		nextID: 1,
	}
}

// Save agrega una nueva tarea al repositorio en memoria.
//
// Comportamiento:
//   - Asigna un ID incremental a la tarea (sobrescribiendo cualquier valor previo).
//   - Agrega la tarea a la colección interna.
//   - No realiza validaciones de negocio (se asume que la entidad ya fue validada).
//
// Limitaciones:
//   - No garantiza concurrencia segura.
//   - No persiste datos más allá del ciclo de vida del proceso.
//
// Ejemplo de uso:
//
//   repo := NewInMemoryRepo()
//   t := &domain.Task{Title: "Mi tarea"}
//   if err := repo.Save(t); err != nil {
//       // manejar error
//   }
func (r *InMemoryRepo) Save(task *domain.Task) error {
	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, *task)

	return nil
}

// GetAll retorna el listado completo de tareas almacenadas en memoria.
//
// Devuelve:
//   - Un slice de domain.Task con todas las tareas hasta el momento.
//   - Un error en caso de fallo (actualmente siempre retorna nil).
//
// NOTA: Se devuelve una copia del slice interno para evitar que el estado
// interno pueda ser modificado de forma no controlada por el consumidor.
func (r *InMemoryRepo) GetAll() ([]domain.Task, error) {
	tasksCopy := make([]domain.Task, len(r.tasks))
	copy(tasksCopy, r.tasks)

	return tasksCopy, nil
}
```