import { useMemo } from "react";
import { getMonthGrid, isSameDay, isToday } from "../lib/date";
import type { EventColor, OccurrenceI } from "../model";

interface YearViewProps {
  currentDate: Date;
  today: Date;
  occurrences: OccurrenceI[];
  programColors: Record<string, EventColor>;
  onDateClick: (date: Date) => void;
  onCardClick: (occurrenceId: string) => void;
}

const MONTH_NAMES = [
  "Janeiro",
  "Fevereiro",
  "Março",
  "Abril",
  "Maio",
  "Junho",
  "Julho",
  "Agosto",
  "Setembro",
  "Outubro",
  "Novembro",
  "Dezembro",
];

export function YearView({
  currentDate,
  occurrences,
  programColors,
  onDateClick,
}: YearViewProps) {
  const year = currentDate.getFullYear();

  const months = useMemo(() => {
    return Array.from({ length: 12 }, (_, m) => {
      const firstDay = new Date(year, m, 1);
      return getMonthGrid(firstDay);
    });
  }, [year]);

  const getOccurrencesForDay = (date: Date) => {
    return occurrences.filter((oc) => {
      const s = new Date(oc.starts_at);
      return isSameDay(s, date);
    });
  };

  return (
    <div className="calendar-scroll flex-1 overflow-auto p-6">
      <div className="grid grid-cols-3 lg:grid-cols-4 gap-6">
        {months.map((weeks, monthIdx) => (
          <div
            key={MONTH_NAMES[monthIdx]}
            className="border border-border/60 rounded-xl p-3 bg-card hover:shadow-sm transition-shadow"
          >
            <div className="text-sm font-semibold text-foreground mb-2">
              {MONTH_NAMES[monthIdx]}
            </div>
            <div className="grid grid-cols-7 gap-px">
              {["D", "S", "T", "Q", "Q", "S", "S"].map((d) => (
                <div
                  key={d}
                  className="text-center text-[10px] text-muted-foreground py-0.5"
                >
                  {d}
                </div>
              ))}
              {weeks.flat().map((day) => {
                const dayOccs = getOccurrencesForDay(day);
                const dayIsToday = isToday(day);
                const isOtherMonth = day.getMonth() !== monthIdx;
                const hasEvents = dayOccs.length > 0;

                return (
                  <div
                    key={day.toISOString()}
                    className={`relative text-center py-1 text-[11px] cursor-pointer rounded transition-colors ${
                      isOtherMonth
                        ? "text-muted-foreground/50"
                        : "text-foreground"
                    } ${dayIsToday ? "bg-primary text-primary-foreground font-medium" : "hover:bg-muted"}`}
                    onClick={() => onDateClick(day)}
                  >
                    {day.getDate()}
                    {hasEvents && !dayIsToday && (
                      <div
                        className="absolute bottom-0.5 left-1/2 -translate-x-1/2 w-1 h-1 rounded-full"
                        style={{
                          background:
                            programColors[dayOccs[0].program_id]?.border ||
                            "var(--primary)",
                        }}
                      />
                    )}
                  </div>
                );
              })}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
