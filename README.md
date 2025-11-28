```go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/learning/go-todo-clean/internal/service"
)

// TaskHandler expone los endpoints HTTP relacionados con tareas (Task).
// Actúa como capa de transporte/entrega, delegando la lógica de negocio
// en el servicio de dominio (TaskService).
type TaskHandler struct {
	service *service.TaskService
}

// NewTaskHandler construye una nueva instancia de TaskHandler.
func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// CreateTask maneja la petición HTTP POST para crear una nueva tarea.
//
// Endpoint esperado:
//   POST /tasks
//
// Request Body (JSON):
//   {
//     "title": "Texto del título de la tarea"
//   }
//
// Response 201 (application/json):
//   {
//     "id": "...",
//     "title": "...",
//     "created_at": "...",
//     ...
//   }
//
// Códigos de estado posibles:
//   201 Created      -> Tarea creada exitosamente
//   400 Bad Request  -> JSON inválido o datos de entrada inválidos
//   500 Internal Server Error -> Error interno al crear la tarea
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// 1. Definimos un DTO temporal para el cuerpo de la petición
	type requestBody struct {
		Title string `json:"title"`
	}

	var body requestBody

	// 2. Decodificamos el JSON recibido
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 3. Validación básica de campos requeridos
	if body.Title == "" {
		http.Error(w, "El campo 'title' es obligatorio", http.StatusBadRequest)
		return
	}

	// 4. Llamamos al servicio de dominio para crear la tarea
	createdTask, err := h.service.CreateTask(body.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 5. Devolvemos la respuesta al cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		// Si fallamos al codificar la respuesta, devolvemos 500
		http.Error(w, "Error al generar la respuesta", http.StatusInternalServerError)
		return
	}
}

// GetAllTasks maneja la petición HTTP GET para obtener el listado de tareas.
//
// Endpoint esperado:
//   GET /tasks
//
// Response 200 (application/json):
//   [
//     {
//       "id": "...",
//       "title": "...",
//       "created_at": "...",
//       ...
//     },
//     ...
//   ]
//
// Códigos de estado posibles:
//   200 OK           -> Listado devuelto correctamente
//   500 Internal Server Error -> Error al obtener las tareas
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetTasks()
	if err != nil {
		http.Error(w, "Error al obtener las tareas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Error al generar la respuesta", http.StatusInternalServerError)
		return
	}
}
```