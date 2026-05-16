import { PencilLine } from 'lucide-react';
import type { InlineEditProps } from './types';
import { cn } from '@/shared/lib/utils';

const InlineEdit = ({
  value,
  onChange,
  isEditEnabled,
  isEditing,
  onStartEdit,
  onFinishEdit,
  multiline = false,
  className = '',
  placeholder = 'Clique para editar...',
}: InlineEditProps) => {

  if (!isEditEnabled) return <span className={className}>{value}</span>

  if (isEditing) {
    return (
      <span
        contentEditable
        suppressContentEditableWarning
        onBlur={(e) => {
          onChange(e.currentTarget.innerText);
          onFinishEdit();
        }}
        onKeyDown={(e) => {
          if (e.key === 'Enter' && !multiline) {
            e.preventDefault();
            onChange(e.currentTarget.innerText);
            onFinishEdit();
          }
          if (e.key === 'Escape') onFinishEdit();
        }}
        className={cn(
          "relative outline-none border-b-2 border-primary/60",
          "bg-primary/5 px-2 py-1 rounded-sm whitespace-pre-wrap",
          multiline ? "block min-h-15" : "min-h-auto inline-block",
          className
        )}
      >
        {value}
      </span>
    );
  }

  return (
    <button
      onClick={(e) => {
        e.stopPropagation();
        onStartEdit();
      }}
      className={cn(
        "group relative inline-flex items-start gap-2 text-left",
        "rounded-md px-3 py-2 transition-all cursor-text",
        "hover:bg-muted/50 border border-dashed border-muted-foreground/25",
        "hover:border-muted-foreground/40",
        className
      )}
    >
      <span className={cn(
        "flex-1",
        !value && "text-muted-foreground italic"
      )}>
        {value ?? placeholder}
      </span>

      <span className={cn(
        "absolute -top-2 -right-2",
        "flex items-center justify-center",
        "h-6 w-6 rounded-md",
        "bg-background border border-muted-foreground/20",
        "shadow-sm cursor-pointer",
        "transition-all duration-200",
        "group-hover:bg-primary group-hover:text-primary-foreground group-hover:border-primary",
        "group-hover:scale-110 group-hover:-translate-y-0.5 group-hover:translate-x-0.5"
      )}>
        <PencilLine className="h-3 w-3" />
      </span>
    </button>
  );
};

export default InlineEdit;