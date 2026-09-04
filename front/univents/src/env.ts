import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

export const env = createEnv({
  server: {
    SERVER_URL: z.url().optional(),
    IDENTITYX_ACCESS_API_KEY: z.string(),
    AUTH_SESSION_PASSWORD: z.string().min(32),

    STORAGE_IMAGE_ALLOWED_TYPES: z.string().optional(),
    STORAGE_IMAGE_MAX_SIZE_BYTES: z.string().optional(),
    STORAGE_IMAGE_UPLOAD_EXPIRES_SECONDS: z.string().optional(),
    STORAGE_IMAGE_MODERATION_MODEL: z.string().optional(),
    STORAGE_IMAGE_MODERATION_PROMPT: z.string().optional(),
    TRACES_OTLP_USER: z.string().optional(),
    TRACES_OTLP_PASSWORD: z.string().optional(),
    TRACES_OTLP_URL: z.string().optional(),
    TRACES_ENABLED: z.string().optional(),
  },

  /**
   * The prefix that client-side variables must have. This is enforced both at
   * a type-level and at runtime.
   */
  clientPrefix: "VITE_",

  client: {
    VITE_TRACING_ENABLED: z.coerce.boolean().default(false),
    VITE_POSTHOG_KEY: z.string(),
    VITE_POSTHOG_HOST: z.url().optional(),

    VITE_APP_TITLE: z.string().min(1).optional(),
    VITE_API_URL: z.url(),
    VITE_AUTH_API_URL: z.url(),
    VITE_STORAGE_URL: z.url(),
    VITE_AUTH_TRANSPORT: z.enum(["bff", "direct"]).default("bff"),
    VITE_TRIEOH_AUTH_PROJECT_ID: z.string(),

    VITE_UPLOAD_MAX_RETRIES: z.coerce.number().int().min(0).default(5),
    VITE_UPLOAD_RETRY_BASE_DELAY_MS: z.coerce
      .number()
      .int()
      .positive()
      .default(1000),
    VITE_UPLOAD_RETRY_MAX_DELAY_MS: z.coerce
      .number()
      .int()
      .positive()
      .default(30000),
  },

  runtimeEnv: {
    ...import.meta.env,
    VITE_TRACING_ENABLED: import.meta.env.VITE_TRACING_ENABLED,
    SERVER_URL: process.env.SERVER_URL,
    IDENTITYX_ACCESS_API_KEY: process.env.IDENTITYX_ACCESS_API_KEY,
    AUTH_SESSION_PASSWORD: process.env.AUTH_SESSION_PASSWORD,
    STORAGE_IMAGE_ALLOWED_TYPES: process.env.STORAGE_IMAGE_ALLOWED_TYPES,
    STORAGE_IMAGE_MAX_SIZE_BYTES: process.env.STORAGE_IMAGE_MAX_SIZE_BYTES,
    STORAGE_IMAGE_UPLOAD_EXPIRES_SECONDS:
      process.env.STORAGE_IMAGE_UPLOAD_EXPIRES_SECONDS,
    STORAGE_IMAGE_MODERATION_MODEL: process.env.STORAGE_IMAGE_MODERATION_MODEL,
    STORAGE_IMAGE_MODERATION_PROMPT:
      process.env.STORAGE_IMAGE_MODERATION_PROMPT,
    TRACES_OTLP_USER: process.env.TRACES_OTLP_USER,
    TRACES_OTLP_PASSWORD: process.env.TRACES_OTLP_PASSWORD,
    TRACES_OTLP_URL: process.env.TRACES_OTLP_URL,
    TRACES_ENABLED: process.env.TRACES_ENABLED,
  },
  onValidationError: (issues) => {
    console.error("Invalid or missing environment variables:");
    issues.forEach((issue) => {
      const path = issue.path?.map(String).join(".");
      console.error(`  → ${path}: ${issue.message}`);
    });
    process.exit(1);
  },
  onInvalidAccess: (key) => {
    console.error(
      `Attempted to access a server variable on the client: ${key}`,
    );
    throw new Error(`Invalid Access: ${key}`);
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
});
