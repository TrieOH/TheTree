import type { JSX } from "@solidjs/web";
import { cn } from "./lib/cn";

export interface LinkProps {
  href: string;
  children: JSX.Element;
  class?: string;
  "aria-label"?: string;
}

export function Link(props: LinkProps) {
  return (
    <a
      href={props.href}
      aria-label={props["aria-label"]}
      class={cn(
        "text-primary underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        props.class,
      )}
    >
      {props.children}
    </a>
  );
}
