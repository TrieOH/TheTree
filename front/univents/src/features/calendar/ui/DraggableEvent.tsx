import { useDraggable } from "@dnd-kit/core";
import { Pencil, Trash2 } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import type { EventColor, OccurrenceI, ProgramI } from "../model";
import { formatTimeRange, getDurationMinutes } from "../lib/date";

interface DraggableEventProps {
  occurrence: OccurrenceI;
  program: ProgramI | undefined;
  color: EventColor;
  onClick: (occurrenceId: string) => void;
  onDelete?: (occurrenceId: string) => void;
  overlapIndex?: number;
  overlapCount?: number;
}

export function DraggableEvent({
  occurrence,
  program,
  color,
  onClick,
  onDelete,
  overlapIndex = 0,
  overlapCount = 1,
}: DraggableEventProps) {
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!menuOpen) return;
    const close = (event: MouseEvent) => {
      if (!menuRef.current?.contains(event.target as Node)) setMenuOpen(false);
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, [menuOpen]);
  const { attributes, listeners, setNodeRef, transform, isDragging } =
    useDraggable({
      id: `event-${occurrence.id}`,
      data: { type: "occurrence", occurrenceId: occurrence.id },
    });

  const start = new Date(occurrence.starts_at);
  const startMins = start.getHours() * 60 + start.getMinutes();
  const durMins = getDurationMinutes(occurrence.starts_at, occurrence.ends_at);
  // The event is rendered inside the slot for its start hour. Only the
  // minute offset belongs here; using the full day offset would add the
  // hour position twice (e.g. 09:00 would render around 18:00).
  const top = startMins % 60;
  const height = Math.max(22, (durMins / 60) * 60 - 2);

  // The event lives inside its source hour slot, so clamp its vertical
  // translation to the 24-hour calendar bounds. This prevents it from
  // visually escaping above midnight or below the final grid row.
  const sourceAbsoluteTop = Math.floor(startMins / 60) * 60 + top;
  const maxTranslateY = 24 * 60 - sourceAbsoluteTop - height;
  const clampedY = transform
    ? Math.max(-top, Math.min(transform.y, maxTranslateY))
    : 0;
  const style = transform
    ? {
        transform: `translate3d(${transform.x}px, ${clampedY}px, 0)`,
        zIndex: 100,
      }
    : {};

  return (
    <div
      ref={setNodeRef}
      {...listeners}
      {...attributes}
      className={`group absolute rounded-md cursor-grab ${menuOpen ? "z-[9999] overflow-visible" : "z-[5] overflow-hidden"} hover:z-10 hover:shadow-lg transition-shadow`}
      style={{
        top: `${top}px`,
        height: `${height}px`,
        left: `calc(${(overlapIndex * 100) / overlapCount}% + 2px)`,
        width: `calc(${100 / overlapCount}% - 4px)`,
        background: color.bg,
        borderLeft: `3px solid ${color.border}`,
        color: color.text,
        padding: "2px 6px",
        opacity: isDragging ? 0 : 1,
        ...style,
      }}
      onClick={(e) => {
        if (!isDragging) {
          e.stopPropagation();
          onClick(occurrence.id);
        }
      }}
      onContextMenu={(e) => {
        e.preventDefault();
        e.stopPropagation();
        setMenuOpen(true);
      }}
    >
      {menuOpen && <div ref={menuRef} className="absolute right-1 top-1 z-30 min-w-28 rounded-lg border border-border bg-popover p-1 text-popover-foreground shadow-xl" onPointerDown={(e) => e.stopPropagation()}>
        <button type="button" className="flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-xs hover:bg-accent" onClick={(e) => { e.stopPropagation(); setMenuOpen(false); onClick(occurrence.id); }}><Pencil className="size-3" />Alterar</button>
        <button type="button" className="flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-xs text-destructive hover:bg-destructive/10" onClick={(e) => { e.stopPropagation(); setMenuOpen(false); onDelete?.(occurrence.id); }}><Trash2 className="size-3" />Excluir</button>
      </div>}
      <div className="text-[11px] font-semibold whitespace-nowrap overflow-hidden text-ellipsis leading-tight">
        {program?.name || "Evento"}
      </div>
      <div className="text-[10px] opacity-80 whitespace-nowrap overflow-hidden text-ellipsis leading-tight">
        {formatTimeRange(occurrence.starts_at, occurrence.ends_at)}
      </div>
      {occurrence.max_capacity && (
        <div className="text-[10px] opacity-65 mt-px">
          Cap: {occurrence.max_capacity}
        </div>
      )}
    </div>
  );
}
