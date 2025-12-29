# Go Auth Test Framework Documentation

A fluent, declarative testing framework for HTTP APIs built on top of `httpexpect`. This framework provides a clean, readable way to write integration tests with minimal boilerplate.

## Table of Contents

1. [Core Concepts](#core-concepts)
2. [Quick Start](#quick-start)
3. [Architecture](#architecture)
4. [Usage Guide](#usage-guide)
5. [Extending the Framework](#extending-the-framework)
6. [Best Practices](#best-practices)
7. [Common Patterns](#common-patterns)

---

## Core Concepts

### Design Philosophy

The framework is built around three key principles:

1. **Fluent Interface**: Chain methods together for readable test code
2. **Declarative Testing**: Express *what* you want to test, not *how*
3. **Type Safety**: Leverage Go's type system to catch errors at compile time

### Key Components

```
TestSuite → Client → RequestBuilder → Response
     ↓         ↓
   User    AuthContext
```

- **TestSuite**: Manages test environment (server, database)
- **Client**: Makes HTTP requests with optional authentication
- **RequestBuilder**: Constructs HTTP requests fluently
- **Response**: Validates responses with chainable assertions
- **User**: High-level abstraction for user workflows (register, login, etc.)
- **AuthContext**: Stores authentication tokens
 
---

## Quick Start

### Basic Test Structure

```go
func TestMyFeature(t *testing.T) {
    suite := NewTestSuite(t)
    
    t.Run("CreateResource", func(t *testing.T) {
        client := suite.Client(t)
        
        client.POST("/resources").
            WithBody(map[string]string{
                "name": "Test Resource",
            }).
            Expect(http.StatusCreated).
            Success("my-module", "Resource created")
    })
}
```

### With Authentication

```go
func TestAuthenticatedEndpoint(t *testing.T) {
    suite := NewTestSuite(t)
    client := suite.Client(t)
    
    // Create and login user
    user := client.User("test@example.com", "Password123!").
        Register().
        Login()
    
    // Make authenticated request
    user.AuthedClient().GET("/protected").
        Expect(http.StatusOK).
        Success("my-module", "Success")
}
```

---

## Architecture

### 1. TestSuite

The `TestSuite` manages your test environment lifecycle.

```go
type TestSuite struct {
    Server *httptest.Server
    DB     *sql.DB
    t      *testing.T
}
```

**Responsibilities:**
- Initializes test server and database
- Provides client factory method
- Handles cleanup automatically

**Usage:**
```go
suite := NewTestSuite(t)  // Setup happens automatically
// Cleanup registered with t.Cleanup()
```

### 2. Client

The `Client` is your main interface for making HTTP requests.

```go
type Client struct {
    expect *httpexpect.Expect
    t      *testing.T
    auth   *AuthContext
}
```

**Methods:**
- `POST(path)`, `GET(path)`, `PATCH(path)`, `DELETE(path)`: Create requests
- `Auth(ctx)`: Return authenticated client
- `User(email, password)`: Create user workflow helper

**Important:** Always create clients with the **current subtest's** `*testing.T`:

```go
t.Run("Subtest", func(t *testing.T) {
    client := suite.Client(t)  // Use subtest's t!
})
```

### 3. RequestBuilder

Fluently construct HTTP requests before sending.

```go
type RequestBuilder struct {
    req *httpexpect.Request
    t   *testing.T
}
```

**Methods:**
- `WithBody(body)`: Add JSON body
- `Expect(status)`: Send request and expect status code
- `ExpectStatus(status)`: Alias for `Expect`

**Example:**
```go
client.POST("/api/endpoint").
    WithBody(map[string]interface{}{
        "field1": "value1",
        "field2": 123,
    }).
    Expect(http.StatusOK)
```

### 4. Response

Chain assertions on HTTP responses.

```go
type Response struct {
    resp   *httpexpect.Response
    t      *testing.T
    status int
}
```

**Methods:**
- `Success(module, message)`: Assert successful response structure
- `Error(module, message)`: Assert error response structure
- `ValidationError(errors...)`: Assert validation error with specific messages
- `Data()`: Get data object from response
- `DataArray()`: Get data array from response
- `Cookies()`: Extract authentication cookies
- `JSON()`: Get raw JSON object

**Example:**
```go
resp := client.GET("/users/123").
    Expect(http.StatusOK).
    Success("users", "User found")

data := resp.Data()
data.Value("email").String().IsEqual("test@example.com")
```

### 5. User

High-level abstraction for user authentication workflows.

```go
type User struct {
    Email    string
    Password string
    client   *Client
    auth     *AuthContext
    t        *testing.T
}
```

**Methods:**
- `Register()`: Register user
- `Login()`: Login and store tokens
- `Logout()`: Logout user
- `Refresh()`: Refresh authentication tokens
- `AuthedClient()`: Get authenticated client

**Example:**
```go
user := client.User("test@example.com", "Pass123!").
    Register().
    Login()

// Now make authenticated requests
user.AuthedClient().GET("/profile").
    Expect(http.StatusOK)
```

---

## Usage Guide

### Making Requests

#### Simple GET Request

```go
client.GET("/users").
    Expect(http.StatusOK)
```

#### POST with Body

```go
client.POST("/users").
    WithBody(map[string]interface{}{
        "name": "John Doe",
        "email": "john@example.com",
        "age": 30,
    }).
    Expect(http.StatusCreated)
```

#### With Authentication

```go
// Method 1: Via User helper
user := client.User("test@example.com", "Pass123!").Login()
user.AuthedClient().GET("/protected").Expect(http.StatusOK)

// Method 2: Manual auth context
authClient := client.Auth(&AuthContext{
    AccessToken: "token123",
    RefreshToken: "refresh456",
})
authClient.GET("/protected").Expect(http.StatusOK)
```

### Response Assertions

#### Standard Success Response

Expects response format:
```json
{
    "module": "my-module",
    "message": "Success message",
    "code": 200,
    "data": { ... }
}
```

```go
client.GET("/resource").
    Expect(http.StatusOK).
    Success("my-module", "Resource retrieved")
```

#### Validation Errors

Expects response format:
```json
{
    "module": "validation",
    "message": "Validation failed",
    "trace": ["error1", "error2"]
}
```

```go
client.POST("/users").
    WithBody(map[string]string{"email": "invalid"}).
    Expect(http.StatusBadRequest).
    ValidationError("valid email address", "password is required")
```

#### Working with Response Data

```go
// Object data
data := client.GET("/user/123").
    Expect(http.StatusOK).
    Data()

data.Value("name").String().IsEqual("John Doe")
data.Value("age").Number().IsEqual(30)

// Array data
arr := client.GET("/users").
    Expect(http.StatusOK).
    DataArray()

arr.Length().IsEqual(5)
arr.Value(0).Object().Value("email").String().Contains("@")
```

### User Workflows

#### Complete Registration and Login Flow

```go
user := client.User("test@example.com", "SecurePass123!").
    Register().
    Login()

// User now has auth tokens stored
user.AuthedClient().GET("/profile").
    Expect(http.StatusOK)
```

#### Testing Token Refresh

```go
user := client.User("test@example.com", "Pass123!").
    Register().
    Login()

oldToken := user.auth.AccessToken

user.Refresh()

if oldToken == user.auth.AccessToken {
    t.Error("Token should have changed")
}
```

#### Multiple Sessions

```go
user := client.User("test@example.com", "Pass123!").Register()

// Create multiple sessions
session1 := suite.Client(t).User(user.Email, user.Password).Login()
session2 := suite.Client(t).User(user.Email, user.Password).Login()
session3 := suite.Client(t).User(user.Email, user.Password).Login()

// Each has independent auth tokens
```

---

## Extending the Framework

### Adding New Request Methods

Add methods to the `Client` struct:

```go
func (c *Client) PUT(path string) *RequestBuilder {
    return c.newRequest("PUT", path)
}

func (c *Client) OPTIONS(path string) *RequestBuilder {
    return c.newRequest("OPTIONS", path)
}
```

### Adding Custom Assertions

Extend the `Response` struct:

```go
// Assert paginated response format
func (r *Response) Pagination(page, perPage, total int) *Response {
    r.t.Helper()
    obj := r.resp.JSON().Object()
    
    meta := obj.Value("meta").Object()
    meta.Value("page").Number().IsEqual(page)
    meta.Value("per_page").Number().IsEqual(perPage)
    meta.Value("total").Number().IsEqual(total)
    
    return r
}

// Usage:
client.GET("/users?page=2").
    Expect(http.StatusOK).
    Pagination(2, 20, 100)
```

### Creating Domain-Specific Helpers

Create helpers for your specific domain:

```go
// Project helper
type Project struct {
    ID     string
    Name   string
    client *Client
    t      *testing.T
}

func (c *Client) Project(name string) *Project {
    return &Project{
        Name:   name,
        client: c,
        t:      c.t,
    }
}

func (p *Project) Create() *Project {
    p.t.Helper()
    
    data := p.client.POST("/projects").
        WithBody(map[string]string{"name": p.Name}).
        Expect(http.StatusCreated).
        Data()
    
    p.ID = data.Value("id").String().Raw()
    return p
}

func (p *Project) Delete() *Project {
    p.t.Helper()
    
    p.client.DELETE("/projects/" + p.ID).
        Expect(http.StatusOK)
    
    return p
}

// Usage:
project := client.Project("My Project").Create()
defer project.Delete()
```

### Adding Custom Request Headers

Extend `RequestBuilder`:

```go
func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
    rb.req = rb.req.WithHeader(key, value)
    return rb
}

func (rb *RequestBuilder) WithBearerToken(token string) *RequestBuilder {
    return rb.WithHeader("Authorization", "Bearer "+token)
}

// Usage:
client.GET("/api/data").
    WithBearerToken("my-token").
    Expect(http.StatusOK)
```

### Creating Test Data Builders

Define reusable test data:

```go
type UserBuilder struct {
    email    string
    password string
    name     string
    age      int
}

func NewUserBuilder() *UserBuilder {
    return &UserBuilder{
        email:    "test@example.com",
        password: "DefaultPass123!",
        name:     "Test User",
        age:      25,
    }
}

func (ub *UserBuilder) WithEmail(email string) *UserBuilder {
    ub.email = email
    return ub
}

func (ub *UserBuilder) WithAge(age int) *UserBuilder {
    ub.age = age
    return ub
}

func (ub *UserBuilder) Build() map[string]interface{} {
    return map[string]interface{}{
        "email":    ub.email,
        "password": ub.password,
        "name":     ub.name,
        "age":      ub.age,
    }
}

// Usage:
userData := NewUserBuilder().
    WithEmail("custom@example.com").
    WithAge(30).
    Build()

client.POST("/users").
    WithBody(userData).
    Expect(http.StatusCreated)
```

---

## Best Practices

### 1. Always Use the Correct `*testing.T`

**❌ Wrong:**
```go
func testFeature(t *testing.T, suite *TestSuite) {
    client := suite.Client(t)  // Parent t
    
    t.Run("Subtest", func(t *testing.T) {
        // Still using parent's client!
        client.GET("/endpoint").Expect(http.StatusOK)
    })
}
```

**✅ Correct:**
```go
func testFeature(t *testing.T, suite *TestSuite) {
    t.Run("Subtest", func(t *testing.T) {
        client := suite.Client(t)  // Subtest's t
        client.GET("/endpoint").Expect(http.StatusOK)
    })
}
```

### 2. Use `t.Helper()` in Custom Helpers

```go
func (u *User) DoSomething() *User {
    u.t.Helper()  // Makes stack traces point to caller
    // ... test logic
    return u
}
```

### 3. Chain for Readability

```go
// Good: Clear workflow
user := client.User("test@example.com", "Pass123!").
    Register().
    Login()

// Bad: Unnecessary intermediate variables
user := client.User("test@example.com", "Pass123!")
user = user.Register()
user = user.Login()
```

### 4. Use Declarative Test Specs for Repetitive Tests

```go
type TestCase struct {
    Name     string
    Input    map[string]string
    Expected int
    Errors   []string
}

var testCases = []TestCase{
    {
        Name:     "InvalidEmail",
        Input:    map[string]string{"email": "invalid"},
        Expected: http.StatusBadRequest,
        Errors:   []string{"valid email"},
    },
    // More cases...
}

for _, tc := range testCases {
    tc := tc
    t.Run(tc.Name, func(t *testing.T) {
        client := suite.Client(t)
        client.POST("/endpoint").
            WithBody(tc.Input).
            Expect(tc.Expected).
            ValidationError(tc.Errors...)
    })
}
```

### 5. Organize Tests by Feature

```go
func TestGoAuth(t *testing.T) {
    suite := NewTestSuite(t)
    
    t.Run("Authentication", func(t *testing.T) {
        testRegister(t, suite)
        testLogin(t, suite)
        testLogout(t, suite)
    })
    
    t.Run("Sessions", func(t *testing.T) {
        testListSessions(t, suite)
        testRevokeSessions(t, suite)
    })
}
```

### 6. Clean Up Resources

```go
t.Run("CreateAndDelete", func(t *testing.T) {
    client := suite.Client(t)
    
    // Create resource
    data := client.POST("/resources").
        WithBody(map[string]string{"name": "Test"}).
        Expect(http.StatusCreated).
        Data()
    
    resourceID := data.Value("id").String().Raw()
    
    // Clean up
    defer func() {
        client.DELETE("/resources/" + resourceID).
            Expect(http.StatusOK)
    }()
    
    // ... test logic
})
```

---

## Common Patterns

### Testing Error Cases

```go
t.Run("NotFound", func(t *testing.T) {
    client := suite.Client(t)
    client.GET("/users/nonexistent").
        Expect(http.StatusNotFound).
        Error("users", "User not found")
})

t.Run("Unauthorized", func(t *testing.T) {
    client := suite.Client(t)
    client.GET("/protected").
        Expect(http.StatusUnauthorized).
        Error("auth", "Authentication required")
})
```

### Testing with Multiple Users

```go
t.Run("UserInteraction", func(t *testing.T) {
    client := suite.Client(t)
    
    alice := client.User("alice@example.com", "Pass123!").
        Register().
        Login()
    
    bob := client.User("bob@example.com", "Pass123!").
        Register().
        Login()
    
    // Alice creates resource
    data := alice.AuthedClient().POST("/resources").
        WithBody(map[string]string{"name": "Alice's Resource"}).
        Expect(http.StatusCreated).
        Data()
    
    resourceID := data.Value("id").String().Raw()
    
    // Bob can't access Alice's resource
    bob.AuthedClient().GET("/resources/" + resourceID).
        Expect(http.StatusForbidden)
})
```

### Testing Pagination

```go
t.Run("Pagination", func(t *testing.T) {
    client := suite.Client(t)
    user := client.User("test@example.com", "Pass123!").
        Register().
        Login()
    
    // Create 25 items
    for i := 0; i < 25; i++ {
        user.AuthedClient().POST("/items").
            WithBody(map[string]string{"name": fmt.Sprintf("Item %d", i)}).
            Expect(http.StatusCreated)
    }
    
    // Test first page
    arr := user.AuthedClient().GET("/items?page=1&per_page=10").
        Expect(http.StatusOK).
        DataArray()
    
    arr.Length().IsEqual(10)
    
    // Test second page
    arr = user.AuthedClient().GET("/items?page=2&per_page=10").
        Expect(http.StatusOK).
        DataArray()
    
    arr.Length().IsEqual(10)
})
```

### Testing Concurrency

```go
t.Run("ConcurrentAccess", func(t *testing.T) {
    client := suite.Client(t)
    user := client.User("test@example.com", "Pass123!").
        Register().
        Login()
    
    var wg sync.WaitGroup
    errors := make(chan error, 10)
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            resp := user.AuthedClient().POST("/items").
                WithBody(map[string]string{
                    "name": fmt.Sprintf("Concurrent Item %d", id),
                }).
                Expect(http.StatusCreated)
            
            // Check for any assertion failures
            // (they'll panic in goroutines)
        }(i)
    }
    
    wg.Wait()
    close(errors)
})
```

### Testing Data Relationships

```go
t.Run("UserWithProjects", func(t *testing.T) {
    client := suite.Client(t)
    user := client.User("test@example.com", "Pass123!").
        Register().
        Login()
    
    // Create projects
    project1ID := user.AuthedClient().POST("/projects").
        WithBody(map[string]string{"name": "Project 1"}).
        Expect(http.StatusCreated).
        Data().
        Value("id").String().Raw()
    
    project2ID := user.AuthedClient().POST("/projects").
        WithBody(map[string]string{"name": "Project 2"}).
        Expect(http.StatusCreated).
        Data().
        Value("id").String().Raw()
    
    // Verify user's projects
    arr := user.AuthedClient().GET("/projects").
        Expect(http.StatusOK).
        DataArray()
    
    arr.Length().IsEqual(2)
    
    // Verify each project belongs to user
    data := user.AuthedClient().GET("/projects/" + project1ID).
        Expect(http.StatusOK).
        Data()
    
    data.Value("owner_id").String().NotEmpty()
})
```

---

## Troubleshooting

### Tests Pass When They Should Fail

**Problem:** Using parent test's `*testing.T` instead of subtest's.

**Solution:**
```go
t.Run("Subtest", func(t *testing.T) {
    client := suite.Client(t)  // Use THIS t, not parent's
})
```

### Nil Pointer Errors

**Problem:** Trying to use `User.auth` before calling `Login()`.

**Solution:**
```go
user := client.User("test@example.com", "Pass123!")
user.Register()  // No auth tokens yet
user.Login()     // Now auth tokens are set
user.AuthedClient()  // Safe to use
```

### Authentication Not Working

**Problem:** Cookies not being sent correctly.

**Solution:** Ensure you're using `AuthedClient()` or `Auth()`:
```go
// Wrong
user.client.GET("/protected")  // No auth!

// Correct
user.AuthedClient().GET("/protected")
```

---

## Summary

This framework provides a clean, expressive way to write integration tests:

- **Fluent API** makes tests read like specifications
- **Type safety** catches errors at compile time
- **Extensible** architecture lets you add domain-specific helpers
- **Proper test isolation** with correct `*testing.T` usage

Start simple, extend as needed, and keep your tests readable!