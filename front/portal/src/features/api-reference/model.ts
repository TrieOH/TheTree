import type { ServiceId } from "@/features/catalog/services";
import type { ApiAuth, ApiModel, ApiOperation, ApiParameter, SchemaNode } from "../../../mdx/api-plugin";

export type { ApiAuth, ApiModel, ApiOperation, ApiParameter, SchemaNode };

// One lazy chunk per service, built from its spec at build time
// (mdx/api-plugin.ts).
const MODELS: Record<ServiceId, () => Promise<{ default: ApiModel | null }>> = {
  identityx: () => import("virtual:api/identityx"),
  payssage: () => import("virtual:api/payssage"),
  univents: () => import("virtual:api/univents"),
  informd: () => import("virtual:api/informd"),
};

export const loadApi = (service: ServiceId) => MODELS[service]().then((module) => module.default);

/** `Wallet[]`, `"open" | "closed"`, `string (uuid) | null`… */
export function typeLabel(node: SchemaNode | null | undefined): string {
  if (!node) return "any";
  const nullable = node.nullable ? " | null" : "";
  if ("ref" in node) return node.ref + nullable;
  if (node.enum) return node.enum.map((value) => JSON.stringify(value)).join(" | ") + nullable;
  const union = node.oneOf ?? node.anyOf;
  if (union) return union.map(typeLabel).join(" | ") + nullable;
  if (node.allOf) return node.allOf.map(typeLabel).join(" & ") + nullable;
  if (node.type === "array") {
    const item = typeLabel(node.items);
    return (item.includes(" ") ? `(${item})[]` : `${item}[]`) + nullable;
  }
  if (node.type === "object" && !node.properties && node.additionalProperties) {
    const value = typeof node.additionalProperties === "object" ? typeLabel(node.additionalProperties) : "any";
    return `Record<string, ${value}>` + nullable;
  }
  return (node.type ?? (node.properties ? "object" : "any")) + (node.format ? ` (${node.format})` : "") + nullable;
}

/** The component a node points at, through arrays: `Wallet[]` → "Wallet". */
export function refOf(node: SchemaNode | null | undefined): string | null {
  if (!node) return null;
  if ("ref" in node) return node.ref;
  return node.type === "array" ? refOf(node.items) : null;
}

export const operationAnchor = (op: ApiOperation) => `op-${op.id}`;
export const schemaAnchor = (name: string) => `model-${name}`;
export const tagAnchor = (tag: string) => `tag-${tag}`;
