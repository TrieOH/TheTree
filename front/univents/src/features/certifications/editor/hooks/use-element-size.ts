import type { RefObject } from "react";
import { useEffect, useRef, useState } from "react";

interface ElementSize {
  width: number;
  height: number;
}

export function useElementSize<T extends HTMLElement>(): {
  ref: RefObject<T | null>;
  size: ElementSize;
} {
  const ref = useRef<T | null>(null);
  const [size, setSize] = useState<ElementSize>({ width: 0, height: 0 });

  useEffect(() => {
    const element = ref.current;
    if (!element) return;

    const observer = new ResizeObserver(([entry]) => {
      setSize({
        width: entry.contentRect.width,
        height: entry.contentRect.height,
      });
    });

    observer.observe(element);
    return () => observer.disconnect();
  }, []);

  return { ref, size };
}
