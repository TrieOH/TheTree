# Sandbox Setup — Manual Steps

The test suite handles most setup automatically, but provider credentials
require a manual DB insert since MP OAuth is not yet available.

## 1. Insert provider credential
```sql
INSERT INTO provider_credentials (id, workspace_id, provider, display_name, credentials)
VALUES (
  uuidv7(),
  '<workspace_id printed in test output>',
  'mercadopago',
  'MP Test Credential',
  '{"access_token": "TEST-your-access-token", "refresh_token": "", "provider_user_id": ""}'
);

SELECT id FROM provider_credentials WHERE workspace_id = '<workspace_id>';
```

## 2. Set MP_CREDENTIAL_ID in .env
```
MP_CREDENTIAL_ID=<credential id from above>
```

Then re-run tests — the `set marketplace config` test will pick it up.

## 3. Once MP OAuth is implemented

The manual steps above will be replaced by `POST /workspaces/{name}/providers/mercadopago/setup`.