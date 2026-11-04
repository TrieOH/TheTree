import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import type { Plugin } from "vite";
import { parse } from "yaml";

// `virtual:problems`: the RFC 9457 problems (ADR-0019) every service shares,
// from api/problems.yml, and each service's own, from api/<svc>/problems.yml,
// read at build time so the pages ship JSON, not a YAML parser.

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
  "x-problem"?: { slug: string; status: number; title: string; description?: string };
  properties?: Record<string, Property>;
}

function problemsIn(file: string): ProblemDoc[] {
  const doc = parse(readFileSync(file, "utf8")) as { components?: { schemas?: Record<string, Schema> } };
  return Object.values(doc.components?.schemas ?? {})
    .filter((schema) => schema["x-problem"])
    .map((schema) => ({
      slug: schema["x-problem"]!.slug,
      status: schema["x-problem"]!.status,
      title: schema["x-problem"]!.title,
      description: schema["x-problem"]!.description ?? "",
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
  const files = [api("problems.yml"), ...SERVICES.map((svc) => api(`${svc}/problems.yml`))];
  return {
    name: "portal:problems",
    resolveId: (id) => (id === ID ? RESOLVED : undefined),
    load(id) {
      if (id !== RESOLVED) return undefined;
      for (const file of files) this.addWatchFile(file);
      const sets: Record<ProblemSet, ProblemDoc[]> = {
        shared: problemsIn(files[0]),
        ...(Object.fromEntries(SERVICES.map((svc) => [svc, problemsIn(api(`${svc}/problems.yml`))])) as Record<
          (typeof SERVICES)[number],
          ProblemDoc[]
        >),
      };
      return `export default ${JSON.stringify(sets)};`;
    },
  };
}
