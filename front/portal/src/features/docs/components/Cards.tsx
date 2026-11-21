import type { JSX } from "@solidjs/web";
import { Card as RegistryCard } from "@trieoh/ui-solid";
import { Show } from "solid-js";
import { routerLink } from "@/shared/ui/router-link";
import ArrowIcon from "~icons/lucide/arrow-up-right";

/** A grid of <Card>s. */
export function Cards(props: { children?: JSX.Element }) {
  return <div class="not-prose my-5 grid gap-3 sm:grid-cols-2">{props.children}</div>;
}

/** `<Card title="Run it locally" href="/docs/onboarding/local">mise, just and the APIs</Card>` */
export function Card(props: { title: string; href?: string; children?: JSX.Element }) {
  const external = () => !!props.href && !props.href.startsWith("/");
  return (
    <RegistryCard href={props.href} external={external()} renderLink={routerLink} class="group p-4">
      <span class="flex items-center gap-1.5 font-medium">
        {props.title}
        <Show when={props.href?.startsWith("http")}>
          <ArrowIcon aria-hidden="true" class="size-3.5 text-muted-foreground" />
        </Show>
      </span>
      <Show when={props.children}>
        <span class="mt-1 block text-[13px] text-muted-foreground [&_p]:m-0">{props.children}</span>
      </Show>
    </RegistryCard>
  );
}
