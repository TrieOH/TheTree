import type { JSX } from "@solidjs/web";
import { animate } from "motion";
import { onCleanup, untrack } from "solid-js";

export function Reveal(props: {
  children: JSX.Element;
  direction?: "up" | "left" | "right";
  delay?: number;
  class?: string;
  /**
   * `false` mounts the content already visible. Lists that re-slice on resize
   * pass it for the rows that appear because of the new width, so a settled
   * layout does not replay the entrance animation.
   */
  animate?: boolean;
}) {
  let element!: HTMLDivElement;
  let observer: IntersectionObserver | undefined;
  let rafId: number | undefined;

  const offset = () =>
    props.direction === "left"
      ? "translateX(-15px)"
      : props.direction === "right"
        ? "translateX(15px)"
        : "translateY(15px)";

  const reveal = () => {
    if (
      !element ||
      typeof IntersectionObserver === "undefined" ||
      window.matchMedia?.("(prefers-reduced-motion: reduce)")?.matches
    ) {
      return;
    }
    element.style.opacity = "0";
    observer = new IntersectionObserver(
      ([entry]) => {
        if (!entry?.isIntersecting) return;
        observer?.disconnect();
        observer = undefined;
        animate(
          element,
          { opacity: [0, 1], transform: [offset(), "translate(0, 0)"] },
          { delay: props.delay ?? 0, duration: 0.4, ease: [0.22, 1, 0.36, 1] },
        );
      },
      { rootMargin: "0px 0px -10% 0px", threshold: 0.08 },
    );
    observer.observe(element);
  };

  onCleanup(() => {
    if (rafId) cancelAnimationFrame(rafId);
    observer?.disconnect();
  });

  return (
    <div
      class={props.class}
      ref={(value) => {
        element = value;
        if (!value || untrack(() => props.animate) === false) return;
        rafId = requestAnimationFrame(reveal);
      }}
    >
      {props.children}
    </div>
  );
}
