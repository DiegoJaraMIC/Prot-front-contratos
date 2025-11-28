package handler

import (
	"encoding/json"
	"net/http"

	"github.com/learning/go-todo-clean/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// CreateTask maneja la petición POST
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// 1. Definimos una estructura temporal para leer el JSON entrante
	var body struct {
		Title string `json:"title"`
	}

	// 2. Decodificamos el body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 3. Llamamos al servicio
	createdTask, err := h.service.CreateTask(body.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Respondemos con éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

// GetAllTasks maneja la petición GET
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, _ := h.service.GetTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
