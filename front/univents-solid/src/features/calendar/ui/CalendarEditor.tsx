import type { JSX } from "@solidjs/web";
import { useQuery, useQueryClient } from "@trieoh/front-core-solid";
import { useNavigate } from "@tanstack/solid-router";
import { For, Show, createMemo, createSignal } from "solid-js";

import ArrowLeftIcon from "~icons/lucide/arrow-left";
import ChevronLeftIcon from "~icons/lucide/chevron-left";
import ChevronRightIcon from "~icons/lucide/chevron-right";
import GripVerticalIcon from "~icons/lucide/grip-vertical";
import PlusIcon from "~icons/lucide/plus";

import { occurrencesQueryOptions, programsQueryOptions } from "@/features/programs/api";
import {
  useCreateOccurrenceMutation,
  useDeleteOccurrenceMutation,
  useUpdateOccurrenceMutation,
} from "@/features/programs/api/mutations";
import type { OccurrenceCreateOutput } from "@/features/programs/model";
import { programKeys } from "@/features/programs/api/query-keys";
import { DesktopOnly } from "@/shared/ui/DesktopOnly";
import { toast } from "@/shared/ui/toast";
import { addDays, formatMonthYear, isSameDay, startOfDay } from "../lib/date";
import { endCalendarDrag, startCalendarDrag } from "../lib/drag-state";
import type { CalendarDragData, CalendarView, EventColor, OccurrenceI, ProgramI } from "../model";
import { EVENT_COLORS, FALLBACK_EVENT_COLOR } from "../model";
import { DayView } from "./DayView";
import { ManageOccurrenceDialog } from "./ManageOccurrenceDialog";
import { MiniCalendar } from "./MiniCalendar";
import { MonthView } from "./MonthView";
import { WeekView } from "./WeekView";
import { YearView } from "./YearView";

const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const ChevronLeft = ChevronLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const ChevronRight = ChevronRightIcon as unknown as (props: { class?: string }) => JSX.Element;
const GripVertical = GripVerticalIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface CalendarEditorProps {
  eventId: string;
  editionId: string;
  initialProgramId?: string;
}

export function CalendarEditor(props: CalendarEditorProps): JSX.Element {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [currentDate, setCurrentDate] = createSignal<Date>(startOfDay(new Date()));
  const [view, setView] = createSignal<CalendarView>("week");

  // Queries
  const programsQuery = useQuery(() => programsQueryOptions(props.editionId));
  const occurrencesQuery = useQuery(() => occurrencesQueryOptions(props.editionId));

  const programs = createMemo(() => (programsQuery().data ?? []) as ProgramI[]);
  const occurrences = createMemo(() => (occurrencesQuery().data ?? []) as OccurrenceI[]);

  // Mutations
  const createMutation = useCreateOccurrenceMutation(() => props.editionId);
  const updateMutation = useUpdateOccurrenceMutation(() => props.editionId);
  const deleteMutation = useDeleteOccurrenceMutation(() => props.editionId);

  // Dialog State
  const [dialogOpen, setDialogOpen] = createSignal(false);
  const [selectedOccurrence, setSelectedOccurrence] = createSignal<OccurrenceI | null>(null);
  const [newOccurrenceSlot, setNewOccurrenceSlot] = createSignal<{
    dateStr: string;
    hour: number;
  } | null>(null);

  // Filter by program
  const [selectedProgramFilter, setSelectedProgramFilter] = createSignal<string | null>(
    props.initialProgramId ?? null,
  );

  // Program color map
  const programColors = createMemo(() => {
    const map: Record<string, EventColor> = {};
    programs().forEach((p, i) => {
      map[p.id] = EVENT_COLORS[i % EVENT_COLORS.length] ?? FALLBACK_EVENT_COLOR;
    });
    return map;
  });

  // Filtered occurrences
  const visibleOccurrences = createMemo(() => {
    if (!selectedProgramFilter()) return occurrences();
    return occurrences().filter((oc) => oc.program_id === selectedProgramFilter());
  });

  // Navigation handlers
  const handlePrev = () => {
    setCurrentDate((d) => {
      switch (view()) {
        case "day":
          return addDays(d, -1);
        case "week":
          return addDays(d, -7);
        case "month": {
          const nd = new Date(d);
          nd.setMonth(nd.getMonth() - 1);
          return nd;
        }
        case "year": {
          const nd = new Date(d);
          nd.setFullYear(nd.getFullYear() - 1);
          return nd;
        }
      }
    });
  };

  const handleNext = () => {
    setCurrentDate((d) => {
      switch (view()) {
        case "day":
          return addDays(d, 1);
        case "week":
          return addDays(d, 7);
        case "month": {
          const nd = new Date(d);
          nd.setMonth(nd.getMonth() + 1);
          return nd;
        }
        case "year": {
          const nd = new Date(d);
          nd.setFullYear(nd.getFullYear() + 1);
          return nd;
        }
      }
    });
  };

  const handleToday = () => {
    setCurrentDate(startOfDay(new Date()));
  };

  // Title text based on current view
  const titleText = createMemo(() => {
    const d = currentDate();
    switch (view()) {
      case "day":
        return d.toLocaleDateString("pt-BR", {
          weekday: "long",
          day: "numeric",
          month: "long",
          year: "numeric",
        });
      case "week": {
        const start = new Date(d);
        start.setDate(d.getDate() - d.getDay());
        const end = new Date(start);
        end.setDate(start.getDate() + 6);
        return `${start.toLocaleDateString("pt-BR", {
          day: "numeric",
          month: "short",
        })} - ${end.toLocaleDateString("pt-BR", {
          day: "numeric",
          month: "short",
          year: "numeric",
        })}`;
      }
      case "month":
        return formatMonthYear(d);
      case "year":
        return `${d.getFullYear()}`;
    }
  });

  // Dialog actions
  const openNewOccurrence = (dateStr?: string, hour?: number) => {
    setSelectedOccurrence(null);
    if (dateStr && hour !== undefined) {
      setNewOccurrenceSlot({ dateStr, hour });
    } else {
      setNewOccurrenceSlot(null);
    }
    setDialogOpen(true);
  };

  const openEditOccurrence = (occurrence: OccurrenceI) => {
    setNewOccurrenceSlot(null);
    setSelectedOccurrence(occurrence);
    setDialogOpen(true);
  };

  const handleSaveOccurrence = async (data: {
    programId: string;
    occurrenceData: OccurrenceCreateOutput;
    id?: string;
  }): Promise<boolean> => {
    try {
      if (data.id) {
        await updateMutation.mutateAsync({
          id: data.id,
          data: data.occurrenceData,
        });
        toast.success("Horário atualizado com sucesso!");
      } else {
        await createMutation.mutateAsync({
          programId: data.programId,
          data: data.occurrenceData,
        });
        toast.success("Horário criado com sucesso!");
      }
      void queryClient.invalidateQueries({
        queryKey: programKeys.occurrences(props.editionId),
      });
      void queryClient.invalidateQueries({
        queryKey: programKeys.byEdition(props.editionId),
      });
      return true;
    } catch {
      toast.error("Erro ao salvar horário.");
      return false;
    }
  };

  const handleDeleteOccurrenceModal = async (occurrenceId: string): Promise<boolean> => {
    try {
      await deleteMutation.mutateAsync(occurrenceId);
      void queryClient.invalidateQueries({
        queryKey: programKeys.occurrences(props.editionId),
      });
      void queryClient.invalidateQueries({
        queryKey: programKeys.byEdition(props.editionId),
      });
      toast.success("Horário excluído!");
      return true;
    } catch {
      toast.error("Erro ao excluir horário.");
      return false;
    }
  };

  const handleDeleteOccurrence = async (occurrenceId: string) => {
    await handleDeleteOccurrenceModal(occurrenceId);
  };

  // Drag & drop handlers
  const handleDropSlot = async (
    dateStr: string,
    hour: number,
    minute: number,
    data: CalendarDragData,
  ) => {
    try {
      if (data.type === "program") {
        const [year, month, day] = dateStr.split("-").map(Number);
        const start = new Date(year, (month || 1) - 1, day || 1, hour, minute, 0, 0);
        const end = new Date(start.getTime() + 60 * 60 * 1000); // 1 hora padrão

        await createMutation.mutateAsync({
          programId: data.programId,
          data: {
            starts_at: start.toISOString(),
            ends_at: end.toISOString(),
          },
        });
        void queryClient.invalidateQueries({
          queryKey: programKeys.occurrences(props.editionId),
        });
        void queryClient.invalidateQueries({
          queryKey: programKeys.byEdition(props.editionId),
        });
        toast.success("Horário adicionado com sucesso!");
      } else if (data.type === "occurrence") {
        const oc = occurrences().find((o) => o.id === data.occurrenceId);
        if (!oc) return;

        const currentStart = new Date(oc.starts_at);
        const durMs = Math.max(
          new Date(oc.ends_at).getTime() - currentStart.getTime(),
          15 * 60 * 1000,
        );
        const [year, month, day] = dateStr.split("-").map(Number);
        const newStart = new Date(year, (month || 1) - 1, day || 1, hour, minute, 0, 0);
        const newEnd = new Date(newStart.getTime() + durMs);

        // Don't update if same start time
        if (currentStart.getTime() === newStart.getTime()) {
          return;
        }

        await updateMutation.mutateAsync({
          id: oc.id,
          data: {
            starts_at: newStart.toISOString(),
            ends_at: newEnd.toISOString(),
            max_capacity: oc.max_capacity ?? undefined,
          },
        });
        void queryClient.invalidateQueries({
          queryKey: programKeys.occurrences(props.editionId),
        });
        void queryClient.invalidateQueries({
          queryKey: programKeys.byEdition(props.editionId),
        });
        toast.success("Horário reagendado!");
      }
    } catch {
      toast.error("Erro ao salvar horário.");
    }
  };

  const handleDropDay = async (day: Date, data: CalendarDragData) => {
    try {
      if (data.type === "program") {
        const start = new Date(day);
        start.setHours(9, 0, 0, 0);
        const end = new Date(start.getTime() + 60 * 60 * 1000);

        await createMutation.mutateAsync({
          programId: data.programId,
          data: {
            starts_at: start.toISOString(),
            ends_at: end.toISOString(),
          },
        });
        void queryClient.invalidateQueries({
          queryKey: programKeys.occurrences(props.editionId),
        });
        void queryClient.invalidateQueries({
          queryKey: programKeys.byEdition(props.editionId),
        });
        toast.success("Horário adicionado com sucesso!");
      } else if (data.type === "occurrence") {
        const oc = occurrences().find((o) => o.id === data.occurrenceId);
        if (!oc) return;

        const origStart = new Date(oc.starts_at);
        const newStart = new Date(day);
        newStart.setHours(origStart.getHours(), origStart.getMinutes(), 0, 0);

        if (isSameDay(origStart, newStart)) return;

        const durMs = Math.max(
          new Date(oc.ends_at).getTime() - origStart.getTime(),
          30 * 60 * 1000,
        );
        const newEnd = new Date(newStart.getTime() + durMs);

        await updateMutation.mutateAsync({
          id: oc.id,
          data: {
            starts_at: newStart.toISOString(),
            ends_at: newEnd.toISOString(),
            max_capacity: oc.max_capacity ?? undefined,
          },
        });
        void queryClient.invalidateQueries({
          queryKey: programKeys.occurrences(props.editionId),
        });
        void queryClient.invalidateQueries({
          queryKey: programKeys.byEdition(props.editionId),
        });
        toast.success("Horário reagendado!");
      }
    } catch {
      toast.error("Erro ao reagendar horário.");
    }
  };

  return (
    <DesktopOnly
      threshold={1200}
      title="Disponível apenas no computador"
      description="O editor de horários da programação foi projetado para telas a partir de 1200px de largura para permitir o uso da grade e organização dos eventos."
      actionLabel="Voltar aos programas"
      onAction={() =>
        navigate({
          to: "/admin/events/$eventId/editions/$editionId/programs",
          params: { eventId: props.eventId, editionId: props.editionId },
        })
      }
    >
      <div
        class="flex flex-col h-dvh w-full overflow-hidden bg-background"
        onDragOver={(e) => e.preventDefault()}
        onDrop={(e) => e.preventDefault()}
      >
        {/* Top Bar / Toolbar - Clean & Organized */}
        <header class="flex h-14 shrink-0 items-center justify-between border-b border-border bg-card px-4 sm:px-6 select-none z-30">
          {/* Left: Navigation & Date Info */}
          <div class="flex items-center gap-3">
            {/* Back Button */}
            <button
              type="button"
              onClick={() =>
                navigate({
                  to: "/admin/events/$eventId/editions/$editionId/programs",
                  params: { eventId: props.eventId, editionId: props.editionId },
                })
              }
              class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-md hover:bg-muted text-xs font-medium text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
            >
              <ArrowLeft class="size-4" />
              <span>Voltar aos programas</span>
            </button>

            <div class="h-4 w-px bg-border" />

            {/* Date Navigation */}
            <div class="flex items-center gap-1">
              <button
                type="button"
                onClick={handleToday}
                class="px-2.5 py-1 text-xs font-medium rounded-md border border-border hover:bg-muted text-foreground transition-colors cursor-pointer"
              >
                Hoje
              </button>
              <div class="flex items-center">
                <button
                  type="button"
                  onClick={handlePrev}
                  title="Anterior"
                  class="size-7 flex items-center justify-center rounded-md hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                >
                  <ChevronLeft class="size-4" />
                </button>
                <button
                  type="button"
                  onClick={handleNext}
                  title="Próximo"
                  class="size-7 flex items-center justify-center rounded-md hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                >
                  <ChevronRight class="size-4" />
                </button>
              </div>
            </div>

            {/* Current View Title */}
            <h1 class="text-sm sm:text-base font-semibold text-foreground tracking-tight capitalize pl-1">
              {titleText()}
            </h1>
          </div>

          {/* Right: View Switcher */}
          <div class="flex items-center">
            <div class="flex items-center rounded-lg border border-border bg-muted/30 p-0.5">
              {(["day", "week", "month", "year"] as const).map((v) => {
                const labels = { day: "Dia", week: "Semana", month: "Mês", year: "Ano" };
                const isActive = () => view() === v;
                return (
                  <button
                    type="button"
                    onClick={() => setView(v)}
                    class={`px-3 py-1 text-xs font-medium rounded-md transition-colors cursor-pointer ${
                      isActive()
                        ? "bg-background text-foreground font-semibold shadow-2xs"
                        : "text-muted-foreground hover:text-foreground"
                    }`}
                  >
                    {labels[v]}
                  </button>
                );
              })}
            </div>
          </div>
        </header>

        {/* Main layout */}
        <div class="flex flex-1 overflow-hidden">
          {/* Sidebar with mini calendar & program list (desktop only) */}
          <div class="hidden lg:flex w-64 flex-col border-r border-border bg-background p-3 gap-4 overflow-y-auto shrink-0 select-none">
            <MiniCalendar
              currentDate={currentDate()}
              onDateClick={(d) => setCurrentDate(d)}
              onPrevMonth={() =>
                setCurrentDate((d) => {
                  const nd = new Date(d);
                  nd.setMonth(nd.getMonth() - 1);
                  return nd;
                })
              }
              onNextMonth={() =>
                setCurrentDate((d) => {
                  const nd = new Date(d);
                  nd.setMonth(nd.getMonth() + 1);
                  return nd;
                })
              }
            />

            <button
              type="button"
              class="inline-flex items-center justify-center gap-1.5 px-4 py-2 rounded-full text-sm font-medium bg-primary text-primary-foreground hover:opacity-90 transition-opacity shadow-xs cursor-pointer"
              onClick={() => openNewOccurrence()}
            >
              <Plus class="size-4" /> Criar ocorrência
            </button>

            <div class="space-y-2 pt-2 border-t border-border">
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold text-foreground uppercase tracking-wider">
                  Programas
                </span>
                <Show when={selectedProgramFilter() !== null}>
                  <button
                    type="button"
                    class="text-[10px] text-primary hover:underline cursor-pointer"
                    onClick={() => setSelectedProgramFilter(null)}
                  >
                    Limpar filtro
                  </button>
                </Show>
              </div>

              <p class="text-[11px] text-muted-foreground leading-snug">
                Arraste um programa para a grade para criar um horário.
              </p>

              <div class="space-y-1 mt-2">
                <For
                  each={programs()}
                  fallback={
                    <p class="text-xs text-muted-foreground/60 py-2">
                      Nenhum programa cadastrado nesta edição.
                    </p>
                  }
                >
                  {(prog) => {
                    const color = () => programColors()[prog.id] ?? FALLBACK_EVENT_COLOR;
                    const isFiltered = () => selectedProgramFilter() === prog.id;

                    return (
                      <div
                        draggable="true"
                        onDragStart={(e) => {
                          startCalendarDrag(e, {
                            type: "program",
                            programId: prog.id,
                          });
                        }}
                        onDragEnd={() => {
                          endCalendarDrag();
                        }}
                        onClick={() => {
                          setSelectedProgramFilter((prev) =>
                            prev === prog.id ? null : prog.id,
                          );
                        }}
                        class={`flex items-center gap-2 px-2.5 py-1.5 rounded-lg border text-xs font-medium cursor-grab active:cursor-grabbing transition-all select-none ${
                          isFiltered()
                            ? "ring-2 ring-primary border-primary font-semibold"
                            : "border-border/60 hover:bg-muted/50"
                        }`}
                        style={{
                          "border-left": `4px solid ${color().border}`,
                          background: isFiltered() ? color().bg : undefined,
                        }}
                      >
                        <GripVertical class="size-3 text-muted-foreground/50 shrink-0" />
                        <span class="truncate flex-1">{prog.name}</span>
                      </div>
                    );
                  }}
                </For>
              </div>
            </div>
          </div>

          {/* View container */}
          <div class="flex flex-1 flex-col overflow-hidden">
            <Show when={view() === "month"}>
              <MonthView
                currentDate={currentDate()}
                occurrences={visibleOccurrences()}
                programs={programs()}
                programColors={programColors()}
                onDateClick={(d) => {
                  setCurrentDate(d);
                  setView("day");
                }}
                onOccurrenceClick={openEditOccurrence}
                onDropDay={handleDropDay}
              />
            </Show>

            <Show when={view() === "week"}>
              <WeekView
                currentDate={currentDate()}
                occurrences={visibleOccurrences()}
                programs={programs()}
                programColors={programColors()}
                onSlotClick={(dateStr, hour) => openNewOccurrence(dateStr, hour)}
                onDayClick={(d) => {
                  setCurrentDate(d);
                  setView("day");
                }}
                onOccurrenceClick={openEditOccurrence}
                onDeleteOccurrence={handleDeleteOccurrence}
                onDropSlot={handleDropSlot}
              />
            </Show>

            <Show when={view() === "day"}>
              <DayView
                currentDate={currentDate()}
                occurrences={visibleOccurrences()}
                programs={programs()}
                programColors={programColors()}
                onSlotClick={(dateStr, hour) => openNewOccurrence(dateStr, hour)}
                onOccurrenceClick={openEditOccurrence}
                onDeleteOccurrence={handleDeleteOccurrence}
                onDropSlot={handleDropSlot}
              />
            </Show>

            <Show when={view() === "year"}>
              <YearView
                currentDate={currentDate()}
                occurrences={visibleOccurrences()}
                programColors={programColors()}
                onDateClick={(d) => {
                  setCurrentDate(d);
                  setView("day");
                }}
              />
            </Show>
          </div>
        </div>

        {/* Create/Edit Occurrence Dialog */}
        <ManageOccurrenceDialog
          open={dialogOpen()}
          onOpenChange={setDialogOpen}
          programs={programs()}
          initialProgramId={selectedProgramFilter() ?? undefined}
          occurrence={selectedOccurrence()}
          initialDate={newOccurrenceSlot()?.dateStr}
          initialHour={newOccurrenceSlot()?.hour}
          onSave={handleSaveOccurrence}
          onDelete={handleDeleteOccurrenceModal}
        />
      </div>
    </DesktopOnly>
  );
}
