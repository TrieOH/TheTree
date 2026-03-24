import { useEffect, useMemo, useRef, useState } from "react"
import { Clock, AlertCircle } from "lucide-react"
import { cn } from "@/shared/lib/utils"

interface TimerProps {
  expiresAt: string | null
  warningThresholdSeconds?: number
  onExpire?: () => void
  label?: string
  className?: string
}

export function Timer({
  expiresAt,
  warningThresholdSeconds = 60,
  onExpire,
  label = "Reserva expira em",
  className = ""
}: TimerProps) {
  const [secondsLeft, setSecondsLeft] = useState(0)
  const onExpireRef = useRef(onExpire)
  const hasExpiredRef = useRef(false)

  onExpireRef.current = onExpire

  const calculateSeconds = () => {
    if (!expiresAt) return 0
    return Math.max(0, Math.floor((new Date(expiresAt).getTime() - Date.now()) / 1000))
  }

  useEffect(() => {
    hasExpiredRef.current = false
    const initial = calculateSeconds()
    setSecondsLeft(initial)

    if (initial === 0) {
      hasExpiredRef.current = true
      onExpireRef.current?.()
      return
    }

    const interval = setInterval(() => {
      const left = calculateSeconds()
      setSecondsLeft(left)

      if (left === 0) {
        clearInterval(interval)
        if (!hasExpiredRef.current) {
          hasExpiredRef.current = true
          onExpireRef.current?.()
        }
      }
    }, 1000)

    return () => { clearInterval(interval); }
  }, [expiresAt])

  const formatted = useMemo(() => {
    const m = Math.floor(secondsLeft / 60).toString().padStart(2, "0")
    const s = (secondsLeft % 60).toString().padStart(2, "0")
    return `${m}:${s}`
  }, [secondsLeft])

  const isWarning = secondsLeft < warningThresholdSeconds && secondsLeft > 0
  const isExpired = secondsLeft === 0

  if (!expiresAt) return null

  return (
    <div className={`w-full min-w-75 ${className}`}>
      <div
        className={cn(
          "flex items-center justify-between px-3 py-2.5 border rounded-xs",
          isExpired ? "bg-destructive border-destructive text-destructive-foreground" :
            isWarning ? "bg-accent border-accent text-accent-foreground" :
              "bg-primary/10 border-primary/20"
        )}
      >
        <div className="flex items-center gap-2">
          {isExpired ? (
            <AlertCircle className="w-4 h-4" />
          ) : (
            <Clock className={`w-4 h-4 ${isWarning && "animate-pulse"}`} />
          )}
          <span className="text-xs font-medium uppercase tracking-wide">
            {isExpired ? "Expirado" : label}
          </span>
        </div>

        <span
          className={cn(
            "text-sm font-bold tabular-nums",
            isExpired ? "text-destructive-foreground" :
              isWarning ? "text-accent-foreground" :
                "text-primary"
          )}
        >
          {isExpired ? "00:00" : formatted}
        </span>
      </div>
    </div>
  )
}