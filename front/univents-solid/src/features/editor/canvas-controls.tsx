import type { JSX } from "@solidjs/web";
import Maximize2Icon from "~icons/lucide/maximize-2";
import MinusIcon from "~icons/lucide/minus";
import PlusIcon from "~icons/lucide/plus";

const Maximize2 = Maximize2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Minus = MinusIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

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
        <Minus class="size-3.5" />
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
        <Plus class="size-3.5" />
      </button>

      <div class="mx-0.5 h-4 w-px bg-border" />

      <button
        type="button"
        class="flex size-7 cursor-pointer items-center justify-center rounded-md text-foreground transition-colors hover:bg-muted"
        aria-label="Ajustar à tela"
        onClick={() => props.onReset()}
      >
        <Maximize2 class="size-3.5" />
      </button>
    </div>
  );
}
