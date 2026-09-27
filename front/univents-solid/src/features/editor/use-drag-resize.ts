import { onCleanup } from "solid-js";
import type { ElementBounds, ResizeHandle } from "./types";

export interface UseDragResizeOptions {
  bounds: () => ElementBounds;
  scale: () => number;
  canvas: () => { width: number; height: number };
  overflowAllowance?: (() => number) | number;
  minWidth?: (() => number) | number;
  minHeight?: (() => number) | number;
  onChange: (bounds: ElementBounds) => void;
  onInteractionEnd?: (bounds: ElementBounds) => void;
}

interface InteractionStart {
  pointerX: number;
  pointerY: number;
  bounds: ElementBounds;
}

function resolveVal(val: (() => number) | number | undefined, defaultVal: number): number {
  if (typeof val === "function") return val();
  return val ?? defaultVal;
}

function clamp(value: number, min: number, max: number): number {
  return Math.min(Math.max(value, min), max);
}

function clampPosition(
  bounds: ElementBounds,
  canvas: { width: number; height: number },
  overflowAllowance: number,
): ElementBounds {
  const allowedX = bounds.width * overflowAllowance;
  const allowedY = bounds.height * overflowAllowance;

  return {
    ...bounds,
    x: clamp(
      bounds.x,
      -allowedX,
      canvas.width - bounds.width + allowedX,
    ),
    y: clamp(
      bounds.y,
      -allowedY,
      canvas.height - bounds.height + allowedY,
    ),
  };
}

export function createDragResize(options: UseDragResizeOptions) {
  let drag: InteractionStart | null = null;
  let resize: (InteractionStart & { handle: ResizeHandle }) | null = null;
  let cleanupInteraction: () => void = () => undefined;
  let latestBounds: ElementBounds | null = null;

  const getScale = () => {
    const s = options.scale();
    return s > 0 ? s : 1;
  };

  const handleDragMove = (event: PointerEvent) => {
    if (!drag) return;

    const normalizedScale = getScale();
    const canvas = options.canvas();
    const allowance = resolveVal(options.overflowAllowance, 0);

    const next = clampPosition(
      {
        ...drag.bounds,
        x: drag.bounds.x + (event.clientX - drag.pointerX) / normalizedScale,
        y: drag.bounds.y + (event.clientY - drag.pointerY) / normalizedScale,
      },
      canvas,
      allowance,
    );

    latestBounds = next;
    options.onChange(next);
  };

  const handleDragEnd = () => {
    drag = null;
    window.removeEventListener("pointermove", handleDragMove);
    window.removeEventListener("pointerup", handleDragEnd);
    cleanupInteraction = () => undefined;
    if (latestBounds) {
      options.onInteractionEnd?.(latestBounds);
    }
  };

  const startDrag = (event: PointerEvent) => {
    event.stopPropagation();
    cleanupInteraction();

    drag = {
      pointerX: event.clientX,
      pointerY: event.clientY,
      bounds: options.bounds(),
    };
    latestBounds = drag.bounds;

    const onPointerUp = () => handleDragEnd();
    window.addEventListener("pointermove", handleDragMove);
    window.addEventListener("pointerup", onPointerUp);
    cleanupInteraction = () => {
      window.removeEventListener("pointermove", handleDragMove);
      window.removeEventListener("pointerup", onPointerUp);
    };
  };

  const handleResizeMove = (event: PointerEvent) => {
    if (!resize) return;

    const normalizedScale = getScale();
    const canvas = options.canvas();
    const minW = resolveVal(options.minWidth, 4);
    const minH = resolveVal(options.minHeight, 4);
    const allowance = resolveVal(options.overflowAllowance, 0);

    const deltaX = (event.clientX - resize.pointerX) / normalizedScale;
    const deltaY = (event.clientY - resize.pointerY) / normalizedScale;

    const next: ElementBounds = { ...resize.bounds };

    if (resize.handle.includes("e")) {
      next.width = Math.max(minW, resize.bounds.width + deltaX);
    }
    if (resize.handle.includes("s")) {
      next.height = Math.max(minH, resize.bounds.height + deltaY);
    }
    if (resize.handle.includes("w")) {
      const maxWidth = resize.bounds.x + resize.bounds.width;
      const desiredWidth = resize.bounds.width - deltaX;
      next.width = clamp(desiredWidth, minW, maxWidth);
      next.x = resize.bounds.x + (resize.bounds.width - next.width);
    }
    if (resize.handle.includes("n")) {
      const maxHeight = resize.bounds.y + resize.bounds.height;
      const desiredHeight = resize.bounds.height - deltaY;
      next.height = clamp(desiredHeight, minH, maxHeight);
      next.y = resize.bounds.y + (resize.bounds.height - next.height);
    }

    const clamped = clampPosition(next, canvas, allowance);
    latestBounds = clamped;
    options.onChange(clamped);
  };

  const handleResizeEnd = () => {
    resize = null;
    window.removeEventListener("pointermove", handleResizeMove);
    window.removeEventListener("pointerup", handleResizeEnd);
    cleanupInteraction = () => undefined;
    if (latestBounds) {
      options.onInteractionEnd?.(latestBounds);
    }
  };

  const startResize = (handle: ResizeHandle, event: PointerEvent) => {
    event.stopPropagation();
    cleanupInteraction();

    resize = {
      handle,
      pointerX: event.clientX,
      pointerY: event.clientY,
      bounds: options.bounds(),
    };
    latestBounds = resize.bounds;

    const onPointerUp = () => handleResizeEnd();
    window.addEventListener("pointermove", handleResizeMove);
    window.addEventListener("pointerup", onPointerUp);
    cleanupInteraction = () => {
      window.removeEventListener("pointermove", handleResizeMove);
      window.removeEventListener("pointerup", onPointerUp);
    };
  };

  onCleanup(() => cleanupInteraction());

  return { startDrag, startResize };
}
