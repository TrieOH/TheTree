import { useCallback } from "react";
import { ChevronLeft, ChevronRight, Plus, Pencil } from "lucide-react";
import { DndContext,  PointerSensor, useSensor, useSensors } from "@dnd-kit/core";
import type {DragEndEvent} from "@dnd-kit/core";
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
import type { StepI } from "../model";
import type { FieldI } from "#/features/fields/model";
import SortableFieldRow from "#/features/fields/ui/sortable-field-row";
import { cn } from "#/shared/lib/utils";

interface StepCardProps {
  step: StepI;
  fields?: FieldI[];
  active?: boolean;
  onClick?: (step: StepI) => void;
  onEdit?: (step: StepI) => void;
  onMoveLeft?: (step: StepI) => void;
  onMoveRight?: (step: StepI) => void;
  onAddField?: (step: StepI) => void;
  onEditField?: (field: FieldI) => void;
  onDeleteField?: (field: FieldI) => void;
  onReorderFields?: (step: StepI, fieldIds: string[]) => void;
  onFieldDragChange?: (dragging: boolean) => void;
  canMoveLeft?: boolean;
  canMoveRight?: boolean;
  className?: string;
}

export function StepCard({
  step,
  fields = [],
  active = false,
  onClick,
  onEdit,
  onMoveLeft,
  onMoveRight,
  onAddField,
  onEditField,
  onDeleteField,
  onReorderFields,
  onFieldDragChange,
  canMoveLeft = true,
  canMoveRight = true,
  className,
}: StepCardProps) {
  const sortedFields = [...fields].sort((a, b) => a.position_hint - b.position_hint);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
  );

  const handleDragEnd = useCallback((event: DragEndEvent) => {
    onFieldDragChange?.(false);

    const { active: activeItem, over } = event;
    if (!over || activeItem.id === over.id) return;

    const oldIndex = sortedFields.findIndex(f => f.id === activeItem.id);
    const newIndex = sortedFields.findIndex(f => f.id === over.id);
    if (oldIndex === -1 || newIndex === -1) return;

    // Reorder the sortedFields array
    const reordered = [...sortedFields];
    const [moved] = reordered.splice(oldIndex, 1);
    reordered.splice(newIndex, 0, moved);

    onReorderFields?.(step, reordered.map(f => f.id));
  }, [sortedFields, step, onReorderFields, onFieldDragChange]);

  return (
    <div
      onClick={() => onClick?.(step)}
      onKeyDown={(e) => {
        if (onClick && (e.key === "Enter" || e.key === " ")) {
          e.preventDefault();
          onClick(step);
        }
      }}
      tabIndex={onClick ? 0 : -1}
      role={onClick ? "button" : undefined}
      aria-current={active ? "true" : undefined}
      className={cn(
        "w-full flex flex-col text-left select-none outline-none",
        "rounded-sm border bg-card transition-all duration-300 ease-in-out",
        active
          ? "border-primary shadow-[0_4px_28px_rgba(var(--primary),0.13)] opacity-100"
          : "border-border opacity-40 scale-[0.96] hover:opacity-55",
        onClick && active && "cursor-default",
        onClick && !active && "cursor-pointer",
        className
      )}
    >
      {/* Header */}
      <div className="flex items-center justify-between gap-2 px-4 pt-4 pb-0">
        <div className="flex items-center gap-1">
          {active && onMoveLeft && (
            <button
              type="button"
              tabIndex={-1}
              onClick={(e) => { e.stopPropagation(); onMoveLeft(step); }}
              disabled={!canMoveLeft}
              className={cn(
                "p-0.5 rounded-xs transition-colors",
                canMoveLeft
                  ? "text-muted-foreground hover:text-primary hover:bg-primary/5 cursor-pointer"
                  : "text-muted-foreground/20 cursor-not-allowed"
              )}
              aria-label="Move step left"
            >
              <ChevronLeft size={14} strokeWidth={2.5} />
            </button>
          )}
          <span
            className={cn(
              "text-[10px] font-bold tracking-[0.14em] uppercase transition-colors duration-300",
              active ? "text-primary" : "text-muted-foreground"
            )}
          >
            Step {String(step.position_hint).padStart(2, "0")}
          </span>
          {active && onMoveRight && (
            <button
              type="button"
              tabIndex={-1}
              onClick={(e) => { e.stopPropagation(); onMoveRight(step); }}
              disabled={!canMoveRight}
              className={cn(
                "p-0.5 rounded-xs transition-colors",
                canMoveRight
                  ? "text-muted-foreground hover:text-primary hover:bg-primary/5 cursor-pointer"
                  : "text-muted-foreground/20 cursor-not-allowed"
              )}
              aria-label="Move step right"
            >
              <ChevronRight size={14} strokeWidth={2.5} />
            </button>
          )}
        </div>

        {active && onEdit && (
          <button
            type="button"
            tabIndex={-1}
            onClick={(e) => { e.stopPropagation(); onEdit(step); }}
            className={cn(
              "p-1 rounded-xs transition-colors",
              "text-muted-foreground hover:text-primary hover:bg-primary/5 cursor-pointer"
            )}
            aria-label="Edit step"
          >
            <Pencil size={13} strokeWidth={2.5} />
          </button>
        )}
      </div>

      {/* Title + description */}
      <div className="px-4 pt-2 pb-4">
        <p className="text-[18px] font-bold text-foreground leading-snug">
          {step.title}
        </p>
        {step.description && (
          <p className="text-xs text-muted-foreground leading-relaxed mt-1">
            {step.description}
          </p>
        )}
      </div>
      {/* Fields table - only when active */}
      {active && (
        <div
          onClick={(e) => e.stopPropagation()}
          className="border-t border-border"
        >
          {/* Rows */}
          <DndContext
            sensors={sensors}
            onDragStart={() => onFieldDragChange?.(true)}
            onDragEnd={handleDragEnd}
          >
            <SortableContext
              items={sortedFields.map(f => f.id)}
              strategy={verticalListSortingStrategy}
            >
              <div className="flex flex-col divide-y divide-border/40">
                {sortedFields.length === 0 ? (
                  <p className="px-4 py-5 text-center text-[11px] text-muted-foreground/40 italic">
                    No fields yet
                  </p>
                ) : (
                  sortedFields.map((field) => (
                    <SortableFieldRow
                      key={field.id}
                      field={field}
                      onEdit={onEditField}
                      onDelete={onDeleteField}
                    />
                  ))
                )}
              </div>
            </SortableContext>
          </DndContext>

          {/* Add field */}
          {onAddField && (
            <div className="border-t border-border/50 px-4 py-3">
              <button
                type="button"
                tabIndex={-1}
                onClick={() => onAddField(step)}
                className="flex items-center gap-1.5 text-muted-foreground/50 hover:text-primary transition-colors duration-150 cursor-pointer"
              >
                <div className="w-4 h-4 rounded-full border border-current flex items-center justify-center shrink-0">
                  <Plus size={9} strokeWidth={3} />
                </div>
                <span className="text-[11px] font-medium">Add Field</span>
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}