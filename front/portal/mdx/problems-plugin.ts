import { existsSync, readFileSync, readdirSync } from "node:fs";
import { fileURLToPath } from "node:url";
import type { Plugin } from "vite";
import { parse } from "yaml";

// `virtual:problems`: the RFC 9457 problems (ADR-0019) every service shares,
// from api/problems.yml, and each service's own, from api/<svc>/problems.yml
// and its features' spec.yml (ADR-0021), read at build time so the pages
// ship JSON, not a YAML parser.

const SERVICES = ["identityx", "payssage", "univents", "informd"] as const;
const ID = "virtual:problems";
const RESOLVED = `\0${ID}`;
const api = (path: string) => fileURLToPath(new URL(`../../../api/${path}`, import.meta.url));

export interface ProblemMember {
  name: string;
  /** JSON Schema type: string, integer, array… */
  type: string;
  description: string;
  example?: unknown;
}

export interface ProblemDoc {
  slug: string;
  status: number;
  title: string;
  description: string;
  /**
   * The detail it answers with, from x-problem's detail: fixed for a named
   * error, a template over the members ({provider}) for a typed one.
   * Problems whose detail comes from a hand-written error have none.
   */
  detail?: string;
  /** Extension members the problem adds to the standard ones. */
  members: ProblemMember[];
}

export type ProblemSet = "shared" | (typeof SERVICES)[number];

interface Property {
  type?: string;
  description?: string;
  example?: unknown;
  $ref?: string;
}

interface Schema {
  "x-problem"?: { slug: string; status: number; title: string; description?: string; detail?: string };
  properties?: Record<string, Property>;
}

/** A service's problem files: its problems.yml, then each feature's spec.yml. */
function filesOf(svc: (typeof SERVICES)[number]): string[] {
  const features = api(`${svc}/internal/features`);
  const specs = readdirSync(features, { withFileTypes: true })
    .filter((entry) => entry.isDirectory())
    .map((entry) => `${features}/${entry.name}/spec.yml`)
    .filter((file) => existsSync(file))
    .sort();
  return [api(`${svc}/problems.yml`), ...specs];
}

function problemsIn(files: readonly string[]): ProblemDoc[] {
  const schemas = files.flatMap((file) => {
    const doc = parse(readFileSync(file, "utf8")) as { components?: { schemas?: Record<string, Schema> } };
    return Object.values(doc.components?.schemas ?? {});
  });
  return schemas
    .filter((schema) => schema["x-problem"])
    .map((schema) => ({
      slug: schema["x-problem"]!.slug,
      status: schema["x-problem"]!.status,
      title: schema["x-problem"]!.title,
      description: schema["x-problem"]!.description ?? "",
      ...(schema["x-problem"]!.detail ? { detail: schema["x-problem"]!.detail } : {}),
      members: Object.entries(schema.properties ?? {}).map(([name, property]) => ({
        name,
        type: property.type ?? (property.$ref ? "object" : "string"),
        description: property.description ?? "",
        ...(property.example === undefined ? {} : { example: property.example }),
      })),
    }))
    .sort((a, b) => a.status - b.status || a.slug.localeCompare(b.slug));
}

export function problems(): Plugin {
  return {
    name: "portal:problems",
    resolveId: (id) => (id === ID ? RESOLVED : undefined),
    load(id) {
      if (id !== RESOLVED) return undefined;
      const shared = [api("problems.yml")];
      const own = Object.fromEntries(SERVICES.map((svc) => [svc, filesOf(svc)]));
      for (const file of [...shared, ...Object.values(own).flat()]) this.addWatchFile(file);
      const sets: Record<ProblemSet, ProblemDoc[]> = {
        shared: problemsIn(shared),
        ...(Object.fromEntries(SERVICES.map((svc) => [svc, problemsIn(own[svc])])) as Record<
          (typeof SERVICES)[number],
          ProblemDoc[]
        >),
      };
      return `export default ${JSON.stringify(sets)};`;
    },
  };
}
