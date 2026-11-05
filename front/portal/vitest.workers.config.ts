import { cloudflareTest, readD1Migrations } from "@cloudflare/vitest-plugin";
import { defineConfig } from "vitest/config";

// Tests for the Worker entrypoint (`src/server.ts`). They run inside workerd
// through Miniflare, so `env` is the real binding set from `wrangler.jsonc`,
// with a local D1 that has the migrations applied.
export default defineConfig(async () => ({
  plugins: [
    cloudflareTest({
      remoteBindings: false,
      wrangler: { configPath: "./wrangler.jsonc" },
      miniflare: {
        bindings: {
          LINEAR_API_KEY: "lin_api_test",
          LINEAR_WEBHOOK_SECRET: "whsec_test",
          GITHUB_WEBHOOK_SECRET: "ghsec_test",
          DEPLOY_HOOK_SECRET: "deploy_test",
          CLOUDFLARE_ANALYTICS_TOKEN: "cf_test",
          NEON_API_KEY: "neon_test",
          GRAFANA_READ_TOKEN: "glc_test",
          // Applied to the local D1 by tests/workers/setup.ts.
          TEST_MIGRATIONS: await readD1Migrations("./migrations"),
        },
      },
    }),
  ],
  resolve: {
    tsconfigPaths: true,
  },
  test: {
    include: ["tests/workers/**/*.test.ts"],
    setupFiles: ["./tests/workers/setup.ts"],
  },
}));
