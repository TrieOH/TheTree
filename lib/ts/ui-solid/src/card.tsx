import type { JSX } from "@solidjs/web";
import { cn } from "./lib/cn";

export interface CardProps {
  children: JSX.Element;
  class?: string;
}

export function Card(props: CardProps) {
  return (
    <article
      class={cn(
        "min-w-0 rounded-xl border border-border bg-card text-card-foreground shadow-xs",
        props.class,
      )}
    >
      {props.children}
    </article>
  );
}
