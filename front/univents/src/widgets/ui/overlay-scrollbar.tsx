import { useEffect, useState } from "react";

interface ScrollbarState {
  visible: boolean;
  top: number;
  height: number;
  trackHeight: number;
}

const MIN_THUMB_HEIGHT = 32;

export function OverlayScrollbar() {
  const [state, setState] = useState<ScrollbarState>({
    visible: false,
    top: 0,
    height: 0,
    trackHeight: 0,
  });

  useEffect(() => {
    const root = document.documentElement;
    const desktop = window.matchMedia("(hover: hover) and (pointer: fine)");
    if (!desktop.matches) return;

    root.classList.add("overlay-scrollbar");
    let frame = 0;

    const update = () => {
      cancelAnimationFrame(frame);
      frame = requestAnimationFrame(() => {
        const trackHeight = window.innerHeight - 8;
        const contentHeight = root.scrollHeight;
        const visible = contentHeight > window.innerHeight;
        const height = Math.max(
          MIN_THUMB_HEIGHT,
          (window.innerHeight / contentHeight) * trackHeight,
        );
        const maxTop = trackHeight - height;
        const maxScroll = contentHeight - window.innerHeight;
        setState({
          visible,
          top: maxScroll ? (window.scrollY / maxScroll) * maxTop + 4 : 4,
          height,
          trackHeight,
        });
      });
    };

    update();
    window.addEventListener("scroll", update, { passive: true });
    window.addEventListener("resize", update);
    const observer = new ResizeObserver(update);
    observer.observe(document.body);

    return () => {
      cancelAnimationFrame(frame);
      observer.disconnect();
      window.removeEventListener("scroll", update);
      window.removeEventListener("resize", update);
      root.classList.remove("overlay-scrollbar");
    };
  }, []);

  if (!state.visible) return null;

  const move = (clientY: number) => {
    const maxTop = state.trackHeight - state.height;
    const maxScroll =
      document.documentElement.scrollHeight - window.innerHeight;
    const nextTop = Math.max(
      4,
      Math.min(clientY - state.height / 2, maxTop + 4),
    );
    window.scrollTo({
      top: ((nextTop - 4) / maxTop) * maxScroll,
      behavior: "auto",
    });
  };

  return (
    <div
      className="pointer-events-auto fixed inset-y-0 right-0 z-40 hidden w-3 md:block"
      onPointerDown={(event) => {
        if (event.target !== event.currentTarget) return;
        const maxTop = state.trackHeight - state.height;
        const maxScroll =
          document.documentElement.scrollHeight - window.innerHeight;
        const targetTop = Math.max(
          0,
          Math.min(event.clientY - state.height / 2, maxTop),
        );
        window.scrollTo({
          top: maxTop ? (targetTop / maxTop) * maxScroll : 0,
          behavior: "smooth",
        });
      }}
    >
      <div
        role="scrollbar"
        tabIndex={0}
        aria-label="Rolagem da página"
        aria-orientation="vertical"
        aria-controls="root"
        aria-valuemin={0}
        aria-valuemax={100}
        aria-valuenow={Math.round(
          (window.scrollY /
            (document.documentElement.scrollHeight - window.innerHeight)) *
            100,
        )}
        className="pointer-events-auto absolute right-0.5 w-2 cursor-grab select-none rounded-full bg-primary/45 shadow-sm shadow-primary/20 outline-none transition-colors duration-200 hover:bg-primary/75 focus-visible:bg-primary/75 active:cursor-grabbing"
        style={{ top: state.top, height: state.height }}
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
            window.scrollTo({ top: document.documentElement.scrollHeight });
          }
        }}
        onPointerDown={(event) => {
          if (event.button !== 0) return;
          event.preventDefault();
          event.stopPropagation();
          event.currentTarget.setPointerCapture(event.pointerId);
          document.body.style.userSelect = "none";
          const movePointer = (moveEvent: PointerEvent) =>
            move(moveEvent.clientY);
          const stop = () => {
            window.removeEventListener("pointermove", movePointer);
            window.removeEventListener("pointerup", stop);
            document.body.style.removeProperty("user-select");
            event.currentTarget.releasePointerCapture(event.pointerId);
          };
          window.addEventListener("pointermove", movePointer);
          window.addEventListener("pointerup", stop);
        }}
      />
    </div>
  );
}
