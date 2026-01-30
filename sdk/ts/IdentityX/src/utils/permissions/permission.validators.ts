const NAMESPACE_RE = /^[a-zA-Z][a-zA-Z0-9_]*$/;
const SPECIFIER_RE = /^([a-zA-Z0-9_.-]+|[0-9]+)$/;
const ACTION_PART_RE = /^[a-zA-Z0-9_]+$/;

export function assertNamespace(ns: string) {
  if (!NAMESPACE_RE.test(ns)) throw new Error(`Invalid namespace "${ns}"`);
  if (ns === "*" || ns === "**") throw new Error(`Namespace cannot be wildcard`);
}

export function assertSpecifier(spec: string) {
  if (spec === "*" || spec === "**") throw new Error(`Specifier cannot be wildcard`);
  if (!SPECIFIER_RE.test(spec)) throw new Error(`Invalid specifier "${spec}"`);
}

export function assertActionPart(part: string) {
  if (part === "*" || part === "**") throw new Error(`Action token cannot be wildcard`);
  if (!ACTION_PART_RE.test(part)) throw new Error(`Invalid action part "${part}"`);
}
