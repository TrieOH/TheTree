import type { JSX } from "@solidjs/web";

export function KpiRow(props: { children: JSX.Element }) {
  return <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">{props.children}</div>;
}
