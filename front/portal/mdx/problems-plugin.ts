import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import type { Plugin } from "vite";
import { parse } from "yaml";

// `virtual:problems`: every service's RFC 9457 problems (ADR-0019), read from
// api/problems.yml (shared) and api/<svc>/problems.yml at build time, so the
// <Problems> doc component ships JSON, not a YAML parser.

const SERVICES = ["identityx", "payssage", "univents", "informd"] as const;
const ID = "virtual:problems";
const RESOLVED = `\0${ID}`;
const api = (path: string) => fileURLToPath(new URL(`../../../api/${path}`, import.meta.url));

export interface ProblemDoc {
  slug: string;
  status: number;
  title: string;
  description: string;
  /** Extension members the problem adds to the standard ones. */
  members: string[];
  shared: boolean;
}

interface Schema {
  "x-problem"?: { slug: string; status: number; title: string; description?: string };
  properties?: Record<string, unknown>;
}

function problemsIn(file: string, shared: boolean): ProblemDoc[] {
  const doc = parse(readFileSync(file, "utf8")) as { components?: { schemas?: Record<string, Schema> } };
  return Object.values(doc.components?.schemas ?? {})
    .filter((schema) => schema["x-problem"])
    .map((schema) => ({
      slug: schema["x-problem"]!.slug,
      status: schema["x-problem"]!.status,
      title: schema["x-problem"]!.title,
      description: schema["x-problem"]!.description ?? "",
      members: Object.keys(schema.properties ?? {}),
      shared,
    }));
}

export function problems(): Plugin {
  const files = [api("problems.yml"), ...SERVICES.map((svc) => api(`${svc}/problems.yml`))];
  return {
    name: "portal:problems",
    resolveId: (id) => (id === ID ? RESOLVED : undefined),
    load(id) {
      if (id !== RESOLVED) return undefined;
      for (const file of files) this.addWatchFile(file);
      const shared = problemsIn(files[0], true);
      const bySlug = Object.fromEntries(
        SERVICES.map((svc) => [
          svc,
          [...shared, ...problemsIn(api(`${svc}/problems.yml`), false)].sort(
            (a, b) => a.status - b.status || a.slug.localeCompare(b.slug),
          ),
        ]),
      );
      return `export default ${JSON.stringify(bySlug)};`;
    },
  };
}
