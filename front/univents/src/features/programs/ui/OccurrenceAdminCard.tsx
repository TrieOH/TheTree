import { CalendarDays, Pencil, Trash2, UserCheck } from "lucide-react";
import { Button } from "@/shared/ui/shadcn/button";
import type { OccurrenceI } from "../model";

export function OccurrenceAdminCard({
  occurrence,
  onEdit,
  onDelete,
  onAttendance,
}: {
  occurrence: OccurrenceI;
  onEdit: () => void;
  onDelete: () => void;
  onAttendance: () => void;
}) {
  return (
    <div className="flex items-center justify-between rounded-xl border border-border/60 bg-card px-4 py-3 transition-colors hover:bg-accent/5">
      <div className="min-w-0 flex-1 space-y-0.5">
        <p className="flex items-center gap-2 text-sm font-medium">
          <CalendarDays className="size-3.5 text-muted-foreground" />
          {new Date(occurrence.starts_at).toLocaleDateString("pt-BR")}
        </p>
        <p className="text-xs text-muted-foreground">
          {new Date(occurrence.starts_at).toLocaleTimeString("pt-BR", {
            hour: "2-digit",
            minute: "2-digit",
          })}{" "}
          –{" "}
          {new Date(occurrence.ends_at).toLocaleTimeString("pt-BR", {
            hour: "2-digit",
            minute: "2-digit",
          })}
          {occurrence.max_capacity
            ? ` · ${occurrence.max_capacity} vagas`
            : " · Vagas ilimitadas"}
        </p>
      </div>
      <div className="flex gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={onAttendance}
        >
          <UserCheck className="mr-2 size-3.5" />
          Presença
        </Button>
        <Button type="button" variant="outline" size="sm" onClick={onEdit}>
          <Pencil className="mr-2 size-3.5" />
          Editar
        </Button>
        <Button
          type="button"
          variant="destructive"
          size="icon"
          onClick={onDelete}
          aria-label="Excluir ocorrência"
        >
          <Trash2 className="size-3.5" />
        </Button>
      </div>
    </div>
  );
}
