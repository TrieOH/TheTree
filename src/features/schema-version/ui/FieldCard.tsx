import { cn } from "@/shared/lib/utils";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical, Pencil, Trash2 } from "lucide-react";
import { getFieldTypeIcon } from "../model/field-type-to-icon";
import type { VersionField } from "../model/types";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";

interface PropsI {
  field: VersionField;
  className?: string;
  overwriteType?: "password";
  isFixed?: boolean;
  onEdit?: (fieldKey: string) => void;
  onDelete?: (fieldKey: string) => void;
  onUpdateField?: (updatedField: VersionField) => void;
  onOpenEditPanel?: (field: VersionField) => void;
}

export default function FieldCard({field, className, overwriteType, isFixed = false, onEdit, onDelete, onUpdateField, onOpenEditPanel}: PropsI) {

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging
  } = useSortable({id: field.key, disabled: isFixed});

  const style = {
    transform: CSS.Transform.toString(transform),
    transition: isDragging ? undefined : transition,
    willChange: "transform",
    zIndex: isDragging ? 999 : "auto",
  }

  const displayType = overwriteType || field.type;
  const Icon = getFieldTypeIcon(displayType);

  return (
    <div 
      ref={setNodeRef}
      style={style}
      className={cn(
        "flex w-full select-none p-2 bg-card",
        "transition-[box-shadow,border-color,background-color] duration-300 ease-out group",
        "border-2 border-border rounded-sm",
        !isFixed && "shadow-[1px_1px_0_0_var(--color-border)]",
        !isFixed && "hover:shadow-[2px_2px_0_0_var(--color-border)] hover:bg-muted/50",
        isDragging && "opacity-50",
        isFixed && "bg-muted-foreground/10 border-dashed cursor-not-allowed",
        className
      )}
    >
      <div className="flex items-center gap-2 w-full">
        <div className="w-6 h-8 flex items-center justify-center shrink-0">
          <GripVertical 
            {...attributes} 
            {...listeners} 
            className={cn(
              "text-muted-foreground cursor-grab active:cursor-grabbing outline-0 w-6 h-6",
              !isFixed && "hover:text-accent duration-150 transition-colors"
            )}
          />
        </div>
        <Icon className="w-8 h-8 p-2 bg-accent rounded-sm text-accent-foreground shrink-0"/>
        <div className="flex-1 min-w-0">
          <h5 className="font-semibold truncate">{field.title}</h5>
          <div className="text-xs text-muted-foreground flex gap-1 items-center truncate">
            <span className="truncate">{field.key}</span>
            <span className="shrink-0">({displayType})</span>
          </div>
        </div>
        {!isFixed && (onEdit || onDelete || onUpdateField) && (
          <div className="flex gap-1 shrink-0">
            {onEdit && (
              <ShadowButton
                onClick={() => onEdit(field.key)}
                leftIcon={<Pencil className="w-4 h-4" />}
                variant="ghost"
                className="p-1 h-auto"
              />
            )}
            {onOpenEditPanel && (
              <ShadowButton
                leftIcon={<Pencil className="w-4 h-4" />}
                variant="ghost"
                className="p-1 h-auto"
                onClick={() => onOpenEditPanel(field)}
              />
            )}
            {onDelete && (
              <ShadowButton
                onClick={() => onDelete(field.key)}
                leftIcon={<Trash2 className="w-4 h-4" />}
                variant="destructive"
                className="p-1 h-auto"
              />
            )}
          </div>
        )}
      </div>
    </div>
  )
}