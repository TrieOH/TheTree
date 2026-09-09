import type { JSX } from "@solidjs/web";
import { Portal } from "@solidjs/web";
import XIcon from "~icons/lucide/x";

const CloseIcon = XIcon as unknown as () => JSX.Element;

export function Drawer(props: {
  open: boolean;
  title: string;
  onClose: () => void;
  children: JSX.Element;
}) {
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
        class={`fixed inset-0 z-50 transition-[visibility] duration-300 ${props.open ? "visible" : "invisible pointer-events-none"}`}
      >
        <button
          type="button"
          aria-label="Fechar"
          class={`absolute inset-0 bg-black/30 backdrop-blur-sm transition-opacity duration-300 ${props.open ? "opacity-100" : "opacity-0"}`}
          onClick={() => props.onClose()}
        />
        <section
          role="dialog"
          aria-modal="true"
          onPointerDown={handlePointerDown}
          onPointerMove={handlePointerMove}
          onPointerUp={handlePointerUp}
          style={{ "touch-action": "none", "overscroll-behavior": "contain" }}
          class={`fixed! bottom-0! left-0! right-0! top-auto! max-h-[85vh] overflow-y-auto rounded-t-2xl border-t border-border bg-card p-4 pb-8 shadow-2xl transition-transform duration-300 ease-out ${props.open ? "translate-y-0" : "translate-y-full"}`}
        >
          <header class="mb-4 flex items-center justify-between border-b border-border pb-4">
            <h2 class="text-base font-semibold">{props.title}</h2>
            <button
              type="button"
              aria-label="Fechar filtros"
              class="flex size-8 items-center justify-center rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground"
              onClick={() => props.onClose()}
            >
              <CloseIcon />
            </button>
          </header>
          {props.children}
        </section>
      </div>
    </Portal>
  );
}
