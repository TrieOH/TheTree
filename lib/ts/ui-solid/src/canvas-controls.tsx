import type { JSX } from "@solidjs/web";

export interface CanvasControlsProps {
  scale: number;
  onZoomIn: () => void;
  onZoomOut: () => void;
  onReset: () => void;
}

export function CanvasControls(props: CanvasControlsProps): JSX.Element {
  return (
    <div
      data-canvas-controls
      class="absolute bottom-4 left-1/2 z-30 flex -translate-x-1/2 select-none items-center gap-1 rounded-lg border border-border bg-card/95 p-1 shadow-lg backdrop-blur-xs"
      onPointerDown={(event) => event.stopPropagation()}
    >
      <button
        type="button"
        class="flex size-7 cursor-pointer items-center justify-center rounded-md text-foreground transition-colors hover:bg-muted"
        aria-label="Diminuir zoom"
        onClick={() => props.onZoomOut()}
      >
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="size-3.5"
          aria-hidden="true"
        >
          <path d="M5 12h14" />
        </svg>
      </button>

      <span class="w-12 text-center text-xs font-medium tabular-nums text-muted-foreground">
        {Math.round(props.scale * 100)}%
      </span>

      <button
        type="button"
        class="flex size-7 cursor-pointer items-center justify-center rounded-md text-foreground transition-colors hover:bg-muted"
        aria-label="Aumentar zoom"
        onClick={() => props.onZoomIn()}
      >
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="size-3.5"
          aria-hidden="true"
        >
          <path d="M5 12h14" />
          <path d="M12 5v14" />
        </svg>
      </button>

      <div class="mx-0.5 h-4 w-px bg-border" />

      <button
        type="button"
        class="flex size-7 cursor-pointer items-center justify-center rounded-md text-foreground transition-colors hover:bg-muted"
        aria-label="Ajustar à tela"
        onClick={() => props.onReset()}
      >
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="size-3.5"
          aria-hidden="true"
        >
          <polyline points="15 3 21 3 21 9" />
          <polyline points="9 21 3 21 3 15" />
          <line x1="21" x2="14" y1="3" y2="10" />
          <line x1="3" x2="10" y1="21" y2="14" />
        </svg>
      </button>
    </div>
  );
}
