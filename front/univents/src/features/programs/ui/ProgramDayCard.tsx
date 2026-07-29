import { Activity, ChevronRight, LogIn } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import type { OccurrenceI, ProgramI } from "../model";

function formatTimeRange(start: string, end: string): string {
  const fmt = (iso: string) =>
    new Date(iso).toLocaleTimeString("pt-BR", {
      hour: "2-digit",
      minute: "2-digit",
    });
  return `${fmt(start)} – ${fmt(end)}`;
}

function formatDayLabel(dateIso: string): string {
  const date = new Date(dateIso);
  const today = new Date();
  const tomorrow = new Date(today);
  tomorrow.setDate(tomorrow.getDate() + 1);

  if (date.toDateString() === today.toDateString()) return "Hoje";
  if (date.toDateString() === tomorrow.toDateString()) return "Amanhã";

  return date.toLocaleDateString("pt-BR", {
    weekday: "long",
    day: "numeric",
    month: "long",
  });
}

function getIconForProgram(program: ProgramI) {
  if (program.kind === "checkpoint")
    return { Icon: LogIn, bg: "bg-primary", text: "text-primary-foreground" };
  return { Icon: Activity, bg: "bg-muted", text: "text-muted-foreground" };
}

interface ProgramDayCardProps {
  date: string;
  items: { program: ProgramI; occurrence: OccurrenceI }[];
  maxItems?: number;
}

export function ProgramDayCard({
  date,
  items,
  maxItems = 3,
}: ProgramDayCardProps) {
  const visible = items.slice(0, maxItems);
  const remaining = items.length - maxItems;

  return (
    <div className="flex flex-col w-80 rounded-2xl p-6">
      {/* Day header */}
      <h3 className="w-full rounded-xl bg-primary/10 px-4 py-2 text-center text-sm font-bold tracking-wide text-primary capitalize mb-6">
        {formatDayLabel(date)}
      </h3>

      {/* Timeline */}
      <div className="flex flex-col">
        {visible.map(({ program, occurrence }, index) => {
          const isLast = index === visible.length - 1 && remaining <= 0;
          const { Icon, bg, text } = getIconForProgram(program);

          return (
            <div key={occurrence.id} className="flex gap-4">
              {/* Timeline column */}
              <div className="flex flex-col items-center shrink-0">
                <div
                  className={cn(
                    "w-9 h-9 rounded-xl flex items-center justify-center shrink-0",
                    bg,
                  )}
                >
                  <Icon className={cn("w-4 h-4", text)} />
                </div>
                {!isLast && (
                  <div className="w-px flex-1 min-h-6 bg-border rounded-full" />
                )}
              </div>

              {/* Content */}
              <div className={cn("flex-1", isLast ? "pb-0" : "pb-6")}>
                {/* Time */}
                <span className="text-xs font-semibold text-primary">
                  {formatTimeRange(occurrence.starts_at, occurrence.ends_at)}
                </span>

                {/* Title */}
                <h4 className="mt-1 text-[15px] font-bold text-card-foreground leading-snug">
                  {program.name}
                </h4>

                {/* Description */}
                {program.description && (
                  <p className="mt-1 text-xs text-muted-foreground line-clamp-2 leading-relaxed">
                    {program.description}
                  </p>
                )}
              </div>
            </div>
          );
        })}

        {/* "+ N" */}
        {remaining > 0 && (
          <div className="flex items-center gap-2 -mt-2 ml-13 text-xs font-medium text-muted-foreground">
            <span>
              +{remaining} ocorrência{remaining > 1 ? "s" : ""}
            </span>
            <ChevronRight className="w-3 h-3" />
          </div>
        )}
      </div>
    </div>
  );
}
