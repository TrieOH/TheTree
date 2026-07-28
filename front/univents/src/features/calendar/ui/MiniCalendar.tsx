import { ChevronLeft, ChevronRight } from "lucide-react";
import { useMemo } from "react";
import { getMonthGrid, isSameDay, isToday, monthName } from "../lib/date";

interface MiniCalendarProps {
  currentDate: Date;
  onDateClick: (date: Date) => void;
  onPrevMonth: () => void;
  onNextMonth: () => void;
}

export function MiniCalendar({
  currentDate,
  onDateClick,
  onPrevMonth,
  onNextMonth,
}: MiniCalendarProps) {
  const weeks = useMemo(() => getMonthGrid(currentDate), [currentDate]);
  const dowLabels = [
    { label: "D", key: "domingo" },
    { label: "S", key: "segunda" },
    { label: "T", key: "terca" },
    { label: "Q", key: "quarta" },
    { label: "Q", key: "quinta" },
    { label: "S", key: "sexta" },
    { label: "S", key: "sabado" },
  ];

  return (
    <div className="p-2">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium text-foreground">
          {monthName(currentDate)} de {currentDate.getFullYear()}
        </span>
        <div className="flex gap-0.5">
          <button
            type="button"
            className="w-7 h-7 flex items-center justify-center rounded-full hover:bg-muted transition-colors text-muted-foreground"
            onClick={onPrevMonth}
          >
            <ChevronLeft size={16} />
          </button>
          <button
            type="button"
            className="w-7 h-7 flex items-center justify-center rounded-full hover:bg-muted transition-colors text-muted-foreground"
            onClick={onNextMonth}
          >
            <ChevronRight size={16} />
          </button>
        </div>
      </div>
      <div className="grid grid-cols-7 gap-1.5 text-center">
        {dowLabels.map((d) => (
          <div
            key={d.key}
            className="text-[11px] font-medium text-muted-foreground py-1"
          >
            {d.label}
          </div>
        ))}
        {weeks.flat().map((day) => {
          const dayIsToday = isToday(day);
          const isSelected = isSameDay(day, currentDate);
          const isOtherMonth = day.getMonth() !== currentDate.getMonth();
          let className =
            "text-xs py-1.5 rounded-full cursor-pointer transition-colors font-variant-numeric tabular-nums ";
          if (dayIsToday)
            className += "bg-primary text-primary-foreground font-medium ";
          else if (isSelected)
            className += "bg-primary/15 text-primary font-medium ";
          else if (isOtherMonth) className += "text-muted-foreground/60 ";
          else
            className +=
              "text-muted-foreground hover:bg-muted hover:text-foreground ";

          return (
            <div
              key={day.toISOString()}
              className={className}
              onClick={() => onDateClick(day)}
            >
              {day.getDate()}
            </div>
          );
        })}
      </div>
    </div>
  );
}
