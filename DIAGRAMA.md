# 🔄 Diagrama de Flujo - Arquitectura Clean

## Flujo Completo: Crear una Tarea

```mermaid
sequenceDiagram
    participant Cliente as 👤 Cliente<br/>(Insomnia)
    participant Main as 🚪 Main<br/>(main.go)
    participant Handler as 📥 Handler<br/>(http.go)
    participant Service as 🧠 Service<br/>(task.go)
    participant Repo as 💾 Repository<br/>(memory.go)
    participant Domain as 📋 Domain<br/>(task.go)

    Cliente->>Main: POST /tasks<br/>{"title": "Aprender Go"}
    Main->>Handler: ruta /tasks detectada<br/>llama CreateTask()
    
    Handler->>Handler: Lee JSON del body<br/>body.Title = "Aprender Go"
    Handler->>Service: CreateTask("Aprender Go")
    
    Service->>Service: Valida: title != ""
    Service->>Domain: Crea nueva Task<br/>{Title: "Aprender Go"}
    Service->>Repo: Save(newTask)
    
    Repo->>Repo: Asigna ID = 1
    Repo->>Repo: Guarda en tasks[]
    Repo-->>Service: return nil (éxito)
    
    Service-->>Handler: return newTask, nil
    Handler->>Handler: Convierte a JSON
    Handler-->>Cliente: 201 Created<br/>{"id": 1, "title": "Aprender Go", ...}
```

## Flujo Completo: Obtener Todas las Tareas

```mermaid
sequenceDiagram
    participant Cliente as 👤 Cliente<br/>(Insomnia)
    participant Main as 🚪 Main<br/>(main.go)
    participant Handler as 📥 Handler<br/>(http.go)
    participant Service as 🧠 Service<br/>(task.go)
    participant Repo as 💾 Repository<br/>(memory.go)

    Cliente->>Main: GET /tasks
    Main->>Handler: ruta /tasks detectada<br/>llama GetAllTasks()
    
    Handler->>Service: GetTasks()
    Service->>Repo: GetAll()
    
    Repo-->>Service: return tasks[]
    Service-->>Handler: return tasks[]
    
    Handler->>Handler: Convierte a JSON
    Handler-->>Cliente: 200 OK<br/>[{"id": 1, ...}, {"id": 2, ...}]
```

## Arquitectura en Capas

```mermaid
graph TB
    subgraph "Capa Externa"
        Handler[📥 Handler<br/>http.go<br/>Recibe HTTP]
    end
    
    subgraph "Capa de Servicio"
        Service[🧠 Service<br/>task.go<br/>Lógica de Negocio]
    end
    
    subgraph "Capa de Datos"
        Repo[💾 Repository<br/>memory.go<br/>Guarda Datos]
    end
    
    subgraph "Capa de Dominio"
        Domain[📋 Domain<br/>task.go<br/>Estructuras e Interfaces]
    end
    
    Handler -->|usa| Service
    Service -->|usa| Repo
    Service -->|usa| Domain
    Repo -->|implementa| Domain
    
    style Handler fill:#e1f5ff
    style Service fill:#fff4e1
    style Repo fill:#e8f5e9
    style Domain fill:#f3e5f5
```

## Inyección de Dependencias

```mermaid
graph LR
    Main[main.go] -->|1. Crea| Repo[InMemoryRepo]
    Main -->|2. Crea con repo| Service[TaskService]
    Main -->|3. Crea con service| Handler[TaskHandler]
    
    Repo -.->|implementa| Interface[TaskRepository<br/>interface]
    Service -.->|usa| Interface
    
    style Main fill:#ffebee
    style Repo fill:#e8f5e9
    style Service fill:#fff4e1
    style Handler fill:#e1f5ff
    style Interface fill:#f3e5f5
```

## Estructura de Datos

```mermaid
classDiagram
    class Task {
        +int ID
        +string Title
        +bool IsCompleted
        +time.Time CreatedAt
        +time.Time UpdatedAt
    }
    
    class TaskRepository {
        <<interface>>
        +Save(task *Task) error
        +GetAll() []Task, error
    }
    
    class InMemoryRepo {
        -[]Task tasks
        -int nextID
        +Save(task *Task) error
        +GetAll() []Task, error
    }
    
    class TaskService {
        -TaskRepository repo
        +CreateTask(title string) *Task, error
        +GetTasks() []Task, error
    }
    
    class TaskHandler {
        -TaskService service
        +CreateTask(w, r) void
        +GetAllTasks(w, r) void
    }
    
    TaskRepository <|.. InMemoryRepo : implements
    TaskService --> TaskRepository : uses
    TaskHandler --> TaskService : uses
    TaskService --> Task : creates
    InMemoryRepo --> Task : stores
```

## Flujo de Datos: Crear Tarea

```
┌─────────────────────────────────────────────────────────────┐
│                    DATOS ENTRANTES                          │
│                                                               │
│  HTTP Request:                                               │
│  POST /tasks                                                 │
│  Content-Type: application/json                             │
│  Body: {"title": "Aprender Go"}                             │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    HANDLER (Conversión)                      │
│                                                               │
│  Input:  HTTP Request                                        │
│  Output: body.Title = "Aprender Go" (string)                │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    SERVICE (Lógica)                         │
│                                                               │
│  Input:  "Aprender Go" (string)                             │
│  Proceso:                                                    │
│    • Valida: title != "" ✅                                 │
│    • Crea: Task{Title: "Aprender Go", IsCompleted: false}   │
│  Output: *Task (puntero)                                    │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    REPOSITORY (Persistencia)                 │
│                                                               │
│  Input:  *Task (sin ID)                                      │
│  Proceso:                                                    │
│    • Asigna ID = 1                                          │
│    • Guarda en []Task                                        │
│  Output: Task guardada con ID                                │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    DATOS SALIENTES                          │
│                                                               │
│  HTTP Response:                                              │
│  Status: 201 Created                                         │
│  Content-Type: application/json                              │
│  Body: {                                                     │
│    "id": 1,                                                  │
│    "title": "Aprender Go",                                   │
│    "is_completed": false,                                    │
│    "created_at": "0001-01-01T00:00:00Z",                    │
│    "updated_at": "0001-01-01T00:00:00Z"                      │
│  }                                                           │
└─────────────────────────────────────────────────────────────┘
```

## Comparación: Con vs Sin Clean Architecture

### ❌ Sin Clean Architecture (Todo mezclado)
```
┌─────────────────────────────────────┐
│         Un solo archivo             │
│                                     │
│  • Recibe HTTP                      │
│  • Valida datos                     │
│  • Guarda en base de datos          │
│  • Responde HTTP                    │
│                                     │
│  Problemas:                         │
│  • Difícil de probar                │
│  • Difícil de cambiar               │
│  • Todo acoplado                    │
└─────────────────────────────────────┘
```

### ✅ Con Clean Architecture (Separado)
```
┌──────────────┐
│   Handler    │ ← Solo HTTP
└──────┬───────┘
       │
┌──────▼───────┐
│   Service    │ ← Solo lógica
└──────┬───────┘
       │
┌──────▼───────┐
│  Repository  │ ← Solo datos
└──────────────┘

Ventajas:
• Fácil de probar cada capa
• Fácil de cambiar implementación
• Código organizado y mantenible
```

---

## 📊 Resumen Visual

```
┌─────────────────────────────────────────────────────────────┐
│                         CLIENTE                              │
│                    (Insomnia/Postman)                        │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTP Request/Response
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  MAIN (main.go)                                             │
│  • Crea todas las dependencias                              │
│  • Registra rutas                                           │
│  • Arranca servidor                                         │
└───────────────────────────┬─────────────────────────────────┘
                            │
        ┌───────────────────┴───────────────────┐
        │                                       │
        ▼                                       ▼
┌───────────────┐                      ┌───────────────┐
│   HANDLER     │                      │   HANDLER     │
│  CreateTask   │                      │  GetAllTasks  │
└───────┬───────┘                      └───────┬───────┘
        │                                       │
        └───────────────┬───────────────────────┘
                        │
                        ▼
                ┌───────────────┐
                │    SERVICE    │
                │  TaskService  │
                └───────┬───────┘
                        │
                        ▼
                ┌───────────────┐
                │  REPOSITORY   │
                │ InMemoryRepo  │
                └───────┬───────┘
                        │
                        ▼
                ┌───────────────┐
                │    DOMAIN     │
                │  Task struct  │
                │  Interface    │
                └───────────────┘
```

---

**Nota:** Estos diagramas muestran cómo las diferentes capas se comunican entre sí. Cada flecha representa una llamada o dependencia entre componentes.


