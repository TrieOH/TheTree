import { createEnv } from '@t3-oss/env-core'
import { z } from 'zod'

export const env = createEnv({
  server: {
    SERVER_URL: z.url().optional(),
    TRIEOH_AUTHZED_TOKEN: z.string(),
    TRIEOH_AUTHZED_ENDPOINT: z.string(),
  },

  /**
   * The prefix that client-side variables must have. This is enforced both at
   * a type-level and at runtime.
   */
  clientPrefix: 'VITE_',

  client: {
    VITE_APP_TITLE: z.string().min(1).optional(),
    VITE_POSTHOG_KEY: z.string(),
    VITE_POSTHOG_HOST: z.url().optional(),
  },
  runtimeEnv: {
    ...import.meta.env,
    SERVER_URL: process.env.SERVER_URL,
    TRIEOH_AUTHZED_TOKEN: process.env.TRIEOH_AUTHZED_TOKEN,
    TRIEOH_AUTHZED_ENDPOINT: process.env.TRIEOH_AUTHZED_ENDPOINT
  },
  onValidationError: (issues) => {
    console.error("Invalid or missing environment variables:")
    issues.forEach((issue) => {
      const path = issue.path?.map(String).join(".")
      console.error(`  → ${path}: ${issue.message}`)
    })
    process.exit(1)
  },
  onInvalidAccess: (key) => {
    console.error(`Attempted to access a server variable on the client: ${key}`)
    throw new Error(`Invalid Access: ${key}`)
  },

  /**
   * By default, this library will feed the environment variables directly to
   * the Zod validator.
   *
   * This means that if you have an empty string for a value that is supposed
   * to be a number (e.g. `PORT=` in a ".env" file), Zod will incorrectly flag
   * it as a type mismatch violation. Additionally, if you have an empty string
   * for a value that is supposed to be a string with a default value (e.g.
   * `DOMAIN=` in an ".env" file), the default value will never be applied.
   *
   * In order to solve these issues, we recommend that all new projects
   * explicitly specify this option as true.
   */
  emptyStringAsUndefined: true,
})
