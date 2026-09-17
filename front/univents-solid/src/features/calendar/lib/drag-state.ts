import { createSignal } from "solid-js";
import type { CalendarDragData } from "../model";

let currentDragData: CalendarDragData | null = null;
let ignoreNextClick = false;

export const [isCalendarDragging, setIsCalendarDragging] = createSignal(false);

export function startCalendarDrag(
  e: DragEvent,
  data: CalendarDragData,
): void {
  currentDragData = data;
  setIsCalendarDragging(true);

  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = "all";
    const serialized = JSON.stringify(data);
    try {
      e.dataTransfer.setData("application/json", serialized);
    } catch {
      // ignore
    }
    try {
      e.dataTransfer.setData("text/plain", serialized);
    } catch {
      // ignore
    }
  }
}

export function endCalendarDrag(): void {
  currentDragData = null;
  setIsCalendarDragging(false);
  ignoreNextClick = true;
  setTimeout(() => {
    ignoreNextClick = false;
  }, 200);
}

export function getCalendarDrag(e?: DragEvent): CalendarDragData | null {
  if (currentDragData) {
    return currentDragData;
  }

  if (e?.dataTransfer) {
    try {
      const raw =
        e.dataTransfer.getData("application/json") ||
        e.dataTransfer.getData("text/plain");
      if (raw) {
        return JSON.parse(raw) as CalendarDragData;
      }
    } catch {
      // ignore
    }
  }

  return null;
}

export function shouldIgnoreClick(): boolean {
  return isCalendarDragging() || ignoreNextClick;
}
