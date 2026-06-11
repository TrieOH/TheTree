import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { cn } from "#/shared/lib/utils";
import FieldRow from "./field-row";
import type { FieldI } from "../model";

interface SortableFieldRowProps {
  field: FieldI;
  onEdit?: (field: FieldI) => void;
  onDelete?: (field: FieldI) => void;
}

export default function SortableFieldRow({
  field,
  onEdit,
  onDelete,
}: SortableFieldRowProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    setActivatorNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: field.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        "relative",
        isDragging && "z-50 opacity-60",
      )}
    >
      <FieldRow
        field={field}
        onEdit={onEdit}
        onDelete={onDelete}
        dragHandleRef={setActivatorNodeRef}
        dragHandleProps={{ ...attributes, ...listeners }}
      />
    </div>
  );
}
