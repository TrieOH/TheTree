# GoAuth

### Generate the keys

`cd keys`
`openssl genrsa -out private.pem 2048`
`openssl rsa -in private.pem -pubout -out public.pem`

## How to run tests
1. `docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from go-auth-test`
2. `docker compose -f docker-compose.test.yml down -v`

Always run`down -v` when running tests
