import { cn } from '@/shared/lib/utils';

interface ShadowInputProps {
  id?: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
  type?: string;
}

export const ShadowInput: React.FC<ShadowInputProps> = ({
  id,
  value,
  onChange,
  placeholder,
  className,
  disabled,
  type = "text",
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
        placeholder={placeholder}
        className={cn(
          "h-9 w-full rounded-sm bg-transparent px-3 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none", // Adjusted padding
          disabled && "cursor-not-allowed"
        )}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
      />
    </div>
  );
};
