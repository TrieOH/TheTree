import { shapeClass, type Shape } from "./lib/shape";
import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "./lib/cn";

export interface CardLinkRenderProps {
  href: string;
  class: string;
  children: JSX.Element;
}

export interface CardProps {
  shape?: Shape;
  children: JSX.Element;
  class?: string;
  /** Turns the whole card into a link, with a hover state. */
  href?: string;
  /** Opens `href` in a new tab. */
  external?: boolean;
  /** Renders the link with a router's own component instead of a plain `<a>`. */
  renderLink?: (props: CardLinkRenderProps) => JSX.Element;
}

export function Card(props: CardProps) {
  const classes = () =>
    cn(
      "min-w-0 rounded-xl border border-border bg-card text-card-foreground shadow-xs",
      props.href &&
        "block transition-colors hover:bg-muted/60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
      shapeClass(props.shape, "surface"),
      props.class,
    );
  return (
    <Show
      when={props.href}
      fallback={<article class={classes()}>{props.children}</article>}
    >
      {(href) => (
        <Show
          when={props.renderLink && !props.external}
          fallback={
            <a
              href={href()}
              class={classes()}
              target={props.external ? "_blank" : undefined}
              rel={props.external ? "noopener noreferrer" : undefined}
            >
              {props.children}
            </a>
          }
        >
          {props.renderLink!({
            href: href(),
            class: classes(),
            children: props.children,
          })}
        </Show>
      )}
    </Show>
  );
}
