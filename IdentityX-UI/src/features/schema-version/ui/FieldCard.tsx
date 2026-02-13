import { cn } from "@/shared/lib/utils";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical } from "lucide-react";
import { getFieldTypeIcon } from "../model/field-type-to-icon";

interface PropsI {
  id: string;
  className?: string;
  isFixed?: boolean;
}

export default function FieldCard({id, className, isFixed}: PropsI) {

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging
  } = useSortable({id: id, disabled: isFixed});

  const style = {
    transform: CSS.Transform.toString(transform),
    transition: isDragging ? undefined : transition,
    willChange: "transform",
    zIndex: isDragging ? 999 : "auto",
  }

  const Icon = getFieldTypeIcon("string");

  return (
    <div 
      ref={setNodeRef}
      style={style}
      className={cn(
        "flex w-full select-none p-2",
        "bg-card",
        "transition-[box-shadow,border-color,background-color] duration-300 ease-out group",
        "border-2 border-border rounded-sm",
        "shadow-[1px_1px_0_0_var(--color-border)] hover:shadow-[2px_2px_0_0_var(--color-border)] hover:bg-muted/50",
        isDragging && "opacity-50",
        isFixed && "bg-muted-foreground/10 border-dashed cursor-not-allowed",
        className
      )}
    >
      <div className="flex items-center gap-4 w-full">
        {!isFixed && (
          <GripVertical 
            {...attributes} 
            {...listeners} 
            className={cn(
              "text-muted-foreground cursor-grab active:cursor-grabbing outline-0 w-8 h-8",
              "hover:text-accent duration-150 transition-colors"
            )}
          />
        )}
        {isFixed && <div className="w-8 h-8"></div>}
        <Icon className="w-8 h-8 p-2 bg-accent rounded-sm text-accent-foreground shrink-0"/>
        <div className="flex-1">
          <h5 className="font-semibold">Label</h5>
          <div className="text-xs text-muted-foreground flex justify-between items-center">
            <span>Field Name</span>
            <span>Field Type</span>
          </div>
        </div>
      </div>
    </div>
  )
}