import {
  createSignal,
  onSettled,
  type Component,
} from "solid-js";
import { cn } from "./lib/cn";

const MIN_THUMB_HEIGHT = 32;
const TRACK_PADDING = 4;

export interface OverlayScrollbarProps {
  class?: string;
}

export const OverlayScrollbar: Component<OverlayScrollbarProps> = (props) => {
  let thumbRef: HTMLDivElement | undefined;

  const [visible, setVisible] = createSignal(false);

  const metrics = {
    trackHeight: 0,
    thumbHeight: 0,
    maxScroll: 0,
  };

  onSettled(() => {
    const desktop = window.matchMedia("(hover: hover) and (pointer: fine)");

    let frame = 0;
    let scrollScheduled = false;
    let observer: ResizeObserver | null = null;
    let active = false;

    const applyPosition = () => {
      const { trackHeight, thumbHeight, maxScroll } = metrics;
      const maxTop = Math.max(0, trackHeight - thumbHeight);

      const progress =
        maxScroll > 0
          ? Math.min(1, Math.max(0, window.scrollY / maxScroll))
          : 0;

      if (!thumbRef) return;
      thumbRef.style.transform = `translate3d(0, ${
        TRACK_PADDING + progress * maxTop
      }px, 0)`;

      thumbRef.setAttribute(
        "aria-valuenow",
        String(Math.round(progress * 100)),
      );
    };

    const measure = () => {
      const root = document.documentElement;
      const docHeight = Math.max(
        root.scrollHeight,
        document.body.scrollHeight,
      );
      const viewportHeight = window.innerHeight;
      const maxScroll = Math.max(0, docHeight - viewportHeight);

      const trackHeight = viewportHeight - TRACK_PADDING * 2;
      const thumbHeight =
        maxScroll > 0
          ? Math.max(
              MIN_THUMB_HEIGHT,
              Math.round((viewportHeight / docHeight) * trackHeight),
            )
          : 0;

      metrics.trackHeight = trackHeight;
      metrics.thumbHeight = thumbHeight;
      metrics.maxScroll = maxScroll;

      const shouldShow = maxScroll > 8;
      setVisible(shouldShow);

      if (thumbRef) {
        thumbRef.style.height = `${thumbHeight}px`;
      }

      applyPosition();
    };

    const onScroll = () => {
      if (!scrollScheduled) {
        scrollScheduled = true;
        frame = window.requestAnimationFrame(() => {
          scrollScheduled = false;
          applyPosition();
        });
      }
    };

    const teardown = () => {
      if (!active) return;
      active = false;
      window.removeEventListener("scroll", onScroll);
      window.removeEventListener("resize", measure);
      observer?.disconnect();
      observer = null;
      window.cancelAnimationFrame(frame);
    };

    const setup = () => {
      if (active) return;
      active = true;
      measure();
      window.addEventListener("scroll", onScroll, { passive: true });
      window.addEventListener("resize", measure, { passive: true });
      observer = new ResizeObserver(measure);
      observer.observe(document.body);
    };

    const handleMedia = (event: MediaQueryList | MediaQueryListEvent) => {
      if (event.matches) {
        setup();
      } else {
        teardown();
        setVisible(false);
      }
    };

    handleMedia(desktop);
    desktop.addEventListener("change", handleMedia);

    return () => {
      desktop.removeEventListener("change", handleMedia);
      teardown();
    };
  });

  return (
    <aside
      class={cn(
        "pointer-events-none fixed top-0 right-0 z-50 h-screen w-3 transition-opacity duration-300",
        visible() ? "opacity-100" : "opacity-0",
        props.class,
      )}
      aria-hidden="true"
    >
      <div
        ref={(el) => (thumbRef = el)}
        class="absolute right-0.5 w-1.5 rounded-full bg-foreground/20 backdrop-blur-sm transition-[width,background-color] duration-150 will-change-transform hover:w-2 hover:bg-foreground/40"
        role="scrollbar"
        aria-orientation="vertical"
        aria-valuemin={0}
        aria-valuemax={100}
        aria-valuenow={0}
      />
    </aside>
  );
};
