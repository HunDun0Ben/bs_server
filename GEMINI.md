# Project Context: bs_server

This document provides an overview of the `bs_server` project, outlining its purpose, technical stack, architecture, and development conventions.

## Project Overview

The `bs_server` is the backend API for an image recognition system, specifically for butterfly recognition. It enables frontend interactions for training models and performing user recognition. A key feature is its support for configurable pipelines, allowing users to combine different feature extraction algorithms, clustering methods, and training strategies to generate corresponding trainers based on their configurations.

**Key Technologies:**

- **Language:** Go (version 1.24)
- **Web Framework:** Gin
- **Image Processing:** GoCV (OpenCV bindings for Go)
- **Database:** MongoDB
- **Caching/Data Store:** Redis
- **Authentication:** JWT (JSON Web Tokens)
- **Configuration:** Viper
- **API Documentation:** Swagger
- **Distributed Tracing:** OpenTelemetry

**Architecture:**
The project follows a typical layered architecture:

- **`api`**: Defines API routes.
- **`internal/handler`**: Handles incoming API requests and interacts with services.
- **`internal/service`**: Contains the core business logic.
- **`app/pkg/data/imongo`**: Provides MongoDB client and data access functionalities.
- **`app/pkg/gocv`**: Encapsulates image processing utilities using GoCV.
- **`app/pkg/conf`**: Manages application configuration.

## Git Repository

- The primary branch for this project is named **`master`**.
- All Git commits should automatically follow the standard format for [Git Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

## GoCV Memory Management

As this project uses the `gocv` library, strict adherence to `gocv.Mat` memory management principles is critical due to its manual memory handling, which is not managed by the Go garbage collector.

1.  **Principle of Ownership:**
    - The function that creates a `gocv.Mat` (e.g., via `gocv.IMRead`, `NewMat`, `Clone`, `IMDecode`) is solely responsible for its destruction. Always pair a creation with a `defer mat.Close()` call to guarantee cleanup.

2.  **Principle of Borrowing (Intra-Service):**
    - When passing a `gocv.Mat` to another function within the same process, always use a pointer (`*gocv.Mat`). This signals that the callee is only "borrowing" the `Mat`.
    - A function that receives a `*gocv.Mat` must **not** close it.
    - If a function returns a new `gocv.Mat`, it transfers ownership to the caller. The caller is now responsible for closing it.

3.  **Principle of Serialization (Inter-Service):**
    - Never attempt to serialize or send a `gocv.Mat` object directly across process boundaries (e.g., via Kafka or an API), as its internal pointer is meaningless elsewhere.
    - The sending service must encode the `Mat` into a byte slice (e.g., JPEG/PNG). The receiving service must decode those bytes into a fresh `Mat` instance, thereby becoming its new owner and responsible for closing it.

## Building and Running

The project utilizes a `Makefile` for streamlined development tasks.

**Common Commands:**

- **`make all`**: Builds the main application, all Go script tools, and generates Swagger documentation. This is the recommended command for a full setup.
- **`make build`**: Compiles the main application (`app/main.go`) and places the executable at `scripts/bin/bs_server`.
- **`make tools`**: Compiles various utility scripts (e.g., JWT token generation tool) and places them in `scripts/bin/`.
- **`make swagger`**: Generates API documentation based on Swagger annotations in `app/main.go`. The output will be in `app/docs/swagger`.
- **`make clean`**: Removes all generated build artifacts and documentation files.
- **`make help`**: Displays a list of available `Makefile` targets and their descriptions.

**To Run the Application:**

1.  **Build:** Execute `make build` to compile the server.
2.  **Execute:** Run the generated binary: `./scripts/bin/bs_server`
    - The server will listen on the port configured in `application.yaml` (defaulting to `localhost:8080`).

## Development Conventions

- **Go Version:** The project requires Go 1.24.
- **Dependency Management:** Go Modules (`go.mod`) are used for managing project dependencies.
- **Code Linting && Formatting:** Maintain high code quality and style standards by using the project's automation tools. Always run `make format` after making modifications to ensure consistency across the codebase. Linting is enforced via project configurations like `golangci-lint`, `prettier`.
- **API Documentation:** Swagger annotations are embedded throughout the Go source code, including `app/main.go` and individual handler files. Any changes to API handlers must be reflected in the corresponding Swagger annotations in their respective source files. Always regenerate the Swagger documentation (make swag) after modifying API definitions to ensure documentation consistency.
- **Tracing:** OpenTelemetry is implemented to provide distributed tracing capabilities for monitoring and debugging.
- **Configuration:** Project configuration is managed using `Viper`, allowing for flexible external configuration.

## Architectural Standards

To ensure maintainability, testability, and scalability, the following architectural standards are mandatory:

1.  **Layered Responsibility:**
    - **Handler Layer (`internal/handler`)**: Responsible for HTTP request parsing, input validation, and mapping responses. It must NOT contain business logic or database queries. It interacts only with the Service layer via interfaces.
    - **Service Layer (`internal/service`)**: Contains the core business logic. It orchestrates domain objects and interacts with the Repository layer via interfaces. Services must be decoupled from transport protocols (like HTTP).
    - **Repository Layer (`internal/repository`)**: Handles data persistence (e.g., MongoDB, Redis). It abstracts the underlying storage details from the Service layer.

2.  **Interface-Based Design:**
    - All cross-layer and cross-package communication MUST be done through interfaces.
    - Interfaces should be defined in the package that provides the implementation.

3.  **Dependency Injection (DI):**
    - Use **Constructor Injection**. Every struct representing a handler, service, or repository must have a `New...` constructor that accepts its dependencies as interfaces.
    - Avoid package-level global variables for stateful objects (e.g., database clients should be passed, not accessed globally).
    - Do NOT instantiate dependencies inside methods using `New...` or by calling package-level functions.

4.  **Handler Implementation:**
    - Handlers must be implemented as methods on a struct (e.g., `type UserHandler struct { ... }`). Standalone functions for handlers are prohibited to allow for proper dependency injection and testing.

5.  **Error Handling:**
    - Use the project's `AppError` for consistent error reporting across all layers.
