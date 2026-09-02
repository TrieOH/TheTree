import { Calendar } from "lucide-react";
import { Badge } from "@/shared/ui/shadcn/badge";
import type { EventI } from "../model";
import { EventVisualCard } from "./EventVisualCard";

const statusConfig = {
  draft: {
    label: "Rascunho",
    className: "bg-amber-500/10 text-amber-700",
  },
  active: {
    label: "Ativo",
    className: "bg-emerald-500/10 text-emerald-700",
  },
  archived: {
    label: "Arquivado",
    className: "bg-slate-500/10 text-slate-700",
  },
  discontinued: {
    label: "Descontinuado",
    className: "bg-rose-500/10 text-rose-700",
  },
} as const;

export function EventOverviewHeader({ event }: { event: EventI | null }) {
  if (!event) return null;
  const status = statusConfig[event.status];
  const createdDate = new Date(event.created_at)
    .toLocaleDateString("pt-BR", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    })
    .replace(".", "");

  return (
    <>
      <EventVisualCard event={event} />
      <div className="space-y-1 px-1 text-center md:text-left">
        <h1 className="text-xl font-medium tracking-tight text-foreground/90">
          {event.full_name}
        </h1>
        <p className="flex items-center justify-center gap-1.5 text-xs text-muted-foreground md:justify-start">
          <Calendar className="size-3.5" />
          Criado em {createdDate}
        </p>
        <Badge
          className={`w-fit border-0 px-2 py-0.5 text-xs font-normal ${status.className}`}
        >
          {status.label}
        </Badge>
      </div>
    </>
  );
}
