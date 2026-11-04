import type { ServiceId } from "@/features/catalog/services";

// Each API's spec as generated for its release (tools/problemgen writes
// api/<svc>/api-spec.gen.yml with the problems bundled in; CI generates it
// before building every front). Read from the repo, never from the live
// /docs/openapi.yml, which would wake the API. One lazy chunk per spec.
const specs = import.meta.glob<string>("../../../../../api/*/api-spec.gen.yml", {
  query: "?raw",
  import: "default",
});

export function loadSpec(service: ServiceId): Promise<string> | null {
  const loader = specs[`../../../../../api/${service}/api-spec.gen.yml`];
  return loader ? loader() : null;
}
