# Project Context: bs_server

This document provides an overview of the `bs_server` project, outlining its purpose, technical stack, architecture, and development conventions.

## Project Overview

The `bs_server` is the backend API for an image recognition system, specifically for butterfly recognition. It enables frontend interactions for training models and performing user recognition. A key feature is its support for configurable pipelines, allowing users to combine different feature extraction algorithms, clustering methods, and training strategies to generate corresponding trainers based on their configurations.

**Key Technologies:**
*   **Language:** Go (version 1.24)
*   **Web Framework:** Gin
*   **Image Processing:** GoCV (OpenCV bindings for Go)
*   **Database:** MongoDB
*   **Caching/Data Store:** Redis
*   **Authentication:** JWT (JSON Web Tokens)
*   **Configuration:** Viper
*   **API Documentation:** Swagger
*   **Distributed Tracing:** OpenTelemetry

**Architecture:**
The project follows a typical layered architecture:
*   **`api`**: Defines API routes.
*   **`internal/handler`**: Handles incoming API requests and interacts with services.
*   **`internal/service`**: Contains the core business logic.
*   **`app/pkg/data/imongo`**: Provides MongoDB client and data access functionalities.
*   **`app/pkg/gocv`**: Encapsulates image processing utilities using GoCV.
*   **`app/pkg/conf`**: Manages application configuration.

## Git Repository

*   The primary branch for this project is named **`master`**.
*   All Git commits should automatically follow the standard format for [Git Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

## GoCV Memory Management

As this project uses the `gocv` library, strict adherence to `gocv.Mat` memory management principles is critical due to its manual memory handling, which is not managed by the Go garbage collector.

1.  **Principle of Ownership:**
    *   The function that creates a `gocv.Mat` (e.g., via `gocv.IMRead`, `NewMat`, `Clone`, `IMDecode`) is solely responsible for its destruction. Always pair a creation with a `defer mat.Close()` call to guarantee cleanup.

2.  **Principle of Borrowing (Intra-Service):**
    *   When passing a `gocv.Mat` to another function within the same process, always use a pointer (`*gocv.Mat`). This signals that the callee is only "borrowing" the `Mat`.
    *   A function that receives a `*gocv.Mat` must **not** close it.
    *   If a function returns a new `gocv.Mat`, it transfers ownership to the caller. The caller is now responsible for closing it.

3.  **Principle of Serialization (Inter-Service):**
    *   Never attempt to serialize or send a `gocv.Mat` object directly across process boundaries (e.g., via Kafka or an API), as its internal pointer is meaningless elsewhere.
    *   The sending service must encode the `Mat` into a byte slice (e.g., JPEG/PNG). The receiving service must decode those bytes into a fresh `Mat` instance, thereby becoming its new owner and responsible for closing it.

## Building and Running

The project utilizes a `Makefile` for streamlined development tasks.

**Common Commands:**

*   **`make all`**: Builds the main application, all Go script tools, and generates Swagger documentation. This is the recommended command for a full setup.
*   **`make build`**: Compiles the main application (`app/main.go`) and places the executable at `scripts/bin/bs_server`.
*   **`make tools`**: Compiles various utility scripts (e.g., JWT token generation tool) and places them in `scripts/bin/`.
*   **`make swagger`**: Generates API documentation based on Swagger annotations in `app/main.go`. The output will be in `app/docs/swagger`.
*   **`make clean`**: Removes all generated build artifacts and documentation files.
*   **`make help`**: Displays a list of available `Makefile` targets and their descriptions.

**To Run the Application:**

1.  **Build:** Execute `make build` to compile the server.
2.  **Execute:** Run the generated binary: `./scripts/bin/bs_server`
    *   The server will listen on the port configured in `application.yaml` (defaulting to `localhost:8080`).

## Development Conventions

*   **Go Version:** The project requires Go 1.24.
*   **Dependency Management:** Go Modules (`go.mod`) are used for managing project dependencies.
*   **Code Linting:** `golangci-lint` is configured via `.golangci.yml` to enforce code quality and style.
*   **Code Formatting:** `.editorconfig` is used to ensure consistent code formatting across the project.
*   **API Documentation:** Swagger annotations are integrated directly into the Go source code (specifically `app/main.go`) for auto-generating API documentation.
*   **Tracing:** OpenTelemetry is implemented to provide distributed tracing capabilities for monitoring and debugging.
*   **Configuration:** Project configuration is managed using `Viper`, allowing for flexible external configuration.
