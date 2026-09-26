import { Show } from "solid-js";

export function ProfileSkeleton(props: { ownProfile: boolean }) {
  return (
    <div aria-hidden="true" class="min-h-dvh animate-pulse bg-background pb-28">
      <section class="w-full border-b border-border bg-card shadow-md">
        <div class="h-40 w-full bg-muted sm:h-48 md:h-56" />
        <div class="md:hidden">
          <div class="relative mx-auto px-4">
            <div class="relative z-10 -mt-12 flex flex-col items-center">
              <div class="size-24 rounded-full border-4 border-background bg-muted shadow-xl" />
            </div>
            <div class="mt-3 space-y-2 pb-4">
              <div class="mx-auto h-7 w-44 rounded bg-muted" />
              <div class="mx-auto h-4 w-28 rounded bg-muted" />
              <div class="mx-auto h-4 w-36 rounded bg-muted" />
            </div>
            <Show when={props.ownProfile}>
              <div class="mb-4 flex gap-2">
                <div class="h-10 flex-1 rounded-md bg-muted" />
                <div class="h-10 w-10 rounded-md bg-muted" />
              </div>
            </Show>
          </div>
        </div>
        <div class="mx-auto hidden max-w-7xl px-4 md:block">
          <div class="relative pb-6">
            <div class="absolute -top-16 left-0 size-32 rounded-full border-4 border-background bg-muted shadow-xl" />
            <Show when={props.ownProfile}>
              <div class="absolute right-0 top-4 flex gap-2">
                <div class="h-10 w-28 rounded-md bg-muted" />
                <div class="h-10 w-10 rounded-md bg-muted" />
              </div>
            </Show>
            <div class="space-y-2 pt-18">
              <div class="h-8 w-56 rounded bg-muted" />
              <div class="h-4 w-36 rounded bg-muted" />
              <div class="h-4 w-44 rounded bg-muted" />
            </div>
          </div>
        </div>
        <nav class="mx-auto flex max-w-7xl gap-1 overflow-hidden px-4">
          <div class="h-12 w-24 shrink-0 border-b-2 border-primary bg-muted/40" />
          <div class="h-12 w-24 shrink-0 rounded-t bg-muted/50" />
          <Show when={props.ownProfile}>
            <div class="h-12 w-28 shrink-0 rounded-t bg-muted/50" />
            <div class="h-12 w-24 shrink-0 rounded-t bg-muted/50" />
          </Show>
        </nav>
      </section>
      <div class="mx-auto mt-4 grid max-w-7xl gap-4 px-4 md:grid-cols-[minmax(0,1fr)_280px] md:gap-5">
        <div class="space-y-5">
          <div class="rounded-md border border-border bg-card p-5 shadow-md shadow-foreground/5">
            <div class="mb-4 h-5 w-36 rounded bg-muted" />
            <div class="h-4 w-full rounded bg-muted" />
            <div class="mt-2 h-4 w-4/5 rounded bg-muted" />
            <div class="mt-2 h-4 w-2/3 rounded bg-muted" />
          </div>
          <div class="rounded-md border border-border bg-card p-5 shadow-md shadow-foreground/5">
            <div class="mb-4 h-5 w-28 rounded bg-muted" />
            <div class="h-4 w-full rounded bg-muted" />
            <div class="mt-2 h-4 w-5/6 rounded bg-muted" />
          </div>
        </div>
        <div class="space-y-5">
          <div class="rounded-md border border-border bg-card p-5 shadow-md shadow-foreground/5">
            <div class="mb-4 h-5 w-24 rounded bg-muted" />
            <div class="h-7 w-32 rounded bg-muted" />
          </div>
          <div class="rounded-md border border-border bg-card p-5 shadow-md shadow-foreground/5">
            <div class="mb-4 h-5 w-24 rounded bg-muted" />
            <div class="grid grid-cols-2 gap-2">
              <div class="h-10 rounded-md bg-muted" />
              <div class="h-10 rounded-md bg-muted" />
              <div class="h-10 rounded-md bg-muted" />
              <div class="h-10 rounded-md bg-muted" />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
