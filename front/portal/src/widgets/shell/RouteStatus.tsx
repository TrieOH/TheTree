import { Link, useRouter, type ErrorComponentProps } from "@tanstack/solid-router";
import { Button } from "@trieoh/ui-solid";

export function NotFoundPage() {
  return (
    <section class="mx-auto flex min-h-[60vh] max-w-md flex-col items-center justify-center gap-4 px-6 py-16 text-center">
      <h1 class="font-heading text-2xl font-semibold tracking-tight">Nothing here</h1>
      <p class="text-muted-foreground">This page doesn't exist, or it moved.</p>
      <Link to="/" class="text-[13px] underline underline-offset-4">Go home</Link>
    </section>
  );
}

export function RouteError(props: ErrorComponentProps) {
  const router = useRouter();
  const retry = () => {
    props.reset();
    void router.invalidate();
  };

  return (
    <section class="mx-auto flex min-h-[60vh] max-w-md flex-col items-center justify-center gap-4 px-6 py-16 text-center">
      <h1 class="font-heading text-2xl font-semibold tracking-tight">Something went wrong</h1>
      <p class="text-muted-foreground">This page failed to load. Try again.</p>
      <Button onClick={retry}>Try again</Button>
    </section>
  );
}
