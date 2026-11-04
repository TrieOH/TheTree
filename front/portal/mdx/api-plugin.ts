import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { validate } from "@scalar/openapi-parser";
import { sample } from "openapi-sampler";
import rehypeStringify from "rehype-stringify";
import remarkGfm from "remark-gfm";
import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import { unified } from "unified";
import type { Plugin } from "vite";
import { parse } from "yaml";

// `virtual:api/<svc>`: one service's API, read from api/<svc>/api-spec.gen.yml
// (tools/problemgen's output, generated before every front build) and turned
// into a small model at build time: operations by tag, schemas by name,
// descriptions as HTML, sample bodies. The page ships JSON, not a YAML parser
// or an OpenAPI viewer. Hosts never appear: the base URL is `{API_URL}`.

const SERVICES = ["identityx", "payssage", "univents", "informd"] as const;
const PREFIX = "virtual:api/";
const METHODS = ["get", "post", "put", "patch", "delete"] as const;
const spec = (svc: string) => fileURLToPath(new URL(`../../../api/${svc}/api-spec.gen.yml`, import.meta.url));

/** A schema with component references kept as names, so it stays small and acyclic. */
export type SchemaNode =
  | { ref: string; nullable?: boolean }
  | {
      type?: string;
      format?: string;
      description?: string;
      enum?: unknown[];
      nullable?: boolean;
      default?: unknown;
      example?: unknown;
      required?: string[];
      properties?: Record<string, SchemaNode>;
      additionalProperties?: SchemaNode | boolean;
      items?: SchemaNode;
      oneOf?: SchemaNode[];
      anyOf?: SchemaNode[];
      allOf?: SchemaNode[];
    };

export interface ApiParameter {
  name: string;
  in: "path" | "query" | "header" | "cookie";
  required: boolean;
  description: string;
  schema: SchemaNode | null;
}

export interface ApiBody {
  contentType: string;
  schema: SchemaNode | null;
  /** JSON sample built from the schema. */
  example: string | null;
}

export interface ApiResponse extends Partial<ApiBody> {
  status: string;
  description: string;
}

export interface ApiOperation {
  id: string;
  method: (typeof METHODS)[number];
  path: string;
  summary: string;
  /** HTML. */
  description: string;
  deprecated: boolean;
  /** Security scheme names; empty when the endpoint is public. */
  auth: string[];
  /** Linked-project scope (x-scope), when the operation has one. */
  scope: string | null;
  parameters: ApiParameter[];
  requestBody: (ApiBody & { required: boolean }) | null;
  /** Successful responses; failures are the problem slugs below. */
  responses: ApiResponse[];
  problems: string[];
}

export interface ApiAuth {
  name: string;
  /** "Bearer token", "API key in header X-API-Key"… */
  summary: string;
  description: string;
  /** The header that carries it, and what goes before the value. */
  header: string;
  prefix: string;
  /** Stands in for the value in snippets: `<token>`, `<api-key>`. */
  placeholder: string;
}

export interface ApiModel {
  title: string;
  version: string;
  /** HTML. */
  description: string;
  baseUrl: string;
  auth: ApiAuth[];
  tags: { name: string; description: string; operations: ApiOperation[] }[];
  schemas: Record<string, SchemaNode>;
}

type Json = Record<string, unknown>;

const markdown = unified().use(remarkParse).use(remarkGfm).use(remarkRehype).use(rehypeStringify);
const html = (text: unknown) => (typeof text === "string" && text.trim() ? String(markdown.processSync(text)) : "");

function resolve(doc: Json, value: unknown): Json {
  let node = value as Json;
  for (let hops = 0; typeof node?.$ref === "string" && hops < 10; hops++) {
    node = (node.$ref as string)
      .replace(/^#\//, "")
      .split("/")
      .map((part) => part.replace(/~1/g, "/").replace(/~0/g, "~"))
      .reduce<unknown>((at, key) => (at as Json)?.[key], doc) as Json;
  }
  return node ?? {};
}

function schemaNode(value: unknown): SchemaNode | null {
  if (!value || typeof value !== "object") return null;
  const raw = value as Json;
  if (typeof raw.$ref === "string") {
    const name = raw.$ref.match(/^#\/components\/schemas\/(.+)$/)?.[1];
    if (name) return { ref: name };
  }
  // OpenAPI 3.1 writes nullable as `type: [x, "null"]`.
  const types = Array.isArray(raw.type) ? (raw.type as string[]) : raw.type ? [raw.type as string] : [];
  const node: SchemaNode = {};
  const type = types.find((t) => t !== "null");
  if (type) node.type = type;
  if (types.includes("null") || raw.nullable === true) node.nullable = true;
  for (const key of ["format", "description", "default", "example"] as const) {
    if (raw[key] !== undefined) (node as Json)[key] = raw[key];
  }
  if (Array.isArray(raw.enum)) node.enum = raw.enum;
  if (Array.isArray(raw.required)) node.required = raw.required as string[];
  if (raw.properties && typeof raw.properties === "object") {
    node.properties = Object.fromEntries(
      Object.entries(raw.properties as Json).map(([name, child]) => [name, schemaNode(child) ?? {}]),
    );
  }
  if (raw.items) node.items = schemaNode(raw.items) ?? {};
  if (typeof raw.additionalProperties === "boolean") node.additionalProperties = raw.additionalProperties;
  else if (raw.additionalProperties) node.additionalProperties = schemaNode(raw.additionalProperties) ?? true;
  for (const key of ["oneOf", "anyOf", "allOf"] as const) {
    if (Array.isArray(raw[key])) node[key] = (raw[key] as unknown[]).map((child) => schemaNode(child) ?? {});
  }
  return node;
}

function exampleOf(doc: Json, schema: unknown): string | null {
  if (!schema) return null;
  try {
    return JSON.stringify(sample(schema as never, { skipReadOnly: true, quiet: true }, doc as never), null, 2);
  } catch {
    return null;
  }
}

function bodyOf(doc: Json, value: unknown): ApiBody | null {
  const content = resolve(doc, value).content as Record<string, Json> | undefined;
  const [contentType, media] = Object.entries(content ?? {})[0] ?? [];
  if (!contentType) return null;
  return {
    contentType,
    schema: schemaNode(media.schema),
    example: contentType.includes("json") ? exampleOf(doc, media.schema) : null,
  };
}

function authSummary(scheme: Json): string {
  if (scheme.type === "http") return scheme.scheme === "bearer" ? "Bearer token" : `HTTP ${scheme.scheme}`;
  if (scheme.type === "apiKey") return `API key in ${scheme.in} ${scheme.name}`;
  return String(scheme.type ?? "");
}

/**
 * Real API hosts become `{API_URL}` everywhere, prose and examples included:
 * one pasted into a browser would wake the API, and whoever runs the API
 * decides its URL.
 */
const HOST = /https?:\/\/(?:api\.[a-z0-9.-]+|localhost(?::\d+)?)/gi;

export function buildModel(doc: Json): ApiModel {
  return JSON.parse(JSON.stringify(rawModel(doc)).replace(HOST, "{API_URL}")) as ApiModel;
}

function rawModel(doc: Json): ApiModel {
  const info = (doc.info ?? {}) as Json;
  const components = (doc.components ?? {}) as Json;
  const schemes = (components.securitySchemes ?? {}) as Record<string, Json>;
  const globalAuth = (doc.security ?? []) as Record<string, unknown>[];
  const tagOrder = ((doc.tags ?? []) as Json[]).map((tag) => String(tag.name));
  const tagText = Object.fromEntries(((doc.tags ?? []) as Json[]).map((tag) => [tag.name, html(tag.description)]));
  const byTag = new Map<string, ApiOperation[]>();

  for (const [path, item] of Object.entries((doc.paths ?? {}) as Record<string, Json>)) {
    const shared = (item.parameters ?? []) as unknown[];
    for (const method of METHODS) {
      const op = item[method] as Json | undefined;
      if (!op) continue;
      const parameters = [...shared, ...((op.parameters ?? []) as unknown[])].map((value): ApiParameter => {
        const param = resolve(doc, value);
        return {
          name: String(param.name),
          in: param.in as ApiParameter["in"],
          required: param.required === true,
          description: html(param.description),
          schema: schemaNode(param.schema),
        };
      });
      const body = op.requestBody ? bodyOf(doc, op.requestBody) : null;
      const responses = Object.entries((op.responses ?? {}) as Json)
        .filter(([status]) => /^[123]/.test(status))
        .map(([status, value]): ApiResponse => ({
          status,
          description: String(resolve(doc, value).description ?? ""),
          ...(bodyOf(doc, value) ?? {}),
        }));
      const security = (op.security ?? globalAuth) as Record<string, unknown>[];
      const operation: ApiOperation = {
        id: String(op.operationId ?? `${method}-${path}`),
        method,
        path,
        summary: String(op.summary ?? ""),
        description: html(op.description),
        deprecated: op.deprecated === true,
        auth: [...new Set(security.flatMap((requirement) => Object.keys(requirement)))],
        scope: typeof op["x-scope"] === "string" ? op["x-scope"] : null,
        parameters,
        requestBody: body ? { ...body, required: resolve(doc, op.requestBody).required === true } : null,
        responses,
        problems: Array.isArray(op["x-problems"]) ? (op["x-problems"] as string[]) : [],
      };
      for (const tag of (op.tags as string[] | undefined) ?? ["other"]) {
        byTag.set(tag, [...(byTag.get(tag) ?? []), operation]);
      }
    }
  }

  const tags = [...byTag.keys()].sort((a, b) => {
    const ia = tagOrder.indexOf(a);
    const ib = tagOrder.indexOf(b);
    return (ia < 0 ? Infinity : ia) - (ib < 0 ? Infinity : ib) || a.localeCompare(b);
  });

  return {
    title: String(info.title ?? ""),
    version: String(info.version ?? ""),
    description: html(info.description),
    baseUrl: "{API_URL}",
    auth: Object.entries(schemes).map(([name, scheme]) => {
      const bearer = scheme.type === "http";
      return {
        name,
        summary: authSummary(scheme),
        description: html(scheme.description),
        header: bearer ? "Authorization" : String(scheme.name ?? "Authorization"),
        prefix: bearer ? `${String(scheme.scheme ?? "bearer").replace(/^./, (c) => c.toUpperCase())} ` : "",
        placeholder: bearer ? "<token>" : "<api-key>",
      };
    }),
    tags: tags.map((name) => ({ name, description: tagText[name] ?? "", operations: byTag.get(name)! })),
    schemas: Object.fromEntries(
      Object.entries((components.schemas ?? {}) as Json)
        .filter(([, value]) => !(value as Json)["x-problem"])
        .map(([name, value]) => [name, schemaNode(value) ?? {}]),
    ),
  };
}

export function apiModels(): Plugin {
  return {
    name: "portal:api",
    resolveId: (id) => (id.startsWith(PREFIX) ? `\0${id}` : undefined),
    async load(id) {
      if (!id.startsWith(`\0${PREFIX}`)) return undefined;
      const svc = id.slice(PREFIX.length + 1);
      if (!(SERVICES as readonly string[]).includes(svc)) return undefined;
      const file = spec(svc);
      this.addWatchFile(file);
      let text: string;
      try {
        text = readFileSync(file, "utf8");
      } catch {
        // Not generated in this checkout (`just gen`): the page says so.
        return "export default null;";
      }
      const doc = parse(text) as Json;
      const { valid, errors } = await validate(doc);
      if (!valid) this.warn(`api-spec.gen.yml for ${svc} has ${errors?.length ?? 0} validation errors`);
      return `export default ${JSON.stringify(buildModel(doc))};`;
    },
  };
}
