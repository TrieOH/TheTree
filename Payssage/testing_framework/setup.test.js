import { describe, test, beforeAll } from "vitest"
import {loginAs, post, put} from "./helpers.js"
import { CreateAPIKeySchema } from "./schemas/api_key.js"
import { WorkspaceSchema } from "./schemas/workspace.js"
import { validate } from "./helpers.js"
import { WebhookEndpointSchema } from "./schemas/webhook_endpoint.js"

const workspaceName = process.env.WORKSPACE_NAME || `trie-payments-test-${Math.random().toString(36).substring(2, 8)}`

let owner
beforeAll(async () => {
    owner = await loginAs(process.env.OWNER_EMAIL, process.env.OWNER_PASSWORD)
})

let workspace
describe("create workspace", () => {
    test("get or create workspace", async () => {
        const workspaces = await get(owner, "/workspaces")
        workspace = workspaces.find(w => w.name === workspaceName)
        if (!workspace) {
            workspace = await post(owner, "/workspaces", { name: workspaceName })
            console.log("\n✅ Workspace created:", workspace.name)
        } else {
            console.log("\n✅ Workspace found:", workspace.name)
        }
        console.log("   ID:", workspace.id)
    })
})

describe("create api keys", () => {
    test("create primary api key", async () => {
        await new Promise(resolve => setTimeout(resolve, 2000))
        const res = await post(owner, `/workspaces/${workspaceName}/keys`, { name: "test-key" })
        console.log("\n🔑 Primary API Key:", res.key)
    })

    test("create secondary api key", async () => {
        const res = await post(owner, `/workspaces/${workspaceName}/keys`, { name: "webhook-test-key" })
        console.log("🔑 Secondary API Key:", res.key)
    })
})

describe("sandbox setup", () => {
    test("enable sandbox", async () => {
        const ws = await post(owner, `/workspaces/${workspaceName}/sandbox/enable`)
        console.log("\n🧪 Sandbox enabled:", ws.sandbox)
    })

    test("register univents webhook endpoint", async () => {
        const endpoint = await post(owner, `/workspaces/${workspaceName}/webhooks`, {
            url: "http://univents:8080/webhooks/payments"
        })
        console.log("🪝 Webhook endpoint:", endpoint.url)
        console.log("🔐 Webhook secret:", endpoint.secret)
    })

    test("set marketplace config", async () => {
        const credentialId = process.env.MP_CREDENTIAL_ID
        if (!credentialId) {
            printCredentialInstructions(workspace.id)
            return
        }
        try {
            const config = await put(owner, `/workspaces/${workspaceName}/marketplace`, {
                credential_id: credentialId,
                fee_bps: 500,
            })
            console.log("💳 Marketplace config set, fee_bps:", config.fee_bps)
        } catch (e) {
            console.warn("\n⚠️  Failed to set marketplace config — credential may not exist in this workspace")
            printCredentialInstructions(workspace.id)
        }
    })
})

describe("print env", () => {
    test("print .env values to copy", async () => {
        console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        console.log("📋 Copy these to your .env:")
        console.log(`   WORKSPACE_NAME=${workspaceName}`)
        console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
    })
})

function printCredentialInstructions(workspaceId) {
    console.warn("\n📋 Run this in the payments DB:")
    console.warn(`
INSERT INTO provider_credentials (id, workspace_id, provider, display_name, credentials)
VALUES (
  uuidv7(),
  '${workspaceId}',
  'mercadopago',
  'MP Test Credential',
  '{"access_token": "TEST-your-access-token", "refresh_token": "", "provider_user_id": ""}'
);

SELECT id FROM provider_credentials WHERE workspace_id = '${workspaceId}';
    `)
    console.warn("   Then update .env:")
    console.warn("   MP_CREDENTIAL_ID=<id from SELECT above>")
    console.warn("\n   Then re-run: pnpm setup:workspace")
}