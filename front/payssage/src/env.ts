import { createEnv } from '@t3-oss/env-core'
import { z } from 'zod'

const supportedProvidersSchema = z
  .string()
  .default('mercadopago')
  .transform((value) => value.split(',').map((provider) => provider.trim()).filter(Boolean))
  .pipe(z.array(z.string().regex(/^[a-z0-9]+(?:_[a-z0-9]+)*$/, 'Provider must be a snake_case slug')).min(1))
  .transform((providers) => [...new Set(providers)])

export const env = createEnv({
  server: {
    SERVER_URL: z.url().optional(),
    IDENTITYX_ACCESS_API_KEY: z.string(),
  },

  /**
   * The prefix that client-side variables must have. This is enforced both at
   * a type-level and at runtime.
   */
  clientPrefix: 'VITE_',

  client: {
    VITE_POSTHOG_KEY: z.string(),
    VITE_POSTHOG_HOST: z.url().optional(),

    VITE_APP_TITLE: z.string().min(1).optional(),
    VITE_API_URL: z.url(),
    VITE_AUTH_API_URL: z.url(),
    VITE_TRIEOH_AUTH_PROJECT_ID: z.string(),

    VITE_SUPPORTED_PROVIDERS: supportedProvidersSchema,
  },

  /**
   * What object holds the environment variables at runtime. This is usually
   * `process.env` or `import.meta.env`.
   */
  runtimeEnv: {
    ...import.meta.env,
    SERVER_URL: process.env.SERVER_URL,
    IDENTITYX_ACCESS_API_KEY: process.env.IDENTITYX_ACCESS_API_KEY,
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
