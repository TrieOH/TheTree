# IdentityX Access SDK TS

TypeScript SDK focused on IdentityX API keys and project capabilities.

## Install

```bash
pnpm add @trieoh/identityx-access-sdk-ts
```

## Usage

```ts
import { createIdentityXAccessClient } from "@trieoh/identityx-access-sdk-ts";

const client = createIdentityXAccessClient({
  baseURL: "https://api.identityx.example",
  apiKey: "your-api-key",
});

const response = await client.apiKeys.create("project-id", {
  name: "Production key",
  capabilities: [],
  env: "prod",
});

const capabilities = await client.capabilities.list("project-id");

const actors = await client.actors.list("project-id");
const actor = await client.actors.getById("project-id", "actor-id");
const actorByEmail = await client.actors.getByEmail(
  "project-id",
  "person@example.com",
);
const createdActor = await client.actors.create("project-id", {
  auth_method: "password",
  type: "human",
  email: "person@example.com",
});
```
