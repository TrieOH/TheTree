import { defineConfig } from "orval"

const services = ["identityx", "informd", "payssage", "univents"] as const

export default defineConfig(
  Object.fromEntries(
    services.map((svc) => [
      svc,
      {
        input: {
          target: `api/${svc}/api-spec.yml`,
          override: {
            transformer: "./tools/orval/transformers/unwrap-envelope.ts",
          },
        },
        output: {
          target: `lib/ts/${svc}/client/endpoints.ts`,
          schemas: `lib/ts/${svc}/client/schemas`,
          client: "react-query",
          httpClient: "fetch",
          clean: true,
          override: {
            mutator: {
              path: "lib/ts/api-client/src/orval-mutator.ts",
              name: "customInstance",
            },
          },
        },
      },
    ]),
  ),
)
