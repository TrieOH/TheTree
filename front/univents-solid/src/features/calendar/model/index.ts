import type { ProgramI, OccurrenceI } from "@/features/programs/model";

export type { ProgramI, OccurrenceI };

export type CalendarView = "day" | "week" | "month" | "year";

export type CalendarDragData =
  | { type: "program"; programId: string }
  | { type: "occurrence"; occurrenceId: string };

export interface EventColor {
  bg: string;
  border: string;
  text: string;
}

export const EVENT_COLORS: EventColor[] = [
  {
    bg: "rgba(99, 102, 241, 0.15)",
    border: "rgba(99, 102, 241, 0.8)",
    text: "inherit",
  },
  {
    bg: "rgba(16, 185, 129, 0.15)",
    border: "rgba(16, 185, 129, 0.8)",
    text: "inherit",
  },
  {
    bg: "rgba(245, 158, 11, 0.15)",
    border: "rgba(245, 158, 11, 0.8)",
    text: "inherit",
  },
  {
    bg: "rgba(236, 72, 153, 0.15)",
    border: "rgba(236, 72, 153, 0.8)",
    text: "inherit",
  },
  {
    bg: "rgba(14, 165, 233, 0.15)",
    border: "rgba(14, 165, 233, 0.8)",
    text: "inherit",
  },
];

export const FALLBACK_EVENT_COLOR: EventColor = EVENT_COLORS[0];
