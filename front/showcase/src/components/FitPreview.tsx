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
    <div ref={(el) => (box = el)} class="fit-box">
      <div
        ref={(el) => (inner = el)}
        class="fit-inner"
        style={{ transform: `translate(-50%, -50%) scale(${scale()})` }}
      >
        {props.children}
      </div>
    </div>
  );
}
