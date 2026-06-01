import { GripVertical, Settings, Trash2, Text, Mail, Hash, ToggleLeft, Calendar, Clock, CalendarClock, ChevronDown, File, Phone, Link } from "lucide-react";
import type { FieldI } from "../model";
import { cn } from "#/shared/lib/utils";
import { Button } from "#/shared/ui/shadcn/button";

interface FieldRowProps {
  field: FieldI;
  onEdit?: (field: FieldI) => void;
  onDelete?: (field: FieldI) => void;
  dragHandleRef?: React.Ref<HTMLButtonElement>;
  dragHandleProps?: Record<string, unknown>;
}

export default function FieldRow({ field, onEdit, onDelete, dragHandleRef, dragHandleProps }: FieldRowProps) {
  return (
    <div className="flex items-center gap-2 px-3.5 py-2.5 hover:bg-muted/20 transition-colors duration-100 group">
      <button
        ref={dragHandleRef}
        type="button"
        tabIndex={-1}
        {...(dragHandleProps as React.HTMLAttributes<HTMLButtonElement>)}
        className="shrink-0 inline-flex cursor-grab active:cursor-grabbing touch-none"
        aria-label={`Drag ${field.title}`}
      >
        <GripVertical
          size={20}
          strokeWidth={2}
          className="text-muted-foreground/20 group-hover:text-muted-foreground/50 transition-colors"
        />
      </button>

      <TypeIcon type={field.type} />

      {/* label + key */}
      <div className="flex items-baseline gap-1.5 flex-1 min-w-0 overflow-hidden">
        <span className="text-[12px] font-medium text-foreground whitespace-nowrap overflow-hidden text-ellipsis shrink-0 max-w-[55%]">
          {field.title}
        </span>
        <span className="font-mono text-[10px] text-muted-foreground/50 whitespace-nowrap overflow-hidden text-ellipsis min-w-0">
          {field.key}
        </span>
      </div>

      {field.required && (
        <span className="shrink-0 text-[9px] font-bold tracking-widest uppercase text-primary/60 leading-none">
          req
        </span>
      )}

      <div className="flex items-center shrink-0">
        {onEdit && (
          <Button
            variant="ghost"
            tabIndex={-1}
            onClick={() => onEdit(field)}
            className="p-1 rounded-xs text-muted-foreground/40 hover:text-primary hover:bg-primary/5 transition-colors cursor-pointer"
            aria-label={`Edit ${field.title}`}
          >
            <Settings size={13} strokeWidth={2} />
          </Button>
        )}
        {onDelete && (
          <Button
            variant="ghost"
            tabIndex={-1}
            onClick={() => onDelete(field)}
            className="p-1 rounded-xs text-muted-foreground/40 hover:text-destructive hover:bg-destructive/5 transition-colors cursor-pointer"
            aria-label={`Delete ${field.title}`}
          >
            <Trash2 size={13} strokeWidth={2} />
          </Button>
        )}
      </div>
    </div>
  );
}

function TypeIcon({ type }: { type: string }) {
  const Icon = TYPE_ICON[type] ?? Text;
  const label = type.charAt(0).toUpperCase() + type.slice(1);
  return (
    <span
      title={label}
      aria-label={label}
      className={cn("shrink-0 inline-flex", TYPE_COLOR[type] ?? "text-muted-foreground/30")}
    >
      <Icon size={13} strokeWidth={2} />
    </span>
  );
}

const TYPE_ICON: Record<string, React.ComponentType<{ size?: number; strokeWidth?: number; title?: string; "aria-label"?: string; className?: string }>> = {
  string: Text,
  email: Mail,
  int: Hash,
  float: Hash,
  bool: ToggleLeft,
  date: Calendar,
  time: Clock,
  datetime: CalendarClock,
  select: ChevronDown,
  file: File,
  phone: Phone,
  url: Link,
};

const TYPE_COLOR: Record<string, string> = {
  string: "text-blue-400 dark:text-blue-500",
  email: "text-violet-400 dark:text-violet-500",
  int: "text-amber-400 dark:text-amber-500",
  float: "text-amber-400 dark:text-amber-500",
  bool: "text-emerald-400 dark:text-emerald-500",
  date: "text-sky-400 dark:text-sky-500",
  time: "text-sky-400 dark:text-sky-500",
  datetime: "text-sky-400 dark:text-sky-500",
  select: "text-fuchsia-400 dark:text-fuchsia-500",
  file: "text-orange-400 dark:text-orange-500",
  phone: "text-teal-400 dark:text-teal-500",
  url: "text-rose-400 dark:text-rose-500",
};