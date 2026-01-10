# 2. Getting Started

This guide will walk you through setting up GoAuth for local development.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Docker and Docker Compose:** For running the application and its database in a containerized environment.
*   **OpenSSL:** For generating the cryptographic keys required for JWT signing.
*   **Git:** For cloning the project repository.

## 1. Clone the Repository

First, clone the GoAuth repository to your local machine:

```bash
git clone https://github.com/TrieOH/GoAuth.git
cd GoAuth
```

## 2. Configure Environment Variables

GoAuth uses environment variables for configuration. The project includes an example file that you can use as a template.

Copy the `.env.example` file to a new file named `.env`:

```bash
cp .env.example .env
```

The default values in the `.env` file are suitable for the standard local setup, so you don't need to change them to get started.

## 3. Generate Cryptographic Keys

The service uses an Ed25519 key pair to sign and verify JSON Web Tokens (JWTs). The `README.md` provides `openssl` commands to generate them.

First, create the `keys` directory if it doesn't exist:
```bash
mkdir -p keys
```

Then, generate the private and public keys:
```bash
openssl genpkey -algorithm ed25519 -out keys/ed25519-private.pem
openssl pkey -in keys/ed25519-private.pem -pubout -out keys/ed25519-public.pem
```
These commands will create `ed25519-private.pem` and `ed25519-public.pem` inside the `keys/` directory. The application is configured via `docker-compose.yml` to use them automatically.

## 4. Run the Service

The entire application stack (the Go service and a PostgreSQL database) can be launched using Docker Compose.

To build and start the services, run:

```bash
docker compose up --build
```

The GoAuth service will be running and accessible at `http://localhost:8080`. The PostgreSQL database will be available on port `5432`.

To stop the services, press `Ctrl+C` in the terminal and then run:
```bash
docker compose down
```

## 5. Run Tests

The project has a separate Docker Compose file for running the integration test suite. This ensures tests are run in a clean, isolated environment.

To run the tests, execute the following command:

```bash
docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from go-auth-test
```

After the tests have completed, it's important to clean up the test-specific containers and volumes to avoid conflicts. Run the following command to do so:

```bash
docker compose -f docker-compose.test.yml down -v
```
