import { createFileRoute } from "@tanstack/react-router";
import { Clock, Lock, Shield } from "lucide-react";
import { useState, useEffect } from "react";
import { useCart } from "@/features/products/hooks/use-cart";
import { PaymentProviderSelector } from "@/features/payments/ui/PaymentProviderSelector";

export const Route = createFileRoute(
  "/events/$eventId/editions/$editionId/checkout",
)({
  component: CheckoutPage,
});

function OrderSummary({ editionId }: { editionId: string }) {
  const { items, totalCents } = useCart(editionId);
  const [timeLeft, setTimeLeft] = useState(15 * 60);

  useEffect(() => {
    if (timeLeft === 0) return;
    const timer = setInterval(() => { setTimeLeft((t) => t - 1); }, 1000);
    return () => { clearInterval(timer); };
  }, [timeLeft]);

  const formatTime = (s: number) =>
    `${Math.floor(s / 60).toString().padStart(2, "0")}:${(s % 60).toString().padStart(2, "0")}`;

  const price = (cents: number) =>
    new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" }).format(cents / 100);

  return (
    <div className="space-y-0 w-full">
      {/* Itens */}
      <div className="divide-y divide-border/50">
        {items.map((item) => (
          <div
            key={item.id}
            className="flex gap-3 py-3"
          >
            <div className="flex-1 min-w-0">
              <h4 className="font-medium text-sm text-foreground truncate leading-tight">
                {item.name}
              </h4>
              <p className="text-xs text-muted-foreground mt-0.5">
                {price(item.price_cents)} un · {item.quantity}x
              </p>
            </div>
            <div className="flex items-center">
              <span className="text-sm font-semibold text-foreground tabular-nums">
                {price(item.price_cents * item.quantity)}
              </span>
            </div>
          </div>
        ))}
      </div>

      {/* Total */}
      <div className="flex justify-between items-center border-t border-border pt-3 mt-3">
        <span className="text-sm font-semibold text-foreground uppercase tracking-wide">
          Total
        </span>
        <span className="text-2xl font-bold text-primary tabular-nums">
          {price(totalCents)}
        </span>
      </div>

      {/* Timer */}
      <div className="flex items-center justify-between mt-4 py-2.5 px-3 bg-accent/10 border border-accent/20">
        <div className="flex items-center gap-2 text-accent-foreground">
          <Clock className="w-4 h-4" />
          <span className="text-xs font-medium uppercase tracking-wide">Reserva expira em</span>
        </div>
        <span className="font-mono text-sm font-bold text-accent-foreground tabular-nums">
          {formatTime(timeLeft)}
        </span>
      </div>
    </div>
  );
}

function CheckoutPage() {
  const { editionId } = Route.useParams();
  const { totalCents, items } = useCart(editionId);

  const price = (cents: number) =>
    new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" }).format(cents / 100);

  return (
    <div className="min-h-screen bg-background py-4 px-4">
      <div className="max-w-4xl mx-auto">
        {/* Header Compacto */}
        <div className="flex items-center justify-between mb-6 pb-3 border-b border-border">
          <div>
            <h1 className="text-lg font-semibold text-foreground tracking-tight">Checkout</h1>
            <p className="text-xs text-muted-foreground">Finalize sua compra</p>
          </div>
          <div className="text-right">
            <p className="text-xs text-muted-foreground uppercase tracking-wide">{items.length} itens</p>
            <p className="text-sm font-bold text-primary tabular-nums">{price(totalCents)}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 lg:gap-8">
          {/* Order Summary */}
          <div className="w-full">
            <div className="flex items-center gap-2 mb-3 pb-2 border-b border-border">
              <Shield className="w-4 h-4 text-primary" />
              <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                Resumo do Pedido
              </h2>
            </div>
            <OrderSummary editionId={editionId} />
          </div>

          {/* Payment */}
          <div className="w-full">
            <div className="flex items-center gap-2 mb-3 pb-2 border-b border-border">
              <Lock className="w-4 h-4 text-primary" />
              <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">Pagamento</h2>
            </div>

            <PaymentProviderSelector amount={totalCents} />
          </div>
        </div>
      </div>
    </div>
  );
}