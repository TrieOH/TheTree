import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

export interface DashboardStatCardProps {
  label: string;
  value: JSX.Element;
  hint: string;
  icon: (props: { class?: string }) => JSX.Element;
  action?: JSX.Element;
}

export function DashboardStatCard(props: DashboardStatCardProps): JSX.Element {
  return (
    <article class="rounded-xl border border-border bg-card p-4 shadow-xs transition-shadow hover:shadow-md">
      <div class="flex items-center justify-between text-muted-foreground">
        <p class="text-xs">{props.label}</p>
        <div class="flex items-center gap-2">
          <Show when={props.action}>
            <div class="shrink-0">{props.action}</div>
          </Show>
          <div class="size-4">
            {props.icon({ class: "size-4" })}
          </div>
        </div>
      </div>
      <p class="mt-2 text-xl font-semibold tracking-tight text-foreground">
        {props.value}
      </p>
      <p class="mt-2 truncate text-xs leading-4 text-muted-foreground">
        {props.hint}
      </p>
    </article>
  );
}
