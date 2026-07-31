import {
  CalendarDays,
  Clock3,
  MoreVertical,
  Pencil,
  ShieldCheck,
  Trash2,
} from "lucide-react";
import { motion } from "motion/react";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/shared/ui/shadcn/context-menu";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/shared/ui/shadcn/dropdown-menu";
import type { OccurrenceI, ProgramI } from "../model";

function gradient(index: number) {
  return [
    "from-violet-500/20 via-fuchsia-500/10 to-background",
    "from-emerald-500/20 via-teal-500/10 to-background",
    "from-amber-500/20 via-orange-500/10 to-background",
    "from-blue-500/20 via-cyan-500/10 to-background",
  ][index % 4];
}

export function ProgramAdminCard({
  program,
  occurrences,
  index,
  onEdit,
  onManageOccurrences,
  onDelete,
}: {
  program: ProgramI;
  occurrences: OccurrenceI[];
  index: number;
  onEdit: () => void;
  onManageOccurrences: () => void;
  onDelete: () => void;
}) {
  const card = (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05, duration: 0.35 }}
      className="group flex w-full flex-col overflow-hidden rounded-2xl bg-card text-left ring-1 ring-foreground/10 shadow-xs transition-all duration-300 hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm"
    >
      <div
        role="button"
        tabIndex={0}
        onClick={onEdit}
        onKeyDown={(event) => {
          if (event.key === "Enter" || event.key === " ") onEdit();
        }}
        className={cn(
          "relative aspect-video overflow-hidden bg-linear-to-br text-left",
          gradient(index),
        )}
      >
        <div className="flex h-full items-center justify-center">
          <div className="flex size-18 items-center justify-center rounded-full border border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
            <CalendarDays className="size-7 text-muted-foreground/45" />
          </div>
        </div>
        <div className="absolute inset-0 bg-linear-to-t from-background/95 via-background/25 to-transparent" />
        <div className="absolute left-3 top-3 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-2 py-0.5 text-[10px] font-medium text-emerald-700">
          <span className="mr-1.5 inline-block size-1.5 rounded-full bg-emerald-500" />
          Ativo
        </div>
        <div className="absolute right-3 top-3">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <button
                  type="button"
                  onClick={(event) => event.stopPropagation()}
                  className="inline-flex size-9 items-center justify-center rounded-full bg-background/85 text-foreground shadow-sm backdrop-blur-sm hover:bg-background"
                  aria-label="Abrir ações do programa"
                />
              }
            >
              <MoreVertical className="size-4" />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuItem onClick={onEdit}>
                <Pencil className="size-4" />
                Editar programa
              </DropdownMenuItem>
              <DropdownMenuItem onClick={onManageOccurrences}>
                <CalendarDays className="size-4" />
                Gerenciar ocorrências
              </DropdownMenuItem>
              <DropdownMenuItem
                variant="destructive"
                onClick={(event) => {
                  event.stopPropagation();
                  onDelete();
                }}
              >
                <Trash2 className="size-4" />
                Excluir programa
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        <div className="absolute inset-x-0 bottom-0 p-4 sm:p-5">
          <p className="text-[10px] font-medium uppercase tracking-wide text-muted-foreground">
            {program.kind === "activity" ? "Atividade" : "Checkpoint"}
          </p>
          <h3 className="line-clamp-1 text-lg font-semibold leading-snug group-hover:text-primary sm:text-xl">
            {program.name}
          </h3>
        </div>
      </div>
      <div className="space-y-3 p-4 sm:p-5 sm:pt-4">
        <p className="line-clamp-2 min-h-10 text-sm text-muted-foreground">
          {program.description || "Sem descrição cadastrada."}
        </p>
        <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
          <span className="inline-flex items-center gap-1.5">
            <Clock3 className="size-3.5" />
            {occurrences.length} ocorrência(s)
          </span>
          {program.staff_only && (
            <span className="inline-flex items-center gap-1.5">
              <ShieldCheck className="size-3.5" />
              Somente equipe
            </span>
          )}
        </div>
        <Button
          type="button"
          variant="outline"
          className="h-9 w-full gap-2"
          onClick={onManageOccurrences}
        >
          Gerenciar ocorrências
        </Button>
      </div>
    </motion.article>
  );

  return (
    <ContextMenu>
      <ContextMenuTrigger render={card} />
      <ContextMenuContent className="w-56">
        <ContextMenuItem onClick={onEdit}>
          <Pencil className="size-4" />
          Editar programa
        </ContextMenuItem>
        <ContextMenuItem onClick={onManageOccurrences}>
          <CalendarDays className="size-4" />
          Gerenciar ocorrências
        </ContextMenuItem>
        <ContextMenuItem
          variant="destructive"
          onClick={(event) => {
            event.stopPropagation();
            onDelete();
          }}
        >
          <Trash2 className="size-4" />
          Excluir programa
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  );
}
