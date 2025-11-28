package main

import (
	"fmt"
	"net/http"

	"github.com/learning/go-todo-clean/internal/handler"
	"github.com/learning/go-todo-clean/internal/repository"
	"github.com/learning/go-todo-clean/internal/service"
)

func main() {
	// 1. Inyección de Dependencias
	repo := repository.NewInMemoryRepo()
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	// 2. Definición de Rutas
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateTask(w, r)
		} else if r.Method == http.MethodGet {
			h.GetAllTasks(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// 3. Arrancar servidor
	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
