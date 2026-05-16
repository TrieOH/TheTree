import { Search, X } from 'lucide-react';
import { cn } from '@/shared/lib/utils';

interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
}

export const SearchInput: React.FC<SearchInputProps> = ({
  value,
  onChange,
  placeholder = 'Search...',
  className,
  disabled,
}) => {
  const handleClear = () => {
    onChange('');
  };

  return (
    <div
      className={cn(
        "relative flex items-center rounded-sm border border-input bg-background shadow-[1px_1px_0_0_var(--color-input)] focus-within:shadow-[2px_2px_0_0_var(--color-input)] transition-all duration-300 ease-out",
        disabled && "opacity-50 cursor-not-allowed shadow-none",
        className
      )}
    >
      <Search size={16} className="absolute left-3 text-muted-foreground" />
      <input
        type="text"
        placeholder={placeholder}
        className={cn(
          "h-9 w-full rounded-sm bg-transparent pl-9 pr-8 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none",
          disabled && "cursor-not-allowed"
        )}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
      />
      {value && (
        <button
          type="button"
          onClick={handleClear}
          className="absolute right-2 p-1 text-muted-foreground hover:text-foreground disabled:cursor-not-allowed"
          disabled={disabled}
          aria-label="Clear search"
        >
          <X size={16} />
        </button>
      )}
    </div>
  );
};
