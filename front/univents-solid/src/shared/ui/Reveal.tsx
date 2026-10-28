import type { JSX } from "@solidjs/web";
import { onCleanup, untrack } from "solid-js";

// One observer for every Reveal on the page instead of one per element
const pending = new WeakMap<Element, () => void>();
let shared: IntersectionObserver | undefined;

function whenVisible(element: Element, onVisible: () => void): () => void {
  shared ??= new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (!entry.isIntersecting) continue;
        const callback = pending.get(entry.target);
        shared?.unobserve(entry.target);
        pending.delete(entry.target);
        callback?.();
      }
    },
    { rootMargin: "0px 0px -10% 0px", threshold: 0.08 },
  );
  pending.set(element, onVisible);
  shared.observe(element);

  return () => {
    shared?.unobserve(element);
    pending.delete(element);
  };
}

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
  let release: (() => void) | undefined;
  let animation: Animation | undefined;

  const offset = () =>
    props.direction === "left"
      ? "translateX(-15px)"
      : props.direction === "right"
        ? "translateX(15px)"
        : "translateY(15px)";

  const play = () => {
    // Native WAAPI: the backwards fill holds the hidden state through the
    // delay, so the inline opacity can go as soon as the animation exists.
    animation = element.animate(
      { opacity: [0, 1], transform: [offset(), "translate(0, 0)"] },
      {
        delay: (props.delay ?? 0) * 1000,
        duration: 400,
        easing: "cubic-bezier(0.22, 1, 0.36, 1)",
        fill: "backwards",
      },
    );
    element.style.removeProperty("opacity");
  };

  const reveal = () => {
    if (
      !element ||
      typeof IntersectionObserver === "undefined" ||
      window.matchMedia?.("(prefers-reduced-motion: reduce)")?.matches
    ) {
      return;
    }
    element.style.opacity = "0";
    release = whenVisible(element, play);
  };

  onCleanup(() => {
    animation?.cancel();
    release?.();
  });

  return (
    <div
      class={props.class}
      ref={(value) => {
        element = value;
        if (!value || untrack(() => props.animate) === false) return;
        // Set the initial opacity before insertion, never after the first paint.
        untrack(reveal);
      }}
    >
      {props.children}
    </div>
  );
}
