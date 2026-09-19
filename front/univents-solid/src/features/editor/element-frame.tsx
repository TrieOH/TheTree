import type { JSX } from "@solidjs/web";
import { For, Show, createMemo } from "solid-js";
import Trash2Icon from "~icons/lucide/trash-2";
import { cn } from "@trieoh/ui-solid";
import { createDragResize } from "./use-drag-resize";
import type { ElementBounds, ResizeHandle } from "./types";

const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;

const HANDLES: readonly ResizeHandle[] = ["nw", "ne", "sw", "se"];

const HANDLE_CLASS: Record<ResizeHandle, string> = {
  nw: "-left-1.5 -top-1.5 cursor-nwse-resize",
  ne: "-right-1.5 -top-1.5 cursor-nesw-resize",
  sw: "-bottom-1.5 -left-1.5 cursor-nesw-resize",
  se: "-bottom-1.5 -right-1.5 cursor-nwse-resize",
};

export interface EditorElementFrameProps {
  bounds: ElementBounds;
  scale: number;
  canvas: { width: number; height: number };
  zIndex: number;
  minSize?: number;
  overflowAllowance?: number;
  selected: boolean;
  editing?: boolean;
  onSelect: () => void;
  onDoubleClick?: () => void;
  onChangeBounds: (bounds: ElementBounds) => void;
  onDelete?: () => void;
  children: JSX.Element;
}

export function EditorElementFrame(props: EditorElementFrameProps): JSX.Element {
  const dragResize = createDragResize({
    bounds: () => props.bounds,
    scale: () => props.scale,
    canvas: () => props.canvas,
    overflowAllowance: () => props.overflowAllowance ?? 0,
    minWidth: () => props.minSize ?? 4,
    minHeight: () => props.minSize ?? 4,
    onChange: (next) => props.onChangeBounds(next),
  });

  const frameStyle = createMemo<JSX.CSSProperties>(() => ({
    position: "absolute",
    left: `${props.bounds.x}px`,
    top: `${props.bounds.y}px`,
    width: `${props.bounds.width}px`,
    height: `${props.bounds.height}px`,
    "z-index": props.zIndex,
  }));

  const cursorClass = createMemo(() =>
    props.selected ? "cursor-move" : "cursor-pointer",
  );
  const outlineClass = createMemo(() =>
    props.selected ? "outline-2 outline-offset-2 outline-ring" : "",
  );
  const showHandles = createMemo(() => Boolean(props.selected && !props.editing));

  return (
    <div
      data-badge-element
      data-editor-element
      style={frameStyle()}
      class={cn("absolute select-none", cursorClass())}
      onPointerDown={(event) => {
        props.onSelect();
        if (!props.editing) {
          event.preventDefault();
          dragResize.startDrag(event);
        }
      }}
      onDblClick={() => props.onDoubleClick?.()}
    >
      <div class={cn("size-full", outlineClass())}>
        {props.children}
      </div>

      <Show when={showHandles()}>
        <For each={HANDLES}>
          {(handle) => (
            <div
              aria-hidden="true"
              class={cn(
                "absolute size-3 touch-none rounded-full border border-ring bg-popover shadow-sm",
                HANDLE_CLASS[handle],
              )}
              onPointerDown={(event) => {
                event.preventDefault();
                dragResize.startResize(handle, event);
              }}
            />
          )}
        </For>

        <Show when={props.onDelete}>
          <button
            type="button"
            class="absolute -top-2.5 -right-2.5 flex size-5 items-center justify-center rounded-full bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/90 cursor-pointer"
            aria-label="Excluir elemento"
            onPointerDown={(event) => event.stopPropagation()}
            onClick={() => props.onDelete?.()}
          >
            <Trash2 class="size-3" />
          </button>
        </Show>
      </Show>
    </div>
  );
}
