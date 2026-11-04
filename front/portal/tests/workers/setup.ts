/// <reference types="@cloudflare/vitest-plugin/types" />
import { applyD1Migrations, type D1Migration } from "cloudflare:test";
import { env } from "cloudflare:workers";

// vitest.workers.config.ts reads migrations/ into TEST_MIGRATIONS; apply them
// to the test D1 once per run.
await applyD1Migrations(env.DB, (env as unknown as { TEST_MIGRATIONS: D1Migration[] }).TEST_MIGRATIONS);
