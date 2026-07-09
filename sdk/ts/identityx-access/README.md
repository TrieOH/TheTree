# IdentityX Access SDK TS

TypeScript SDK focused on IdentityX API keys and project capabilities.

## Install

```bash
npm install @trieoh/identityx-access-sdk-ts
```

## Usage

```ts
import { createIdentityXAccessClient } from "@trieoh/identityx-access-sdk-ts";

const client = createIdentityXAccessClient({
  baseURL: "https://api.identityx.example",
  apiKey: "your-api-key",
});

const keys = await client.apiKeys.list("project-id");
const response = await client.apiKeys.create("project-id", {
  name: "Production key",
  capabilities: [],
  env: "prod",
});

const capabilities = await client.capabilities.list("project-id");
```
