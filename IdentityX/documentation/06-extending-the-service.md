# 6. Extending the Service: A Tutorial

This tutorial will guide you through the process of adding a new feature to GoAuth. We will add a new API endpoint to fetch a user's profile by their ID. This is a common requirement and serves as a perfect example of how to work with the hexagonal architecture.

**Our Goal:** Create a new authenticated endpoint `GET /users/{user_id}`.

The process can be broken down into these steps:
1.  Define the Application Use Case.
2.  Create a new HTTP Handler.
3.  Register the new route.
4.  Add an integration test.

For this tutorial, we are lucky: the persistence layer we need already exists! The `UserRepository` port (`internal/ports/outbound/user_repository.go`) already defines a `GetUserByID` method, and the persistence adapter (`internal/adapters/persistence/user_repo.go`) already implements it. We will just be building on top of that existing foundation.

---

## Step 1: Define the Application Use Case

Our first step is to define the business logic. We'll create a new use case for user-related operations.

1.  Create a new directory: `internal/application/user`.
2.  Inside this directory, create `contract.go`:

    ```go
    // file: internal/application/user/contract.go
    package user

    import (
        "GoAuth/internal/domain/user"
        "context"
        "github.com/google/uuid"
    )

    type Usecase interface {
        Get(ctx context.Context, userID uuid.UUID) (*user.User, error)
    }
    ```

3.  Now, create the use case implementation in `usecase.go`:

    ```go
    // file: internal/application/user/usecase.go
    package user

    import (
        "GoAuth/internal/domain/user"
        "GoAuth/internal/ports/outbound"
        "context"
        "github.com/google/uuid"
    )

    type usecase struct {
        userRepo outbound.UserRepository
    }

    func New(userRepo outbound.UserRepository) Usecase {
        return &usecase{userRepo: userRepo}
    }

    func (u *usecase) Get(ctx context.Context, userID uuid.UUID) (*user.User, error) {
        // The core logic is simple: just call the repository.
        // More complex logic, like authorization checks, would go here.
        return u.userRepo.GetUserByID(ctx, userID)
    }
    ```

---

## Step 2: Create the HTTP Handler

Next, we need an inbound adapter to handle HTTP requests. We will create a new handler for user-related endpoints.

1.  Create a new file: `internal/adapters/http/user.go`.
2.  Add the following code to create the handler:

    ```go
    // file: internal/adapters/http/user.go
    package http

    import (
        "GoAuth/internal/adapters/http/dto"
        "GoAuth/internal/application/user"
        "encoding/json"
        "net/http"

    	"github.com/go-chi/chi/v5"
        "github.com/google/uuid"
    )

    type UserHandler struct {
        userUC user.Usecase
    }

    func NewUserHandler(userUC user.Usecase) *UserHandler {
        return &UserHandler{userUC: userUC}
    }

    func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
        userIDStr := chi.URLParam(r, "user_id")
        userID, err := uuid.Parse(userIDStr)
        if err != nil {
            // In a real implementation, you'd use a proper error response helper.
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        foundUser, err := h.userUC.Get(r.Context(), userID)
        if err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        // We use a DTO to control what data is exposed in the API.
        response := dto.UserResponse{
            ID:        foundUser.ID,
            Email:     foundUser.Email,
            UserType:  foundUser.UserType,
            CreatedAt: foundUser.CreatedAt,
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(response)
    }
    ```
    *Note: We would also need to create the `dto.UserResponse` struct in a new file like `internal/adapters/http/dto/user.go`.*

---

## Step 3: Register the New Route

Now we need to initialize our new handler and register the route.

1.  Open `internal/adapters/http/router/register_routes.go`.
2.  In the `registerRoutes` function, instantiate the `UserUsecase` and `UserHandler`:

    ```go
    // Near the other use case initializations
    userUC := user.New(userRepo)

    // Near the other handler initializations
    userHandler := http2.NewUserHandler(userUC)
    ```

3.  Create a new function to register user-specific routes and call it from `registerRoutes`:

    ```go
    func registerRoutes(...) {
        // ...
        registerAuthRoutes(mux, authHandler, authMW)
        registerSessionRoutes(mux, sessionHandler, authMW)
        registerProjectRoutes(mux, projectHandler, authMW)
        registerUserRoutes(mux, userHandler, authMW) // Add this line

        return mux
    }

    func registerUserRoutes(
        r chi.Router,
        h *http2.UserHandler,
        authMW *middleware.AuthMiddleware,
    ) {
        r.Group(func(r chi.Router) {
            r.Use(authMW.Auth())
    
            r.Get("/users/{user_id}", h.GetUser)
        }
    }
    ```

---

## Step 4: Add an Integration Test

Finally, let's add a test to ensure our new endpoint works correctly.

1.  Create a new test file: `testing_framework/user_test.go`.
2.  Add a test case that:
    *   Creates a user by calling the register endpoint.
    *   Logs in as that user to get an auth token.
    *   Calls the new `GET /users/{user_id}` endpoint with the auth token.
    *   Asserts that the response is `200 OK` and the user data matches.

    ```go
    // file: testing_framework/user_test.go
    package testing_framework

    import (
        "testing"
        "github.com/gavv/httpexpect/v2"
    )

    func Test_GetUserByID(t *testing.T) {
        e, c, shutdown := New(t) // Assuming a test setup helper
        defer shutdown()

        // 1. Register a new user
        user := NewUserBuilder().Build()
        e.POST("/auth/register").
            WithJSON(user).
            Expect().
            Status(httpexpect.StatusCreated)

        // 2. Login to get tokens
        loginResp := e.POST("/auth/login").
            WithJSON(map[string]interface{}{"email": user.Email, "password": user.Password}).
            Expect().
            Status(httpexpect.StatusOK).
            JSON().Object()

        accessToken := loginResp.Value("access_token").String().Raw()
        userID := loginResp.Value("user").Object().Value("id").String().Raw()

        // 3. Call the new endpoint
        e.GET("/users/{user_id}", userID).
            WithHeader("Authorization", "Bearer "+accessToken).
            Expect().
            Status(httpexpect.StatusOK).
            JSON().Object().
            ValueEqual("id", userID).
            ValueEqual("email", user.Email)
    }
    ```

And that's it! You have successfully extended the service by adding a new, fully tested feature, following the clean architecture principles of the project.
