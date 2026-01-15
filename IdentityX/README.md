# GoAuth

## Requirements
1. Postgres 18+

## Invariants
1. Only apierr package is allowed to wrap errors

### Generate the keys

`mkdir keys && cd keys`
`openssl genpkey -algorithm ed25519 -out ed25519-private.pem`
`openssl pkey -in ed25519-private.pem -pubout -out ed25519-public.pem`

## How to run tests
1. `docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from go-auth-test`
2. `docker compose -f docker-compose.test.yml down -v`

Always run`down -v` when running tests