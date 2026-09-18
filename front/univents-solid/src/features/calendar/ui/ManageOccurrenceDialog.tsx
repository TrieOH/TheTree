import type { JSX } from "@solidjs/web";
import { Portal } from "@solidjs/web";
import { Show, createEffect, createMemo, createSignal, onCleanup } from "solid-js";

import ClockIcon from "~icons/lucide/clock";
import GiftIcon from "~icons/lucide/gift";
import TrashIcon from "~icons/lucide/trash-2";
import UserCheckIcon from "~icons/lucide/user-check";
import UsersIcon from "~icons/lucide/users";
import XIcon from "~icons/lucide/x";

import type { OccurrenceCreateOutput } from "@/features/programs/model";
import { Combobox, type ComboboxOption } from "@/shared/ui/Combobox";
import { Button } from "@trieoh/ui-solid";
import { AlertModal } from "@/widgets/ui/AlertModal";
import { toISODate } from "../lib/date";
import type { OccurrenceI, ProgramI } from "../model";

const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Gift = GiftIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash = TrashIcon as unknown as (props: { class?: string }) => JSX.Element;
const UserCheck = UserCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;
const X = XIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface ManageOccurrenceDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  programs: ProgramI[];
  occurrence?: OccurrenceI | null;
  initialProgramId?: string;
  initialDate?: string; // YYYY-MM-DD
  initialHour?: number;
  onSave: (data: {
    programId: string;
    occurrenceData: OccurrenceCreateOutput;
    id?: string;
  }) => Promise<boolean>;
  onDelete?: (occurrenceId: string) => Promise<boolean>;
  onOpenAttendance?: (occurrence: OccurrenceI) => void;
  onOpenDraw?: (occurrence: OccurrenceI) => void;
}

interface FormValues {
  programId: string;
  date: string; // YYYY-MM-DD
  startTime: string; // HH:mm
  endTime: string; // HH:mm
  maxCapacity: string;
}

function padZero(n: number): string {
  return String(n).padStart(2, "0");
}

function initValues(
  programs: ProgramI[],
  occurrence?: OccurrenceI | null,
  initialProgramId?: string,
  initialDate?: string,
  initialHour?: number,
): FormValues {
  if (occurrence) {
    const s = new Date(occurrence.starts_at);
    const e = new Date(occurrence.ends_at);
    return {
      programId: occurrence.program_id,
      date: toISODate(s),
      startTime: `${padZero(s.getHours())}:${padZero(s.getMinutes())}`,
      endTime: `${padZero(e.getHours())}:${padZero(e.getMinutes())}`,
      maxCapacity: occurrence.max_capacity != null ? String(occurrence.max_capacity) : "",
    };
  }

  const progId = initialProgramId || programs[0]?.id || "";
  const dateStr = initialDate ?? toISODate(new Date());
  const h = initialHour ?? 9;

  const nextH = (h + 1) % 24;

  return {
    programId: progId,
    date: dateStr,
    startTime: `${padZero(h)}:00`,
    endTime: `${padZero(nextH)}:00`,
    maxCapacity: "",
  };
}

/**
 * Google Calendar-style floating quick card to create or edit an occurrence.
 * Features inline Combobox for program search/selection, clean date/time inputs,
 * quick Presença & Sorteio actions, and quick delete/save.
 */
export function ManageOccurrenceDialog(props: ManageOccurrenceDialogProps): JSX.Element {
  const [values, setValues] = createSignal<FormValues>({
    programId: "",
    date: "",
    startTime: "09:00",
    endTime: "10:00",
    maxCapacity: "",
  });
  const [submitting, setSubmitting] = createSignal(false);
  const [deleting, setDeleting] = createSignal(false);
  const [confirmDeleteOpen, setConfirmDeleteOpen] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<"programId" | "time", string>>>({});

  // Reset form when opened or target changed
  createEffect(
    () =>
      [
        props.open,
        props.occurrence,
        props.initialProgramId,
        props.initialDate,
        props.initialHour,
        props.programs,
      ] as const,
    ([open, occ, initProg, initDate, initH, progs]) => {
      if (!open) return;
      setValues(initValues(progs, occ, initProg, initDate, initH));
      setErrors({});
    },
  );

  // Auto-select first program if empty
  createEffect(
    () =>
      [
        props.open,
        values().programId,
        props.programs,
        props.initialProgramId,
      ] as const,
    ([open, currentProgId, progs, initProg]) => {
      if (open && !currentProgId && progs.length > 0) {
        setValues((prev) => ({
          ...prev,
          programId: initProg || progs[0]?.id || "",
        }));
      }
    },
  );

  // Keyboard escape listener
  if (typeof window !== "undefined") {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (props.open && e.key === "Escape" && !submitting() && !deleting()) {
        props.onOpenChange(false);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    onCleanup(() => window.removeEventListener("keydown", handleKeyDown));
  }

  const update = <K extends keyof FormValues>(key: K, value: FormValues[K]) => {
    setValues((prev) => ({ ...prev, [key]: value }));
    setErrors({});
  };

  const isEditing = () => Boolean(props.occurrence);

  const programOptions = createMemo((): ComboboxOption[] => {
    return props.programs.map((p) => ({
      value: p.id,
      label: p.name,
      description: p.kind === "checkpoint" ? "Checkpoint" : "Atividade",
    }));
  });

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    const current = values();
    const errs: Partial<Record<"programId" | "time", string>> = {};

    if (!current.programId) {
      errs.programId = "Selecione um programa para vincular.";
    }

    if (!current.date || !current.startTime || !current.endTime) {
      errs.time = "Preencha a data e os horários de início e término.";
    }

    const startDate = new Date(`${current.date}T${current.startTime}:00`);
    const endDate = new Date(`${current.date}T${current.endTime}:00`);

    if (isNaN(startDate.getTime()) || isNaN(endDate.getTime())) {
      errs.time = "Data ou horário inválido.";
    } else if (endDate <= startDate) {
      errs.time = "O horário de término deve ser posterior ao início.";
    }

    if (Object.keys(errs).length > 0) {
      setErrors(errs);
      return;
    }

    const cap = current.maxCapacity.trim() ? parseInt(current.maxCapacity, 10) : undefined;

    setSubmitting(true);
    try {
      const ok = await props.onSave({
        programId: current.programId,
        id: props.occurrence?.id,
        occurrenceData: {
          starts_at: startDate.toISOString(),
          ends_at: endDate.toISOString(),
          max_capacity: cap,
        },
      });

      if (ok) {
        props.onOpenChange(false);
      }
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (!props.occurrence || !props.onDelete) return;
    setDeleting(true);
    try {
      const ok = await props.onDelete(props.occurrence.id);
      if (ok) {
        props.onOpenChange(false);
      }
    } finally {
      setDeleting(false);
    }
  };

  return (
    <Show when={props.open}>
      <Portal>
        <div class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6 select-none animate-in fade-in duration-100">
          {/* Lightweight non-intrusive backdrop (Google Calendar style) */}
          <div
            class="fixed inset-0 bg-black/30 backdrop-blur-[1px] transition-opacity"
            onClick={() => !submitting() && !deleting() && props.onOpenChange(false)}
          />

          {/* Google Agenda Floating Card */}
          <div class="relative w-full max-w-md rounded-2xl border border-border bg-card shadow-2xl z-10 overflow-hidden animate-in zoom-in-95 duration-100">
            {/* Top Bar: Title & Quick Actions */}
            <div class="flex items-center justify-between px-4 py-3 border-b border-border/50 bg-muted/20">
              <div class="flex items-center gap-2">
                <span class="size-2.5 rounded-full bg-primary" />
                <span class="text-xs font-semibold text-foreground uppercase tracking-wider">
                  {isEditing() ? "Editar Ocorrência" : "Novo Horário"}
                </span>
              </div>

              <div class="flex items-center gap-1">
                <Show when={isEditing() && props.onDelete}>
                  <button
                    type="button"
                    onClick={() => setConfirmDeleteOpen(true)}
                    disabled={submitting() || deleting()}
                    title="Excluir ocorrência"
                    class="size-7 flex items-center justify-center rounded-lg text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors cursor-pointer"
                  >
                    <Trash class="size-4" />
                  </button>
                </Show>
                <button
                  type="button"
                  onClick={() => props.onOpenChange(false)}
                  disabled={submitting() || deleting()}
                  title="Fechar"
                  class="size-7 flex items-center justify-center rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted transition-colors cursor-pointer"
                >
                  <X class="size-4" />
                </button>
              </div>
            </div>

            {/* Quick Actions for Existing Occurrences: Presença & Sorteio */}
            <Show when={isEditing()}>
              <div class="grid grid-cols-2 gap-2 px-5 pt-3.5">
                <button
                  type="button"
                  onClick={() => {
                    if (!props.occurrence) return;
                    props.onOpenChange(false);
                    props.onOpenAttendance?.(props.occurrence);
                  }}
                  class="flex items-center justify-center gap-2 py-2 px-3 rounded-xl border border-border/70 bg-muted/30 hover:bg-muted/70 text-xs font-medium text-foreground transition-colors cursor-pointer"
                >
                  <UserCheck class="size-4 text-primary" />
                  <span>Presença</span>
                </button>

                <button
                  type="button"
                  onClick={() => {
                    if (!props.occurrence) return;
                    props.onOpenChange(false);
                    props.onOpenDraw?.(props.occurrence);
                  }}
                  class="flex items-center justify-center gap-2 py-2 px-3 rounded-xl border border-border/70 bg-muted/30 hover:bg-muted/70 text-xs font-medium text-foreground transition-colors cursor-pointer"
                >
                  <Gift class="size-4 text-amber-500" />
                  <span>Sorteio</span>
                </button>
              </div>
            </Show>

            {/* Quick Form */}
            <form onSubmit={handleSubmit} class="p-5 space-y-4">
              {/* Program Selector using Combobox */}
              <div class="space-y-1.5">
                <label class="text-xs font-medium text-foreground">
                  Programa vinculado
                </label>
                <Combobox
                  value={values().programId}
                  options={programOptions()}
                  placeholder="Selecione um programa..."
                  searchPlaceholder="Buscar por nome..."
                  disabled={submitting() || isEditing()}
                  onChange={(val) => update("programId", val)}
                />
                <Show when={errors().programId}>
                  <p class="text-xs text-destructive">{errors().programId}</p>
                </Show>
              </div>

              {/* Date & Time Row */}
              <div class="space-y-1.5">
                <label class="text-xs font-medium text-foreground flex items-center gap-1.5">
                  <Clock class="size-3.5 text-muted-foreground" />
                  Data e Horário
                </label>

                <div class="grid grid-cols-12 gap-2">
                  {/* Date Input */}
                  <div class="col-span-12 sm:col-span-6">
                    <input
                      type="date"
                      value={values().date}
                      onInput={(e) => update("date", e.currentTarget.value)}
                      disabled={submitting()}
                      class="h-9 w-full rounded-lg border border-border bg-background px-3 text-xs text-foreground focus:border-primary focus:outline-none"
                    />
                  </div>

                  {/* Start Time */}
                  <div class="col-span-6 sm:col-span-3">
                    <input
                      type="time"
                      value={values().startTime}
                      onInput={(e) => update("startTime", e.currentTarget.value)}
                      disabled={submitting()}
                      class="h-9 w-full rounded-lg border border-border bg-background px-2 text-xs text-foreground text-center focus:border-primary focus:outline-none"
                    />
                  </div>

                  {/* End Time */}
                  <div class="col-span-6 sm:col-span-3">
                    <input
                      type="time"
                      value={values().endTime}
                      onInput={(e) => update("endTime", e.currentTarget.value)}
                      disabled={submitting()}
                      class="h-9 w-full rounded-lg border border-border bg-background px-2 text-xs text-foreground text-center focus:border-primary focus:outline-none"
                    />
                  </div>
                </div>

                <Show when={errors().time}>
                  <p class="text-xs text-destructive">{errors().time}</p>
                </Show>
              </div>

              {/* Max Capacity */}
              <div class="space-y-1.5">
                <label class="text-xs font-medium text-foreground flex items-center gap-1.5">
                  <Users class="size-3.5 text-muted-foreground" />
                  Capacidade de vagas (opcional)
                </label>
                <input
                  type="number"
                  min="1"
                  placeholder="Deixe em branco para ilimitado"
                  value={values().maxCapacity}
                  onInput={(e) => update("maxCapacity", e.currentTarget.value)}
                  disabled={submitting()}
                  class="h-9 w-full rounded-lg border border-border bg-background px-3 text-xs text-foreground placeholder:text-muted-foreground/60 focus:border-primary focus:outline-none"
                />
              </div>

              {/* Action Buttons */}
              <div class="flex items-center justify-end gap-2 pt-2 border-t border-border/50">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => props.onOpenChange(false)}
                  disabled={submitting() || deleting()}
                  class="text-xs font-medium cursor-pointer"
                >
                  Cancelar
                </Button>
                <Button
                  type="submit"
                  size="sm"
                  disabled={submitting() || deleting()}
                  class="text-xs font-semibold cursor-pointer"
                >
                  {submitting()
                    ? "Salvando..."
                    : isEditing()
                    ? "Salvar alterações"
                    : "Adicionar"}
                </Button>
              </div>
            </form>
          </div>
        </div>

        <AlertModal
          open={confirmDeleteOpen()}
          onOpenChange={setConfirmDeleteOpen}
          title="Excluir ocorrência"
          description="Tem certeza que deseja excluir esta ocorrência? Participantes inscritos perderão o agendamento."
          confirmLabel="Excluir ocorrência"
          variant="destructive"
          loading={deleting()}
          onConfirm={() => {
            void handleDelete();
            setConfirmDeleteOpen(false);
          }}
        />
      </Portal>
    </Show>
  );
}
