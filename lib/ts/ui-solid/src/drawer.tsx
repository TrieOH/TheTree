import { shapeClass, type Shape } from "./lib/shape";
import type { JSX } from "@solidjs/web";
import { Portal } from "@solidjs/web";
import { cn } from "./lib/cn";

export interface DrawerProps {
  shape?: Shape;
  open: boolean;
  title?: JSX.Element;
  onClose: () => void;
  children: JSX.Element;
  class?: string;
}

export function Drawer(props: DrawerProps): JSX.Element {
  let startY = 0;
  let dragging = false;

  const handlePointerDown = (event: PointerEvent) => {
    if (event.pointerType === "mouse") return;
    event.stopPropagation();
    startY = event.clientY;
    dragging = true;
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
  };

  const handlePointerMove = (event: PointerEvent) => {
    if (!dragging) return;
    event.preventDefault();
    event.stopPropagation();
  };

  const handlePointerUp = (event: PointerEvent) => {
    if (!dragging) return;
    event.stopPropagation();
    dragging = false;
    (event.currentTarget as HTMLElement).releasePointerCapture(event.pointerId);
    if (event.clientY - startY > 80) props.onClose();
  };

  return (
    <Portal>
      <div
        class={cn(
          "fixed inset-0 z-50 transition-[visibility] duration-300",
          props.open ? "visible" : "invisible pointer-events-none",
        )}
      >
        <button
          type="button"
          aria-label="Fechar"
          class={cn(
            "absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity duration-300 cursor-pointer",
            props.open ? "opacity-100" : "opacity-0",
          )}
          onClick={() => props.onClose()}
        />
        <section
          role="dialog"
          aria-modal="true"
          onPointerDown={handlePointerDown}
          onPointerMove={handlePointerMove}
          onPointerUp={handlePointerUp}
          style={{ "touch-action": "none", "overscroll-behavior": "contain" }}
          class={cn(
            "fixed! bottom-0! left-0! right-0! top-auto! max-h-[85vh] overflow-y-auto rounded-t-2xl border-t border-border bg-card p-4 pb-8 shadow-2xl transition-transform duration-300 ease-out",
            props.open ? "translate-y-0" : "translate-y-full",
            shapeClass(props.shape, "surface"),
            props.class,
          )}
        >
          {/* Subtle drag handle bar */}
          <div class="mx-auto -mt-1 mb-3 h-1.5 w-12 rounded-full bg-muted-foreground/30" />

          <header class="mb-4 flex items-center justify-between border-b border-border/60 pb-3">
            <h2 class="text-base font-semibold text-foreground">
              {props.title}
            </h2>
            <button
              type="button"
              aria-label="Fechar gaveta"
              class="flex size-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-muted hover:text-foreground cursor-pointer"
              onClick={() => props.onClose()}
            >
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="size-4"
                aria-hidden="true"
              >
                <path d="M18 6 6 18M6 6l12 12" />
              </svg>
            </button>
          </header>
          {props.children}
        </section>
      </div>
    </Portal>
  );
}
