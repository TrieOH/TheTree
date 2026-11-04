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

declare module "virtual:problems" {
  import type { ProblemDoc, ProblemSet } from "../mdx/problems-plugin";

  /** "shared" holds the problems every service answers; each service, its own. */
  const problems: Record<ProblemSet, ProblemDoc[]>;
  export default problems;
}

declare module "virtual:api/*" {
  import type { ApiModel } from "../mdx/api-plugin";

  /** null when this checkout has no generated spec (`just gen`). */
  const model: ApiModel | null;
  export default model;
}
