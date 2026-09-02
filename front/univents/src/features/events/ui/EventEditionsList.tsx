import { Link } from "@tanstack/react-router";
import { format } from "date-fns";
import { ptBR } from "date-fns/locale";
import { ChevronRight, Layers3 } from "lucide-react";
import type { EditionI } from "@/features/editions/model";

const statusConfig = {
  draft: {
    label: "Rascunho",
    className: "border-amber-500/20 bg-amber-500/10 text-amber-700",
    dot: "bg-amber-500",
  },
  future: {
    label: "Futura",
    className: "border-sky-500/20 bg-sky-500/10 text-sky-700",
    dot: "bg-sky-500",
  },
  active: {
    label: "Ativa",
    className: "border-emerald-500/20 bg-emerald-500/10 text-emerald-700",
    dot: "bg-emerald-500",
  },
  past: {
    label: "Encerrada",
    className: "border-slate-500/20 bg-slate-500/10 text-slate-600",
    dot: "bg-slate-500",
  },
} as const;

export function EventEditionsList({
  eventId,
  editions,
}: {
  eventId: string;
  editions: EditionI[];
}) {
  return (
    <section className="order-7 space-y-3">
      <div className="flex min-w-0 items-center justify-between gap-3 px-1">
        <div className="flex min-w-0 flex-1 items-center gap-3">
          <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
            <Layers3 className="size-5" />
          </div>
          <div className="min-w-0 flex-1">
            <h2 className="text-base font-semibold tracking-tight">Edições</h2>
            <p className="truncate text-xs text-muted-foreground">
              Acesse rapidamente cada edição do evento.
            </p>
          </div>
        </div>
        <span className="flex size-8 shrink-0 items-center justify-center rounded-full border border-border/60 bg-muted/50 p-0 text-xs font-medium text-muted-foreground">
          {editions.length}
        </span>
      </div>
      <div className="overflow-hidden rounded-lg border border-border/60 bg-card/95 shadow-sm">
        {editions.length === 0 ? (
          <p className="p-5 text-sm text-muted-foreground">
            Nenhuma edição cadastrada.
          </p>
        ) : (
          editions.slice(0, 6).map((edition, index) => {
            const status = statusConfig[edition.status];
            return (
              <Link
                key={edition.id}
                to="/admin/events/$eventId/editions/$editionId"
                params={{ eventId, editionId: edition.id }}
                className={`group flex items-center justify-between gap-4 px-5 py-4 transition-colors hover:bg-muted/30 ${index > 0 ? "border-t border-border/60" : ""}`}
              >
                <div className="flex min-w-0 items-center gap-3">
                  <span className="flex size-8 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">
                    {index + 1}
                  </span>
                  <div className="min-w-0">
                    <p className="truncate text-sm font-semibold">
                      {edition.name}
                    </p>
                    <p className="mt-1 truncate text-xs text-muted-foreground">
                      {format(new Date(edition.starts_at), "dd MMM yyyy", {
                        locale: ptBR,
                      })}
                      {edition.location_name
                        ? ` · ${edition.location_name}`
                        : ""}
                    </p>
                  </div>
                </div>
                <div className="flex shrink-0 items-center gap-3">
                  <span
                    className={`hidden items-center gap-1.5 rounded-full border px-2 py-0.5 text-[11px] sm:inline-flex ${status.className}`}
                  >
                    <span className={`size-1.5 rounded-full ${status.dot}`} />
                    {status.label}
                  </span>
                  <ChevronRight className="size-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
                </div>
              </Link>
            );
          })
        )}
      </div>
    </section>
  );
}
