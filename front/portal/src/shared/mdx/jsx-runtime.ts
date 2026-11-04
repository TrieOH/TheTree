import { Dynamic } from "@solidjs/web";
import type { JSX as SolidJSX } from "@solidjs/web";
import { createComponent, type Component } from "solid-js";

// The automatic JSX runtime MDX compiles docs against (vite.config.ts sets
// `jsxImportSource: "@/shared/mdx"`). Docs are static, so props are plain
// values: an intrinsic tag goes through <Dynamic>, anything else is a Solid
// component, e.g. one a doc imports or a page passes in `components`.

type Props = Record<string, unknown>;
type ElementType = string | Component<Props>;

export function jsx(type: ElementType, props: Props): SolidJSX.Element {
  if (typeof type === "string") {
    return createComponent(Dynamic, { ...props, component: type } as never);
  }
  return createComponent(type, props);
}

export const jsxs = jsx;

export function Fragment(props: { children?: SolidJSX.Element }): SolidJSX.Element {
  // oxlint-disable-next-line solid/reactivity -- docs are static
  return props.children;
}

export declare namespace JSX {
  type Element = import("@solidjs/web").JSX.Element;
  type IntrinsicElements = import("@solidjs/web").JSX.IntrinsicElements;
}
