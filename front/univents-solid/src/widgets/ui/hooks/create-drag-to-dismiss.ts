import type { JSX } from "@solidjs/web";
import { createSignal, onCleanup, type Accessor } from "solid-js";

export interface DragToDismissOptions {
  threshold?: number;
  onDismiss: () => void;
}

export interface DragToDismissReturn {
  dragY: Accessor<number>;
  isDragging: Accessor<boolean>;
  dragProps: {
    onPointerDown: (e: PointerEvent | MouseEvent) => void;
    onTouchStart: (e: TouchEvent) => void;
    onMouseDown: (e: MouseEvent) => void;
  };
  style: () => JSX.CSSProperties;
}

const getClientY = (e: TouchEvent | MouseEvent | PointerEvent): number => {
  if ("touches" in e && e.touches && e.touches.length > 0) {
    return e.touches[0].clientY ?? 0;
  }
  if ("changedTouches" in e && e.changedTouches && e.changedTouches.length > 0) {
    return e.changedTouches[0].clientY ?? 0;
  }
  if ("clientY" in e && typeof e.clientY === "number") {
    return e.clientY;
  }
  return 0;
};

export function createDragToDismiss(options: DragToDismissOptions): DragToDismissReturn {
  const threshold = options.threshold ?? 70;
  const [dragY, setDragY] = createSignal(0);
  const [isDragging, setIsDragging] = createSignal(false);
  let touchStartY = 0;

  const moveDrag = (clientY: number) => {
    if (!isDragging()) return;
    const delta = clientY - touchStartY;
    setDragY(delta > 0 ? delta : 0);
  };

  const endDrag = () => {
    if (!isDragging()) return;
    setIsDragging(false);
    cleanupGlobalListeners();
    if (dragY() > threshold) {
      options.onDismiss();
    }
    setDragY(0);
  };

  const onGlobalPointerMove = (e: PointerEvent) => moveDrag(getClientY(e));
  const onGlobalPointerUp = () => endDrag();
  const onGlobalTouchMove = (e: TouchEvent) => moveDrag(getClientY(e));
  const onGlobalTouchEnd = () => endDrag();
  const onGlobalMouseMove = (e: MouseEvent) => moveDrag(getClientY(e));
  const onGlobalMouseUp = () => endDrag();

  const cleanupGlobalListeners = () => {
    if (typeof window === "undefined") return;
    window.removeEventListener("pointermove", onGlobalPointerMove);
    window.removeEventListener("pointerup", onGlobalPointerUp);
    window.removeEventListener("touchmove", onGlobalTouchMove);
    window.removeEventListener("touchend", onGlobalTouchEnd);
    window.removeEventListener("mousemove", onGlobalMouseMove);
    window.removeEventListener("mouseup", onGlobalMouseUp);
  };

  const startDrag = (clientY: number) => {
    touchStartY = clientY;
    setIsDragging(true);
    if (typeof window !== "undefined") {
      window.addEventListener("pointermove", onGlobalPointerMove);
      window.addEventListener("pointerup", onGlobalPointerUp);
      window.addEventListener("touchmove", onGlobalTouchMove);
      window.addEventListener("touchend", onGlobalTouchEnd);
      window.addEventListener("mousemove", onGlobalMouseMove);
      window.addEventListener("mouseup", onGlobalMouseUp);
    }
  };

  const handlePointerDown = (e: PointerEvent | MouseEvent) => {
    if (e.button !== undefined && e.button !== 0 && (e as PointerEvent).pointerType === "mouse") {
      return;
    }
    startDrag(getClientY(e));
  };

  const handleTouchStart = (e: TouchEvent) => {
    startDrag(getClientY(e));
  };

  onCleanup(() => {
    cleanupGlobalListeners();
  });

  const style = (): JSX.CSSProperties => ({
    transform: dragY() > 0 ? `translateY(${dragY()}px)` : undefined,
    transition: isDragging()
      ? "none"
      : "transform 0.2s cubic-bezier(0.16, 1, 0.3, 1)",
  });

  return {
    dragY,
    isDragging,
    dragProps: {
      onPointerDown: handlePointerDown,
      onTouchStart: handleTouchStart,
      onMouseDown: handlePointerDown,
    },
    style,
  };
}
