import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

export function Kpi(props: { label: string; children: JSX.Element; note?: string }) {
  return (
    <div class="rounded-xl border border-border bg-card p-3.5">
      <span class="block text-[12px] text-muted-foreground">{props.label}</span>
      <b class="mt-1 block text-lg font-semibold tabular-nums">{props.children}</b>
      <Show when={props.note}>
        <small class="block text-[12px] text-muted-foreground">{props.note}</small>
      </Show>
    </div>
  );
}

export function KpiRow(props: { children: JSX.Element }) {
  return <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">{props.children}</div>;
}
