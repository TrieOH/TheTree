import { CalendarDays, Gift, Pencil, Trash2, UserCheck } from "lucide-react";
import { Button } from "@/shared/ui/shadcn/button";
import type { OccurrenceI } from "../model";

export function OccurrenceAdminCard({
  occurrence,
  onEdit,
  onDelete,
  onAttendance,
  onDraw,
}: {
  occurrence: OccurrenceI;
  onEdit: () => void;
  onDelete: () => void;
  onAttendance: () => void;
  onDraw?: () => void;
}) {
  return (
    <div className="flex min-w-0 flex-col gap-3 rounded-xl border border-border/60 bg-card px-3 py-3 transition-colors hover:bg-accent/5 sm:flex-row! sm:items-center! sm:justify-between! sm:px-4!">
      <div className="min-w-0 flex-1 space-y-0.5">
        <p className="flex min-w-0 items-center gap-2 text-sm font-medium">
          <CalendarDays className="size-3.5 text-muted-foreground" />
          <span className="truncate">
            {new Date(occurrence.starts_at).toLocaleDateString("pt-BR")}
          </span>
        </p>
        <p className="truncate text-xs text-muted-foreground">
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
      <div className="flex min-w-0 flex-wrap gap-1.5 sm:shrink-0 sm:gap-2">
        {onDraw ? (
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="min-w-0 px-2 text-xs sm:px-3 sm:text-sm"
            onClick={onDraw}
          >
            <Gift className="mr-1.5 size-3.5 shrink-0 sm:mr-2" />
            <span className="truncate">Sortear</span>
          </Button>
        ) : null}
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="min-w-0 px-2 text-xs sm:px-3 sm:text-sm"
          onClick={onAttendance}
        >
          <UserCheck className="mr-1.5 size-3.5 shrink-0 sm:mr-2" />
          <span className="truncate">Presença</span>
        </Button>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="px-2 text-xs sm:px-3 sm:text-sm"
          onClick={onEdit}
        >
          <Pencil className="mr-1.5 size-3.5 sm:mr-2" />
          <span>Editar</span>
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
