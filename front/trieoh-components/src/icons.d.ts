// Solid 2 DOM types live in @solidjs/web, rather than solid-js.
declare module "~icons/*" {
  import type { JSX } from "@solidjs/web";
  const Icon: (props: JSX.IntrinsicElements["svg"]) => JSX.Element;
  export default Icon;
}
