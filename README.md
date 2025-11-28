```go
package main

import (
	"log"
	"net/http"

	"github.com/learning/go-todo-clean/internal/handler"
	"github.com/learning/go-todo-clean/internal/repository"
	"github.com/learning/go-todo-clean/internal/service"
)

// main es el punto de entrada de la aplicación HTTP.
// Se encarga de:
//   1. Instanciar las dependencias (repositorio, servicios, handlers).
//   2. Definir las rutas HTTP y asociarlas a los handlers correspondientes.
//   3. Iniciar el servidor HTTP.
func main() {
	// 1. Inyección de dependencias
	repo := repository.NewInMemoryRepo()
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	// 2. Definición de rutas
	// Se utiliza la ruta base /tasks para exponer el recurso de tareas.
	//   - POST   /tasks -> Crear una nueva tarea
	//   - GET    /tasks -> Listar todas las tareas
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateTask(w, r)
		case http.MethodGet:
			h.GetAllTasks(w, r)
		default:
			http.Error(
				w,
				"Método no permitido",
				http.StatusMethodNotAllowed,
			)
		}
	})

	// 3. Arranque del servidor
	addr := ":8080"
	log.Printf("Servidor iniciado en http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("error al iniciar el servidor: %v", err)
	}
}
```