import {
  addDays as dfAddDays,
  isSameDay as dfIsSameDay,
  isToday as dfIsToday,
  startOfDay as dfStartOfDay,
  eachDayOfInterval,
  endOfMonth,
  endOfWeek,
  format,
  startOfMonth,
  startOfWeek,
} from "date-fns";
import { ptBR } from "date-fns/locale";

export {
  dfAddDays as addDays,
  dfStartOfDay as startOfDay,
  dfIsSameDay as isSameDay,
};
export { dfIsToday as isToday };

export function toISODate(date: Date): string {
  // Calendar dates are wall-clock dates. Using toISOString() here converts
  // midnight to UTC and can move the date to the previous day in negative
  // timezones (for example, São Paulo).
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
}

export function pad(n: number): string {
  return String(n).padStart(2, "0");
}

export function formatTime(iso: string): string {
  const d = new Date(iso);
  let h = d.getHours();
  const m = pad(d.getMinutes());
  const ampm = h >= 12 ? "PM" : "AM";
  h = h % 12 || 12;
  return `${h}:${m} ${ampm}`;
}

export function formatTimeRange(start: string, end: string): string {
  return `${formatTime(start)} – ${formatTime(end)}`;
}

export function makeDateStr(date: Date, hour: number, minute: number): string {
  const d = new Date(date);
  d.setHours(hour, minute, 0, 0);
  return d.toISOString();
}

export function getDurationMinutes(start: string, end: string): number {
  return (new Date(end).getTime() - new Date(start).getTime()) / (1000 * 60);
}

export function formatMonthYear(date: Date): string {
  return format(date, "MMM. yyyy", { locale: ptBR });
}

export function formatFullDate(date: Date): string {
  return format(date, "EEEE, d 'de' MMMM", { locale: ptBR });
}

export function monthName(date: Date): string {
  return format(date, "MMMM", { locale: ptBR });
}

export function shortDayName(date: Date): string {
  return format(date, "EEE", { locale: ptBR });
}

export function toISODateTimeLocal(date: Date): string {
  const p = pad;
  return `${date.getFullYear()}-${p(date.getMonth() + 1)}-${p(date.getDate())}T${p(date.getHours())}:${p(date.getMinutes())}`;
}

export function getNowPosition(): number {
  const now = new Date();
  const mins = now.getHours() * 60 + now.getMinutes();
  return (mins / 60) * 60;
}

export function getMonthGrid(date: Date): Date[][] {
  const start = startOfWeek(startOfMonth(date), { weekStartsOn: 0 });
  const end = endOfWeek(endOfMonth(date), { weekStartsOn: 0 });
  const days = eachDayOfInterval({ start, end });
  const weeks: Date[][] = [];
  for (let i = 0; i < days.length; i += 7) {
    weeks.push(days.slice(i, i + 7));
  }
  return weeks;
}

export function getWeekDays(date: Date): Date[] {
  const start = startOfWeek(date, { weekStartsOn: 0 });
  return Array.from({ length: 7 }, (_, i) => dfAddDays(start, i));
}
