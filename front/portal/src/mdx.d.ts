declare module "*.md" {
  import type { Component } from "solid-js";
  import type { DocMeta } from "@/features/docs/types";

  export const meta: DocMeta;
  const Content: Component<{ components?: Record<string, unknown> }>;
  export default Content;
}

declare module "*.mdx" {
  import type { Component } from "solid-js";
  import type { DocMeta } from "@/features/docs/types";

  export const meta: DocMeta;
  const Content: Component<{ components?: Record<string, unknown> }>;
  export default Content;
}
