import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

export function PageHeader(props: { title: JSX.Element; description?: JSX.Element; children?: JSX.Element }) {
  return (
    <header class="mb-8 space-y-2">
      <h1 class="font-heading text-[28px] font-semibold leading-tight tracking-tight">{props.title}</h1>
      <Show when={props.description}>
        <p class="text-muted-foreground">{props.description}</p>
      </Show>
      {props.children}
    </header>
  );
}

/** A titled block of a page, with an optional link or note on the right. */
export function Section(props: { title: string; aside?: JSX.Element; children: JSX.Element }) {
  return (
    <section class="mb-10">
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

/** Placeholder for data a later phase connects. */
export function NotConnected(props: { source: string; children?: JSX.Element }) {
  return (
    <div class="rounded-xl border border-dashed border-border px-4 py-6 text-center text-[13px] text-muted-foreground">
      {props.children ?? <>Not connected yet. This fills in from {props.source}.</>}
    </div>
  );
}
