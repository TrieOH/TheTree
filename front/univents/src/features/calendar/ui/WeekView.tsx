import { useDroppable } from "@dnd-kit/core";
import { useMemo } from "react";
import {
  addDays,
  getNowPosition,
  isSameDay,
  shortDayName,
  toISODate,
} from "../lib/date";
import {
  type EventColor,
  FALLBACK_EVENT_COLOR,
  type OccurrenceI,
  type ProgramI,
} from "../model";
import { DraggableEvent } from "./DraggableEvent";

const HOURS = Array.from({ length: 24 }, (_, i) => i);

interface WeekViewProps {
  currentDate: Date;
  today: Date;
  occurrences: OccurrenceI[];
  programs: ProgramI[];
  programColors: Record<string, EventColor>;
  onSlotClick: (date: string, hour: number) => void;
  onDayClick: (date: Date) => void;
  onCardClick: (occurrenceId: string) => void;
  onDelete?: (occurrenceId: string) => void;
}

function TimeSlot({
  dateStr,
  hour,
  children,
  onClick,
}: {
  dateStr: string;
  hour: number;
  children: React.ReactNode;
  onClick: (date: string, hour: number) => void;
}) {
  const { isOver, setNodeRef } = useDroppable({
    id: `slot-${dateStr}-${hour}`,
    data: { date: dateStr, hour },
  });

  return (
    <div
      ref={setNodeRef}
      className={`h-15 border-b border-border/60 relative transition-colors ${isOver ? "bg-primary/8" : "hover:bg-muted/40"}`}
      onClick={() => onClick(dateStr, hour)}
    >
      {children}
    </div>
  );
}

export function WeekView({
  currentDate,
  today,
  occurrences,
  programs,
  programColors,
  onSlotClick,
  onDayClick,
  onCardClick,
  onDelete,
}: WeekViewProps) {
  const days = useMemo(() => {
    return Array.from({ length: 7 }, (_, i) => addDays(currentDate, i));
  }, [currentDate]);

  const nowPos = getNowPosition();
  const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;

  const getOccurrencesForDay = (date: Date) => {
    const dayStart = new Date(
      date.getFullYear(),
      date.getMonth(),
      date.getDate(),
    );
    const dayEnd = new Date(dayStart);
    dayEnd.setDate(dayEnd.getDate() + 1);
    return occurrences
      .filter(
        (oc) =>
          new Date(oc.starts_at) < dayEnd && new Date(oc.ends_at) > dayStart,
      )
      .map((oc) => {
        const start = new Date(oc.starts_at);
        return start < dayStart
          ? { ...oc, starts_at: dayStart.toISOString() }
          : oc;
      });
  };

  const formatHourLabel = (h: number): string => {
    if (h === 0) return "12 AM";
    if (h < 12) return `${h} AM`;
    if (h === 12) return "12 PM";
    return `${h - 12} PM`;
  };

  return (
    <div className="calendar-scroll min-w-0 flex-1 overflow-auto">
      <div className="flex min-w-full">
        <div className="w-19 shrink-0 border-r border-border sticky left-0 z-15 bg-background">
          <div className="sticky top-0 z-30 flex h-19 items-center justify-center border-b border-border bg-background px-1 text-center text-[9px] uppercase leading-3 text-muted-foreground break-all">
            {timeZone}
          </div>
          <div className="h-15" />
          {HOURS.map((h) => (
            <div
              key={h}
              className="h-15 flex items-start justify-center pt-0.5 text-[10px] text-muted-foreground relative"
            >
              <span className="-translate-y-1/2 bg-background px-1">
                {formatHourLabel(h)}
              </span>
            </div>
          ))}
        </div>

        {days.map((day) => {
          const dayOccs = getOccurrencesForDay(day);
          const dayIsToday = isSameDay(day, today);
          const dateStr = toISODate(day);

          return (
            <div
              key={dateStr}
              className="min-w-30 flex-1 border-r border-border relative"
            >
              <div
                className="sticky top-0 z-20 flex h-19 flex-col items-center justify-center border-b border-border bg-background text-center"
                onClick={() => onDayClick(day)}
              >
                <div
                  className={`text-[11px] font-medium uppercase tracking-wide ${dayIsToday ? "text-primary" : "text-muted-foreground"}`}
                >
                  {shortDayName(day)}
                </div>
                <div
                  className={`w-10 h-10 mx-auto mt-1 flex items-center justify-center rounded-full text-2xl font-normal ${dayIsToday ? "bg-primary text-primary-foreground font-medium" : "text-foreground"}`}
                >
                  {day.getDate()}
                </div>
              </div>

              <div className="h-7 border-b border-border" />

              {HOURS.map((h) => {
                const slotOccs = dayOccs.filter(
                  (oc) => new Date(oc.starts_at).getHours() === h,
                );
                return (
                  <TimeSlot
                    key={h}
                    dateStr={dateStr}
                    hour={h}
                    onClick={onSlotClick}
                  >
                    {slotOccs.map((oc) => (
                      <DraggableEvent
                        key={oc.id}
                        occurrence={oc}
                        program={programs.find((p) => p.id === oc.program_id)}
                        color={
                          programColors[oc.program_id] ?? FALLBACK_EVENT_COLOR
                        }
                        onClick={onCardClick}
                        onDelete={onDelete}
                        overlapIndex={slotOccs.indexOf(oc)}
                        overlapCount={slotOccs.length}
                      />
                    ))}
                  </TimeSlot>
                );
              })}

              {dayIsToday && (
                <div
                  className="absolute left-0 right-0 h-0.5 z-12 pointer-events-none"
                  style={{
                    top: `${84 + nowPos}px`,
                    background: "var(--destructive)",
                  }}
                >
                  <div
                    className="absolute -left-1.25 -top-1 w-2.5 h-2.5 rounded-full"
                    style={{ background: "var(--destructive)" }}
                  />
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
