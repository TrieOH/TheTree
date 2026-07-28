import { useMemo } from "react";
import type { EventColor, OccurrenceI, ProgramI } from "../model";
import { formatTime, getMonthGrid, isSameDay, isToday } from "../lib/date";

interface MonthViewProps {
  currentDate: Date;
  today: Date;
  occurrences: OccurrenceI[];
  programs: ProgramI[];
  programColors: Record<string, EventColor>;
  onDateClick: (date: Date) => void;
  onCardClick: (occurrenceId: string) => void;
}

export function MonthView({
  currentDate,
  occurrences,
  programs,
  programColors,
  onDateClick,
  onCardClick,
}: MonthViewProps) {
  const weeks = useMemo(() => getMonthGrid(currentDate), [currentDate]);
  const dowLabels = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"];

  const getOccurrencesForDay = (date: Date) => {
    const dayStart = new Date(date.getFullYear(), date.getMonth(), date.getDate());
    const dayEnd = new Date(dayStart);
    dayEnd.setDate(dayEnd.getDate() + 1);
    return occurrences.filter((oc) => {
      const start = new Date(oc.starts_at);
      const end = new Date(oc.ends_at);
      return start < dayEnd && end > dayStart;
    });
  };

  return (
    <div className="calendar-scroll flex-1 overflow-auto p-4">
      <div className="grid grid-cols-7 border-t border-l border-border/60">
        {dowLabels.map((d, index) => (
          <div
            key={`${d}-${index}`}
            className="text-center py-2 text-xs font-medium text-muted-foreground border-r border-b border-border/60"
          >
            {d}
          </div>
        ))}
        {weeks.flat().map((day) => {
          const dayOccs = getOccurrencesForDay(day);
          const dayIsToday = isToday(day);
          const isOtherMonth = day.getMonth() !== currentDate.getMonth();
          return (
            <div
              key={day.toISOString()}
              className={`min-h-[100px] border-r border-b border-border/60 p-1 cursor-pointer transition-colors hover:bg-muted/30 ${isOtherMonth ? "bg-muted/20" : ""}`}
              onClick={() => onDateClick(day)}
            >
              <div
                className={`w-7 h-7 flex items-center justify-center rounded-full text-sm mx-auto ${dayIsToday ? "bg-primary text-primary-foreground font-medium" : "text-foreground"}`}
              >
                {day.getDate()}
              </div>
              <div className="mt-1 flex flex-col gap-0.5">
                {dayOccs.slice(0, 3).map((oc) => {
                  const prog = programs.find((p) => p.id === oc.program_id);
                  const color = programColors[oc.program_id];
                  return (
                    <div
                      key={oc.id}
                      className="text-[10px] px-1.5 py-0.5 rounded cursor-pointer truncate hover:opacity-80 transition-opacity"
                      style={{
                        background: color.bg,
                        color: color.text,
                        borderLeft: `2px solid ${color.border}`,
                      }}
                      onClick={(e) => {
                        e.stopPropagation();
                        onCardClick(oc.id);
                      }}
                    >
                      {isSameDay(new Date(oc.starts_at), day)
                        ? formatTime(oc.starts_at)
                        : "Continua"} {prog?.name}
                    </div>
                  );
                })}
                {dayOccs.length > 3 && (
                  <div className="text-[10px] text-muted-foreground px-1.5">
                    +{dayOccs.length - 3} mais
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
