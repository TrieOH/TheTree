import assert from "node:assert/strict"
import { test } from "node:test"

import { envelopeFetch, toEnvelope } from "./envelope.ts"

const enveloped = {
  code: 200,
  message: "",
  module: "IdentityX",
  timestamp: "2026-10-29T00:00:00Z",
  data: { id: "a1" },
}

test("an envelope passes through", () => {
  assert.deepEqual(toEnvelope(enveloped, 200), enveloped)
})

test("a bare payload becomes data, with the HTTP status as code", () => {
  assert.deepEqual(toEnvelope({ id: "a1" }, 201), { code: 201, data: { id: "a1" } })
  assert.deepEqual(toEnvelope({ items: [] }, 200), { code: 200, data: { items: [] } })
})

test("no body has no data", () => {
  assert.deepEqual(toEnvelope(undefined, 204), { code: 204 })
})

test("a bare error object becomes the failure envelope", () => {
  const error = {
    code: "VALIDATION_ERROR",
    message: "invalid request body",
    fields: [{ field: "name", message: "this field is required" }],
  }
  assert.deepEqual(toEnvelope(error, 400), {
    code: 400,
    message: "invalid request body",
    error_id: "VALIDATION_ERROR",
    error,
  })
})

test("a problem becomes the failure envelope", () => {
  const problem = {
    type: "about:blank",
    status: 400,
    title: "Bad Request",
    detail: "invalid request body",
    code: "VALIDATION_ERROR",
    errors: [{ field: "name", message: "this field is required" }],
  }
  assert.deepEqual(toEnvelope(problem, 400, "application/problem+json"), {
    code: 400,
    message: "invalid request body",
    error_id: "VALIDATION_ERROR",
    error: {
      code: "VALIDATION_ERROR",
      message: "invalid request body",
      fields: [{ field: "name", message: "this field is required" }],
    },
  })
})

test("envelopeFetch rewrites bare responses and keeps the real status", async () => {
  const fetchImpl = async (input: RequestInfo | URL) =>
    String(input) === "/empty"
      ? new Response(null, { status: 204 })
      : Response.json({ id: "a1" }, { status: 201 })
  const wrapped = envelopeFetch(fetchImpl as typeof fetch)

  const created = await wrapped("/actors")
  assert.equal(created.status, 201)
  assert.deepEqual(await created.json(), { code: 201, data: { id: "a1" } })

  const empty = await wrapped("/empty")
  assert.equal(empty.status, 200)
  assert.deepEqual(await empty.json(), { code: 204 })
})

test("envelopeFetch returns an enveloped response untouched", async () => {
  const original = Response.json(enveloped)
  const wrapped = envelopeFetch((async () => original) as typeof fetch)
  assert.equal(await wrapped("/actors"), original)
})
