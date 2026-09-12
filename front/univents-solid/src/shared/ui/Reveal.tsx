import { animate } from "motion/mini";
import type { JSX } from "@solidjs/web";
import { untrack } from "solid-js";

export function Reveal(props: {
  children: JSX.Element;
  direction?: "up" | "left" | "right";
  delay?: number;
  /**
   * `false` mounts the content already visible. Lists that re-slice on resize
   * pass it for the rows that appear because of the new width, so a settled
   * layout does not replay the entrance animation.
   */
  animate?: boolean;
}) {
  let element!: HTMLDivElement;
  const offset = () =>
    props.direction === "left"
      ? "translateX(-15px)"
      : props.direction === "right"
        ? "translateX(15px)"
        : "translateY(15px)";

  const reveal = () => {
    if (
      !element ||
      window.matchMedia("(prefers-reduced-motion: reduce)").matches
    ) {
      return;
    }
    element.style.opacity = "0";
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (!entry?.isIntersecting) return;
        observer.disconnect();
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

  return (
    <div
      ref={(value) => {
        element = value;
        if (untrack(() => props.animate) === false) return;
        requestAnimationFrame(reveal);
      }}
    >
      {props.children}
    </div>
  );
}
