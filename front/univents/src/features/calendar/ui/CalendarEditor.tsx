import {
  closestCenter,
  DndContext,
  type DragEndEvent,
  DragOverlay,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import {
  ArrowLeft,
  ChevronLeft,
  ChevronRight,
  Monitor,
  Plus,
} from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import {
  useDeleteOccurrenceMutation,
  useOccurrenceMutation,
} from "@/features/programs/api/mutations";
import { Button } from "@/shared/ui/shadcn/button";
import { addDays, formatMonthYear, startOfDay, toISODate } from "../lib/date";
import type { CalendarView, EventColor, OccurrenceI, ProgramI } from "../model";
import { CalendarCombobox } from "./CalendarCombobox";
import { DayView } from "./DayView";
import { EventModal } from "./EventModal";
import { MiniCalendar } from "./MiniCalendar";
import { MonthView } from "./MonthView";
import { ProgramList } from "./ProgramList";
import { WeekView } from "./WeekView";
import { YearView } from "./YearView";

const COLORS: EventColor[] = [
  {
    bg: "oklch(0.85 0.08 264.05 / 0.25)",
    border: "oklch(0.55 0.15 264.05)",
    text: "var(--foreground)",
  },
  {
    bg: "oklch(0.85 0.08 150 / 0.25)",
    border: "oklch(0.5 0.15 150)",
    text: "var(--foreground)",
  },
  {
    bg: "oklch(0.85 0.08 30 / 0.25)",
    border: "oklch(0.6 0.15 30)",
    text: "var(--foreground)",
  },
  {
    bg: "oklch(0.85 0.08 280 / 0.25)",
    border: "oklch(0.55 0.15 280)",
    text: "var(--foreground)",
  },
  {
    bg: "oklch(0.85 0.08 200 / 0.25)",
    border: "oklch(0.55 0.15 200)",
    text: "var(--foreground)",
  },
];

const VIEW_LABELS: Record<CalendarView, string> = {
  day: "Dia",
  week: "Semana",
  month: "Mês",
  year: "Ano",
};

export function CalendarEditor({
  eventId,
  editionId,
}: {
  eventId: string;
  editionId: string;
}) {
  const navigate = useNavigate();
  const today = useMemo(() => startOfDay(new Date()), []);
  const [currentDate, setCurrentDate] = useState<Date>(new Date(today));
  const [view, setView] = useState<CalendarView>("week");
  const { data: loadedPrograms = [] } = useQuery(
    programsQueryOptions(editionId),
  );
  const { data: loadedOccurrences = [] } = useQuery(
    occurrencesQueryOptions(editionId),
  );
  const occurrenceMutation = useOccurrenceMutation();
  const deleteOccurrenceMutation = useDeleteOccurrenceMutation();
  const [occurrences, setOccurrences] = useState<OccurrenceI[]>([]);
  const programs = loadedPrograms as ProgramI[];
  const visibleOccurrences = loadedOccurrences.length
    ? loadedOccurrences
    : occurrences;
  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<"create" | "edit">("create");
  const [modalOccurrence, setModalOccurrence] = useState<OccurrenceI | null>(
    null,
  );
  const [modalInitial, setModalInitial] = useState<{
    programId?: string;
    date?: string;
    hour?: number;
  }>({});
  const [activeDragId, setActiveDragId] = useState<string | null>(null);

  const programColors = useMemo(() => {
    const map: Record<string, EventColor> = {};
    programs.forEach((p, i) => {
      map[p.id] = COLORS[i % COLORS.length];
    });
    return map;
  }, [programs]);

  const titleText = useMemo(() => {
    if (view === "day") return formatMonthYear(currentDate);
    if (view === "week") {
      const days = Array.from({ length: 7 }, (_, i) => addDays(currentDate, i));
      return `${formatMonthYear(days[0])} – ${formatMonthYear(days[6])}`;
    }
    if (view === "month") return formatMonthYear(currentDate);
    return String(currentDate.getFullYear());
  }, [currentDate, view]);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 5 },
    }),
  );

  const handlePrev = useCallback(() => {
    setCurrentDate((d) => {
      const nd = new Date(d);
      if (view === "day") nd.setDate(nd.getDate() - 1);
      else if (view === "week") nd.setDate(nd.getDate() - 7);
      else if (view === "month") nd.setMonth(nd.getMonth() - 1);
      else if (view === "year") nd.setFullYear(nd.getFullYear() - 1);
      return nd;
    });
  }, [view]);

  const handleNext = useCallback(() => {
    setCurrentDate((d) => {
      const nd = new Date(d);
      if (view === "day") nd.setDate(nd.getDate() + 1);
      else if (view === "week") nd.setDate(nd.getDate() + 7);
      else if (view === "month") nd.setMonth(nd.getMonth() + 1);
      else if (view === "year") nd.setFullYear(nd.getFullYear() + 1);
      return nd;
    });
  }, [view]);

  const handleToday = useCallback(() => {
    setCurrentDate(new Date(today));
  }, [today]);

  const handleMiniPrevMonth = useCallback(() => {
    setCurrentDate((d) => {
      const nd = new Date(d);
      nd.setMonth(nd.getMonth() - 1);
      return nd;
    });
  }, []);

  const handleMiniNextMonth = useCallback(() => {
    setCurrentDate((d) => {
      const nd = new Date(d);
      nd.setMonth(nd.getMonth() + 1);
      return nd;
    });
  }, []);

  const handleMiniDateClick = useCallback((date: Date) => {
    setCurrentDate(date);
  }, []);

  const handleYearDateClick = useCallback((date: Date) => {
    setCurrentDate(date);
    setView("month");
  }, []);

  const handleMonthDateClick = useCallback((date: Date) => {
    setCurrentDate(date);
    setView("week");
  }, []);

  const handleWeekDateClick = useCallback((date: Date) => {
    setCurrentDate(date);
    setView("day");
  }, []);

  const handleDragEnd = useCallback(
    (event: DragEndEvent) => {
      const { active, over } = event;
      setActiveDragId(null);

      if (!over) return;

      const activeData = active.data.current;
      const overData = over.data.current;

      if (!activeData || !overData) return;

      if (activeData.type === "program" && overData.date !== undefined) {
        const [year, month, day] = String(overData.date).split("-").map(Number);
        const start = new Date(
          year,
          month - 1,
          day,
          Number(overData.hour),
          0,
          0,
          0,
        );
        const end = new Date(start.getTime() + 60 * 60 * 1000);
        occurrenceMutation.mutate({
          programId: activeData.programId as string,
          data: { starts_at: start.toISOString(), ends_at: end.toISOString() },
        });
        setOccurrences((prev) => [
          ...prev,
          {
            id: `oc_${Math.random().toString(36).slice(2, 9)}`,
            program_id: activeData.programId as string,
            edition_id: "e1",
            starts_at: start.toISOString(),
            ends_at: end.toISOString(),
            max_capacity: null,
            created_at: new Date().toISOString(),
            updated_at: null,
            deleted_at: null,
          },
        ]);
      } else if (
        activeData.type === "occurrence" &&
        overData.date !== undefined
      ) {
        const occurrenceId = activeData.occurrenceId as string;
        const occurrence = visibleOccurrences.find(
          (oc) => oc.id === occurrenceId,
        );
        if (!occurrence) return;
        const durMs =
          new Date(occurrence.ends_at).getTime() -
          new Date(occurrence.starts_at).getTime();
        const [year, month, day] = String(overData.date).split("-").map(Number);
        const hour = Math.max(0, Math.min(23, Number(overData.hour)));
        const currentStart = new Date(occurrence.starts_at);
        if (
          toISODate(currentStart) === String(overData.date) &&
          currentStart.getHours() === hour
        ) {
          return;
        }
        const newStart = new Date(year, month - 1, day, hour, 0, 0, 0);
        const newEnd = new Date(newStart.getTime() + durMs);
        occurrenceMutation.mutate({
          id: occurrenceId,
          data: {
            starts_at: newStart.toISOString(),
            ends_at: newEnd.toISOString(),
          },
        });
        setOccurrences((prev) =>
          prev.map((oc) => {
            if (oc.id !== occurrenceId) return oc;
            return {
              ...oc,
              starts_at: newStart.toISOString(),
              ends_at: newEnd.toISOString(),
              updated_at: new Date().toISOString(),
            };
          }),
        );
      }
    },
    [occurrenceMutation],
  );

  const handleSlotClick = useCallback(
    (date: string, hour: number) => {
      setModalMode("create");
      setModalInitial({ programId: programs[0]?.id, date, hour });
      setModalOccurrence(null);
      setModalOpen(true);
    },
    [programs],
  );

  const handleCardClick = useCallback(
    (occurrenceId: string) => {
      const oc = visibleOccurrences.find((o) => o.id === occurrenceId);
      if (!oc) return;
      setModalMode("edit");
      setModalOccurrence(oc);
      setModalInitial({});
      setModalOpen(true);
    },
    [visibleOccurrences],
  );

  const handleModalSave = useCallback(
    (data: {
      program_id: string;
      starts_at: string;
      ends_at: string;
      max_capacity?: number;
      id?: string;
    }) => {
      if (modalMode === "edit" && data.id) {
        occurrenceMutation.mutate({ id: data.id, data });
        setOccurrences((prev) =>
          prev.map((oc) =>
            oc.id === data.id
              ? {
                  ...oc,
                  program_id: data.program_id,
                  starts_at: data.starts_at,
                  ends_at: data.ends_at,
                  max_capacity: data.max_capacity ?? null,
                  updated_at: new Date().toISOString(),
                }
              : oc,
          ),
        );
      } else {
        occurrenceMutation.mutate({ programId: data.program_id, data });
        setOccurrences((prev) => [
          ...prev,
          {
            id: `oc_${Math.random().toString(36).slice(2, 9)}`,
            program_id: data.program_id,
            edition_id: "e1",
            starts_at: data.starts_at,
            ends_at: data.ends_at,
            max_capacity: data.max_capacity ?? null,
            created_at: new Date().toISOString(),
            updated_at: null,
            deleted_at: null,
          },
        ]);
      }
      setModalOpen(false);
    },
    [modalMode, occurrenceMutation],
  );

  const handleModalDelete = useCallback(
    (occurrenceId: string) => {
      deleteOccurrenceMutation.mutate(occurrenceId);
      setOccurrences((prev) =>
        prev.map((oc) =>
          oc.id === occurrenceId
            ? { ...oc, deleted_at: new Date().toISOString() }
            : oc,
        ),
      );
      setModalOpen(false);
    },
    [deleteOccurrenceMutation],
  );

  const activeDragProgram = useMemo(() => {
    if (!activeDragId?.startsWith("program-")) return null;
    const pid = activeDragId.replace("program-", "");
    return programs.find((p) => p.id === pid);
  }, [activeDragId, programs]);

  const activeDragOccurrence = useMemo(() => {
    if (!activeDragId?.startsWith("event-")) return null;
    const oid = activeDragId.replace("event-", "");
    return visibleOccurrences.find((o) => o.id === oid);
  }, [activeDragId, visibleOccurrences]);

  const activeDragOccurrenceProgram = useMemo(() => {
    if (!activeDragOccurrence) return null;
    return programs.find((p) => p.id === activeDragOccurrence.program_id);
  }, [activeDragOccurrence, programs]);

  const activeDragOccurrenceColor = useMemo(() => {
    if (!activeDragOccurrence) return null;
    return programColors[activeDragOccurrence.program_id];
  }, [activeDragOccurrence, programColors]);

  const activeOccurrences = useMemo(
    () => visibleOccurrences.filter((oc) => !oc.deleted_at),
    [visibleOccurrences],
  );

  return (
    <div className="h-dvh min-h-0 bg-background text-foreground">
      <div className="flex h-full flex-col items-center justify-center gap-3 px-6 text-center lg:hidden!">
        <Monitor className="size-10 text-muted-foreground" />
        <p className="max-w-xs text-sm text-muted-foreground">
          O calendário foi desenvolvido para telas maiores. Abra esta página em
          um computador para acessar a agenda completa.
        </p>
      </div>

      <div className="hidden h-full min-h-0 min-w-0 flex-col overflow-hidden bg-background lg:flex">
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragStart={(e) => setActiveDragId(e.active.id as string)}
          onDragEnd={handleDragEnd}
        >
          {/* Toolbar */}
          <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border h-14 shrink-0 bg-background">
            <Button
              type="button"
              variant="ghost"
              className="h-9 gap-2 px-3"
              onClick={() =>
                void navigate({
                  to: "/admin/events/$eventId/editions/$editionId/programs",
                  params: { eventId, editionId },
                })
              }
            >
              <ArrowLeft className="size-4" />
              Voltar
            </Button>

            <div className="flex items-center gap-0.5 ml-1">
              <button
                type="button"
                className="w-9 h-9 flex items-center justify-center rounded-full hover:bg-muted transition-colors text-foreground"
                onClick={handlePrev}
              >
                <ChevronLeft size={18} />
              </button>
              <button
                type="button"
                className="w-9 h-9 flex items-center justify-center rounded-full hover:bg-muted transition-colors text-foreground"
                onClick={handleNext}
              >
                <ChevronRight size={18} />
              </button>
            </div>

            <button
              type="button"
              className="px-4 py-2 border border-border rounded-full text-sm font-medium text-foreground hover:bg-muted transition-colors"
              onClick={handleToday}
            >
              Hoje
            </button>

            <div className="text-[22px] font-normal text-foreground ml-4 whitespace-nowrap tracking-tight">
              {titleText}
            </div>

            <div className="flex-1" />

            <CalendarCombobox
              value={view}
              onChange={(value) => setView(value as CalendarView)}
              options={(Object.keys(VIEW_LABELS) as CalendarView[]).map(
                (v) => ({ value: v, label: VIEW_LABELS[v] }),
              )}
              placeholder="Visualização"
              compact
            />
          </div>

          {/* Body */}
          <div className="flex min-w-0 flex-1 min-h-0 overflow-hidden">
            {/* Sidebar */}
            <div className="w-65 shrink-0 border-r border-border flex flex-col overflow-hidden p-3 gap-4 bg-background">
              <MiniCalendar
                currentDate={currentDate}
                onDateClick={handleMiniDateClick}
                onPrevMonth={handleMiniPrevMonth}
                onNextMonth={handleMiniNextMonth}
              />
              <button
                type="button"
                className="inline-flex items-center justify-center gap-1.5 px-4 py-2 rounded-full text-sm font-medium bg-primary text-primary-foreground hover:opacity-90 transition-opacity shadow-sm"
                onClick={() => {
                  setModalMode("create");
                  setModalInitial({
                    programId: programs[0]?.id,
                    date: toISODate(currentDate),
                    hour: 9,
                  });
                  setModalOccurrence(null);
                  setModalOpen(true);
                }}
              >
                <Plus size={16} /> Criar ocorrência
              </button>
              <div className="calendar-scroll flex-1 overflow-y-auto">
                <ProgramList
                  programs={programs}
                  programColors={programColors}
                />
              </div>
            </div>

            {/* Views */}
            {view === "day" && (
              <DayView
                currentDate={currentDate}
                today={today}
                occurrences={activeOccurrences}
                programs={programs}
                programColors={programColors}
                onSlotClick={handleSlotClick}
                onCardClick={handleCardClick}
                onDelete={handleModalDelete}
              />
            )}
            {view === "week" && (
              <WeekView
                currentDate={currentDate}
                today={today}
                occurrences={activeOccurrences}
                programs={programs}
                programColors={programColors}
                onSlotClick={handleSlotClick}
                onDayClick={handleWeekDateClick}
                onCardClick={handleCardClick}
                onDelete={handleModalDelete}
              />
            )}
            {view === "month" && (
              <MonthView
                currentDate={currentDate}
                today={today}
                occurrences={activeOccurrences}
                programs={programs}
                programColors={programColors}
                onDateClick={handleMonthDateClick}
                onCardClick={handleCardClick}
              />
            )}
            {view === "year" && (
              <YearView
                currentDate={currentDate}
                today={today}
                occurrences={activeOccurrences}
                programColors={programColors}
                onDateClick={handleYearDateClick}
                onCardClick={handleCardClick}
              />
            )}
          </div>

          {/* Drag Overlay */}
          <DragOverlay>
            {activeDragProgram && (
              <div className="flex items-center gap-2 px-3 py-2 rounded-lg bg-card border border-border shadow-xl text-sm text-foreground opacity-90">
                <span
                  className="w-2.5 h-2.5 rounded-full"
                  style={{
                    background: programColors[activeDragProgram.id]?.border,
                  }}
                />
                <span>{activeDragProgram.name}</span>
              </div>
            )}
            {activeDragOccurrence &&
              activeDragOccurrenceProgram &&
              activeDragOccurrenceColor && (
                <div
                  className="rounded-md px-3 py-2 shadow-xl opacity-90 min-w-30"
                  style={{
                    background: activeDragOccurrenceColor.bg,
                    borderLeft: `3px solid ${activeDragOccurrenceColor.border}`,
                    color: activeDragOccurrenceColor.text,
                  }}
                >
                  <div className="text-xs font-semibold">
                    {activeDragOccurrenceProgram.name}
                  </div>
                </div>
              )}
          </DragOverlay>
        </DndContext>

        {/* Modal */}
        {modalOpen && (
          <EventModal
            mode={modalMode}
            programs={programs}
            occurrence={modalOccurrence}
            initialProgramId={modalInitial.programId}
            initialDate={modalInitial.date}
            initialHour={modalInitial.hour}
            onSave={handleModalSave}
            onDelete={handleModalDelete}
            onClose={() => setModalOpen(false)}
          />
        )}
      </div>
    </div>
  );
}
