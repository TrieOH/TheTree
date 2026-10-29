/** A full-screen editor's frame while its document loads: top bar, tool column and canvas. */
export function EditorSkeleton() {
  return (
    <div aria-hidden="true" class="flex h-dvh flex-col bg-background">
      <div class="h-14 shrink-0 animate-pulse border-b border-border bg-muted/40" />
      <div class="flex min-h-0 flex-1">
        <div class="hidden w-72 shrink-0 animate-pulse border-r border-border bg-muted/30 md:block" />
        <div class="grid flex-1 place-items-center p-6">
          <div class="aspect-3/2 w-full max-w-3xl animate-pulse rounded-lg bg-muted/50" />
        </div>
      </div>
    </div>
  );
}
