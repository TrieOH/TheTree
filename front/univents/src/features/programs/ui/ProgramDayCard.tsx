import {
  Activity,
  Check,
  ChevronRight,
  LogIn,
  ShoppingCart,
} from "lucide-react";
import { useCart } from "@/features/products/hooks/use-cart";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import type { OccurrenceI, ProgramI } from "../model";

function formatTimeRange(start: string, end: string): string {
  const fmt = (iso: string) =>
    new Date(iso).toLocaleTimeString("pt-BR", {
      hour: "2-digit",
      minute: "2-digit",
    });
  return `${fmt(start)} – ${fmt(end)}`;
}

function formatDayLabel(dateIso: string): string {
  const date = new Date(dateIso);
  const today = new Date();
  const tomorrow = new Date(today);
  tomorrow.setDate(tomorrow.getDate() + 1);

  if (date.toDateString() === today.toDateString()) return "Hoje";
  if (date.toDateString() === tomorrow.toDateString()) return "Amanhã";

  return date.toLocaleDateString("pt-BR", {
    weekday: "long",
    day: "numeric",
    month: "long",
  });
}

function getIconForProgram(program: ProgramI) {
  if (program.kind === "checkpoint")
    return { Icon: LogIn, bg: "bg-primary", text: "text-primary-foreground" };
  return { Icon: Activity, bg: "bg-muted", text: "text-muted-foreground" };
}

interface ProgramDayCardProps {
  date: string;
  items: { program: ProgramI; occurrence: OccurrenceI }[];
  maxItems?: number;
  editionId?: string;
}

export function ProgramDayCard({
  date,
  items,
  maxItems = 3,
  editionId,
}: ProgramDayCardProps) {
  const { addItem, items: cartItems } = useCart(editionId ?? "");
  const visible = items.slice(0, maxItems);
  const remaining = items.length - maxItems;

  return (
    <div className="flex flex-col w-80 rounded-2xl p-6">
      {/* Day header */}
      <h3 className="w-full rounded-xl bg-primary/10 px-4 py-2 text-center text-sm font-bold tracking-wide text-primary capitalize mb-6">
        {formatDayLabel(date)}
      </h3>

      {/* Timeline */}
      <div className="flex flex-col">
        {visible.map(({ program, occurrence }, index) => {
          const isLast = index === visible.length - 1 && remaining <= 0;
          const { Icon, bg, text } = getIconForProgram(program);

          return (
            <div key={occurrence.id} className="flex gap-4">
              {/* Timeline column */}
              <div className="flex flex-col items-center shrink-0 relative">
                <div
                  className={cn(
                    "w-9 h-9 rounded-xl z-10 flex items-center justify-center shrink-0",
                    bg,
                  )}
                >
                  <Icon className={cn("w-4 h-4", text)} />
                </div>
                <div
                  className={cn(
                    "w-px rounded-full bg-border",
                    isLast
                      ? "absolute h-[calc(100%+0.5rem)]"
                      : "min-h-6 flex-1",
                  )}
                />
              </div>

              {/* Content */}
              <div className={cn("flex-1", isLast ? "pb-0" : "pb-6")}>
                {/* Time */}
                <span className="text-xs font-semibold text-primary">
                  {formatTimeRange(occurrence.starts_at, occurrence.ends_at)}
                </span>

                {/* Title */}
                <h4 className="mt-1 text-[15px] font-bold text-card-foreground leading-snug">
                  {program.name}
                </h4>

                {/* Description */}
                {program.description && (
                  <p className="mt-1 text-xs text-muted-foreground line-clamp-2 leading-relaxed">
                    {program.description}
                  </p>
                )}
                {editionId &&
                  (() => {
                    const inCart = cartItems.find(
                      (item) =>
                        item.id === occurrence.id && item.type === "activity",
                    );
                    return (
                      <Button
                        size="sm"
                        variant={inCart ? "secondary" : "default"}
                        className="mt-4 h-9 w-full gap-2 text-xs font-semibold shadow-sm"
                        onClick={() =>
                          addItem(
                            {
                              id: occurrence.id,
                              type: "activity",
                              name: program.name,
                              price_cents: 0,
                              inventory_remaining: 999,
                              has_inventory: false,
                            },
                            1,
                          )
                        }
                      >
                        {inCart ? (
                          <Check className="h-3.5 w-3.5" />
                        ) : (
                          <ShoppingCart className="h-3.5 w-3.5" />
                        )}
                        {inCart
                          ? `Adicionado (${inCart.quantity})`
                          : "Adicionar ao carrinho"}
                      </Button>
                    );
                  })()}
              </div>
            </div>
          );
        })}

        {/* "+ N" */}
        {remaining > 0 && (
          <div className="flex items-center gap-2 -mt-2 ml-13 text-xs font-medium text-muted-foreground">
            <span>
              +{remaining} ocorrência{remaining > 1 ? "s" : ""}
            </span>
            <ChevronRight className="w-3 h-3" />
          </div>
        )}
      </div>
    </div>
  );
}
