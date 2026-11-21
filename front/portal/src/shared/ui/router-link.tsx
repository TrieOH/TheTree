import { Link } from "@tanstack/solid-router";
import type { CardLinkRenderProps } from "@trieoh/ui-solid";

/** Lets registry link components (like Card) navigate with the router instead of reloading. */
export function routerLink(props: CardLinkRenderProps) {
  // Paths are plain strings here; the router resolves them at runtime.
  return (
    <Link to={props.href as never} class={props.class}>
      {props.children}
    </Link>
  );
}
