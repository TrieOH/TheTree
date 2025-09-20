# GoAuth

### Generate the keys

`cd keys`
`openssl genrsa -out private.pem 2048`
`openssl rsa -in private.pem -pubout -out public.pem`
