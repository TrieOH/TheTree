import { cn } from '@/shared/lib/utils';

interface ShadowInputProps {
  id?: string;
  value: string;
  onChange: (value: string) => void;
  onKeyDown?: (e: React.KeyboardEvent<HTMLInputElement> | React.FocusEvent<HTMLInputElement>) => void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
  type?: string;
  onBlur?: (event: React.FocusEvent<HTMLInputElement>) => void;
  inputRef?: React.Ref<HTMLInputElement>;
}

export const ShadowInput: React.FC<ShadowInputProps> = ({
  id,
  value,
  onChange,
  onKeyDown,
  onBlur,
  placeholder,
  className,
  disabled,
  type = "text",
  inputRef
}) => {
  return (
    <div
      className={cn(
        "relative flex items-center rounded-sm border border-input bg-background shadow-[1px_1px_0_0_var(--color-input)] focus-within:shadow-[2px_2px_0_0_var(--color-input)] transition-all duration-300 ease-out",
        disabled && "opacity-50 cursor-not-allowed shadow-none",
        className
      )}
    >
      <input
        id={id}
        type={type}
        ref={inputRef}
        placeholder={placeholder}
        className={cn(
          "h-9 w-full rounded-sm bg-transparent px-3 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none", // Adjusted padding
          disabled && "cursor-not-allowed"
        )}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={onKeyDown}
        onBlur={onBlur}
        disabled={disabled}
      />
    </div>
  );
};
