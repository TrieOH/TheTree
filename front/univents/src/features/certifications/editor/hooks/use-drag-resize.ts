import type { PointerEvent as ReactPointerEvent } from "react";
import { useCallback, useEffect, useRef } from "react";
import { clampCertificateValue } from "../utils";

export type ResizeHandle = "nw" | "ne" | "sw" | "se";

export interface ElementBounds {
  x: number;
  y: number;
  width: number;
  height: number;
}

interface UseDragResizeOptions {
  bounds: ElementBounds;
  scale: number;
  canvas: { width: number; height: number };
  overflowAllowance: number;
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

export function useDragResize({
  bounds,
  scale,
  canvas,
  overflowAllowance,
  minWidth,
  minHeight,
  onChange,
  onInteractionEnd,
}: UseDragResizeOptions) {
  const drag = useRef<InteractionStart | null>(null);
  const resize = useRef<(InteractionStart & { handle: ResizeHandle }) | null>(
    null,
  );
  const cleanupInteraction = useRef<() => void>(() => undefined);
  const latestBounds = useRef(bounds);
  latestBounds.current = bounds;

  const normalizedScale = scale > 0 ? scale : 1;

  const handleDragMove = useCallback(
    (event: PointerEvent) => {
      const start = drag.current;
      if (!start) return;

      const next = clampPosition(
        {
          ...start.bounds,
          x:
            start.bounds.x + (event.clientX - start.pointerX) / normalizedScale,
          y:
            start.bounds.y + (event.clientY - start.pointerY) / normalizedScale,
        },
        canvas,
        overflowAllowance,
      );

      latestBounds.current = next;
      onChange(next);
    },
    [canvas, normalizedScale, onChange, overflowAllowance],
  );

  const handleDragEnd = useCallback(() => {
    drag.current = null;
    window.removeEventListener("pointermove", handleDragMove);
    window.removeEventListener("pointerup", handleDragEnd);
    cleanupInteraction.current = () => undefined;
    onInteractionEnd?.(latestBounds.current);
  }, [handleDragMove, onInteractionEnd]);

  const startDrag = useCallback(
    (event: ReactPointerEvent) => {
      event.stopPropagation();
      cleanupInteraction.current();
      latestBounds.current = bounds;
      drag.current = {
        pointerX: event.clientX,
        pointerY: event.clientY,
        bounds,
      };
      window.addEventListener("pointermove", handleDragMove);
      window.addEventListener("pointerup", handleDragEnd);
      cleanupInteraction.current = () => {
        window.removeEventListener("pointermove", handleDragMove);
        window.removeEventListener("pointerup", handleDragEnd);
      };
    },
    [bounds, handleDragEnd, handleDragMove],
  );

  const handleResizeMove = useCallback(
    (event: PointerEvent) => {
      const interaction = resize.current;
      if (!interaction) return;

      const dx = (event.clientX - interaction.pointerX) / normalizedScale;
      const dy = (event.clientY - interaction.pointerY) / normalizedScale;
      const start = interaction.bounds;
      const right = start.x + start.width;
      const bottom = start.y + start.height;
      let { x, y, width, height } = start;

      if (interaction.handle === "se") {
        width = Math.max(minWidth, start.width + dx);
        height = Math.max(minHeight, start.height + dy);
      } else if (interaction.handle === "nw") {
        width = Math.max(minWidth, start.width - dx);
        height = Math.max(minHeight, start.height - dy);
        x = right - width;
        y = bottom - height;
      } else if (interaction.handle === "ne") {
        width = Math.max(minWidth, start.width + dx);
        height = Math.max(minHeight, start.height - dy);
        y = bottom - height;
      } else {
        width = Math.max(minWidth, start.width - dx);
        height = Math.max(minHeight, start.height + dy);
        x = right - width;
      }

      const next = clampPosition(
        { x, y, width, height },
        canvas,
        overflowAllowance,
      );
      latestBounds.current = next;
      onChange(next);
    },
    [canvas, minHeight, minWidth, normalizedScale, onChange, overflowAllowance],
  );

  const handleResizeEnd = useCallback(() => {
    resize.current = null;
    window.removeEventListener("pointermove", handleResizeMove);
    window.removeEventListener("pointerup", handleResizeEnd);
    cleanupInteraction.current = () => undefined;
    onInteractionEnd?.(latestBounds.current);
  }, [handleResizeMove, onInteractionEnd]);

  const startResize = useCallback(
    (handle: ResizeHandle) => (event: ReactPointerEvent) => {
      event.stopPropagation();
      cleanupInteraction.current();
      latestBounds.current = bounds;
      resize.current = {
        pointerX: event.clientX,
        pointerY: event.clientY,
        bounds,
        handle,
      };
      window.addEventListener("pointermove", handleResizeMove);
      window.addEventListener("pointerup", handleResizeEnd);
      cleanupInteraction.current = () => {
        window.removeEventListener("pointermove", handleResizeMove);
        window.removeEventListener("pointerup", handleResizeEnd);
      };
    },
    [bounds, handleResizeEnd, handleResizeMove],
  );

  useEffect(() => () => cleanupInteraction.current(), []);

  return { startDrag, startResize };
}
