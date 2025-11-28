```go
package domain

import "time"

// Task representa una tarea dentro del dominio de la aplicación.
// Define la estructura principal utilizada en la capa de negocio.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	// IsCompleted indica si la tarea ha sido marcada como completada.
	IsCompleted bool      `json:"is_completed"`
	// CreatedAt es la fecha/hora en que la tarea fue creada.
	CreatedAt   time.Time `json:"created_at"`
	// UpdatedAt es la fecha/hora de la última modificación de la tarea.
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskRepository define el contrato que debe cumplir cualquier
// implementación de repositorio para la entidad Task.
//
// Esta interfaz permite desacoplar la capa de dominio de los detalles
// de persistencia (por ejemplo, base de datos SQL, NoSQL, memoria, etc.).
type TaskRepository interface {
	// Save persiste una tarea. Debe crear una nueva entrada si la tarea
	// no existe, o actualizarla si ya está almacenada.
	Save(task *Task) error

	// GetAll devuelve el listado completo de tareas almacenadas.
	// El slice devuelto no debe ser nil; en ausencia de resultados,
	// debe retornarse un slice vacío.
	GetAll() ([]Task, error)
}
```