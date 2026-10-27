import type { JSX } from "@solidjs/web";
import { Show, createSignal, onSettled } from "solid-js";
import { CanvasControls } from "@trieoh/ui-solid";

const MIN_ZOOM = 0.2;
const MAX_ZOOM = 2;
const STEP = 0.1;

/**
 * Renders children inside a simulated screen of a fixed size (say 1920×1080),
 * scaled to fit the preview. `zoom` is relative to that fit: 1 is "fit to the
 * box", and the canvas controls change it.
 */
export function DeviceFrame(props: {
  width: number;
  height: number;
  zoom: number;
  onZoomChange: (zoom: number) => void;
  /** Pin the container's width instead of measuring it (used in catalog cards). */
  fixedWidth?: number;
  /** The frame never gets taller than this once scaled, so a phone stays a phone. */
  maxHeight?: number;
  controls?: boolean;
  children: JSX.Element;
}): JSX.Element {
  const [measured, setMeasured] = createSignal(0);

  let box: HTMLDivElement | undefined;

  onSettled(() => {
    if (!box) return;
    const element = box;
    // The ref runs before the element is in the page, so measure once it is,
    // then follow its width.
    setMeasured(element.clientWidth);
    const observer = new ResizeObserver((entries) =>
      setMeasured(entries[0]?.contentRect.width ?? element.clientWidth),
    );
    observer.observe(element);
    return () => observer.disconnect();
  });

  const available = () => props.fixedWidth ?? measured();
  const fit = () =>
    available() > 0
      ? Math.min(
          available() / props.width,
          (props.maxHeight ?? 640) / props.height,
        )
      : 0;
  const scale = () => fit() * props.zoom;
  const clamp = (value: number) =>
    Math.round(Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, value)) * 100) / 100;

  return (
    <div
      ref={(element) => (box = element)}
      data-wide
      class="relative w-full overflow-hidden rounded-xl border border-border bg-background"
      style={{
        width: props.fixedWidth ? `${props.fixedWidth}px` : undefined,
        height: `${props.height * fit()}px`,
      }}
    >
      <div class="absolute inset-0 overflow-auto [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
        <div
          class="mx-auto"
          style={{
            width: `${props.width * scale()}px`,
            height: `${props.height * scale()}px`,
          }}
        >
          <div
            class="relative origin-top-left overflow-hidden"
            style={{
              width: `${props.width}px`,
              height: `${props.height}px`,
              transform: `scale(${scale()})`,
            }}
          >
            {props.children}
          </div>
        </div>
      </div>
      <Show when={props.controls ?? true}>
        <CanvasControls
          scale={props.zoom}
          onZoomIn={() => props.onZoomChange(clamp(props.zoom + STEP))}
          onZoomOut={() => props.onZoomChange(clamp(props.zoom - STEP))}
          onReset={() => props.onZoomChange(1)}
        />
      </Show>
      <span class="pointer-events-none absolute left-3 top-2 rounded bg-card/90 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
        {props.width}×{props.height}
      </span>
    </div>
  );
}
