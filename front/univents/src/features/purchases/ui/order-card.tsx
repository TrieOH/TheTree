import type { LucideIcon } from "lucide-react";
import { CreditCard, FileText, Ticket, Wallet, Zap } from "lucide-react";
import { cn } from "@/shared/lib/utils";

export type OrderStatus = "approved" | "pending" | "refunded" | "cancelled";
export type Order = {
  purchaseId: string;
  status: OrderStatus;
  paymentIcon: LucideIcon;
  paymentText: string;
  date: string;
  total: string;
  items: Array<{
    icon?: LucideIcon;
    image?: string;
    title: string;
    description: string;
    price: string;
  }>;
};

const STATUS_CONFIG: Record<OrderStatus, { label: string; className: string }> =
  {
    approved: {
      label: "Aprovado",
      className: "bg-secondary text-secondary-foreground",
    },
    pending: {
      label: "Pendente",
      className: "bg-accent/20 text-accent-foreground",
    },
    cancelled: {
      label: "Cancelado",
      className: "bg-muted text-muted-foreground",
    },
    refunded: {
      label: "Reembolsado",
      className: "bg-violet-500/10 text-violet-700",
    },
  };

function LineItem({
  icon: Icon = Ticket,
  image,
  title,
  description,
  price,
}: Order["items"][number]) {
  return (
    <div className="flex items-center gap-3">
      {image ? (
        <img
          src={image}
          alt=""
          className="size-14 shrink-0 rounded-xl object-cover sm:size-16"
        />
      ) : (
        <div className="flex size-14 shrink-0 items-center justify-center rounded-xl bg-muted text-muted-foreground sm:size-16">
          <Icon className="size-5" />
        </div>
      )}
      <div className="min-w-0 flex-1">
        <p className="line-clamp-1 text-[14px] font-semibold leading-snug text-foreground sm:text-[15px]">
          {title}
        </p>
        <p className="mt-0.5 line-clamp-2 text-[12.5px] leading-snug text-muted-foreground sm:text-[13px]">
          {description}
        </p>
      </div>
      <p className="shrink-0 text-[13px] font-semibold tabular-nums text-foreground sm:text-[14px]">
        {price}
      </p>
    </div>
  );
}

export function OrderCard({
  order,
  className,
}: {
  order: Order;
  className?: string;
}) {
  const PaymentIcon = order.paymentIcon;
  const config = STATUS_CONFIG[order.status];
  return (
    <article
      className={cn(
        "flex h-full flex-col rounded-lg border border-border bg-popover p-4 shadow-[0_1px_2px_rgba(16,24,40,0.04)] transition-colors hover:border-primary/40 sm:p-5",
        className,
      )}
    >
      <header className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <span
            className={cn(
              "inline-block rounded-full px-3 py-1 text-[11px] font-bold uppercase tracking-wide",
              config.className,
            )}
          >
            {config.label}
          </span>
          <div className="mt-2 flex flex-wrap items-center gap-x-2 gap-y-1 text-[12.5px] text-muted-foreground">
            <span>{order.date}</span>
            <span aria-hidden="true">·</span>
            <span className="flex items-center gap-1.5">
              <PaymentIcon size={13} />
              {order.paymentText}
            </span>
          </div>
        </div>
        <div className="shrink-0 text-right">
          <p className="text-[11px] uppercase tracking-wide text-muted-foreground">
            Total
          </p>
          <p
            className={cn(
              "text-lg font-bold tabular-nums sm:text-xl",
              order.status === "cancelled" || order.status === "refunded"
                ? "text-muted-foreground line-through"
                : "text-primary",
            )}
          >
            {order.total}
          </p>
        </div>
      </header>
      <div className="mt-3.5 flex flex-col gap-3 border-t border-border pt-3.5 pb-3 sm:mt-4 sm:gap-3.5 sm:pt-4">
        {order.items.map((item) => (
          <LineItem key={`${item.title}-${item.price}`} {...item} />
        ))}
      </div>
      <footer className="mt-auto flex items-center justify-between border-t border-border pt-3.5 text-xs text-muted-foreground sm:pt-4">
        <span>
          {order.items.length} {order.items.length === 1 ? "item" : "itens"}
        </span>
        <span className="font-medium text-primary">
          Ver detalhes <span aria-hidden="true">→</span>
        </span>
      </footer>
    </article>
  );
}

export const orderIcons = { CreditCard, FileText, Ticket, Wallet, Zap };
