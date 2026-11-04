// Turns a fenced ```mermaid block into <Mermaid chart="…" />, before shiki
// would highlight it as code. The component draws it in the browser.

interface MdNode {
  type: string;
  lang?: string | null;
  value?: string;
  children?: MdNode[];
  name?: string;
  attributes?: unknown[];
}

function walk(node: MdNode) {
  node.children?.forEach((child, i) => {
    if (child.type === "code" && child.lang === "mermaid") {
      node.children![i] = {
        type: "mdxJsxFlowElement",
        name: "Mermaid",
        attributes: [{ type: "mdxJsxAttribute", name: "chart", value: child.value ?? "" }],
        children: [],
      };
    } else {
      walk(child);
    }
  });
}

export function remarkMermaid() {
  return (tree: MdNode) => walk(tree);
}
