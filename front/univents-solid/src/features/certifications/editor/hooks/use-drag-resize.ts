import { onCleanup } from "solid-js";
import { clampCertificateValue } from "../utils";

export type ResizeHandle = "nw" | "ne" | "sw" | "se";

export interface ElementBounds {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface UseDragResizeOptions {
  bounds: () => ElementBounds;
  scale: () => number;
  canvas: () => { width: number; height: number };
  overflowAllowance: () => number;
  minWidth: number;
  minHeight: number;
  onChange: (bounds: ElementBounds) => void;
  onInteractionEnd?: (bounds: ElementBounds) => void;
}

interface InteractionStart {
  pointerX: number;
  pointerY: number;
  bounds: ElementBounds;
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
    x: clampCertificateValue(
      bounds.x,
      -allowedX,
      canvas.width - bounds.width + allowedX,
    ),
    y: clampCertificateValue(
      bounds.y,
      -allowedY,
      canvas.height - bounds.height + allowedY,
    ),
  };
}

export function useDragResize(options: UseDragResizeOptions) {
  let drag: InteractionStart | null = null;
  let resize: (InteractionStart & { handle: ResizeHandle }) | null = null;
  let cleanupInteraction: () => void = () => undefined;
  let latestBounds: ElementBounds = options.bounds();

  const handleDragMove = (event: PointerEvent) => {
    if (!drag) return;
    const currentScale = options.scale() > 0 ? options.scale() : 1;
    const currentCanvas = options.canvas();

    const next = clampPosition(
      {
        ...drag.bounds,
        x: drag.bounds.x + (event.clientX - drag.pointerX) / currentScale,
        y: drag.bounds.y + (event.clientY - drag.pointerY) / currentScale,
      },
      currentCanvas,
      options.overflowAllowance(),
    );

    latestBounds = next;
    options.onChange(next);
  };

  const handleDragEnd = () => {
    drag = null;
    window.removeEventListener("pointermove", handleDragMove);
    window.removeEventListener("pointerup", handleDragEnd);
    cleanupInteraction = () => undefined;
    options.onInteractionEnd?.(latestBounds);
  };

  const startDrag = (event: PointerEvent) => {
    event.stopPropagation();
    cleanupInteraction();
    latestBounds = options.bounds();
    drag = {
      pointerX: event.clientX,
      pointerY: event.clientY,
      bounds: options.bounds(),
    };
    window.addEventListener("pointermove", handleDragMove);
    window.addEventListener("pointerup", handleDragEnd);
    cleanupInteraction = () => {
      window.removeEventListener("pointermove", handleDragMove);
      window.removeEventListener("pointerup", handleDragEnd);
    };
  };

  const handleResizeMove = (event: PointerEvent) => {
    if (!resize) return;
    const currentScale = options.scale() > 0 ? options.scale() : 1;
    const currentCanvas = options.canvas();

    const dx = (event.clientX - resize.pointerX) / currentScale;
    const dy = (event.clientY - resize.pointerY) / currentScale;
    const start = resize.bounds;
    const right = start.x + start.width;
    const bottom = start.y + start.height;
    let { x, y, width, height } = start;

    if (resize.handle === "se") {
      width = Math.min(
        currentCanvas.width * (1 + options.overflowAllowance()),
        Math.max(options.minWidth, start.width + dx),
      );
      height = Math.min(
        currentCanvas.height * (1 + options.overflowAllowance()),
        Math.max(options.minHeight, start.height + dy),
      );
    } else if (resize.handle === "nw") {
      width = Math.min(
        currentCanvas.width * (1 + options.overflowAllowance()),
        Math.max(options.minWidth, start.width - dx),
      );
      height = Math.min(
        currentCanvas.height * (1 + options.overflowAllowance()),
        Math.max(options.minHeight, start.height - dy),
      );
      x = right - width;
      y = bottom - height;
    } else if (resize.handle === "ne") {
      width = Math.min(
        currentCanvas.width * (1 + options.overflowAllowance()),
        Math.max(options.minWidth, start.width + dx),
      );
      height = Math.min(
        currentCanvas.height * (1 + options.overflowAllowance()),
        Math.max(options.minHeight, start.height - dy),
      );
      y = bottom - height;
    } else {
      width = Math.min(
        currentCanvas.width * (1 + options.overflowAllowance()),
        Math.max(options.minWidth, start.width - dx),
      );
      height = Math.min(
        currentCanvas.height * (1 + options.overflowAllowance()),
        Math.max(options.minHeight, start.height + dy),
      );
      x = right - width;
    }

    const next = clampPosition(
      { x, y, width, height },
      currentCanvas,
      options.overflowAllowance(),
    );
    latestBounds = next;
    options.onChange(next);
  };

  const handleResizeEnd = () => {
    resize = null;
    window.removeEventListener("pointermove", handleResizeMove);
    window.removeEventListener("pointerup", handleResizeEnd);
    cleanupInteraction = () => undefined;
    options.onInteractionEnd?.(latestBounds);
  };

  const startResize = (handle: ResizeHandle) => (event: PointerEvent) => {
    event.stopPropagation();
    cleanupInteraction();
    latestBounds = options.bounds();
    resize = {
      pointerX: event.clientX,
      pointerY: event.clientY,
      bounds: options.bounds(),
      handle,
    };
    window.addEventListener("pointermove", handleResizeMove);
    window.addEventListener("pointerup", handleResizeEnd);
    cleanupInteraction = () => {
      window.removeEventListener("pointermove", handleResizeMove);
      window.removeEventListener("pointerup", handleResizeEnd);
    };
  };

  onCleanup(() => {
    cleanupInteraction();
  });

  return { startDrag, startResize };
}
