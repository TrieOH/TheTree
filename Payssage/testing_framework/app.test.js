import { describe, test, beforeAll, expect } from "vitest"
import { loginAs, post, get, deleteReq, validate, createAPIKeyClient } from "./helpers.js"
import { createWorkspace } from "./fixtures/workspaces/create.js"
import { createAPIKey } from "./fixtures/api_keys/create.js"
import { createIntent, createIntentNoMetadata } from "./fixtures/intents/create.js"
import { WorkspaceSchema } from "./schemas/workspace.js"
import { APIKeySchema, CreateAPIKeySchema } from "./schemas/api_key.js"
import { IntentSchema } from "./schemas/intent.js"

let owner
beforeAll(async () => {
    owner = await loginAs(process.env.OWNER_EMAIL, process.env.OWNER_PASSWORD)
})

let workspace
describe("workspaces", () => {
    test("create workspace", async () => {
        workspace = await post(owner, "/workspaces", createWorkspace)
        expect(validate(WorkspaceSchema, workspace)).toBe(true)
        expect(workspace.name).toBe(createWorkspace.name)
    })

    test("list workspaces", async () => {
        const workspaces = await get(owner, "/workspaces")
        expect(Array.isArray(workspaces)).toBe(true)
        expect(workspaces.some(w => w.id === workspace.id)).toBe(true)
    })
})

let apiKey
let rawKey
let apiKeyClient
describe("api keys", () => {
    beforeAll(async () => {
        await new Promise(resolve => setTimeout(resolve, 2000))
    })

    test("create api key", async () => {
        const res = await post(owner, `/workspaces/${workspace.name}/keys`, createAPIKey)
        expect(validate(CreateAPIKeySchema, res)).toBe(true)
        expect(res.key).toMatch(/^tp_/)
        apiKey = res
        rawKey = res.key
        apiKeyClient = createAPIKeyClient(rawKey)
    })

    test("list api keys", async () => {
        const keys = await get(owner, `/workspaces/${workspace.name}/keys`)
        expect(Array.isArray(keys)).toBe(true)
        expect(keys.some(k => k.id === apiKey.id)).toBe(true)
    })
})

let intent
describe("intents", () => {
    test("create intent", async () => {
        intent = await post(apiKeyClient, "/intents", createIntent)
        expect(validate(IntentSchema, intent)).toBe(true)
        expect(intent.status).toBe("pending")
        expect(intent.amount).toBe(createIntent.amount)
        expect(intent.currency).toBe(createIntent.currency)
    })

    test("create intent without metadata", async () => {
        const i = await post(apiKeyClient, "/intents", createIntentNoMetadata)
        expect(validate(IntentSchema, i)).toBe(true)
        expect(i.status).toBe("pending")
    })

    test("get intent by id", async () => {
        const fetched = await get(apiKeyClient, `/intents/${intent.id}`)
        expect(validate(IntentSchema, fetched)).toBe(true)
        expect(fetched.id).toBe(intent.id)
    })

    test("list intents via api key", async () => {
        const intents = await get(apiKeyClient, "/intents")
        expect(Array.isArray(intents)).toBe(true)
        expect(intents.some(i => i.id === intent.id)).toBe(true)
    })

    test("list intents via user session", async () => {
        const intents = await get(owner, "/intents")
        expect(Array.isArray(intents)).toBe(true)
        expect(intents.some(i => i.id === intent.id)).toBe(true)
    })

    test("cancel intent", async () => {
        const cancelled = await post(apiKeyClient, `/intents/${intent.id}/cancel`)
        expect(validate(IntentSchema, cancelled)).toBe(true)
        expect(cancelled.status).toBe("cancelled")
    })

    test("cancel already cancelled intent fails", async () => {
        try {
            await post(apiKeyClient, `/intents/${intent.id}/cancel`)
            expect(true).toBe(false) // should not reach here
        } catch (e) {
            expect(e.response.status).toBe(404)
        }
    })
})

describe("api key revocation", () => {
    test("revoke api key", async () => {
        await deleteReq(owner, `/workspaces/${workspace.name}/keys/${apiKey.id}`)
    })

    test("revoked key cannot create intents", async () => {
        try {
            await post(apiKeyClient, "/intents", createIntent)
            expect(true).toBe(false)
        } catch (e) {
            expect(e.response.status).toBe(401)
        }
    })
})