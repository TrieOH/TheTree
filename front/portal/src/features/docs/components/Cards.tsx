import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { Match, Show, Switch } from "solid-js";
import ArrowIcon from "~icons/lucide/arrow-up-right";

/** A grid of <Card>s. */
export function Cards(props: { children?: JSX.Element }) {
  return <div class="not-prose my-5 grid gap-3 sm:grid-cols-2">{props.children}</div>;
}

const cardClass =
  "group block rounded-xl border border-border bg-card p-4 transition-colors hover:bg-foreground/5";

/** `<Card title="Run it locally" href="/docs/onboarding/local">mise, just and the APIs</Card>` */
export function Card(props: { title: string; href?: string; children?: JSX.Element }) {
  const body = (
    <>
      <span class="flex items-center gap-1.5 font-medium">
        {props.title}
        <Show when={props.href?.startsWith("http")}>
          <ArrowIcon aria-hidden="true" class="size-3.5 text-muted-foreground" />
        </Show>
      </span>
      <Show when={props.children}>
        <span class="mt-1 block text-[13px] text-muted-foreground [&_p]:m-0">{props.children}</span>
      </Show>
    </>
  );

  return (
    <Switch fallback={<div class={cardClass}>{body}</div>}>
      <Match when={props.href?.startsWith("/")}>
        {/* Doc links are plain strings; the router resolves them at runtime. */}
        <Link to={props.href as never} class={cardClass}>
          {body}
        </Link>
      </Match>
      <Match when={props.href}>
        <a href={props.href} target="_blank" rel="noopener" class={cardClass}>
          {body}
        </a>
      </Match>
    </Switch>
  );
}
