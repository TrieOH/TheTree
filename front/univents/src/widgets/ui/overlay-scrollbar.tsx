import { useEffect, useRef } from "react";

const MIN_THUMB_HEIGHT = 32;
const TRACK_PADDING = 4;

export function OverlayScrollbar() {
  const trackRef = useRef<HTMLDivElement>(null);
  const thumbRef = useRef<HTMLDivElement>(null);
  const metrics = useRef({
    trackHeight: 0,
    thumbHeight: 0,
    maxScroll: 0,
  });

  useEffect(() => {
    const desktop = window.matchMedia("(hover: hover) and (pointer: fine)");
    const track = trackRef.current;
    const thumb = thumbRef.current;

    if (!track || !thumb) return;

    let frame = 0;
    let scrollScheduled = false;
    let observer: ResizeObserver | null = null;
    let active = false;

    const applyPosition = () => {
      const { trackHeight, thumbHeight, maxScroll } = metrics.current;
      const maxTop = Math.max(0, trackHeight - thumbHeight);
      const progress =
        maxScroll > 0
          ? Math.min(1, Math.max(0, window.scrollY / maxScroll))
          : 0;

      thumb.style.transform = `translate3d(0, ${
        TRACK_PADDING + progress * maxTop
      }px, 0)`;

      thumb.setAttribute("aria-valuenow", String(Math.round(progress * 100)));
    };

    const measure = () => {
      const root = document.documentElement;
      const trackHeight = window.innerHeight - TRACK_PADDING * 2;
      const contentHeight = root.scrollHeight;
      const maxScroll = Math.max(0, contentHeight - window.innerHeight);
      const visible = maxScroll > 0;

      const thumbHeight = visible
        ? Math.max(
            MIN_THUMB_HEIGHT,
            (window.innerHeight / contentHeight) * trackHeight,
          )
        : 0;

      metrics.current = {
        trackHeight,
        thumbHeight,
        maxScroll,
      };

      track.style.display = visible ? "block" : "none";
      thumb.style.height = `${thumbHeight}px`;

      applyPosition();
    };

    const onScroll = () => {
      if (scrollScheduled) return;

      scrollScheduled = true;

      frame = requestAnimationFrame(() => {
        applyPosition();
        scrollScheduled = false;
      });
    };

    const onResize = () => {
      cancelAnimationFrame(frame);
      frame = requestAnimationFrame(measure);
    };

    const enable = () => {
      if (active) return;
      active = true;

      document.documentElement.classList.add("overlay-scrollbar");
      measure();

      window.addEventListener("scroll", onScroll, { passive: true });
      window.addEventListener("resize", onResize);

      observer = new ResizeObserver(onResize);
      observer.observe(document.body);
    };

    const disable = () => {
      if (!active) return;
      active = false;

      cancelAnimationFrame(frame);
      observer?.disconnect();
      observer = null;

      window.removeEventListener("scroll", onScroll);
      window.removeEventListener("resize", onResize);

      document.documentElement.classList.remove("overlay-scrollbar");
      track.style.display = "none";
    };

    const onMediaChange = () => {
      if (desktop.matches) enable();
      else disable();
    };

    if (desktop.matches) enable();
    desktop.addEventListener("change", onMediaChange);

    return () => {
      desktop.removeEventListener("change", onMediaChange);
      disable();
    };
  }, []);

  const move = (clientY: number) => {
    const { trackHeight, thumbHeight, maxScroll } = metrics.current;
    if (maxScroll <= 0) return;

    const maxTop = Math.max(0, trackHeight - thumbHeight);

    const nextTop = Math.max(
      TRACK_PADDING,
      Math.min(clientY - thumbHeight / 2, maxTop + TRACK_PADDING),
    );

    const progress = maxTop > 0 ? (nextTop - TRACK_PADDING) / maxTop : 0;

    window.scrollTo({
      top: progress * maxScroll,
      behavior: "auto",
    });
  };

  return (
    <div
      ref={trackRef}
      className="pointer-events-auto fixed inset-y-0 right-0 z-40 hidden w-3 md:block"
      style={{ display: "none" }}
      onPointerDown={(event) => {
        if (event.target !== event.currentTarget) return;

        const { trackHeight, thumbHeight, maxScroll } = metrics.current;
        if (maxScroll <= 0) return;

        const maxTop = Math.max(0, trackHeight - thumbHeight);
        const targetTop = Math.max(
          0,
          Math.min(event.clientY - thumbHeight / 2, maxTop),
        );

        window.scrollTo({
          top: maxTop ? (targetTop / maxTop) * maxScroll : 0,
          behavior: "smooth",
        });
      }}
    >
      <div
        ref={thumbRef}
        role="scrollbar"
        tabIndex={0}
        aria-label="Rolagem da página"
        aria-orientation="vertical"
        aria-controls="root"
        aria-valuemin={0}
        aria-valuemax={100}
        aria-valuenow={0}
        className="pointer-events-auto absolute right-0.5 w-2 cursor-grab select-none rounded-full bg-primary/45 shadow-sm shadow-primary/20 outline-none transition-colors duration-200 hover:bg-primary/75 focus-visible:bg-primary/75 active:cursor-grabbing will-change-transform"
        onKeyDown={(event) => {
          if (event.key === "ArrowDown") {
            event.preventDefault();
            window.scrollBy({ top: 80 });
          }
          if (event.key === "ArrowUp") {
            event.preventDefault();
            window.scrollBy({ top: -80 });
          }
          if (event.key === "PageDown") {
            event.preventDefault();
            window.scrollBy({ top: window.innerHeight });
          }
          if (event.key === "PageUp") {
            event.preventDefault();
            window.scrollBy({ top: -window.innerHeight });
          }
          if (event.key === "Home") {
            event.preventDefault();
            window.scrollTo({ top: 0 });
          }
          if (event.key === "End") {
            event.preventDefault();
            window.scrollTo({
              top: document.documentElement.scrollHeight,
            });
          }
        }}
        onPointerDown={(event) => {
          if (event.button !== 0) return;

          event.preventDefault();
          event.stopPropagation();
          event.currentTarget.setPointerCapture(event.pointerId);
          document.body.style.userSelect = "none";

          const movePointer = (moveEvent: PointerEvent) => {
            move(moveEvent.clientY);
          };

          const stop = () => {
            window.removeEventListener("pointermove", movePointer);
            window.removeEventListener("pointerup", stop);
            window.removeEventListener("pointercancel", stop);
            document.body.style.removeProperty("user-select");
          };

          window.addEventListener("pointermove", movePointer);
          window.addEventListener("pointerup", stop);
          window.addEventListener("pointercancel", stop);
        }}
      />
    </div>
  );
}
