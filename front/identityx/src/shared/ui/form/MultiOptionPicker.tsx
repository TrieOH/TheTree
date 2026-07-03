import { cn } from "@/shared/lib/utils"
import { Check } from "lucide-react"
import { useFieldContext } from "@/shared/lib/forms"
import { useMemo, useState } from "react"
import type { FieldOption, RuleStatus } from "./types"

interface PropsI {
  label: string
  value: string[]
  onChange: (value: string[]) => void
  options: FieldOption[]
  required?: boolean
  error?: string
  getRulesStatus?: (value: unknown) => RuleStatus[]
  submitted?: boolean
}

export default function MultiOptionPicker({
  label,
  value,
  onChange,
  options,
  required,
  error,
}: PropsI) {
  const field = useFieldContext<string[]>()
  const [filter, setFilter] = useState('')

  const selected = Array.isArray(value) ? value : []
  const filteredOptions = useMemo(() => {
    const search = filter.toLowerCase().trim()
    if (!search) return options
    return options.filter((option) => {
      return option.label.toLowerCase().includes(search) || option.value.toLowerCase().includes(search)
    })
  }, [filter, options])

  return (
    <div className="flex flex-col gap-1.5 w-full mt-2">
      <span className="text-xs font-semibold text-foreground/80">
        {required ? `${label} *` : label}
      </span>
      <input
        value={filter}
        onChange={(e) => setFilter(e.target.value)}
        placeholder={`Search ${label.toLowerCase()}...`}
        className="h-9 w-full rounded-sm border border-input bg-background px-3 text-sm outline-none placeholder:text-muted-foreground focus:ring-1 focus:ring-primary"
      />
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 max-h-72 overflow-auto pr-1">
        {filteredOptions.map((option) => {
          const isSelected = selected.includes(option.value)
          const Icon = option.icon

          return (
            <button
              key={option.value}
              type="button"
              onClick={() => {
                if (isSelected) onChange(selected.filter((item) => item !== option.value))
                else onChange([...selected, option.value])
                field.handleBlur()
              }}
              className={cn(
                "flex items-center gap-2.5 p-2.5 rounded-sm border text-left transition-all duration-300 cursor-pointer select-none",
                "active:translate-x-px active:translate-y-px",
                isSelected
                  ? "border-primary bg-primary/10 text-foreground ring-1 ring-primary"
                  : "border-input bg-card text-muted-foreground hover:bg-muted/50 hover:text-foreground"
              )}
            >
              <div className={cn(
                "shrink-0 size-4.5 rounded-full border flex items-center justify-center",
                isSelected ? "border-primary bg-primary text-primary-foreground" : "border-input bg-background"
              )}>
                <Check className={cn("size-2.5", isSelected ? "opacity-100" : "opacity-0")} />
              </div>
              {Icon && <Icon className={cn("size-3.5", isSelected ? "text-primary" : "text-muted-foreground")} />}
              <span className="text-[11px] font-medium flex-1 leading-tight">{option.label}</span>
            </button>
          )
        })}
        {filteredOptions.length === 0 && (
          <div className="rounded-sm border border-dashed border-input px-3 py-4 text-xs text-muted-foreground">
            No options found.
          </div>
        )}
      </div>
      {error && <span className="text-[10px] text-destructive mt-0.5">{error}</span>}
    </div>
  )
}
