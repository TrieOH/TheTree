import { ChevronLeft, ChevronRight, Pencil } from "lucide-react";
import type { StepI } from "../model";
import { cn } from "#/shared/lib/utils";

interface StepCardProps {
  step: StepI;
  active?: boolean;
  onClick?: (step: StepI) => void;
  onEdit?: (step: StepI) => void;
  onMoveLeft?: (step: StepI) => void;
  onMoveRight?: (step: StepI) => void;
  canMoveLeft?: boolean;
  canMoveRight?: boolean;
  className?: string;
}

export function StepCard({
  step,
  active = false,
  onClick,
  onEdit,
  onMoveLeft,
  onMoveRight,
  canMoveLeft = true,
  canMoveRight = true,
  className,
}: StepCardProps) {
  return (
    <div
      onClick={() => onClick?.(step)}
      onKeyDown={(e) => {
        if (onClick && (e.key === "Enter" || e.key === " ")) {
          e.preventDefault();
          onClick(step);
        }
      }}
      tabIndex={onClick ? 0 : -1}
      role={onClick ? "button" : undefined}
      aria-current={active ? "true" : undefined}
      className={cn(
        "w-full min-h-44 flex flex-col gap-3.5 text-left select-none outline-none",
        "rounded-sm border bg-card px-4.5 py-5 transition-all duration-300 ease-in-out",
        active
          ? "border-primary shadow-[0_4px_28px_rgba(var(--primary),0.13)] opacity-100 scale-100"
          : "border-border opacity-40 scale-[0.96] hover:opacity-55",
        onClick && "cursor-pointer",
        className
      )}
    >
      {/* Header */}
      <div className="flex items-center justify-between w-full">
        <div className="flex items-center gap-1">
          {active && onMoveLeft && (
            <button
              type="button"
              tabIndex={-1}
              onClick={(e) => { e.stopPropagation(); onMoveLeft(step); }}
              disabled={!canMoveLeft}
              className={cn(
                "p-0.5 rounded-xs transition-colors",
                canMoveLeft
                  ? "text-muted-foreground hover:text-primary hover:bg-primary/5 cursor-pointer"
                  : "text-muted-foreground/15 cursor-not-allowed"
              )}
              aria-label="Move step left"
            >
              <ChevronLeft size={14} strokeWidth={2.5} />
            </button>
          )}
          <span
            className={cn(
              "text-[10px] font-bold tracking-[0.13em] uppercase transition-colors duration-300",
              active ? "text-primary" : "text-muted-foreground"
            )}
          >
            Step {String(step.position_hint).padStart(2, "0")}
          </span>
          {active && onMoveRight && (
            <button
              type="button"
              tabIndex={-1}
              onClick={(e) => { e.stopPropagation(); onMoveRight(step); }}
              disabled={!canMoveRight}
              className={cn(
                "p-0.5 rounded-xs transition-colors",
                canMoveRight
                  ? "text-muted-foreground hover:text-primary hover:bg-primary/5 cursor-pointer"
                  : "text-muted-foreground/15 cursor-not-allowed"
              )}
              aria-label="Move step right"
            >
              <ChevronRight size={14} strokeWidth={2.5} />
            </button>
          )}
        </div>

        {active && onEdit && (
          <button
            type="button"
            tabIndex={-1}
            onClick={(e) => { e.stopPropagation(); onEdit(step); }}
            className={cn(
              "p-1 rounded-xs transition-colors",
              "text-muted-foreground hover:text-primary hover:bg-primary/5 cursor-pointer"
            )}
            aria-label="Edit step"
          >
            <Pencil size={14} strokeWidth={2.5} />
          </button>
        )}
      </div>

      {/* Body */}
      <div className="flex-1">
        <p className="text-[17px] font-bold text-foreground leading-snug mb-1">
          {step.title}
        </p>
        {step.description && (
          <p className="text-xs text-muted-foreground leading-relaxed">
            {step.description}
          </p>
        )}
      </div>
    </div>
  );
}
