# 📚 Clase Completa: Arquitectura Clean en Go

## 🎯 ¿Qué es Clean Architecture?

Imagina que tu aplicación es como una cebolla con capas:
- **Capa externa (Handler)**: Recibe peticiones HTTP
- **Capa de servicio**: Contiene la lógica de negocio
- **Capa de repositorio**: Guarda y obtiene datos
- **Capa de dominio**: Define qué es una tarea (el núcleo)

Cada capa solo conoce a la capa de adentro, nunca a la de afuera. Esto hace el código más fácil de mantener y probar.

---

## 📁 Estructura del Proyecto

```
go-todo-clean/
├── cmd/api/main.go          ← Punto de entrada (arranca todo)
├── internal/
│   ├── domain/task.go       ← Define QUÉ es una tarea
│   ├── repository/memory.go ← Guarda tareas en memoria
│   ├── service/task.go      ← Lógica de negocio
│   └── handler/http.go     ← Recibe peticiones HTTP
└── go.mod                   ← Dependencias del proyecto
```

---

## 📄 1. DOMAIN (domain/task.go) - El Núcleo

**¿Qué hace?** Define qué es una "Tarea" y qué operaciones se pueden hacer con ella.

### Estructura Task

```go
type Task struct {
    ID          int       // Identificador único
    Title       string    // Título de la tarea
    IsCompleted bool      // ¿Está completada?
    CreatedAt   time.Time // Cuándo se creó
    UpdatedAt   time.Time // Cuándo se actualizó
}
```

**Analogía:** Es como el "molde" o "plantilla" de una tarea. Define qué información tiene una tarea.

**Los tags `json:"..."`** le dicen a Go cómo convertir la estructura a JSON cuando respondes al cliente.

### Interfaz TaskRepository

```go
type TaskRepository interface {
    Save(task *Task) error
    GetAll() ([]Task, error)
}
```

**¿Qué es una interfaz?** Es como un "contrato" que dice: "Cualquiera que quiera ser un repositorio DEBE tener estos métodos".

**¿Por qué usar una interfaz?** 
- Puedes cambiar cómo guardas las tareas (memoria, base de datos, archivo) sin cambiar el resto del código
- Facilita las pruebas (puedes crear un repositorio "falso" para probar)

**Analogía:** Es como decir "cualquier caja que tenga una ranura para guardar y otra para sacar puede ser un repositorio".

---

## 📄 2. REPOSITORY (repository/memory.go) - El Almacén

**¿Qué hace?** Guarda las tareas en la memoria de la computadora (cuando reinicias el servidor, se pierden).

### Estructura InMemoryRepo

```go
type InMemoryRepo struct {
    tasks  []domain.Task  // Lista de tareas guardadas
    nextID int            // Próximo ID a asignar
}
```

**Analogía:** Es como una caja donde guardas tus tareas. `tasks` es la lista de tareas y `nextID` es un contador para darle un número único a cada tarea.

### Función NewInMemoryRepo()

```go
func NewInMemoryRepo() *InMemoryRepo {
    return &InMemoryRepo{
        tasks:  []domain.Task{},  // Lista vacía
        nextID: 1,                 // Empezamos en 1
    }
}
```

**¿Qué hace?** Crea una nueva "caja" vacía lista para guardar tareas.

**¿Por qué devuelve un puntero `*InMemoryRepo`?** 
- Para que todos compartan la misma caja
- Si devolviera un valor, cada vez que lo uses se copiaría

**Ejemplo:**
```go
repo := NewInMemoryRepo()  // Crea la caja
// repo.tasks = [] (vacía)
// repo.nextID = 1
```

### Función Save()

```go
func (r *InMemoryRepo) Save(task *domain.Task) error {
    task.ID = r.nextID      // Le asigna un ID único
    r.nextID++              // Prepara el siguiente ID
    r.tasks = append(r.tasks, *task)  // Agrega la tarea a la lista
    return nil
}
```

**¿Qué hace?**
1. Le da un ID único a la tarea (1, 2, 3, ...)
2. Incrementa el contador para la próxima tarea
3. Agrega la tarea a la lista

**`(r *InMemoryRepo)`** significa que esta función "pertenece" a `InMemoryRepo`. Es como un método de la estructura.

**Ejemplo paso a paso:**
```go
repo := NewInMemoryRepo()
// repo.tasks = []
// repo.nextID = 1

tarea := &domain.Task{Title: "Aprender Go"}
repo.Save(tarea)
// tarea.ID = 1
// repo.nextID = 2
// repo.tasks = [tarea con ID=1]
```

### Función GetAll()

```go
func (r *InMemoryRepo) GetAll() ([]domain.Task, error) {
    return r.tasks, nil
}
```

**¿Qué hace?** Devuelve todas las tareas guardadas.

**Ejemplo:**
```go
tareas, _ := repo.GetAll()
// tareas = [tarea1, tarea2, tarea3, ...]
```

---

## 📄 3. SERVICE (service/task.go) - El Cerebro

**¿Qué hace?** Contiene la lógica de negocio. Decide qué se puede hacer y qué no.

### Estructura TaskService

```go
type TaskService struct {
    repo domain.TaskRepository  // Necesita un repositorio para guardar
}
```

**Analogía:** El servicio es como un "gerente" que necesita acceso al almacén (repositorio) para hacer su trabajo.

### Función NewTaskService()

```go
func NewTaskService(repo domain.TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}
```

**¿Qué hace?** Crea un nuevo servicio y le pasa un repositorio.

**¿Por qué recibe una interfaz y no un tipo específico?**
- Puede trabajar con cualquier repositorio (memoria, base de datos, etc.)
- Facilita las pruebas

**Ejemplo:**
```go
repo := repository.NewInMemoryRepo()
servicio := service.NewTaskService(repo)
// Ahora servicio tiene acceso al repositorio
```

### Función CreateTask()

```go
func (s *TaskService) CreateTask(title string) (*domain.Task, error) {
    // 1. Validación
    if title == "" {
        return nil, errors.New("el título no puede estar vacío")
    }

    // 2. Crear la tarea
    newTask := &domain.Task{
        Title:       title,
        IsCompleted: false,
    }

    // 3. Guardar en el repositorio
    err := s.repo.Save(newTask)
    if err != nil {
        return nil, err
    }

    // 4. Devolver la tarea creada
    return newTask, nil
}
```

**¿Qué hace?**
1. **Valida** que el título no esté vacío (regla de negocio)
2. **Crea** una nueva tarea con el título
3. **Guarda** la tarea usando el repositorio
4. **Devuelve** la tarea creada

**¿Por qué la validación está aquí y no en el handler?**
- La lógica de negocio debe estar en el servicio
- Si cambias la interfaz (móvil, web, API), la validación sigue funcionando

**Ejemplo:**
```go
servicio := service.NewTaskService(repo)

// Caso exitoso
tarea, err := servicio.CreateTask("Aprender Go")
// tarea = {ID: 1, Title: "Aprender Go", IsCompleted: false}
// err = nil

// Caso con error
tarea, err := servicio.CreateTask("")
// tarea = nil
// err = "el título no puede estar vacío"
```

### Función GetTasks()

```go
func (s *TaskService) GetTasks() ([]domain.Task, error) {
    return s.repo.GetAll()
}
```

**¿Qué hace?** Obtiene todas las tareas del repositorio.

**¿Por qué no va directo al repositorio desde el handler?**
- El servicio puede agregar lógica adicional (filtros, ordenamiento, etc.)
- Mantiene la separación de responsabilidades

---

## 📄 4. HANDLER (handler/http.go) - El Recepcionista

**¿Qué hace?** Recibe las peticiones HTTP, las convierte a datos de Go, llama al servicio, y devuelve respuestas HTTP.

### Estructura TaskHandler

```go
type TaskHandler struct {
    service *service.TaskService  // Necesita el servicio para hacer el trabajo
}
```

**Analogía:** Es como un recepcionista que recibe a los clientes (peticiones HTTP) y los dirige al servicio correcto.

### Función NewTaskHandler()

```go
func NewTaskHandler(s *service.TaskService) *TaskHandler {
    return &TaskHandler{service: s}
}
```

**¿Qué hace?** Crea un nuevo handler y le pasa un servicio.

### Función CreateTask()

```go
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    // 1. Definir estructura para leer JSON
    var body struct {
        Title string `json:"title"`
    }

    // 2. Leer el JSON del body
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }

    // 3. Llamar al servicio
    createdTask, err := h.service.CreateTask(body.Title)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 4. Responder con éxito
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdTask)
}
```

**¿Qué hace paso a paso?**

1. **Define una estructura temporal** para leer solo el campo `title` del JSON
2. **Decodifica el JSON** del body de la petición
   - Si hay error, responde 400 (Bad Request)
3. **Llama al servicio** para crear la tarea
   - Si hay error, responde 400 con el mensaje
4. **Responde con éxito** (201 Created) y envía la tarea creada en JSON

**Parámetros:**
- `w http.ResponseWriter`: Para escribir la respuesta HTTP
- `r *http.Request`: La petición HTTP entrante

**Ejemplo de flujo:**
```
Cliente envía: POST /tasks {"title": "Aprender Go"}
↓
Handler recibe la petición
↓
Handler lee: body.Title = "Aprender Go"
↓
Handler llama: servicio.CreateTask("Aprender Go")
↓
Servicio crea y guarda la tarea
↓
Handler responde: 201 Created {"id": 1, "title": "Aprender Go", ...}
```

### Función GetAllTasks()

```go
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    tasks, _ := h.service.GetTasks()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}
```

**¿Qué hace?**
1. Obtiene todas las tareas del servicio
2. Configura el header para indicar que es JSON
3. Convierte las tareas a JSON y las envía

**Nota:** El `_` ignora el error (no es la mejor práctica, pero funciona para este ejemplo).

---

## 📄 5. MAIN (cmd/api/main.go) - El Orquestador

**¿Qué hace?** Conecta todas las piezas y arranca el servidor.

### Función main()

```go
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
```

**¿Qué hace paso a paso?**

1. **Inyección de Dependencias:**
   - Crea el repositorio
   - Crea el servicio (pasándole el repositorio)
   - Crea el handler (pasándole el servicio)
   
   **Analogía:** Es como construir una casa: primero la base (repo), luego las paredes (service), y finalmente el techo (handler).

2. **Define las rutas:**
   - Cuando alguien hace POST a `/tasks` → llama a `CreateTask`
   - Cuando alguien hace GET a `/tasks` → llama a `GetAllTasks`
   - Cualquier otro método → responde 405 (Method Not Allowed)

3. **Arranca el servidor:**
   - Escucha en el puerto 8080
   - Espera peticiones HTTP

**Flujo completo de creación de tarea:**
```
1. Cliente → POST http://localhost:8080/tasks {"title": "Aprender Go"}
   ↓
2. main.go → Detecta POST en /tasks
   ↓
3. handler.CreateTask() → Lee el JSON
   ↓
4. service.CreateTask() → Valida y crea la tarea
   ↓
5. repository.Save() → Guarda en memoria
   ↓
6. Respuesta: 201 Created {"id": 1, "title": "Aprender Go", ...}
```

---

## 🔄 Diagrama de Comunicación

```
┌─────────────────────────────────────────────────────────────┐
│                    CLIENTE (Insomnia/Postman)                │
│                                                               │
│  POST /tasks {"title": "Aprender Go"}                        │
└───────────────────────────┬─────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                      MAIN (main.go)                         │
│                                                               │
│  1. Crea: repo → service → handler                          │
│  2. Registra ruta: /tasks                                   │
│  3. Escucha en puerto 8080                                  │
└───────────────────────────┬─────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    HANDLER (http.go)                        │
│                                                               │
│  CreateTask():                                              │
│    • Lee JSON del body                                      │
│    • Llama a service.CreateTask()                          │
│    • Responde con JSON                                      │
└───────────────────────────┬─────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    SERVICE (task.go)                        │
│                                                               │
│  CreateTask():                                              │
│    • Valida que title no esté vacío                         │
│    • Crea nueva Task                                        │
│    • Llama a repo.Save()                                    │
│    • Devuelve la tarea creada                               │
└───────────────────────────┬─────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                  REPOSITORY (memory.go)                     │
│                                                               │
│  Save():                                                    │
│    • Asigna ID único                                        │
│    • Guarda en slice []Task                                 │
│    • Devuelve error (si hay)                               │
└───────────────────────────┬─────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    DOMAIN (task.go)                         │
│                                                               │
│  • Define struct Task                                       │
│  • Define interface TaskRepository                          │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎓 Conceptos Clave Explicados

### 1. Inyección de Dependencias

**¿Qué es?** Pasar las dependencias (objetos que necesita) a una función o estructura en lugar de crearlas dentro.

**Ejemplo:**
```go
// ❌ MAL: Crea su propia dependencia
func NewTaskService() *TaskService {
    repo := repository.NewInMemoryRepo()  // Creado aquí
    return &TaskService{repo: repo}
}

// ✅ BIEN: Recibe la dependencia
func NewTaskService(repo domain.TaskRepository) *TaskService {
    return &TaskService{repo: repo}  // Recibida como parámetro
}
```

**Ventajas:**
- Fácil de probar (puedes pasar un repositorio "falso")
- Fácil de cambiar (puedes usar base de datos en lugar de memoria)

### 2. Interfaces

**¿Qué es?** Un contrato que dice "cualquier tipo que tenga estos métodos cumple con esta interfaz".

**Ejemplo:**
```go
// Interfaz: define QUÉ métodos debe tener
type TaskRepository interface {
    Save(task *Task) error
    GetAll() ([]Task, error)
}

// InMemoryRepo: implementa la interfaz
func (r *InMemoryRepo) Save(task *domain.Task) error { ... }
func (r *InMemoryRepo) GetAll() ([]domain.Task, error) { ... }

// Ahora InMemoryRepo "es un" TaskRepository
```

**Ventajas:**
- Puedes cambiar la implementación sin cambiar el código que la usa
- Facilita las pruebas

### 3. Punteros vs Valores

**Valor (`Task`):**
```go
tarea := domain.Task{Title: "Aprender Go"}
// tarea es una COPIA
```

**Puntero (`*Task`):**
```go
tarea := &domain.Task{Title: "Aprender Go"}
// tarea es una REFERENCIA (dirección) a la tarea
```

**¿Cuándo usar cada uno?**
- **Puntero**: Cuando quieres modificar el original o evitar copias grandes
- **Valor**: Cuando solo necesitas leer datos pequeños

### 4. Métodos de Estructura

```go
func (r *InMemoryRepo) Save(task *domain.Task) error {
    // Este método "pertenece" a InMemoryRepo
}
```

**Sintaxis:**
- `(r *InMemoryRepo)`: El "receptor" - la estructura a la que pertenece
- `r`: Nombre de la variable (puede ser cualquier nombre)
- `*InMemoryRepo`: Tipo (puntero a InMemoryRepo)

**Uso:**
```go
repo := NewInMemoryRepo()
repo.Save(tarea)  // Llama al método Save de repo
```

---

## 🧪 Ejemplo Completo: Crear una Tarea

Vamos a seguir una petición completa paso a paso:

### Paso 1: Cliente envía petición
```bash
POST http://localhost:8080/tasks
Content-Type: application/json

{"title": "Aprender Go"}
```

### Paso 2: main.go recibe la petición
```go
// main.go detecta POST en /tasks
http.HandleFunc("/tasks", func(...) {
    h.CreateTask(w, r)  // Llama al handler
})
```

### Paso 3: Handler procesa
```go
// handler/http.go
func (h *TaskHandler) CreateTask(w, r) {
    // Lee: body.Title = "Aprender Go"
    // Llama: h.service.CreateTask("Aprender Go")
}
```

### Paso 4: Service valida y crea
```go
// service/task.go
func (s *TaskService) CreateTask(title string) {
    // Valida: "Aprender Go" != "" ✅
    // Crea: newTask = {Title: "Aprender Go", IsCompleted: false}
    // Guarda: s.repo.Save(newTask)
}
```

### Paso 5: Repository guarda
```go
// repository/memory.go
func (r *InMemoryRepo) Save(task *domain.Task) {
    // Asigna: task.ID = 1
    // Guarda: r.tasks = [tarea con ID=1]
}
```

### Paso 6: Respuesta
```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": 1,
  "title": "Aprender Go",
  "is_completed": false,
  "created_at": "0001-01-01T00:00:00Z",
  "updated_at": "0001-01-01T00:00:00Z"
}
```

---

## 🎯 Resumen de Responsabilidades

| Capa | Responsabilidad | Ejemplo |
|------|----------------|---------|
| **Handler** | Recibir HTTP, convertir JSON | Leer `{"title": "..."}` del body |
| **Service** | Lógica de negocio, validaciones | Validar que title no esté vacío |
| **Repository** | Guardar/obtener datos | Guardar en memoria o base de datos |
| **Domain** | Definir estructuras e interfaces | Qué es una Task |

---

## 💡 Preguntas Frecuentes

**¿Por qué tantas capas?**
- Separación de responsabilidades: cada capa hace una cosa
- Fácil de cambiar: puedes cambiar el repositorio sin tocar el servicio
- Fácil de probar: puedes probar cada capa por separado

**¿Por qué usar interfaces?**
- Permite cambiar la implementación sin cambiar el código que la usa
- Facilita las pruebas (puedes crear implementaciones "falsas")

**¿Por qué punteros?**
- Eficiencia: no copias datos grandes
- Modificar originales: los cambios se reflejan en el original

**¿Qué pasa si quiero usar una base de datos?**
- Solo necesitas crear un nuevo repositorio que implemente `TaskRepository`
- El resto del código no cambia

---

¡Espero que esta explicación te haya ayudado a entender cómo funciona tu proyecto! 🚀


