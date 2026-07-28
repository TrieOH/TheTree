import { FileText } from "lucide-react";
import { useEffect, useState } from "react";
import { toISODateTimeLocal } from "../lib/date";
import type { OccurrenceI, ProgramI } from "../model";
import { CalendarCombobox } from "./CalendarCombobox";

interface EventModalProps {
  mode: "create" | "edit";
  programs: ProgramI[];
  occurrence?: OccurrenceI | null;
  initialProgramId?: string;
  initialDate?: string;
  initialHour?: number;
  onSave: (data: {
    program_id: string;
    starts_at: string;
    ends_at: string;
    max_capacity?: number;
    id?: string;
  }) => void;
  onDelete?: (occurrenceId: string) => void;
  onClose: () => void;
}

export function EventModal({
  mode,
  programs,
  occurrence,
  initialProgramId,
  initialDate,
  initialHour,
  onSave,
  onDelete,
  onClose,
}: EventModalProps) {
  const [programId, setProgramId] = useState(
    initialProgramId || programs[0]?.id || "",
  );
  const [start, setStart] = useState("");
  const [end, setEnd] = useState("");
  const [cap, setCap] = useState<string>("");

  useEffect(() => {
    if (mode === "edit" && occurrence) {
      setProgramId(occurrence.program_id);
      setStart(toISODateTimeLocal(new Date(occurrence.starts_at)));
      setEnd(toISODateTimeLocal(new Date(occurrence.ends_at)));
      setCap(occurrence.max_capacity?.toString() || "");
    } else if (mode === "create") {
      // YYYY-MM-DD is a calendar date, not an UTC timestamp. Parsing it with
      // new Date(string) shifts it to the previous day in São Paulo.
      const d = initialDate
        ? (() => {
            const [year, month, day] = initialDate.split("-").map(Number);
            return new Date(year, month - 1, day);
          })()
        : new Date();
      const h = initialHour ?? 9;
      const p = (n: number) => String(n).padStart(2, "0");
      const dateStr = `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}`;
      setStart(`${dateStr}T${p(h)}:00`);
      setEnd(`${dateStr}T${p(h + 1)}:00`);
      setCap("");
    }
  }, [mode, occurrence, initialDate, initialHour]);

  const prog = programs.find((p) => p.id === programId);

  const handleSave = () => {
    if (!start || !end) return;
    const startDate = new Date(start);
    const endDate = new Date(end);
    if (endDate <= startDate) endDate.setDate(endDate.getDate() + 1);
    onSave({
      program_id: programId,
      starts_at: startDate.toISOString(),
      ends_at: endDate.toISOString(),
      max_capacity: cap ? parseInt(cap, 10) : undefined,
      id: occurrence?.id,
    });
  };

  return (
    <div
      className="fixed inset-0 z-100 flex items-center justify-center"
      style={{ background: "rgba(0,0,0,0.45)" }}
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div className="relative w-160 max-w-[calc(100vw-2rem)] overflow-hidden rounded-3xl border border-border/70 bg-popover shadow-2xl animate-in fade-in zoom-in-95 duration-200">
        {/* Body */}
        <div className="flex flex-col gap-2.5 px-6 py-7">
          {/* Date/Time */}
          <div>
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-[1fr_auto_1fr] sm:items-end">
              <label className="flex flex-col gap-1.5">
                <span className="text-xs font-medium text-muted-foreground">
                  Início
                </span>
                <input
                  type="datetime-local"
                  value={start}
                  onChange={(e) => setStart(e.target.value)}
                  className="w-full rounded-xl border border-border/60 bg-background px-3 py-2.5 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/15"
                />
              </label>
              <span className="hidden pb-3 text-center text-muted-foreground sm:block">
                até
              </span>
              <label className="flex flex-col gap-1.5">
                <span className="text-xs font-medium text-muted-foreground">
                  Término
                </span>
                <input
                  type="datetime-local"
                  value={end}
                  onChange={(e) => setEnd(e.target.value)}
                  className="w-full rounded-xl border border-border/60 bg-background px-3 py-2.5 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/15"
                />
              </label>
            </div>
          </div>

          {/* Program */}
          <div className="block space-y-1.5">
            <span className="text-xs font-medium text-muted-foreground">
              Programa
            </span>
            <CalendarCombobox
              value={programId}
              onChange={setProgramId}
              disabled={mode === "edit"}
              placeholder="Selecionar programa"
              options={programs.map((p) => ({ value: p.id, label: p.name }))}
            />
          </div>

          {/* Capacity */}
          <div>
            <label className="block space-y-1.5">
              <span className="text-xs font-medium text-muted-foreground">
                Capacidade máxima{" "}
                <span className="font-normal">(opcional)</span>
              </span>
              <input
                type="number"
                min={1}
                placeholder="Ex.: 100"
                value={cap}
                onChange={(e) => setCap(e.target.value)}
                className="w-full rounded-xl border border-border/60 bg-background px-3 py-2.5 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/15"
              />
            </label>
          </div>

          {/* Description */}
          {(prog?.description || prog?.price || prog?.staff_only) && (
            <div className="rounded-2xl border border-border/60 bg-muted/20 p-4">
              {prog?.description && (
                <div className="flex items-start gap-3">
                  <FileText className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
                  <p className="text-sm leading-relaxed text-muted-foreground">
                    {prog.description}
                  </p>
                </div>
              )}
              {(prog?.price || prog?.staff_only) && (
                <div className="mt-3 flex flex-wrap gap-2 text-xs text-muted-foreground">
                  {prog.price && (
                    <span className="rounded-full bg-background px-2.5 py-1">
                      R$ {prog.price}
                    </span>
                  )}
                  {prog.staff_only && (
                    <span className="rounded-full bg-background px-2.5 py-1">
                      Somente equipe
                    </span>
                  )}
                </div>
              )}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center gap-3 border-t border-border/60 bg-muted/20 px-6 py-4">
          {mode === "edit" && occurrence && (
            <button
              type="button"
              className="mr-auto rounded-xl px-3 py-2 text-sm font-medium text-destructive transition-colors hover:bg-destructive/10"
              onClick={() => occurrence && onDelete?.(occurrence.id)}
            >
              Excluir
            </button>
          )}
          <button
            type="button"
            className="rounded-xl px-4 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
            onClick={onClose}
          >
            Cancelar
          </button>
          <button
            type="button"
            className="rounded-xl bg-primary px-5 py-2.5 text-sm font-semibold text-primary-foreground shadow-sm transition-opacity hover:opacity-90"
            onClick={handleSave}
          >
            {mode === "edit" ? "Salvar" : "Criar"}
          </button>
        </div>
      </div>
    </div>
  );
}
