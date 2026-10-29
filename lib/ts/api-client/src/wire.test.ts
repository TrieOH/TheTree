import assert from "node:assert/strict"
import { test } from "node:test"

import { itemsOf, readFailure, toFailure, toSuccess } from "./wire.ts"

test("a success body is the payload", () => {
  assert.deepEqual(toSuccess({ id: "a1" }).data, { id: "a1" })
})

test("a 204 has no data", () => {
  assert.equal(toSuccess(undefined).data, undefined)
})

test("an error object is the failure, with the HTTP status as code", () => {
  const failure = toFailure({ code: "VALIDATION_ERROR", message: "invalid request body" }, 400)
  assert.equal(failure.code, 400)
  assert.equal(failure.error_id, "VALIDATION_ERROR")
  assert.equal(failure.message, "invalid request body")
})

test("an error object lists its fields", () => {
  const fields = [{ field: "name", message: "this field is required" }]
  assert.deepEqual(readFailure({ code: "VALIDATION_ERROR", message: "invalid", fields }).fields, fields)
})

test("a problem reads its detail, code and errors", () => {
  const errors = [{ field: "name", message: "this field is required" }]
  assert.deepEqual(
    readFailure({ type: "about:blank", status: 400, title: "Bad Request", detail: "invalid", code: "VALIDATION_ERROR", errors }),
    { code: "VALIDATION_ERROR", message: "invalid", fields: errors },
  )
  assert.equal(readFailure({ status: 404, title: "Not Found" }).message, "Not Found")
})

test("a body that is no error still fails with its status", () => {
  const failure = toFailure(undefined, 502)
  assert.equal(failure.code, 502)
  assert.equal(failure.error_id, "UNKNOWN_ERROR")
})

test("itemsOf reads { items }", () => {
  assert.deepEqual(itemsOf({ items: [1, 2] }), [1, 2])
  assert.deepEqual(itemsOf(undefined), [])
})
