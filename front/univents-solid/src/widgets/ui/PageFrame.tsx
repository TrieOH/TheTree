import { useMatches } from "@tanstack/solid-router";
import { Dynamic } from "@solidjs/web";
import type { ParentProps } from "solid-js";

declare module "@tanstack/solid-router" {
  interface StaticDataRouteOption {
    /**
     * The route renders its own `<main>` (the SDK's auth screens do), so the
     * frame steps back to a plain element instead of nesting a second landmark.
     */
    ownsMain?: boolean;
  }
}

/**
 * The single `<main>` landmark of the page, and the target of the skip link.
 * Exactly one is exposed to assistive tech; routes that bring their own say so
 * with `staticData: { ownsMain: true }`.
 */
export function PageFrame(props: ParentProps) {
  const ownsMain = useMatches({
    select: (matches) => matches.some((match) => match.staticData?.ownsMain),
  });

  return (
    <Dynamic component={ownsMain() ? "div" : "main"} id="main-content" tabindex="-1">
      {props.children}
    </Dynamic>
  );
}
