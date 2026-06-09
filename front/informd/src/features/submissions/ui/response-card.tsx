import { Mail, Clock, Layers } from "lucide-react";
import type { FullStepI, SubmissionSummaryI } from "../model";
import { cn } from "#/shared/lib/utils";

interface ResponseCardProps {
  data: SubmissionSummaryI;
  steps: FullStepI[];
  isSelected: boolean;
  onClick: () => void;
}

export default function ResponseCard({ data, steps, isSelected, onClick }: ResponseCardProps) {
  const getStepName = (stepId: string) => {
    const step = steps.find((s) => s.step.id === stepId);
    return step?.step.title ?? stepId;
  };

  return (
    <div
      onClick={onClick}
      className={cn(
        "group relative flex flex-col md:flex-row md:items-center",
        "justify-between gap-3 md:gap-4 p-4 md:p-5 rounded-xl border transition-all cursor-pointer",
        isSelected
          ? "bg-secondary/10 border-secondary/30 shadow-sm ring-1 ring-secondary/20 md:bg-accent/5 md:border-accent/50 md:shadow-md md:ring-0"
          : "bg-card border-border hover:border-accent/50 hover:shadow-md hover:bg-accent/5"
      )}
    >
      {/* Main Info */}
      <div className="flex items-center gap-3 md:gap-4 flex-1 min-w-0">
        {/* Avatar - desktop only */}
        <div
          className={cn(
            "hidden md:flex shrink-0 w-12 h-12 rounded-full items-center justify-center text-base font-bold transition-colors",
            isSelected
              ? "bg-primary text-primary-foreground"
              : "bg-muted text-muted-foreground group-hover:bg-accent/20 group-hover:text-accent"
          )}
        >
          {data.responder.charAt(0).toUpperCase()}
        </div>

        <div className="min-w-0 flex-1 space-y-1">
          <div className="flex items-center gap-2">
            <Mail
              className={cn(
                "w-4 h-4 shrink-0",
                isSelected ? "text-secondary md:text-accent" : "text-muted-foreground"
              )}
            />
            <span className="text-sm font-semibold text-foreground truncate">
              {data.responder}
            </span>
          </div>

          <div className="flex flex-wrap items-center gap-x-4 gap-y-1">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Clock className="w-3.5 h-3.5 shrink-0" />
              <span className="tabular-nums">
                {new Date(data.completed_at).toLocaleString(undefined, {
                  month: "short",
                  day: "numeric",
                  year: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </span>
            </div>

            <div className="flex items-center gap-1.5 text-xs">
              <Layers className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
              <span
                className={cn(
                  "px-2 py-0.5 rounded-md font-medium border",
                  isSelected
                    ? "bg-secondary/20 border-secondary/30 text-secondary-foreground md:bg-muted md:border-border md:text-muted-foreground md:group-hover:border-accent/30 md:group-hover:text-accent"
                    : "bg-muted border-border text-muted-foreground group-hover:border-accent/30 group-hover:text-accent"
                )}
              >
                {getStepName(data.step_id)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}