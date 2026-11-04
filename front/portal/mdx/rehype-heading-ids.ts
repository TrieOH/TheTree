// Gives h2–h4 a GitHub-style id ("Run it locally" → "run-it-locally") so
// sections can be linked to and searched; repeats get -1, -2.

interface HastNode {
  type: string;
  tagName?: string;
  value?: string;
  properties?: Record<string, unknown>;
  children?: HastNode[];
}

const text = (node: HastNode): string => node.value ?? (node.children ?? []).map(text).join("");

export function slugify(value: string): string {
  return value
    .toLowerCase()
    .trim()
    .replace(/[^\p{L}\p{N}\s-]/gu, "")
    .replace(/\s+/g, "-");
}

export function rehypeHeadingIds() {
  return (tree: HastNode) => {
    const seen = new Map<string, number>();
    const visit = (node: HastNode) => {
      if (node.type === "element" && /^h[2-4]$/.test(node.tagName ?? "") && !node.properties?.id) {
        const base = slugify(text(node)) || "section";
        const count = seen.get(base) ?? 0;
        seen.set(base, count + 1);
        node.properties = { ...node.properties, id: count ? `${base}-${count}` : base };
      }
      node.children?.forEach(visit);
    };
    visit(tree);
  };
}
