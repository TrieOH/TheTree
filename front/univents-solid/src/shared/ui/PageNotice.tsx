import { Show } from "solid-js";

/** A full-page, centered message for "nothing to show here" outcomes (not found, not available). */
export function PageNotice(props: { title: string; description?: string }) {
  return (
    <div class="grid min-h-[60vh] place-items-center px-6 text-center">
      <div class="space-y-2">
        <h1 class="text-xl font-semibold">{props.title}</h1>
        <Show when={props.description}>
          <p class="text-sm text-muted-foreground">{props.description}</p>
        </Show>
      </div>
    </div>
  );
}
