/**
 * Orval input transformer: unwrap the `fun.Response` envelope.
 *
 * Every successful JSON response in the TrieOH specs is shaped
 * `allOf: [Response, { properties: { data: <payload> } }]`. The generated
 * clients and react-query hooks should expose the payload directly — the
 * envelope is unwrapped at runtime by the orval mutator
 * (lib/ts/api-client/src/orval-mutator.ts) — so this transformer replaces
 * each enveloped response schema with its `data` schema.
 */
import { defineTransformer } from "orval"

const RESPONSE_REF = "#/components/schemas/Response"

export default defineTransformer((spec) => {
  for (const pathItem of Object.values(spec.paths ?? {})) {
    for (const operation of Object.values(pathItem)) {
      if (!operation || typeof operation !== "object" || !("responses" in operation)) {
        continue
      }
      const responses = (operation as { responses?: Record<string, unknown> }).responses
      if (!responses) continue
      for (const [code, response] of Object.entries(responses)) {
        if (!/^[23]\d\d$/.test(code)) continue
        const content = (response as { content?: Record<string, { schema?: unknown }> })?.content
        const json = content?.["application/json"]
        const schema = json?.schema
        if (!schema) continue
        const unwrapped = unwrap(schema)
        if (unwrapped !== undefined) json.schema = unwrapped
      }
    }
  }
  return spec
})

function unwrap(schema: unknown): unknown | undefined {
  if (!schema || typeof schema !== "object") return undefined
  const s = schema as { allOf?: unknown[]; $ref?: string }

  if (Array.isArray(s.allOf)) {
    const enveloped = s.allOf.some(
      (member) =>
        member && typeof member === "object" && (member as { $ref?: string }).$ref === RESPONSE_REF,
    )
    if (!enveloped) return undefined
    const dataMember = s.allOf.find(
      (member) =>
        member &&
        typeof member === "object" &&
        (member as { properties?: unknown }).properties &&
        (member as { properties?: { data?: unknown } }).properties?.data !== undefined,
    )
    const data = dataMember
      ? (dataMember as { properties: { data?: unknown } }).properties.data
      : undefined
    return data ?? {}
  }

  if (s.$ref === RESPONSE_REF) return {}
  return undefined
}
