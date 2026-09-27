import type { JSX } from "@solidjs/web";

export function ProfileCard(props: { title: string; children: JSX.Element }) {
  return (
    <section class="rounded-md border border-border bg-card p-5 shadow-md shadow-foreground/5">
      <h2 class="mb-4 text-base font-semibold">{props.title}</h2>
      {props.children}
    </section>
  );
}
