// Lifts a doc's header into `export const meta = {...}` so pages can render
// the title, number, status and date themselves. The repo's ADRs open with
// `# ADR-0019 — Title` and a list of `**Status:** accepted`-style fields;
// both are removed from the body once read.

interface MdNode {
  type: string;
  depth?: number;
  value?: string;
  children?: MdNode[];
  data?: Record<string, unknown>;
}

const FIELD = /^([A-Za-z][\w -]*):$/;
const NUMBERED = /^([A-Z]+-\d+)\s+[—–-]\s+(.+)$/;

function text(node: MdNode): string {
  if (node.value !== undefined) return node.value;
  return (node.children ?? []).map(text).join("");
}

/** `- **Status:** accepted` → ["status", "accepted"], else undefined. */
function field(item: MdNode): [string, string] | undefined {
  const paragraph = item.children?.[0];
  const [label, ...rest] = paragraph?.children ?? [];
  if (item.children?.length !== 1 || label?.type !== "strong") return undefined;
  const match = FIELD.exec(text(label).trim());
  if (!match) return undefined;
  return [match[1].toLowerCase().replace(/\s+/g, "_"), rest.map(text).join("").trim()];
}

function literal(value: string) {
  return { type: "Literal", value, raw: JSON.stringify(value) };
}

function metaExport(meta: Record<string, string>): MdNode {
  const value = `export const meta = ${JSON.stringify(meta)};`;
  const program = {
    type: "Program",
    sourceType: "module",
    body: [
      {
        type: "ExportNamedDeclaration",
        specifiers: [],
        source: null,
        declaration: {
          type: "VariableDeclaration",
          kind: "const",
          declarations: [
            {
              type: "VariableDeclarator",
              id: { type: "Identifier", name: "meta" },
              init: {
                type: "ObjectExpression",
                properties: Object.entries(meta).map(([key, val]) => ({
                  type: "Property",
                  kind: "init",
                  method: false,
                  shorthand: false,
                  computed: false,
                  key: literal(key),
                  value: literal(val),
                })),
              },
            },
          ],
        },
      },
    ],
  };
  return { type: "mdxjsEsm", value, data: { estree: program } };
}

export function remarkAdr() {
  return (tree: MdNode) => {
    const children = tree.children ?? [];
    const meta: Record<string, string> = {};

    const heading = children[0];
    if (heading?.type === "heading" && heading.depth === 1) {
      const title = text(heading).trim();
      const numbered = NUMBERED.exec(title);
      meta.title = numbered ? numbered[2] : title;
      if (numbered) meta.number = numbered[1];
      children.shift();
    }

    const list = children[0];
    if (list?.type === "list") {
      const fields = (list.children ?? []).map(field);
      if (fields.length > 0 && fields.every((entry) => entry !== undefined)) {
        for (const [key, value] of fields as [string, string][]) meta[key] = value;
        children.shift();
      }
    }

    children.unshift(metaExport(meta));
    tree.children = children;
  };
}
