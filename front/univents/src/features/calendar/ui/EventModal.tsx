import { Clock, FileText, Tag, Users, X } from "lucide-react";
import { useEffect, useState } from "react";
import { formatFullDate, toISODateTimeLocal } from "../lib/date";
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

  const displayDate = start ? new Date(start) : new Date();

  return (
    <div
      className="fixed inset-0 z-100 flex items-center justify-center"
      style={{ background: "rgba(0,0,0,0.45)" }}
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div className="bg-popover border border-border rounded-2xl w-110 max-w-[92vw] shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-200">
        {/* Header */}
        <div className="px-6 pt-5 flex items-start justify-between">
          <div className="flex-1 text-[22px] font-normal text-foreground pb-2">
            {mode === "edit" ? prog?.name || "Evento" : "Novo evento"}
          </div>
          <button
            type="button"
            onClick={onClose}
            className="ml-3 p-1.5 rounded-full hover:bg-muted transition-colors text-muted-foreground"
          >
            <X size={20} />
          </button>
        </div>

        {/* Body */}
        <div className="px-6 py-4 flex flex-col gap-3.5">
          {/* Date/Time */}
          <div className="flex items-start gap-3.5">
            <div className="w-5 h-5 mt-0.5 flex items-center justify-center text-muted-foreground shrink-0">
              <Clock size={18} />
            </div>
            <div className="flex-1 flex items-center gap-2 flex-wrap text-sm text-foreground">
              <span className="text-muted-foreground">
                {formatFullDate(displayDate)}
              </span>
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
              <span className="text-muted-foreground">–</span>
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
          <div className="flex items-start gap-3.5">
            <div className="w-5 h-5 mt-0.5 flex items-center justify-center text-muted-foreground shrink-0">
              <Tag size={18} />
            </div>
            <div className="flex-1">
              <CalendarCombobox
                value={programId}
                onChange={setProgramId}
                disabled={mode === "edit"}
                placeholder="Selecionar programa"
                options={programs.map((p) => ({ value: p.id, label: p.name }))}
              />
            </div>
          </div>

          {/* Capacity */}
          <div className="flex items-start gap-3.5">
            <div className="w-5 h-5 mt-0.5 flex items-center justify-center text-muted-foreground shrink-0">
              <Users size={18} />
            </div>
            <div className="flex-1">
              <label className="block space-y-1.5">
                <span className="text-xs font-medium text-muted-foreground">
                  Capacidade máxima{" "}
                  <span className="font-normal">(opcional)</span>
                </span>
                <input
                  type="number"
                  min={0}
                  placeholder="Ex.: 100"
                  value={cap}
                  onChange={(e) => setCap(e.target.value)}
                  className="w-full rounded-xl border border-border/60 bg-background px-3 py-2.5 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/15"
                />
              </label>
            </div>
          </div>

          {/* Description */}
          <div className="flex items-start gap-3.5">
            <div className="w-5 h-5 mt-0.5 flex items-center justify-center text-muted-foreground shrink-0">
              <FileText size={18} />
            </div>
            <div className="flex-1 text-[13px] text-muted-foreground">
              {prog?.description || ""}
              {prog?.price ? ` · R$ ${prog.price}` : ""}
              {prog?.staff_only ? " · Staff only" : ""}
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-2 px-6 pb-5 pt-2">
          {mode === "edit" && occurrence && (
            <button
              type="button"
              className="mr-auto px-4 py-2 text-sm font-medium text-destructive rounded-md hover:bg-destructive/10 transition-colors"
              onClick={() => occurrence && onDelete?.(occurrence.id)}
            >
              Excluir
            </button>
          )}
          <button
            type="button"
            className="px-4 py-2 text-sm font-medium text-foreground rounded-md hover:bg-muted transition-colors"
            onClick={onClose}
          >
            Cancelar
          </button>
          <button
            type="button"
            className="px-5 py-2 text-sm font-medium text-primary-foreground bg-primary rounded-full hover:opacity-90 transition-opacity"
            onClick={handleSave}
          >
            {mode === "edit" ? "Salvar" : "Criar"}
          </button>
        </div>
      </div>
    </div>
  );
}
