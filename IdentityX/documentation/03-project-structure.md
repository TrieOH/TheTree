# 3. Project Structure

GoAuth is organized following the principles of **Hexagonal Architecture** (also known as Ports and Adapters). This architectural style isolates the core application logic from outside concerns (like databases, HTTP handlers, etc.), making the codebase clean, modular, and easy to test and maintain.

The project root contains configuration files like `Dockerfile`, `docker-compose.yml`, `go.mod`, and the main entrypoint (`main.go`). The core logic is housed entirely within the `internal/` directory.

Here is a breakdown of the key directories:

```text
/
├── documentation/      # All project documentation.
├── internal/           # All private application code resides here.
│   ├── adapters/       # Concrete implementations (adapters) for ports.
│   │   ├── http/       # The HTTP API layer (handlers, router, middleware).
│   │   └── persistence/# The database layer (repositories).
│   ├── application/    # Contains the application's use cases (business logic).
│   ├── apierr/         # Custom error handling framework.
│   ├── database/       # Database migrations and SQL queries.
│   ├── domain/         # Core business models and entities.
│   ├── infrastructure/ # Foundational services like telemetry and logging setup.
│   ├── ports/          # Interfaces (ports) for outbound communication.
│   └── utils/          # Shared utility functions.
├── keys/               # Stores cryptographic keys (ignored by Git).
└── testing_framework/  # The integration testing suite.
```

## Detailed Directory Descriptions

### `internal/`
This is the heart of the application. The `internal` package in Go prevents this code from being imported by other projects, enforcing its role as non-reusable, application-specific code.

### `internal/domain`
This directory contains the core business entities of the application, such as `User`, `Project`, and `Session`. These are plain Go structs with no dependencies on external frameworks. They represent the "what" of the application, not the "how".

### `internal/application`
This directory holds the application's use cases. Each use case orchestrates the flow of data, calling domain objects and using interfaces (ports) to interact with external systems. For example, `application/auth/usecase.go` contains the logic for user registration and login.

### `internal/ports`
This directory defines the interfaces (the "ports") that the application layer uses to communicate with the outside world. A key example is `ports/outbound/user_repository.go`, which defines the contract for how the application stores and retrieves user data, without knowing the specific database implementation.

### `internal/adapters`
This directory provides the concrete implementations (the "adapters") for the ports defined in `internal/ports`.
*   **`adapters/http/`**: This is the primary "inbound" adapter. It handles incoming HTTP requests, validates them, calls the appropriate application use case, and formats the response. It contains the router, middleware, and request/response DTOs (Data Transfer Objects).
*   **`adapters/persistence/`**: This is a primary "outbound" adapter. It implements the repository interfaces from `internal/ports` using a specific database technology (in this case, PostgreSQL with `sqlc`).

### `internal/database`
This directory is dedicated to database concerns.
*   **`migrations/`**: Contains the SQL scripts for creating and altering the database schema, managed by `goose`.
*   **`queries/`**: Contains the raw SQL queries that `sqlc` uses to generate type-safe Go code for database access.

### `internal/apierr`
A custom package for creating structured, consistent API errors throughout the application.

### `testing_framework/`
Contains a complete integration test suite. This framework acts as an external client to the API, making real HTTP requests to a test instance of the application and asserting the responses. This provides a high degree of confidence that the system works end-to-end.
