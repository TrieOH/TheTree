import type { RuleStatus } from '@/shared/ui/form/types'
import type { ZodAny, ZodObject } from 'zod'

export function isFieldRequired<S extends Record<string, ZodAny>>(
  schema: ZodObject<S>,
  fieldName: keyof S & string
): boolean {
  const fieldSchema = schema.shape[fieldName]
  return !fieldSchema.safeParse(undefined).success
}

function getAllRuleMessages(fieldSchema: ZodAny): string[] {
  const probes: unknown[] = ['', 'a'.repeat(1000)]
  const messages = new Set<string>()

  for (const probe of probes) {
    const result = fieldSchema.safeParse(probe)
    if (!result.success)
      for (const issue of result.error.issues) messages.add(issue.message)
  }

  return Array.from(messages)
}

export function getFieldRulesStatus<S extends Record<string, ZodAny>>(
  schema: ZodObject<S>,
  fieldName: keyof S & string,
  value: unknown
): RuleStatus[] {
  const fieldSchema = schema.shape[fieldName]
  const allMessages = getAllRuleMessages(fieldSchema)
  if (allMessages.length === 0) return []

  const result = fieldSchema.safeParse(value)
  const failedMessages = result.success
    ? new Set<string>()
    : new Set(result.error.issues.map((i) => i.message))

  return allMessages.map((message) => ({
    message,
    passed: !failedMessages.has(message),
  }))
}