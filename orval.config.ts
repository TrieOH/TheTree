import { defineConfig } from "orval"

// The generated clients are framework-free on purpose: they expose plain
// `fetch`-based functions and types, so both the React and the Solid apps (and
// Workers/tests) can import them. Data-layer hooks live with each app's feature
// code, not here — see `tools/check-package-boundaries.mjs`.
const services = ["identityx", "informd", "payssage", "univents"] as const

export default defineConfig(
  Object.fromEntries(
    services.map((svc) => [
      svc,
      {
        input: {
          target: `api/${svc}/api-spec.yml`,
          override: {
            transformer: "./lib/ts/orval/transformers/unwrap-envelope.ts",
          },
        },
        output: {
          target: `lib/ts/${svc}/client/endpoints.ts`,
          schemas: `lib/ts/${svc}/client/schemas`,
          client: "fetch",
          httpClient: "fetch",
          clean: true,
          override: {
            mutator: {
              path: "lib/ts/api-client/src/orval-mutator.ts",
              name: "customInstance",
            },
          },
        },
        // Drops type imports the generator emits but never uses — see the
        // script for why it is needed.
        hooks: {
          afterAllFilesWrite: ["node tools/orval-prune-unused-imports.mjs"],
        },
      },
    ]),
  ),
)
