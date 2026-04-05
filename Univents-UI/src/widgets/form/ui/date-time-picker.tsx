import { Calendar as CalendarIcon, Clock, ChevronUp, ChevronDown } from "lucide-react"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { useEffect, useRef, useState } from "react"
import { motion, AnimatePresence } from "motion/react"
import { cn } from "@/shared/lib/utils"
import { Button } from "@/shared/ui/shadcn/button"
import { Calendar } from "@/shared/ui/shadcn/calendar"

interface DateTimePickerProps {
  id?: string
  value?: string
  onChange: (value: string) => void
  disabled?: boolean
  className?: string
  placeholder?: string
  min?: string
  max?: string
  error?: boolean
}

function useScrollSpin(onDelta: (delta: number) => void, isDisabled: boolean) {
  const ref = useRef<HTMLInputElement>(null)

  useEffect(() => {
    const el = ref.current
    if (!el) return

    const handler = (e: WheelEvent) => {
      if (isDisabled) return
      e.preventDefault()
      requestAnimationFrame(() => { onDelta(e.deltaY < 0 ? 1 : -1); })
    }

    el.addEventListener("wheel", handler, { passive: false })
    return () => { el.removeEventListener("wheel", handler); }
  }, [onDelta, isDisabled])

  return ref
}

export function DateTimePicker({
  id,
  value,
  onChange,
  disabled,
  className,
  placeholder = "Selecione data e hora",
  min,
  max,
  error,
}: DateTimePickerProps) {
  const [isOpen, setIsOpen] = useState(false)
  const dateValue = value ? new Date(value) : undefined

  const [hours, setHours] = useState(dateValue ? dateValue.getHours() : 12)
  const [minutes, setMinutes] = useState(dateValue ? dateValue.getMinutes() : 0)

  useEffect(() => {
    if (dateValue) {
      setHours(dateValue.getHours())
      setMinutes(dateValue.getMinutes())
    }
  }, [value])

  const commitChange = (selectedDate: Date | undefined, h: number, m: number) => {
    if (!selectedDate) {
      onChange("")
      return
    }
    const newDate = new Date(selectedDate)
    newDate.setHours(h, m, 0, 0)
    onChange(newDate.toISOString())
  }

  const handleSelectDate = (selectedDate: Date | undefined) => {
    commitChange(selectedDate, hours, minutes)
  }

  const handleHoursChange = (delta: number) => {
    const newH = (hours + delta + 24) % 24
    setHours(newH)
    if (dateValue) commitChange(dateValue, newH, minutes)
  }

  const handleMinutesChange = (delta: number) => {
    const newM = (minutes + delta + 60) % 60
    setMinutes(newM)
    if (dateValue) commitChange(dateValue, hours, newM)
  }

  const handleHoursInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = Math.max(0, Math.min(23, parseInt(e.target.value) || 0))
    setHours(val)
    if (dateValue) commitChange(dateValue, val, minutes)
  }

  const handleMinutesInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = Math.max(0, Math.min(59, parseInt(e.target.value) || 0))
    setMinutes(val)
    if (dateValue) commitChange(dateValue, hours, val)
  }

  const handleNow = () => {
    const now = new Date()
    setHours(now.getHours())
    setMinutes(now.getMinutes())
    const newDate = new Date(now)
    newDate.setSeconds(0, 0)
    onChange(newDate.toISOString())
  }

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation()
    onChange("")
    setHours(12)
    setMinutes(0)
  }

  const displayValue = dateValue
    ? format(dateValue, "dd/MM/yyyy 'às' HH:mm", { locale: ptBR })
    : ""

  const minDate = min ? new Date(min) : undefined
  const maxDate = max ? new Date(max) : undefined

  const hoursRef = useScrollSpin(handleHoursChange, !dateValue)
  const minutesRef = useScrollSpin(handleMinutesChange, !dateValue)

  const spinInput = cn(
    "w-10 h-9 text-center text-base font-bold rounded-lg border bg-background transition-all",
    "focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary",
    "disabled:opacity-40 disabled:cursor-not-allowed",
    "[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none",
    "cursor-ns-resize select-none"
  )

  const spinBtn = cn(
    "p-0.5 text-muted-foreground hover:text-primary disabled:opacity-30 disabled:pointer-events-none transition-colors"
  )

  return (
    <div className={cn("flex flex-col", className)}>
      <button
        id={id}
        type="button"
        disabled={disabled}
        onClick={() => { setIsOpen(!isOpen); }}
        className={cn(
          "w-full px-3.5 py-2 border bg-background text-sm transition-all",
          "hover:bg-accent/5 group flex items-center gap-3 text-left",
          "focus:outline-none focus:ring-2 focus:ring-primary/20",
          error ? "border-destructive" : "border-input focus:border-primary",
          "disabled:opacity-50 disabled:cursor-not-allowed",
          isOpen ? "rounded-t-xl border-b-transparent" : "rounded-xl"
        )}
      >
        <div className={cn(
          "p-1.5 rounded-lg transition-colors",
          dateValue ? "bg-primary/10 text-primary" : "bg-muted text-muted-foreground"
        )}>
          <CalendarIcon className="h-4 w-4 shrink-0" />
        </div>
        <div className="flex flex-col flex-1 truncate">
          {!dateValue ? (
            <span className="text-muted-foreground font-medium">{placeholder}</span>
          ) : (
            <div className="flex flex-col leading-tight">
              <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-tight">Data e Hora</span>
              <span className="text-foreground font-semibold">{displayValue}</span>
            </div>
          )}
        </div>
        <ChevronDown className={cn(
          "h-4 w-4 text-muted-foreground transition-transform duration-200",
          isOpen && "rotate-180"
        )} />
      </button>

      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2, ease: "easeInOut" }}
            className="border-x border-b rounded-b-xl bg-card overflow-hidden"
          >
            <div className="p-3 flex flex-col gap-4">
              <div className="flex flex-col sm:flex-row gap-4">
                <div className="flex-1 flex justify-center">
                  <Calendar
                    mode="single"
                    selected={dateValue}
                    onSelect={handleSelectDate}
                    disabled={(date) => {
                      if (minDate && date < new Date(new Date(minDate).setHours(0, 0, 0, 0))) return true
                      if (maxDate && date > new Date(new Date(maxDate).setHours(23, 59, 59, 999))) return true
                      return false
                    }}
                    className="p-0 scale-95 origin-top"
                  />
                </div>

                <div className="flex flex-col items-center justify-center gap-2 p-3 border-t sm:border-t-0 sm:border-l border-border bg-muted/20 rounded-xl min-w-35">
                  <div className="flex items-center gap-1.5 text-[10px] font-bold text-muted-foreground uppercase tracking-wider">
                    <Clock className="h-3 w-3" />
                    <span>Horário</span>
                  </div>

                  <div className="flex items-center gap-1">
                    <div className="flex flex-col items-center">
                      <button type="button" onClick={() => { handleHoursChange(1); }} className={spinBtn} disabled={!dateValue}>
                        <ChevronUp className="h-4 w-4" />
                      </button>
                      <input
                        ref={hoursRef}
                        type="number"
                        min={0}
                        max={23}
                        value={String(hours).padStart(2, "0")}
                        onChange={handleHoursInput}
                        disabled={!dateValue}
                        className={spinInput}
                      />
                      <button type="button" onClick={() => { handleHoursChange(-1); }} className={spinBtn} disabled={!dateValue}>
                        <ChevronDown className="h-4 w-4" />
                      </button>
                    </div>

                    <span className="font-bold text-muted-foreground/30 mb-1">:</span>

                    <div className="flex flex-col items-center">
                      <button type="button" onClick={() => { handleMinutesChange(1); }} className={spinBtn} disabled={!dateValue}>
                        <ChevronUp className="h-4 w-4" />
                      </button>
                      <input
                        ref={minutesRef}
                        type="number"
                        min={0}
                        max={59}
                        value={String(minutes).padStart(2, "0")}
                        onChange={handleMinutesInput}
                        disabled={!dateValue}
                        className={spinInput}
                      />
                      <button type="button" onClick={() => { handleMinutesChange(-1); }} className={spinBtn} disabled={!dateValue}>
                        <ChevronDown className="h-4 w-4" />
                      </button>
                    </div>
                  </div>

                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={handleNow}
                    className="w-full h-7 text-[10px] font-bold uppercase tracking-wider hover:bg-primary/5 hover:text-primary transition-all"
                  >
                    Agora
                  </Button>
                </div>
              </div>

              <div className="flex gap-2 pt-2 border-t border-border/50">
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={handleClear}
                  className="flex-1 h-8 text-xs font-bold text-destructive hover:bg-destructive/10"
                >
                  LIMPAR
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={() => { setIsOpen(false); }}
                  className="flex-1 h-8 text-xs font-bold rounded-lg"
                >
                  CONCLUIR
                </Button>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
