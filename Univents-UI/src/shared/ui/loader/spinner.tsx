import styles from "./spinner.module.css"
import { cn } from "@/shared/lib/utils"

interface SpinnerProps {
  size?: number | string
  activeColor?: string
  trackColor?: string
  duration?: string
  className?: string
}

export function Spinner({
  size = "3rem",
  activeColor = "#7627a3",
  trackColor = "#f2d4fe",
  duration = "8s",
  className,
}: SpinnerProps) {
  return (
    <svg
      viewBox="0 0 384 384"
      style={
        {
          "--size": typeof size === "number" ? `${size}px` : size,
          "--active": activeColor,
          "--track": trackColor,
          "--duration": duration,
        } as React.CSSProperties
      }
      className={cn(
        "origin-center overflow-visible -rotate-90 animate-[spin_2s_linear_infinite]",
        styles.loader, className
      )}
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle
        className={styles.active}
        pathLength={360}
        fill="transparent"
        strokeWidth={32}
        cx={192}
        cy={192}
        r={176}
      />

      <circle
        className={styles.track}
        pathLength={360}
        fill="transparent"
        strokeWidth={32}
        cx={192}
        cy={192}
        r={176}
      />
    </svg>
  )
}