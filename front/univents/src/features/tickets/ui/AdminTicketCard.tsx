import {
  Coins,
  Hash,
  Infinity as InfinityIcon,
  MoreVertical,
  PencilLine,
  Ticket,
} from "lucide-react";
import { motion } from "motion/react";
import type React from "react";
import type { TicketI } from "@/features/tickets/model";
import { cn } from "@/shared/lib/utils";
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

interface AdminTicketCardProps {
  ticket: TicketI;
  index?: number;
  onManage: (ticket: TicketI) => void;
}

function MenuItems({
  isContext = false,
  onManage,
}: {
  isContext?: boolean;
  onManage: () => void;
}) {
  const Item = isContext ? ContextMenuItem : DropdownMenuItem;
  const stop =
    (action?: () => void) => (e: React.MouseEvent | React.KeyboardEvent) => {
      e.preventDefault();
      e.stopPropagation();
      action?.();
    };

  return (
    <Item onClick={stop(onManage)}>
      <PencilLine className="size-4" />
      <span>Editar ticket</span>
    </Item>
  );
}

export default function AdminTicketCard({
  ticket,
  index = 0,
  onManage,
}: AdminTicketCardProps) {
  const handleEdit = () => onManage(ticket);

  const Article = (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{
        delay: index * 0.05,
        duration: 0.35,
        ease: [0.25, 0.1, 0.25, 1],
      }}
      className={cn(
        "group relative flex w-full min-w-62.5 max-w-full flex-col overflow-hidden rounded-2xl bg-card text-left",
        "ring-1 ring-foreground/10 shadow-xs",
        "transform-gpu will-change-transform",
        "transition-all duration-300 ease-out",
        "hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm",
        "focus:outline-none focus-visible:outline-none focus-visible:ring-0",
      )}
      role="button"
      tabIndex={0}
      onClick={handleEdit}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          handleEdit();
        }
      }}
    >
      <div className="relative overflow-hidden bg-linear-to-br from-muted via-background to-muted/40 px-4 py-4">
        <div className="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-primary/70 via-primary to-cyan-500/70" />

        <div className="absolute right-3 top-3">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <button
                  type="button"
                  onClick={(e) => e.stopPropagation()}
                  className={cn(
                    "inline-flex size-9 items-center justify-center rounded-full",
                    "bg-background/85 text-foreground shadow-sm backdrop-blur-sm",
                    "transition-colors hover:bg-background",
                  )}
                  aria-label={`Abrir ações de ${ticket.name}`}
                >
                  <MoreVertical className="size-4" />
                </button>
              }
            />
            <DropdownMenuContent align="end" className="w-56">
              <MenuItems onManage={handleEdit} />
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div className="flex items-start justify-between gap-3 pr-10">
          <div className="min-w-0 space-y-2">
            <div className="flex flex-wrap items-center gap-1.5">
              <span className="inline-flex items-center gap-1 rounded-full border border-sky-500/20 bg-sky-500/10 px-2 py-0.5 text-[10px] font-medium text-sky-700 backdrop-blur-sm">
                <Ticket className="size-3" />
                <span className="max-w-28 truncate">Ticket</span>
              </span>

              <span className="inline-flex items-center gap-1 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-2 py-0.5 text-[10px] font-medium text-emerald-700 backdrop-blur-sm">
                <Coins className="size-3" />

                <span className="max-w-28 truncate">
                  {ticket.price_cents > 0
                    ? `R$ ${(ticket.price_cents / 100).toFixed(2)}`
                    : "Gratuito"}
                </span>
              </span>
            </div>

            <div className="min-w-0 space-y-1">
              <h3 className="line-clamp-2 text-balance text-lg font-semibold leading-snug text-foreground transition-colors duration-300 group-hover:text-primary sm:text-xl">
                {ticket.name}
              </h3>
              {ticket.description && (
                <p className="line-clamp-1 text-sm text-muted-foreground">
                  {ticket.description}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="flex items-center justify-between gap-3 p-4 pt-3 sm:p-5 sm:pt-4">
        <div className="min-w-0 flex-1 space-y-1.5">
          <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
            <span className="inline-flex items-center gap-1.5">
              <Hash className="size-3.5 shrink-0" />
              <span>Nível de acesso: {ticket.access_level}</span>
            </span>

            <span className="inline-flex items-center gap-1.5">
              {ticket.max_quantity ? (
                <>
                  <Ticket className="size-3.5 shrink-0" />
                  <span>Máx. {ticket.max_quantity} ingressos</span>
                </>
              ) : (
                <>
                  <InfinityIcon className="size-3.5 shrink-0" />
                  <span>Quantidade ilimitada</span>
                </>
              )}
            </span>
          </div>

          {ticket.created_at && (
            <p className="text-[11px] text-muted-foreground/70">
              Criado em{" "}
              {new Date(ticket.created_at).toLocaleDateString("pt-BR")}
            </p>
          )}
        </div>
      </div>
    </motion.article>
  );

  return (
    <ContextMenu>
      <ContextMenuTrigger render={Article} />
      <ContextMenuContent className="w-56">
        <MenuItems isContext onManage={handleEdit} />
      </ContextMenuContent>
    </ContextMenu>
  );
}
