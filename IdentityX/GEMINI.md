# Instructions

## Overview
This is a Golang Identity Provider that's meant to run as a SaaS.

It uses sqlc for typesafe sql, and it uses the hexagonal/ports&adapters architecture

## Commands
- To run tests: `docker compose -f docker-compose.test.yml down -v && docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from go-auth-test --attach go-auth-test`
- To check for errors before running tests `go vet ./...` ignore `package GoAuth/docs is not in std` it should be the only error
- To generate sqlc code `~/go/bin/sqlc generate`

## Code Style
- Always adhere to ports&adapters/hexagonal architecture.
- When creating errors always wrap them in a service specific error with cause for better traceability, for example:
```go
err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password))
 err != nil {
	return nil, inbounds.ErrInvalidCredentials{Cause: err}
}
```

## Never
- Never try to run the `go test` command, it will not work, this project is containerized
