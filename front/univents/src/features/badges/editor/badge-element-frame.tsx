import { Trash2 } from "lucide-react";
import type { ReactNode } from "react";
import { cn } from "@/shared/lib/utils";
import type {
  ElementBounds,
  ResizeHandle,
} from "../../certifications/editor/hooks/use-drag-resize";
import { useDragResize } from "../../certifications/editor/hooks/use-drag-resize";

const HANDLES: ResizeHandle[] = ["nw", "ne", "sw", "se"];
const HANDLE_CLASS: Record<ResizeHandle, string> = {
  nw: "-left-1.5 -top-1.5 cursor-nwse-resize",
  ne: "-right-1.5 -top-1.5 cursor-nesw-resize",
  sw: "-bottom-1.5 -left-1.5 cursor-nesw-resize",
  se: "-bottom-1.5 -right-1.5 cursor-nwse-resize",
};

interface BadgeElementFrameProps {
  bounds: ElementBounds;
  scale: number;
  zIndex: number;
  canvas: { width: number; height: number };
  minSize?: number;
  overflowAllowance?: number;
  selected: boolean;
  editing?: boolean;
  onSelect: () => void;
  onDoubleClick?: () => void;
  onChangeBounds: (bounds: ElementBounds) => void;
  onDelete?: () => void;
  children: ReactNode;
}

export function BadgeElementFrame({
  bounds,
  scale,
  zIndex,
  canvas,
  minSize = 4,
  overflowAllowance = 0,
  selected,
  editing = false,
  onSelect,
  onDoubleClick,
  onChangeBounds,
  onDelete,
  children,
}: BadgeElementFrameProps) {
  const { startDrag, startResize } = useDragResize({
    bounds,
    scale,
    canvas,
    overflowAllowance,
    minWidth: minSize,
    minHeight: minSize,
    onChange: onChangeBounds,
  });

  return (
    <div
      data-badge-element
      className={cn(
        "absolute select-none",
        selected ? "cursor-move" : "cursor-pointer",
      )}
      style={{
        left: bounds.x,
        top: bounds.y,
        width: bounds.width,
        height: bounds.height,
        zIndex,
      }}
      onPointerDown={(event) => {
        onSelect();
        if (!editing) {
          event.preventDefault();
          startDrag(event);
        }
      }}
      onDoubleClick={onDoubleClick}
    >
      <div
        className={cn(
          "size-full",
          selected && "outline-2 outline-offset-2 outline-ring",
        )}
      >
        {children}
      </div>
      {selected && !editing ? (
        <>
          {HANDLES.map((handle) => (
            <div
              key={handle}
              aria-hidden="true"
              className={cn(
                "absolute size-3 touch-none rounded-full border border-ring bg-popover shadow-sm",
                HANDLE_CLASS[handle],
              )}
              onPointerDown={startResize(handle)}
            />
          ))}
          {onDelete ? (
            <button
              type="button"
              className="absolute -top-2.5 -right-2.5 flex size-5 items-center justify-center rounded-full bg-destructive text-destructive-foreground shadow-sm"
              aria-label="Excluir elemento"
              onPointerDown={(event) => event.stopPropagation()}
              onClick={onDelete}
            >
              <Trash2 className="size-3" />
            </button>
          ) : null}
        </>
      ) : null}
    </div>
  );
}
