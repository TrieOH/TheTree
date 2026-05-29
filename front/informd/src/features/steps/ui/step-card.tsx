import { FileText } from "lucide-react";
import type { StepI } from "../model";
import { cn } from "#/shared/lib/utils";

interface StepCardProps {
  step: StepI;
  active?: boolean;
  onClick?: (step: StepI) => void;
  className?: string;
}

export function StepCard({ step, active = false, onClick, className }: StepCardProps) {
  return (
    <button
      type="button"
      onClick={() => onClick?.(step)}
      aria-current={active ? "true" : undefined}
      className={cn(
        "w-full min-h-44 flex flex-col gap-3.5 text-left select-none outline-none",
        "rounded-sm border bg-card px-4.5 py-5 transition-all duration-300 ease-in-out",
        active
          ? "border-primary shadow-[0_4px_28px_rgba(var(--primary),0.13)] opacity-100 scale-100 cursor-pointer"
          : "border-border opacity-40 scale-[0.96] cursor-pointer hover:opacity-55",
        className
      )}
    >
      {/* Header */}
      <div className="flex items-center justify-between w-full">
        <span
          className={cn(
            "text-[10px] font-bold tracking-[0.13em] uppercase transition-colors duration-300",
            active ? "text-primary" : "text-muted-foreground"
          )}
        >
          Step {String(step.position_hint).padStart(2, "0")}
        </span>

        <span
          aria-hidden="true"
          className={cn(
            "flex items-center transition-colors duration-300",
            active ? "text-primary" : "text-muted-foreground/50"
          )}
        >
          <FileText size={16} strokeWidth={1.6} />
        </span>
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

      {/* Decorative lines */}
      <div className="w-full flex flex-col gap-1.5">
        <div
          className={cn(
            "h-px rounded-sm w-full transition-colors duration-300",
            active ? "bg-primary/20" : "bg-border"
          )}
        />
        <div
          className={cn(
            "h-px rounded-sm w-[58%] transition-colors duration-300",
            active ? "bg-primary/20" : "bg-border"
          )}
        />
      </div>
    </button>
  );
}
