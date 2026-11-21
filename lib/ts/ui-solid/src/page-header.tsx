import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "./lib/cn";

export interface PageHeaderProps {
  title: JSX.Element;
  description?: JSX.Element;
  /** Buttons or links next to the title. */
  actions?: JSX.Element;
  /** Anything under the description, like a row of filters. */
  children?: JSX.Element;
  class?: string;
}

/** The title block at the top of a page. */
export function PageHeader(props: PageHeaderProps) {
  return (
    <header class={cn("mb-8 space-y-2", props.class)}>
      <div class="flex flex-wrap items-start justify-between gap-3">
        <h1 class="min-w-0 font-heading text-[28px] font-semibold leading-tight tracking-tight">
          {props.title}
        </h1>
        <Show when={props.actions}>
          <div class="flex shrink-0 items-center gap-2">{props.actions}</div>
        </Show>
      </div>
      <Show when={props.description}>
        <p class="text-muted-foreground">{props.description}</p>
      </Show>
      {props.children}
    </header>
  );
}

export interface PageSectionProps {
  title: JSX.Element;
  /** A note or link on the right of the title. */
  aside?: JSX.Element;
  id?: string;
  children: JSX.Element;
  class?: string;
}

/** A titled block of a page. */
export function PageSection(props: PageSectionProps) {
  return (
    <section id={props.id} class={cn("mb-10", props.class)}>
      <header class="mb-3 flex items-baseline justify-between gap-4">
        <h2 class="font-heading text-base font-semibold">{props.title}</h2>
        <Show when={props.aside}>
          <div class="text-[13px] text-muted-foreground">{props.aside}</div>
        </Show>
      </header>
      {props.children}
    </section>
  );
}
