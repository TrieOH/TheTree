import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/shared/ui/shadcn/select'
import { cn } from '../lib/class-utils'
import { useMemo } from 'react'

interface PropsI {
  placeholder: string
  options: string[]
  value: string
  onChange: (value: string | null) => void
  disabled?: boolean
}

export default function CustomSelect({
  placeholder,
  options,
  value,
  onChange,
  disabled
}: PropsI) {
  const optionsKey = useMemo(() => {
    return options.join(',')
  }, [options])
  return (
    <Select
      key={optionsKey}
      value={value}
      onValueChange={onChange}
      disabled={disabled || options.length <= 0}
    >
      <SelectTrigger
        className={cn(
          "h-9! w-full rounded-md! border border-input bg-background px-2.5",
          "text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring transition-shadow"
        )}
      >
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>

      <SelectContent>
        {options.map(opt => (
          <SelectItem key={opt} value={opt}>{opt}</SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}