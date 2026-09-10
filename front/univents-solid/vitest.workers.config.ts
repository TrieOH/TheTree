import { cloudflareTest } from "@cloudflare/vitest-plugin";
import { defineConfig } from "vitest/config";

// Tests for the Worker entrypoint (`src/server.ts`) and its handlers. They run
// inside workerd through Miniflare, so `env` is the real binding set from
// `wrangler.jsonc` instead of the `{} as never` the old jsdom test had to pass.
export default defineConfig({
  plugins: [
    cloudflareTest({
      // Never touch the network. `remoteBindings: false` turns the `AI` binding
      // (declared `remote: true` in wrangler.jsonc) into a local stub, so tests
      // stay offline, deterministic and can spy on `env.AI.run` per case.
      remoteBindings: false,
      wrangler: { configPath: "./wrangler.jsonc" },
      miniflare: {
        // Dummy S3 settings so `validateEnv` passes and the request reaches the
        // validation/moderation branches. Nothing ever leaves the isolate.
        bindings: {
          MINIO_ENDPOINT: "http://minio.test",
          BUCKET_NAME: "univents-test",
          MINIO_ACCESS_KEY: "test-access-key",
          MINIO_SECRET_KEY: "test-secret-key",
        },
      },
    }),
  ],
  resolve: {
    // Same `@/*` -> `./src/*` mapping the app config gets from tsconfig.json.
    tsconfigPaths: true,
  },
  test: {
    include: ["tests/workers/**/*.test.ts"],
  },
});
