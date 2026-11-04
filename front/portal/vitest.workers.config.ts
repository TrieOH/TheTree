import { cloudflareTest } from "@cloudflare/vitest-plugin";
import { defineConfig } from "vitest/config";

// Tests for the Worker entrypoint (`src/server.ts`). They run inside workerd
// through Miniflare, so `env` is the real binding set from `wrangler.jsonc`,
// with a local D1 that has the migrations applied.
export default defineConfig({
  plugins: [
    cloudflareTest({
      remoteBindings: false,
      wrangler: { configPath: "./wrangler.jsonc" },
    }),
  ],
  resolve: {
    tsconfigPaths: true,
  },
  test: {
    include: ["tests/workers/**/*.test.ts"],
  },
});
