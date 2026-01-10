# 5. API Reference

This document provides a reference for the GoAuth RESTful API endpoints.

## Authentication

These endpoints are for managing user identity and authentication.

### `POST /auth/register`

-   **Description:** Registers a new master user (a "client" user) in the system.
-   **Request Body:** `{"name": "...", "email": "...", "password": "..."}`
-   **Response:** `201 Created` with the new user's details.

### `POST /auth/login`

-   **Description:** Authenticates a master user and returns JWT access and refresh tokens.
-   **Request Body:** `{"email": "...", "password": "..."}`
-   **Response:** `200 OK` with `access_token` and `refresh_token`.

### `POST /auth/refresh`

-   **Description:** Issues a new set of tokens in exchange for a valid refresh token.
-   **Request Body:** `{"refresh_token": "..."}`
-   **Response:** `200 OK` with a new `access_token` and `refresh_token`.

### `POST /auth/logout`

-   **Description:** Revokes the current session's refresh token. Requires authentication.
-   **Authentication:** Requires a valid `access_token`.
-   **Response:** `204 No Content`.

### `GET /.well-known/jwks.json`

-   **Description:** Exposes the public key set (JWKS) for verifying the master JWTs issued by GoAuth. This is a public endpoint.

---

## Project-specific Authentication

These endpoints are for registering and authenticating users scoped to a specific project.

### `POST /projects/{project_id}/register`

-   **Description:** Registers a new "end user" within a specific project.
-   **Request Body:** `{"name": "...", "email": "...", "password": "..."}`
-   **Response:** `201 Created` with the new user's details.

### `POST /projects/{project_id}/login`

-   **Description:** Authenticates a user within a specific project and returns project-specific JWTs.
-   **Request Body:** `{"email": "...", "password": "..."}`
-   **Response:** `200 OK` with `access_token` and `refresh_token`.

---

## Session Management

These endpoints are for viewing and managing active user sessions. All endpoints require authentication.

### `GET /sessions`

-   **Description:** Lists all active sessions for the currently authenticated user.
-   **Response:** `200 OK` with a list of session details.

### `GET /sessions/me`

-   **Description:** Returns details about the current session (the one associated with the access token making the request).
-   **Response:** `200 OK` with session details.

### `DELETE /sessions/{session_id}`

-   **Description:** Revokes a specific session by its ID. Users can only revoke their own sessions.
-   **Response:** `204 No Content`.

### `DELETE /sessions/others`

-   **Description:** Revokes all sessions for the authenticated user *except* for the current one.
-   **Response:** `204 No Content`.

### `DELETE /sessions`

-   **Description:** Revokes all sessions for the authenticated user, including the current one.
-   **Response:** `204 No Content`.

---

## Project Management

These endpoints are for managing projects. They are restricted to "client" users (master users). All endpoints require authentication.

### `POST /projects`

-   **Description:** Creates a new project.
-   **Request Body:** `{"name": "..."}`
-   **Response:** `201 Created` with the new project's details.

### `GET /projects`

-   **Description:** Lists all projects owned by the authenticated client.
-   **Response:** `200 OK` with a list of projects.

### `GET /projects/{project_id}`

-   **Description:** Retrieves the details of a specific project.
-   **Response:** `200 OK` with project details.

### `PATCH /projects/{project_id}`

-   **Description:** Updates the details of a specific project (e.g., its name).
-   **Request Body:** `{"name": "..."}`
-   **Response:** `200 OK` with the updated project details.

### `DELETE /projects/{project_id}`

-   **Description:** Deletes a project.
-   **Response:** `204 No Content`.

### `GET /projects/{project_id}/.well-known/jwks.json`

-   **Description:** Exposes the project-specific public key set (JWKS) for verifying JWTs issued for this project. This is a public endpoint.
