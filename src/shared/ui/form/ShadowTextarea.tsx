import { cn } from '@/shared/lib/utils';

interface ShadowTextareaProps {
  id?: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
  rows?: number;
  onBlur?: (event: React.FocusEvent<HTMLTextAreaElement>) => void;
}

export const ShadowTextarea: React.FC<ShadowTextareaProps> = ({
  id,
  value,
  onChange,
  onBlur,
  placeholder,
  className,
  disabled,
  rows = 3,
}) => {
  return (
    <div
      className={cn(
        "relative flex items-center rounded-sm border border-input bg-background shadow-[1px_1px_0_0_var(--color-input)] focus-within:shadow-[2px_2px_0_0_var(--color-input)] transition-all duration-300 ease-out",
        disabled && "opacity-50 cursor-not-allowed shadow-none",
        className
      )}
    >
      <textarea
        id={id}
        rows={rows}
        placeholder={placeholder}
        className={cn(
          "h-auto min-h-16 w-full rounded-sm bg-transparent px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none resize-y",
          disabled && "cursor-not-allowed"
        )}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onBlur={onBlur}
        disabled={disabled}
      />
    </div>
  );
};
