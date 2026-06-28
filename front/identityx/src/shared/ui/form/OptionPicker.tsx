import { cn } from "@/shared/lib/utils"
import type { FieldOption } from "./types"

interface OptionPickerProps {
  label: string
  value: string
  onChange: (value: string) => void
  options: FieldOption[]
  required?: boolean
  error?: string
}

export default function OptionPicker({
  label,
  value,
  onChange,
  options,
  required,
  error
}: OptionPickerProps) {
  return (
    <div className="flex flex-col gap-1.5 w-full mt-2">
      <span className="text-xs font-semibold text-foreground/80">
        {required ? `${label} *` : label}
      </span>
      <div className="grid grid-cols-3 gap-2">
        {options.map((option) => {
          const isSelected = value === option.value
          const Icon = option.icon

          return (
            <button
              key={option.value}
              type="button"
              onClick={() => onChange(option.value)}
              className={cn(
                "flex flex-col items-center justify-center p-3 rounded-sm border text-center transition-all duration-300 cursor-pointer select-none",
                "active:translate-x-px active:translate-y-px",
                isSelected
                  ? "border-primary bg-primary/10 text-foreground ring-1 ring-primary"
                  : "border-input bg-card text-muted-foreground hover:bg-muted/50 hover:text-foreground"
              )}
            >
              {Icon && <Icon className={cn("size-5 mb-1.5", isSelected ? "text-primary" : "text-muted-foreground")} />}
              <span className="text-xs font-medium">{option.label}</span>
            </button>
          )
        })}
      </div>
      {error && <span className="text-[10px] text-destructive mt-0.5">{error}</span>}
    </div>
  )
}
