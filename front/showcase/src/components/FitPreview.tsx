import { createSignal, onSettled, type ParentProps } from "solid-js";

const PADDING = 16;

/**
 * Shows a component at its natural size, scaled down to fit the box it sits
 * in, so large components (a sidebar, a dialog) share a card with small ones.
 * It never scales up.
 */
export function FitPreview(props: ParentProps) {
  let box: HTMLDivElement | undefined;
  let inner: HTMLDivElement | undefined;
  const [scale, setScale] = createSignal(1);

  onSettled(() => {
    const fit = () => {
      if (!box || !inner) return;
      const width = inner.offsetWidth;
      const height = inner.offsetHeight;
      if (!width || !height) return;
      setScale(
        Math.min(
          1,
          (box.clientWidth - PADDING * 2) / width,
          (box.clientHeight - PADDING * 2) / height,
        ),
      );
    };
    const observer = new ResizeObserver(fit);
    if (box) observer.observe(box);
    if (inner) observer.observe(inner);
    return () => observer.disconnect();
  });

  return (
    <div ref={(el) => (box = el)} class="absolute inset-0">
      <div
        ref={(el) => (inner = el)}
        class="absolute top-1/2 left-1/2 flex w-max min-w-55 max-w-none origin-center items-center justify-center"
        style={{ transform: `translate(-50%, -50%) scale(${scale()})` }}
      >
        {props.children}
      </div>
    </div>
  );
}
