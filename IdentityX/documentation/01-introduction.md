# 1. Introduction to GoAuth

Welcome to GoAuth!

GoAuth is a robust, high-performance authentication and authorization service built with Go. It is designed to serve as a centralized identity provider for your applications, handling user registration, login, session management, and access control.

## Core Features

*   **User Authentication:** Securely manage user accounts with password hashing and token-based authentication.
*   **JWT-based Authorization:** Utilizes JSON Web Tokens (JWTs) with public-key cryptography (Ed25519) to issue and validate access tokens.
*   **Project-based Multi-tenancy:** Allows users to be partitioned into different "projects," providing a level of isolation for different applications or services that rely on GoAuth.
*   **Session Management:** Keeps track of active user sessions and provides mechanisms for token refreshing and revocation.
*   **RESTful API:** Exposes a clean and simple RESTful API for easy integration.

## Who is this for?

This service is ideal for developers who need a standalone, reliable authentication solution that can be easily integrated into a microservices architecture or as the auth backbone for a monolithic application.

## About this Documentation

This documentation is a hands-on guide for developers. It will walk you through the process of setting up the project, understanding its architecture, using its API, and extending it with new functionality.

---

## Table of Contents

1.  **Introduction** (You are here)
2.  [Getting Started](./02-getting-started.md)
3.  [Project Structure](./03-project-structure.md)
4.  [Architecture](./04-architecture.md)
5.  [API Reference](./05-api-reference.md)
6.  [Extending the Service: A Tutorial](./06-extending-the-service.md)
