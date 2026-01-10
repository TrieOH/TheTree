# 4. Architecture

As mentioned in the Project Structure guide, GoAuth is built upon the **Hexagonal Architecture** pattern (also known as "Ports and Adapters"). This is a powerful design pattern that helps create a loosely coupled, highly maintainable, and easily testable application.

## The Hexagon

The core idea is to separate the application into two main regions: the **"inside"** and the **"outside"**.

*   **The Inside (The Hexagon):** This is the core of the application. It contains the business logic and has no knowledge of the outside world. In GoAuth, this region is composed of the `domain` and `application` layers.

*   **The Outside:** This includes everything that the application interacts with but is not core to its business logic. Examples include the database, the web server, third-party APIs, and even the test suite. In GoAuth, this region is the `adapters` layer.

## Ports and Adapters

The "inside" and "outside" communicate with each other through **Ports** and **Adapters**.

*   **Ports:** These are interfaces defined by the "inside" (the application layer). They specify a contract for a certain type of interaction. For example, `internal/ports/outbound/user_repository.go` is a port that defines methods like `SaveUser` and `FindUserByEmail`. The application layer depends on this interface, not on any specific database implementation.

*   **Adapters:** These are concrete implementations of the ports. They live in the "outside" world.
    *   **Inbound Adapters (Driving Adapters):** These drive the application. The HTTP handler at `internal/adapters/http/auth.go` is an inbound adapter. It receives an HTTP request and calls a method on an application use case to execute business logic.
    *   **Outbound Adapters (Driven Adapters):** These are driven by the application. The persistence layer at `internal/adapters/persistence/user_repo.go` is an outbound adapter. It implements the `UserRepository` port and contains the actual database code to save a user.

This separation ensures that the core application logic can remain unchanged even if we decide to swap out external components. For example, we could switch from a REST API to a gRPC API by simply writing a new inbound adapter, with no changes to the `application` or `domain` layers. Similarly, we could switch from PostgreSQL to another database by writing a new persistence adapter.

## Anatomy of a Request: User Registration

To make this concrete, let's trace the journey of a "user registration" request through the different layers of GoAuth.

1.  **The Router (`adapters/http/router`)**
    The request first hits the HTTP router. A route is matched (e.g., `POST /api/v1/auth/register`), which is mapped to a specific handler function.

2.  **The HTTP Handler (`adapters/http/auth.go`)**
    The `Register` handler function is an **inbound adapter**. Its responsibilities are:
    *   Parse and validate the incoming HTTP request body into a Data Transfer Object (DTO), e.g., `dto.RegisterRequest`.
    *   Call the appropriate use case in the application layer, passing the necessary data (e.g., `app.AuthUsecase.Register(...)`).
    *   Receive the result (or an error) from the use case.
    *   Format the result into an HTTP response (e.g., a `201 Created` with the new user's ID).

3.  **The Use Case (`application/auth/usecase.go`)**
    The `Register` use case function contains the core business logic for this operation. It is part of the **inside** of the hexagon. Its responsibilities are:
    *   Check if a user with the given email already exists. To do this, it calls the `UserRepository` **port**.
    *   If the user doesn't exist, it creates a new `domain.User` entity.
    *   It hashes the user's password using a crypto utility.
    *   It calls the `UserRepository` **port** again to save the new user.
    *   It returns the newly created user's data.
    *   Crucially, the use case *does not know* it's talking to a PostgreSQL database. It only knows about the `UserRepository` interface.

4.  **The Repository Port (`ports/outbound/user_repository.go`)**
    This is the **port**—an interface that defines the contract for user persistence. The use case depends on this interface.

5.  **The Repository Implementation (`adapters/persistence/user_repo.go`)**
    This repository is an **outbound adapter**. It provides the concrete implementation for the `UserRepository` port.
    *   The `SaveUser` function in this struct takes the `domain.User` object.
    *   It converts the domain object into a format suitable for the database (using `sqlc`'s generated models).
    *   It executes the `INSERT` SQL command against the PostgreSQL database.
    *   It handles any database-specific errors and returns them.

This flow, from adapter to use case and back out to another adapter, is central to the design of GoAuth. It keeps the core logic pure and insulated from the details of infrastructure.
