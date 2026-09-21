import { createSignal, onCleanup } from "solid-js";

export interface ElementSize {
  width: number;
  height: number;
}

export function useElementSize<T extends HTMLElement>() {
  const [size, setSize] = createSignal<ElementSize>({ width: 0, height: 0 });
  let currentEl: T | null = null;
  let observer: ResizeObserver | null = null;

  const ref = (el: T | null) => {
    if (currentEl === el) return;
    if (observer) {
      observer.disconnect();
      observer = null;
    }
    currentEl = el;
    if (el) {
      setSize({
        width: el.clientWidth,
        height: el.clientHeight,
      });
      observer = new ResizeObserver(([entry]) => {
        if (entry) {
          setSize({
            width: entry.contentRect.width,
            height: entry.contentRect.height,
          });
        }
      });
      observer.observe(el);
    }
  };

  onCleanup(() => {
    if (observer) {
      observer.disconnect();
      observer = null;
    }
  });

  return { ref, size };
}
