import { useNavigate } from "@tanstack/react-router";
import { Mail, Pencil } from "lucide-react";
import { motion } from "motion/react";
import type { EventI, EventStatusI } from "@/features/events/model";
import { cn } from "@/shared/lib/utils";

const statusConfig: Record<
  EventStatusI,
  {
    label: string;
    dot: string;
    pill: string;
  }
> = {
  draft: {
    label: "Rascunho",
    dot: "bg-amber-500",
    pill: "bg-amber-500/10 text-amber-700 border-amber-500/20",
  },
  active: {
    label: "Ativo",
    dot: "bg-emerald-500",
    pill: "bg-emerald-500/10 text-emerald-700 border-emerald-500/20",
  },
  discontinued: {
    label: "Descontinuado",
    dot: "bg-rose-500",
    pill: "bg-rose-500/10 text-rose-700 border-rose-500/20",
  },
};

interface AdminEventCardProps {
  event: EventI;
  index?: number;
  onEdit?: (event: EventI) => void;
}

export default function AdminEventCard({
  event,
  index = 0,
  onEdit,
}: AdminEventCardProps) {
  const navigate = useNavigate();
  const status = statusConfig[event.status];
  const hasVisual = Boolean(event.banner_url ?? event.logo_url);
  const handleEdit = () => onEdit?.(event);
  const handleOpenDashboard = () => {
    void navigate({
      to: "/admin/events/$eventId",
      params: { eventId: event.id },
    });
  };
  return (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{
        delay: index * 0.05,
        duration: 0.35,
        ease: [0.25, 0.1, 0.25, 1],
      }}
      className={cn(
        "group relative flex w-full min-w-60 max-w-full flex-col overflow-hidden rounded-2xl bg-card text-left",
        "ring-1 ring-foreground/10 shadow-xs",
        "transform-gpu will-change-transform",
        "transition-all duration-300 ease-out",
        "hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm",
        "focus:outline-none focus-visible:outline-none focus-visible:ring-0",
      )}
      role="button"
      tabIndex={0}
      onClick={handleOpenDashboard}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          handleOpenDashboard();
        }
      }}
    >
      <div className="relative aspect-video overflow-hidden bg-muted">
        {hasVisual ? (
          <img
            src={event.banner_url ?? event.logo_url ?? ""}
            alt="Representação Visual do Evento"
            className={cn(
              "h-full w-full object-cover transition-transform",
              "duration-700 ease-out group-hover:scale-105",
            )}
            loading={index < 4 ? "eager" : "lazy"}
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-linear-to-br from-muted via-background to-muted/40">
            <div className="flex size-20 items-center justify-center rounded-full border border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
              <span className="text-2xl font-semibold text-muted-foreground/40">
                {event.acronym ?? event.full_name.charAt(0)}
              </span>
            </div>
          </div>
        )}

        <div className="absolute left-4 top-4 flex flex-wrap items-center gap-2">
          <span
            className={cn(
              "inline-flex items-center gap-1 rounded-full border border-white/30 bg-black/75 px-2.5 py-1 text-[11px] font-semibold text-white shadow-lg backdrop-blur-sm drop-shadow-[0_1px_2px_rgba(0,0,0,0.9)]",
            )}
          >
            <span className={cn("size-1.5 rounded-full", status.dot)} />
            {status.label}
          </span>
        </div>

        <div className="absolute right-4 top-4 z-20">
          {onEdit && (
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                handleEdit();
              }}
              className="inline-flex items-center gap-1.5 rounded-lg bg-background/85 px-3 py-1.5 text-xs font-medium text-foreground shadow-sm backdrop-blur-sm transition-colors hover:bg-background"
            >
              <Pencil className="size-3.5" />
              Editar
            </button>
          )}
        </div>

        <div className="absolute inset-x-0 bottom-0 flex items-end justify-between gap-3 p-4 sm:p-5">
          <div className="min-w-0 space-y-1">
            <h3 className="line-clamp-2 text-balance text-lg font-semibold leading-snug text-white drop-shadow-[0_1px_3px_rgba(0,0,0,0.9)] transition-colors duration-300 group-hover:text-white sm:text-xl">
              {event.full_name}
            </h3>
            {event.description && (
              <p className="line-clamp-1 max-w-[min(65%,32rem)] truncate text-xs text-white/90 drop-shadow-[0_1px_2px_rgba(0,0,0,0.9)]">
                {event.description}
              </p>
            )}
          </div>
        </div>
      </div>

      <div className="flex items-center justify-between gap-3 p-4 pt-3 sm:p-5 sm:pt-4">
        <div className="min-w-0 flex-1 space-y-1">
          <div className="flex min-w-0 items-center text-xs text-muted-foreground">
            <span
              className="inline-flex min-w-0 max-w-full items-center gap-1.5"
              title={event.contact_email ?? "Sem contato cadastrado"}
            >
              <Mail className="size-3.5 shrink-0" />
              <span className="truncate">
                {event.contact_email ?? "Sem contato"}
              </span>
            </span>
          </div>
          <code className="block truncate text-[11px] font-mono text-muted-foreground/80">
            {event.slug}
          </code>
        </div>
      </div>
    </motion.article>
  );
}
