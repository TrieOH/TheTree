import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "@trieoh/ui-solid";

export const SIGNATURE_CANVAS_WIDTH = 1200;
export const SIGNATURE_CANVAS_HEIGHT = 420;

export interface SignatureCanvasRef {
  canvas: () => HTMLCanvasElement | undefined;
  clear: () => void;
  isEmpty: () => boolean;
  toBlob: (type?: string, quality?: number) => Promise<Blob | null>;
  toDataURL: () => string;
}

export interface SignatureCanvasProps {
  class?: string;
  canvasClass?: string;
  strokeColor?: string;
  lineWidth?: number;
  showBaseline?: boolean;
  baselineLabel?: string;
  onStroke?: () => void;
  ref?: (instance: SignatureCanvasRef) => void;
}

export function SignatureCanvas(props: SignatureCanvasProps): JSX.Element {
  let canvasEl: HTMLCanvasElement | undefined;
  let isDrawing = false;
  let hasStrokes = false;
  let lastPoint: { x: number; y: number } | null = null;

  const initContext = () => {
    if (!canvasEl) return;
    const ctx = canvasEl.getContext("2d");
    if (!ctx) return;

    canvasEl.width = SIGNATURE_CANVAS_WIDTH;
    canvasEl.height = SIGNATURE_CANVAS_HEIGHT;
    ctx.lineWidth = props.lineWidth ?? 3.5;
    ctx.lineCap = "round";
    ctx.lineJoin = "round";
    ctx.strokeStyle = props.strokeColor ?? "#0f172a";
  };

  const setupRef = (el: HTMLCanvasElement) => {
    canvasEl = el;
    initContext();
    if (props.ref) {
      props.ref({
        canvas: () => canvasEl,
        clear: () => {
          if (!canvasEl) return;
          const ctx = canvasEl.getContext("2d");
          if (!ctx) return;
          ctx.clearRect(0, 0, canvasEl.width, canvasEl.height);
          hasStrokes = false;
        },
        isEmpty: () => !hasStrokes,
        toBlob: (type = "image/png", quality = 0.95) => {
          return new Promise<Blob | null>((resolve) => {
            if (!canvasEl) {
              resolve(null);
              return;
            }
            canvasEl.toBlob(resolve, type, quality);
          });
        },
        toDataURL: () => {
          return canvasEl ? canvasEl.toDataURL("image/png") : "";
        },
      });
    }
  };

  const getPoint = (event: PointerEvent) => {
    if (!canvasEl) return null;
    const rect = canvasEl.getBoundingClientRect();
    return {
      x: (event.clientX - rect.left) * (canvasEl.width / rect.width),
      y: (event.clientY - rect.top) * (canvasEl.height / rect.height),
    };
  };

  const handlePointerDown = (event: PointerEvent) => {
    const pt = getPoint(event);
    if (!pt) return;
    isDrawing = true;
    lastPoint = pt;
    if (canvasEl) {
      canvasEl.setPointerCapture(event.pointerId);
    }
  };

  const handlePointerMove = (event: PointerEvent) => {
    if (!isDrawing || !lastPoint || !canvasEl) return;
    const pt = getPoint(event);
    if (!pt) return;

    const ctx = canvasEl.getContext("2d");
    if (!ctx) return;

    ctx.beginPath();
    ctx.moveTo(lastPoint.x, lastPoint.y);
    ctx.lineTo(pt.x, pt.y);
    ctx.stroke();

    lastPoint = pt;
    hasStrokes = true;
    props.onStroke?.();
  };

  const handlePointerUp = (event: PointerEvent) => {
    isDrawing = false;
    lastPoint = null;
    if (canvasEl && canvasEl.hasPointerCapture(event.pointerId)) {
      canvasEl.releasePointerCapture(event.pointerId);
    }
  };

  return (
    <div
      class={cn(
        "relative w-full aspect-[1200/420] overflow-hidden select-none bg-white",
        props.class,
      )}
    >
      <canvas
        ref={setupRef}
        class={cn(
          "size-full touch-none block cursor-crosshair",
          props.canvasClass,
        )}
        onPointerDown={handlePointerDown}
        onPointerMove={handlePointerMove}
        onPointerUp={handlePointerUp}
        onPointerCancel={handlePointerUp}
      />
      <Show when={props.showBaseline ?? true}>
        <div class="pointer-events-none absolute inset-x-6 sm:inset-x-12 bottom-[24%] flex items-center gap-2 border-b border-dashed border-slate-300">
          <span class="text-xs text-slate-400 font-mono pb-0.5 select-none">✕</span>
        </div>
        <div class="pointer-events-none absolute bottom-[6%] inset-x-0 text-center">
          <span class="text-[10px] sm:text-[11px] text-slate-400 select-none">
            {props.baselineLabel ?? "Assine sobre a linha"}
          </span>
        </div>
      </Show>
    </div>
  );
}
