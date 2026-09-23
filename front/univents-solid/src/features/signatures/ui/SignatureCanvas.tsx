import type { JSX } from "@solidjs/web";
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
  strokeColor?: string;
  lineWidth?: number;
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
    <canvas
      ref={setupRef}
      class={cn(
        "touch-none select-none rounded-lg bg-white shadow-inner cursor-crosshair",
        props.class,
      )}
      style={{
        "aspect-ratio": `${SIGNATURE_CANVAS_WIDTH} / ${SIGNATURE_CANVAS_HEIGHT}`,
        width: "100%",
      }}
      onPointerDown={handlePointerDown}
      onPointerMove={handlePointerMove}
      onPointerUp={handlePointerUp}
      onPointerCancel={handlePointerUp}
    />
  );
}
