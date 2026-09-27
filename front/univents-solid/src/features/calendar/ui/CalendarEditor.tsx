import type { JSX } from "@solidjs/web";
import { useQuery, useQueryClient } from "@trieoh/front-core-solid";
import { useNavigate } from "@tanstack/solid-router";
import { For, Show, createMemo, createSignal, onCleanup } from "solid-js";

import ArrowLeftIcon from "~icons/lucide/arrow-left";
import ChevronLeftIcon from "~icons/lucide/chevron-left";
import ChevronRightIcon from "~icons/lucide/chevron-right";
import FilterIcon from "~icons/lucide/filter";
import GripVerticalIcon from "~icons/lucide/grip-vertical";
import PlusIcon from "~icons/lucide/plus";
import SearchIcon from "~icons/lucide/search";
import { AlertModal } from "@/widgets/ui/AlertModal";
import { Combobox, type ComboboxOption } from "@/shared/ui/Combobox";
import { toISODate } from "../lib/date";

import { occurrencesQueryOptions, programsQueryOptions } from "@/features/programs/api";
import {
  useCreateOccurrenceMutation,
  useDeleteOccurrenceMutation,
  useUpdateOccurrenceMutation,
} from "@/features/programs/api/mutations";
import type { OccurrenceCreateOutput } from "@/features/programs/model";
import { programKeys } from "@/features/programs/api/query-keys";
import { OccurrenceAttendanceModal } from "@/features/programs/ui/OccurrenceAttendanceModal";
import { OccurrenceDrawModal } from "@/features/programs/ui/OccurrenceDrawModal";
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
const Filter = FilterIcon as unknown as (props: { class?: string }) => JSX.Element;
const GripVertical = GripVerticalIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;
const Search = SearchIcon as unknown as (props: { class?: string }) => JSX.Element;

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

  // Attendance & Draw modals
  const [attendanceModalOpen, setAttendanceModalOpen] = createSignal(false);
  const [attendanceOccurrence, setAttendanceOccurrence] = createSignal<OccurrenceI | null>(null);
  const [drawModalOpen, setDrawModalOpen] = createSignal(false);
  const [drawOccurrence, setDrawOccurrence] = createSignal<OccurrenceI | null>(null);

  // Filter by program & advanced filters
  const [selectedProgramFilter, setSelectedProgramFilter] = createSignal<string | null>(
    props.initialProgramId ?? null,
  );
  const [programSearch, setProgramSearch] = createSignal("");
  const [programKindFilter, setProgramKindFilter] = createSignal<"all" | "activity" | "checkpoint">("all");
  const [dateFilter, setDateFilter] = createSignal<string | null>(null);

  // Occurrence deletion confirmation state
  const [occurrenceToDelete, setOccurrenceToDelete] = createSignal<OccurrenceI | null>(null);
  const [deletingOccurrence, setDeletingOccurrence] = createSignal(false);

  // Filter Popover state
  const [filterMenuOpen, setFilterMenuOpen] = createSignal(false);
  let filterMenuRef: HTMLDivElement | undefined;
  let filterPopoverEl: HTMLDivElement | undefined;

  const clampCalendarFilter = (el?: HTMLElement | null) => {
    const target = el ?? filterPopoverEl;
    if (!target || typeof window === "undefined") return;
    requestAnimationFrame(() => {
      target.style.transform = "";
      const rect = target.getBoundingClientRect();
      const vw = window.innerWidth;
      if (rect.left < 8) {
        target.style.transform = `translateX(${8 - rect.left}px)`;
      } else if (rect.right > vw - 8) {
        target.style.transform = `translateX(${vw - 8 - rect.right}px)`;
      }
    });
  };

  if (typeof document !== "undefined") {
    const handlePointerDown = (event: PointerEvent) => {
      const target = event.target as Element | null;
      if (
        filterMenuOpen() &&
        filterMenuRef &&
        !filterMenuRef.contains(target) &&
        !target?.closest?.('[role="listbox"], [role="option"], .z-70, [data-portal]')
      ) {
        setFilterMenuOpen(false);
      }
    };
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape" && filterMenuOpen()) {
        setFilterMenuOpen(false);
      }
    };
    document.addEventListener("pointerdown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);
    onCleanup(() => {
      document.removeEventListener("pointerdown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
    });
  }

  const activeFilterCount = createMemo(() => {
    let count = 0;
    if (programSearch().trim().length > 0) count++;
    if (programKindFilter() !== "all") count++;
    if (dateFilter() !== null) count++;
    return count;
  });

  // Program color map
  const programColors = createMemo(() => {
    const map: Record<string, EventColor> = {};
    programs().forEach((p, i) => {
      map[p.id] = EVENT_COLORS[i % EVENT_COLORS.length] ?? FALLBACK_EVENT_COLOR;
    });
    return map;
  });

  // Unique dates with occurrences
  const occurrenceDates = createMemo(() => {
    const datesMap = new Map<string, number>();
    for (const oc of occurrences()) {
      const key = toISODate(new Date(oc.starts_at));
      datesMap.set(key, (datesMap.get(key) ?? 0) + 1);
    }
    return Array.from(datesMap.entries())
      .sort((a, b) => a[0].localeCompare(b[0]))
      .map(([dateKey, count]) => {
        const [y, m, d] = dateKey.split("-").map(Number);
        const dateObj = new Date(y, m - 1, d);
        const weekday = dateObj.toLocaleDateString("pt-BR", { weekday: "short" });
        return {
          dateKey,
          dateObj,
          dayMonth: `${String(d).padStart(2, "0")}/${String(m).padStart(2, "0")}`,
          weekday,
          count,
        };
      });
  });

  const calendarDateOptions = createMemo((): ComboboxOption[] => {
    const list: ComboboxOption[] = [
      { value: "all", label: "Todos os dias com horários" },
    ];
    for (const item of occurrenceDates()) {
      list.push({
        value: item.dateKey,
        label: `${item.dayMonth} (${item.weekday})`,
        description: `${item.count} horário${item.count > 1 ? "s" : ""}`,
      });
    }
    return list;
  });

  // Filtered programs for sidebar
  const filteredPrograms = createMemo(() => {
    const all = programs();
    const currentKind = programKindFilter();
    const q = programSearch().trim().toLowerCase();
    const targetDate = dateFilter();
    const allOccurrences = occurrences();

    const result: ProgramI[] = [];
    for (const p of all) {
      if (currentKind !== "all" && p.kind !== currentKind) {
        continue;
      }
      if (
        q &&
        !p.name.toLowerCase().includes(q) &&
        !(p.description?.toLowerCase().includes(q) ?? false)
      ) {
        continue;
      }
      if (targetDate) {
        let hasDate = false;
        for (const oc of allOccurrences) {
          if (oc.program_id === p.id && toISODate(new Date(oc.starts_at)) === targetDate) {
            hasDate = true;
            break;
          }
        }
        if (!hasDate) continue;
      }
      result.push(p);
    }
    return result;
  });

  // Filtered occurrences for the calendar grid
  const visibleOccurrences = createMemo(() => {
    const all = occurrences();
    const progId = selectedProgramFilter();
    const currentKind = programKindFilter();
    const q = programSearch().trim().toLowerCase();
    const targetDate = dateFilter();
    const allPrograms = programs();

    const result: OccurrenceI[] = [];
    for (const oc of all) {
      if (progId && oc.program_id !== progId) {
        continue;
      }
      if (currentKind !== "all") {
        const prog = allPrograms.find((p) => p.id === oc.program_id);
        if (!prog || prog.kind !== currentKind) continue;
      }
      if (q) {
        const prog = allPrograms.find((p) => p.id === oc.program_id);
        if (!prog || !prog.name.toLowerCase().includes(q)) continue;
      }
      if (targetDate && toISODate(new Date(oc.starts_at)) !== targetDate) {
        continue;
      }
      result.push(oc);
    }
    return result;
  });

  const hasActiveFilters = createMemo(() => {
    return (
      selectedProgramFilter() !== null ||
      programSearch().trim().length > 0 ||
      programKindFilter() !== "all" ||
      dateFilter() !== null
    );
  });

  const clearAllFilters = () => {
    setSelectedProgramFilter(null);
    setProgramSearch("");
    setProgramKindFilter("all");
    setDateFilter(null);
  };

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

  const handleOpenAttendance = (occ: OccurrenceI) => {
    setAttendanceOccurrence(occ);
    setAttendanceModalOpen(true);
  };

  const handleOpenDraw = (occ: OccurrenceI) => {
    setDrawOccurrence(occ);
    setDrawModalOpen(true);
  };

  const targetProgramName = (occ: OccurrenceI | null) => {
    if (!occ) return undefined;
    return programs().find((p) => p.id === occ.program_id)?.name;
  };

  const targetProgramKind = (occ: OccurrenceI | null) => {
    if (!occ) return undefined;
    return programs().find((p) => p.id === occ.program_id)?.kind;
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

  const confirmDeleteOccurrence = async () => {
    const occ = occurrenceToDelete();
    if (!occ) return;
    setDeletingOccurrence(true);
    try {
      await deleteMutation.mutateAsync(occ.id);
      void queryClient.invalidateQueries({
        queryKey: programKeys.occurrences(props.editionId),
      });
      void queryClient.invalidateQueries({
        queryKey: programKeys.byEdition(props.editionId),
      });
      toast.success("Horário excluído!");
      setOccurrenceToDelete(null);
    } catch {
      toast.error("Erro ao excluir horário.");
    } finally {
      setDeletingOccurrence(false);
    }
  };

  const handleRequestDeleteOccurrence = (occurrenceId: string) => {
    const occ = occurrences().find((o) => o.id === occurrenceId);
    setOccurrenceToDelete(occ ?? ({ id: occurrenceId } as OccurrenceI));
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
        <header class="relative z-40 flex h-14 shrink-0 items-center justify-between border-b border-border bg-card px-4 sm:px-6 select-none">
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
              class="inline-flex h-8 items-center gap-1.5 px-3 rounded-lg hover:bg-muted text-xs font-medium text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
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
                class="inline-flex h-8 items-center px-3 text-xs font-medium rounded-lg border border-border hover:bg-muted text-foreground transition-colors cursor-pointer"
              >
                Hoje
              </button>
              <div class="flex items-center">
                <button
                  type="button"
                  onClick={handlePrev}
                  class="size-8 flex items-center justify-center rounded-lg hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                  title="Anterior"
                >
                  <ChevronLeft class="size-4" />
                </button>
                <button
                  type="button"
                  onClick={handleNext}
                  class="size-8 flex items-center justify-center rounded-lg hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                  title="Próximo"
                >
                  <ChevronRight class="size-4" />
                </button>
              </div>
            </div>

            {/* Current View Title */}
            <h1 class="text-sm font-semibold text-foreground capitalize">
              {titleText()}
            </h1>
          </div>

          {/* Right: Filter Popover & View Switcher */}
          <div class="flex items-center gap-2">
            {/* Filter Popover */}
            <div class="relative z-50" ref={(el) => (filterMenuRef = el)}>
              <button
                type="button"
                onClick={() => setFilterMenuOpen((prev) => !prev)}
                class={`relative inline-flex h-8 items-center gap-1.5 px-3 text-xs font-medium rounded-lg border transition-all cursor-pointer ${activeFilterCount() > 0
                  ? "border-primary/50 bg-primary/10 text-primary hover:bg-primary/20"
                  : "border-border bg-background hover:bg-muted text-foreground"
                  }`}
                title="Filtrar calendário"
              >
                <Filter class="size-3.5" />
                <span>Filtrar</span>
                <Show when={activeFilterCount() > 0}>
                  <span class="absolute -top-1.5 -right-1.5 size-4 rounded-full bg-primary text-primary-foreground text-[10px] flex items-center justify-center font-bold shadow-xs">
                    {activeFilterCount()}
                  </span>
                </Show>
              </button>

              <Show when={filterMenuOpen()}>
                <div
                  ref={(el) => {
                    filterPopoverEl = el;
                    clampCalendarFilter(el);
                  }}
                  class="absolute left-0 sm:left-auto sm:right-0 mt-2 w-72 max-w-[calc(100vw-1rem)] rounded-xl border border-border bg-popover p-3 shadow-xl z-50 text-xs space-y-3"
                >
                  <div class="flex items-center justify-between pb-1.5 border-b border-border/60">
                    <div class="flex items-center gap-1.5 font-semibold text-foreground">
                      <Filter class="size-3.5 text-primary" />
                      <span>Filtros do calendário</span>
                    </div>
                    <Show when={hasActiveFilters()}>
                      <button
                        type="button"
                        onClick={clearAllFilters}
                        class="text-[11px] text-primary hover:underline cursor-pointer font-medium"
                      >
                        Limpar todos
                      </button>
                    </Show>
                  </div>

                  {/* Search */}
                  <div class="space-y-1">
                    <label class="text-[11px] font-medium text-muted-foreground">Buscar programa</label>
                    <div class="relative">
                      <Search class="absolute left-2.5 top-1/2 -translate-y-1/2 size-3.5 text-muted-foreground" />
                      <input
                        type="text"
                        placeholder="Nome ou descrição..."
                        value={programSearch()}
                        onInput={(e) => setProgramSearch(e.currentTarget.value)}
                        class="w-full rounded-md border border-border bg-background pl-8 pr-3 py-1.5 text-xs text-foreground placeholder:text-muted-foreground focus:outline-hidden focus:ring-1 focus:ring-primary"
                      />
                    </div>
                  </div>

                  {/* Kind */}
                  <div class="space-y-1">
                    <label class="text-[11px] font-medium text-muted-foreground">Tipo de atividade</label>
                    <div class="flex gap-1 rounded-lg border border-border bg-muted/40 p-1 text-[11px]">
                      <button
                        type="button"
                        onClick={() => setProgramKindFilter("all")}
                        class={`flex-1 py-1 rounded text-center font-medium transition-colors cursor-pointer ${programKindFilter() === "all"
                          ? "bg-background text-foreground shadow-xs font-semibold border border-border/80"
                          : "text-muted-foreground hover:text-foreground border border-transparent"
                          }`}
                      >
                        Todos
                      </button>
                      <button
                        type="button"
                        onClick={() => setProgramKindFilter("activity")}
                        class={`flex-1 py-1 rounded text-center font-medium transition-colors cursor-pointer ${programKindFilter() === "activity"
                          ? "bg-background text-foreground shadow-xs font-semibold border border-border/80"
                          : "text-muted-foreground hover:text-foreground border border-transparent"
                          }`}
                      >
                        Atividades
                      </button>
                      <button
                        type="button"
                        onClick={() => setProgramKindFilter("checkpoint")}
                        class={`flex-1 py-1 rounded text-center font-medium transition-colors cursor-pointer ${programKindFilter() === "checkpoint"
                          ? "bg-background text-foreground shadow-xs font-semibold border border-border/80"
                          : "text-muted-foreground hover:text-foreground border border-transparent"
                          }`}
                      >
                        Checkpoints
                      </button>
                    </div>
                  </div>

                  {/* Day Select */}
                  <Show when={occurrenceDates().length > 0}>
                    <div class="space-y-1">
                      <label class="text-[11px] font-medium text-muted-foreground">Dia específico</label>
                      <Combobox
                        value={dateFilter() ?? "all"}
                        options={calendarDateOptions()}
                        placeholder="Todos os dias com horários"
                        searchPlaceholder="Buscar dia..."
                        onChange={(val) => {
                          const nextVal = val === "all" ? null : val;
                          setDateFilter(nextVal);
                          if (nextVal) {
                            const [y, m, d] = nextVal.split("-").map(Number);
                            setCurrentDate(new Date(y, m - 1, d, 12, 0, 0));
                            setView("day");
                            setFilterMenuOpen(false);
                          }
                        }}
                        class="w-full"
                        triggerClass="h-8 text-xs"
                      />
                    </div>
                  </Show>
                </div>
              </Show>
            </div>

            {/* View Switcher Segmented Control */}
            <div class="inline-flex h-8 items-center gap-1 p-1 rounded-lg border border-border bg-muted/40">
              {(["day", "week", "month", "year"] as const).map((v) => {
                const labels: Record<CalendarView, string> = {
                  day: "Dia",
                  week: "Semana",
                  month: "Mês",
                  year: "Ano",
                };
                const isActive = () => view() === v;
                return (
                  <button
                    type="button"
                    onClick={() => setView(v)}
                    class={`inline-flex h-full items-center px-3 text-xs font-medium rounded-md transition-all cursor-pointer ${isActive()
                      ? "bg-background text-foreground shadow-xs font-semibold border border-border/80"
                      : "text-muted-foreground hover:text-foreground border border-transparent"
                      }`}
                  >
                    {labels[v]}
                  </button>
                );
              })}
            </div>
          </div>
        </header>

        {/* Calendar Body */}
        <div class="flex flex-1 overflow-hidden">
          {/* Left sidebar */}
          <div class="w-64 shrink-0 border-r border-border bg-card p-4 flex flex-col gap-4 overflow-y-auto select-none">
            {/* Mini Calendar */}
            <MiniCalendar
              currentDate={currentDate()}
              onDateClick={(d: Date) => {
                setCurrentDate(d);
                if (view() === "year") setView("day");
              }}
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

            {/* Program Drag & Select List */}
            <div class="space-y-2 pt-2 border-t border-border flex-1 flex flex-col">
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold text-foreground uppercase tracking-wider">
                  Programas ({filteredPrograms().length})
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

              <div class="space-y-1.5 mt-2 flex-1 overflow-y-auto p-1">
                <For
                  each={filteredPrograms()}
                  fallback={
                    <p class="text-xs text-muted-foreground/60 py-2">
                      {hasActiveFilters()
                        ? "Nenhum programa corresponde aos filtros."
                        : "Nenhum programa cadastrado nesta edição."}
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
                        onDragEnd={endCalendarDrag}
                        onClick={() => {
                          setSelectedProgramFilter((prev) =>
                            prev === prog.id ? null : prog.id,
                          );
                        }}
                        class={`flex items-center gap-2.5 px-3 py-2 rounded-lg border text-xs font-medium cursor-grab active:cursor-grabbing transition-all select-none ${isFiltered()
                          ? "border-primary bg-primary/10 shadow-xs ring-2 ring-primary/30"
                          : "border-border/60 hover:border-border hover:bg-muted/70 bg-card"
                          }`}
                        style={{
                          "border-left-width": "4px",
                          "border-left-color": color().border,
                        }}
                      >
                        <GripVertical class="size-3 text-muted-foreground/50 shrink-0" />
                        <span class="truncate flex-1">{prog.name}</span>
                        <Show when={prog.kind === "checkpoint"}>
                          <span class="text-[9px] px-1 py-0.2 rounded bg-muted text-muted-foreground font-mono shrink-0">
                            CP
                          </span>
                        </Show>
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
                onDeleteOccurrence={handleRequestDeleteOccurrence}
                onOpenAttendance={handleOpenAttendance}
                onOpenDraw={handleOpenDraw}
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
                onDeleteOccurrence={handleRequestDeleteOccurrence}
                onOpenAttendance={handleOpenAttendance}
                onOpenDraw={handleOpenDraw}
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

        {/* Create/Edit Occurrence Dialog (Google Calendar style with Presença & Sorteio) */}
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
          onOpenAttendance={handleOpenAttendance}
          onOpenDraw={handleOpenDraw}
        />

        {/* Attendance Modal */}
        <OccurrenceAttendanceModal
          open={attendanceModalOpen()}
          onOpenChange={setAttendanceModalOpen}
          occurrence={attendanceOccurrence()}
          programName={targetProgramName(attendanceOccurrence())}
          programKind={targetProgramKind(attendanceOccurrence())}
        />

        {/* Draw Modal with Projector Mode */}
        <OccurrenceDrawModal
          open={drawModalOpen()}
          onOpenChange={setDrawModalOpen}
          occurrence={drawOccurrence()}
          programName={targetProgramName(drawOccurrence())}
          programKind={targetProgramKind(drawOccurrence())}
        />

        {/* Delete Occurrence Confirmation Modal */}
        <AlertModal
          open={occurrenceToDelete() !== null}
          onOpenChange={(open) => !open && setOccurrenceToDelete(null)}
          title="Excluir ocorrência"
          description="Tem certeza que deseja remover esta ocorrência da grade de horários? Esta ação não pode ser desfeita."
          confirmLabel="Excluir ocorrência"
          variant="destructive"
          loading={deletingOccurrence()}
          onConfirm={confirmDeleteOccurrence}
        />
      </div>
    </DesktopOnly>
  );
}
