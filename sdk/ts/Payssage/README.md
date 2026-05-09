# Payssage SDK TS

Official TypeScript SDK for Payssage

## Installation

```bash
npm install @trieoh/payssage-sdk-ts
# or
yarn add @trieoh/payssage-sdk-ts
# or
bun add @trieoh/payssage-sdk-ts
```

## Configuration

The SDK can be configured using environment variables or by passing parameters to the client factory (if you were to use one directly, though the SDK currently exports pre-configured services).

### Environment Variables

Copy `.env.example` to `.env` and fill in your values:

```bash
TRIEOH_PAY_BASE_URL=https://api.triepayments.trieoh.com
TRIEOH_PAY_SECRET_KEY=your_secret_key_here
```

## Usage

### Workspace Service

Connect or disconnect payment providers to your workspace.

```typescript
import { workspaceService } from "@trieoh/payssage-sdk-ts";

// Connect a provider (e.g., Mercado Pago)
const response = await workspaceService.connectProvider(
  "my-workspace",
  "mercadopago",
  {
    final_redirect_url: "https://your-app.com/success",
    provider_redirect_url: "https://your-app.com/callback"
  }
);

if (response.data) {
  console.log("Connect URL:", response.data.url);
}

// Disconnect a provider
const disconnectResponse = await workspaceService.disconnectProvider(
  "my-workspace",
  "credential-id"
);
```

### Intent Service

Manage payment intents.

```typescript
import { intentService } from "@trieoh/payssage-sdk-ts";

// Get all intents for a workspace
const response = await intentService.getWorkspaceIntents("my-workspace");

if (response.data) {
  console.log("Intents:", response.data);
}
```

## Development

This project uses [Bun](https://bun.sh).

```bash
# Install dependencies
bun install

# Build the project
bun run build
```
