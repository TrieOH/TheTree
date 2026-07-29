import { useDraggable } from "@dnd-kit/core";
import { GripVertical } from "lucide-react";
import type { EventColor, ProgramI } from "../model";

interface DraggableProgramProps {
  program: ProgramI;
  color: EventColor;
}

export function DraggableProgram({ program, color }: DraggableProgramProps) {
  const { attributes, listeners, setNodeRef, transform, isDragging } =
    useDraggable({
      id: `program-${program.id}`,
      data: { type: "program", programId: program.id },
    });

  const style = transform
    ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
        zIndex: 100,
        opacity: 0.8,
      }
    : {};

  return (
    <div
      ref={setNodeRef}
      {...listeners}
      {...attributes}
      className={`flex items-center gap-2 px-2 py-1.5 rounded-md cursor-grab hover:bg-muted transition-colors text-[13px] text-foreground active:cursor-grabbing group ${isDragging ? "opacity-60 scale-105 shadow-lg" : ""}`}
      style={style}
    >
      <GripVertical
        size={14}
        className="text-muted-foreground/40 group-hover:text-muted-foreground"
      />
      <span
        className="w-2.5 h-2.5 rounded-full flex-shrink-0"
        style={{ background: color.border }}
      />
      <span className="flex-1 whitespace-nowrap overflow-hidden text-ellipsis">
        {program.name}
      </span>
    </div>
  );
}
