import type React from "react";
import { cn } from "#/shared/lib/utils";

interface OptionPickerProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  options: {
    label: string;
    value: string | number;
    icon?: React.ComponentType<{ className?: string }>;
  }[];
  required?: boolean;
  error?: string;
}

export default function OptionPicker({
  label,
  value,
  onChange,
  options,
  required,
  error,
}: OptionPickerProps) {
  return (
    <div className="mt-2 flex w-full flex-col gap-1.5">
      <span className="text-xs font-semibold text-foreground/80">
        {required ? `${label} *` : label}
      </span>
      <div className="grid grid-cols-3 gap-2">
        {options.map((option) => {
          const isSelected = value === String(option.value);
          const Icon = option.icon;

          return (
            <button
              key={String(option.value)}
              type="button"
              onClick={() => onChange(String(option.value))}
              className={cn(
                "flex cursor-pointer select-none flex-col items-center justify-center rounded-sm border p-3 text-center transition-all duration-300",
                "active:translate-x-px active:translate-y-px",
                isSelected
                  ? "border-primary bg-primary/10 text-foreground ring-1 ring-primary"
                  : "border-input bg-card text-muted-foreground hover:bg-muted/50 hover:text-foreground",
              )}
            >
              {Icon && (
                <Icon
                  className={cn(
                    "mb-1.5 size-5",
                    isSelected ? "text-primary" : "text-muted-foreground",
                  )}
                />
              )}
              <span className="text-xs font-medium">{option.label}</span>
            </button>
          );
        })}
      </div>
      {error && (
        <span className="mt-0.5 text-[10px] text-destructive">{error}</span>
      )}
    </div>
  );
}
