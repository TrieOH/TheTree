import { Dynamic } from "@solidjs/web";
import type { JSX } from "@solidjs/web";
import { onCleanup, untrack } from "solid-js";
import { cn } from "./lib/cn";

/** How the content enters: `slide` fades while travelling, the rest stay in place. */
export type RevealVariant = "slide" | "fade" | "scale" | "blur";
/** Side the content travels in from (`slide` only). */
export type RevealDirection = "up" | "down" | "left" | "right";
export type RevealTag =
  | "div"
  | "section"
  | "article"
  | "header"
  | "footer"
  | "aside"
  | "nav"
  | "ul"
  | "ol"
  | "li"
  | "span"
  | "p";

export interface RevealProps {
  children: JSX.Element;
  /** Entrance effect. */
  variant?: RevealVariant;
  /** Side the content travels in from. Only used by the `slide` variant. */
  direction?: RevealDirection;
  /** Travel distance in px for the `slide` variant. */
  distance?: number;
  /** Animation length in seconds. */
  duration?: number;
  /** Wait before the animation starts, in seconds. */
  delay?: number;
  /** Any CSS easing: a keyword or a `cubic-bezier(...)`. */
  easing?: string;
  /**
   * Animate each direct child in sequence instead of the wrapper, waiting this
   * many seconds between one child and the next. Leave unset to animate the
   * wrapper as a single block.
   */
  stagger?: number;
  /**
   * `true` plays once, the first time the content scrolls into view. `false`
   * hides it again when it leaves and replays on every return.
   */
  once?: boolean;
  /** Share of the content (0-1) that must be visible to trigger. */
  threshold?: number;
  /** Margin around the viewport used to trigger, in CSS `rootMargin` syntax. */
  rootMargin?: string;
  /**
   * `false` mounts the content already visible. Lists that re-slice on resize
   * pass it for the rows that appear because of the new width, so a settled
   * layout does not replay the entrance animation.
   */
  animate?: boolean;
  /** Element rendered around the content. */
  as?: RevealTag;
  class?: string;
  /** Called each time the entrance starts. */
  onReveal?: () => void;
}

interface Watcher {
  once: boolean;
  enter: () => void;
  leave: () => void;
}

interface SharedObserver {
  observer: IntersectionObserver;
  watchers: WeakMap<Element, Watcher>;
}

// One observer per trigger configuration, shared by every Reveal on the page,
// instead of one per element. A landing page mounts ~10 of them, and an
// observer per element that disconnected itself from inside its own callback
// left each visit's whole subtree reachable (~300 DOM nodes per visit).
const observers = new Map<string, SharedObserver>();

function watch(
  element: Element,
  options: { threshold: number; rootMargin: string },
  watcher: Watcher,
): () => void {
  const key = `${options.threshold}|${options.rootMargin}`;
  let shared = observers.get(key);

  if (!shared) {
    const watchers = new WeakMap<Element, Watcher>();
    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          const target = watchers.get(entry.target);
          if (!target) continue;

          if (entry.isIntersecting) {
            if (target.once) {
              observer.unobserve(entry.target);
              watchers.delete(entry.target);
            }
            target.enter();
          } else if (!target.once) {
            target.leave();
          }
        }
      },
      { threshold: options.threshold, rootMargin: options.rootMargin },
    );
    shared = { observer, watchers };
    observers.set(key, shared);
  }

  shared.watchers.set(element, watcher);
  shared.observer.observe(element);

  return () => {
    shared.observer.unobserve(element);
    shared.watchers.delete(element);
  };
}

const slideOffset: Record<RevealDirection, (distance: number) => string> = {
  up: (distance) => `translate(0, ${distance}px)`,
  down: (distance) => `translate(0, ${-distance}px)`,
  left: (distance) => `translate(${-distance}px, 0)`,
  right: (distance) => `translate(${distance}px, 0)`,
};

function keyframes(
  variant: RevealVariant,
  direction: RevealDirection,
  distance: number,
): Keyframe[] {
  switch (variant) {
    case "fade":
      return [{ opacity: 0 }, { opacity: 1 }];
    case "scale":
      return [
        { opacity: 0, transform: "scale(0.94)" },
        { opacity: 1, transform: "scale(1)" },
      ];
    case "blur":
      return [
        { opacity: 0, filter: "blur(8px)" },
        { opacity: 1, filter: "blur(0)" },
      ];
    default:
      return [
        { opacity: 0, transform: slideOffset[direction](distance) },
        { opacity: 1, transform: "translate(0, 0)" },
      ];
  }
}

const prefersReducedMotion = () =>
  window.matchMedia?.("(prefers-reduced-motion: reduce)")?.matches ?? false;

/**
 * Plays an entrance animation when its content scrolls into view. Uses the
 * native Web Animations API and one shared `IntersectionObserver`, so a page
 * full of them costs no animation library. Content stays visible without
 * JavaScript, with `prefers-reduced-motion`, or with `animate={false}`.
 */
export function Reveal(props: RevealProps): JSX.Element {
  const variant = () => props.variant ?? "slide";
  const direction = () => props.direction ?? "up";
  const distance = () => props.distance ?? 15;
  const duration = () => props.duration ?? 0.4;
  const delay = () => props.delay ?? 0;
  const easing = () => props.easing ?? "cubic-bezier(0.22, 1, 0.36, 1)";
  const once = () => props.once ?? true;
  const threshold = () => props.threshold ?? 0.08;
  const rootMargin = () => props.rootMargin ?? "0px 0px -10% 0px";
  const tag = () => props.as ?? "div";

  let element!: HTMLElement;
  let release: (() => void) | undefined;
  let animations: Animation[] = [];

  const stop = () => {
    for (const animation of animations) animation.cancel();
    animations = [];
  };

  const play = () => {
    stop();
    // The backwards fill holds the hidden pose through the delay, so the
    // hidden marker can go as soon as the animations exist.
    const stagger = props.stagger;
    const targets =
      stagger === undefined ? [element] : Array.from(element.children);
    const frames = keyframes(variant(), direction(), distance());

    targets.forEach((target, index) => {
      if (typeof target.animate !== "function") return;
      animations.push(
        target.animate(frames, {
          delay: Math.round((delay() + index * (stagger ?? 0)) * 1000),
          duration: Math.round(duration() * 1000),
          easing: easing(),
          fill: "backwards",
        }),
      );
    });
    element.removeAttribute("data-reveal");
    props.onReveal?.();
  };

  const hide = () => {
    stop();
    element.setAttribute("data-reveal", "hidden");
  };

  const arm = () => {
    if (
      !element ||
      typeof IntersectionObserver === "undefined" ||
      prefersReducedMotion()
    ) {
      return;
    }
    // Marked before insertion, never after the first paint.
    element.setAttribute("data-reveal", "hidden");
    release = watch(
      element,
      { threshold: threshold(), rootMargin: rootMargin() },
      { once: once(), enter: play, leave: hide },
    );
  };

  onCleanup(() => {
    release?.();
    stop();
  });

  return (
    <Dynamic
      component={tag()}
      class={cn(
        props.stagger === undefined
          ? "data-[reveal=hidden]:opacity-0"
          : "data-[reveal=hidden]:*:opacity-0",
        props.class,
      )}
      ref={(value: HTMLElement) => {
        element = value;
        if (!value || untrack(() => props.animate) === false) return;
        untrack(arm);
      }}
    >
      {props.children}
    </Dynamic>
  );
}
