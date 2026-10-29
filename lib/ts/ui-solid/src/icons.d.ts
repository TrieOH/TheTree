// Icons are compiled by the consuming app's unplugin-icons Solid plugin.
declare module "~icons/*" {
  import type { JSX } from "@solidjs/web";
  const Icon: (props: JSX.IntrinsicElements["svg"]) => JSX.Element;
  export default Icon;
}
