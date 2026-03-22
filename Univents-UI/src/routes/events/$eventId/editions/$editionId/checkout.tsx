import { createFileRoute } from "@tanstack/react-router";
import { CheckCircle, Clock, Lock, Shield, XCircle } from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useAuth } from "@soramux/node-auth-sdk/react";
import { useQuery } from "@tanstack/react-query";
import type { BuyRequestItem } from "@/features/products/hooks/use-checkout-socket";
import type { SubmitPaymentPayloadI } from "@/features/payments/model";
import { useCart } from "@/features/products/hooks/use-cart";
import { PaymentProviderSelector } from "@/features/payments/ui/PaymentProviderSelector";
import { useCheckoutSocket } from "@/features/products/hooks/use-checkout-socket";
import { env } from "@/env";
import { editionQueryOptions } from "@/features/editions/api";

export const Route = createFileRoute(
  "/events/$eventId/editions/$editionId/checkout",
)({
  component: CheckoutPage,
});

function useCountdown(expiresAt: string | null) {
  const [secondsLeft, setSecondsLeft] = useState(0);

  useEffect(() => {
    if (!expiresAt) return;
    const tick = () => {
      setSecondsLeft(
        Math.max(0, Math.floor((new Date(expiresAt).getTime() - Date.now()) / 1000))
      );
    };
    tick();
    const id = setInterval(tick, 1000);
    return () => { clearInterval(id); };
  }, [expiresAt]);

  const formatted = useMemo(() => {
    const m = Math.floor(secondsLeft / 60).toString().padStart(2, "0");
    const s = (secondsLeft % 60).toString().padStart(2, "0");
    return `${m}:${s}`;
  }, [secondsLeft]);

  return { secondsLeft, formatted };
}

function OrderSummary({ editionId, expiresAt }: { editionId: string; expiresAt: string | null }) {
  const { items, totalCents } = useCart(editionId);
  const { secondsLeft, formatted } = useCountdown(expiresAt);

  const price = (cents: number) =>
    new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" }).format(cents / 100);

  return (
    <div className="w-full space-y-0">
      <div className="divide-y divide-border/50">
        {items.map((item) => (
          <div key={item.id} className="flex gap-3 py-3">
            <div className="flex-1 min-w-0">
              <h4 className="text-sm font-medium text-foreground truncate leading-tight">{item.name}</h4>
              <p className="text-xs text-muted-foreground mt-0.5">
                {price(item.price_cents)} un · {item.quantity}x
              </p>
            </div>
            <span className="text-sm font-semibold text-foreground tabular-nums self-center">
              {price(item.price_cents * item.quantity)}
            </span>
          </div>
        ))}
      </div>

      <div className="flex justify-between items-center border-t border-border pt-3 mt-3">
        <span className="text-sm font-semibold text-foreground uppercase tracking-wide">Total</span>
        <span className="text-2xl font-bold text-primary tabular-nums">{price(totalCents)}</span>
      </div>

      {expiresAt && (
        <div className={`flex items-center justify-between mt-4 py-2.5 px-3 border rounded-md ${secondsLeft < 60 ? "bg-destructive/10 border-destructive/20" : "bg-accent/10 border-accent/20"
          }`}>
          <div className="flex items-center gap-2 text-accent-foreground">
            <Clock className="w-4 h-4" />
            <span className="text-xs font-medium uppercase tracking-wide">Reserva expira em</span>
          </div>
          <span className="font-mono text-sm font-bold text-accent-foreground tabular-nums">
            {formatted}
          </span>
        </div>
      )}
    </div>
  );
}

function CheckoutPage() {
  const { eventId, editionId } = Route.useParams();
  const { totalCents, items } = useCart(editionId);
  const { auth } = useAuth();
  const userInfo = auth.profile();

  const { data: edition } = useQuery(editionQueryOptions(eventId, editionId));

  if (!userInfo) throw Route.redirect({ to: "/events/$eventId/editions/$editionId/products" });

  const wsUrl = `${env.VITE_API_URL.replace("http", "ws")}events/${eventId}/editions/${editionId}/products/purchase`;

  const { state, buyRequest, confirmPartial, cancelReservation, submitPayment, reset } =
    useCheckoutSocket(wsUrl);

  const { phase, errorMessage, partialData, orderId, pendingMessage, reservationExpiresAt } = state;

  const getCartItems = useCallback(
    (): BuyRequestItem[] => items.map((item) => ({ product_id: item.id, quantity: item.quantity })),
    [items],
  );

  useEffect(() => {
    if (edition && (!edition.trie_payments_credential_id || !edition.trie_payments_provider_public_key))
      return;
    if (items.length > 0 && phase === "idle" && edition?.trie_payments_credential_id) {
      void buyRequest(getCartItems());
    }
  }, [buyRequest, getCartItems, phase, items.length, edition]);

  const price = (cents: number) =>
    new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" }).format(cents / 100);

  const handleSubmitPayment = useCallback(
    (payload: SubmitPaymentPayloadI) => { submitPayment(payload); },
    [submitPayment],
  );

  const paymentUnavailable = edition && (
    !edition.trie_payments_credential_id || !edition.trie_payments_provider_public_key
  );

  return (
    <div className="min-h-screen bg-background py-4 px-4">
      <div className="max-w-4xl mx-auto">

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

          <div className="w-full">
            <div className="flex items-center gap-2 mb-3 pb-2 border-b border-border">
              <Shield className="w-4 h-4 text-primary" />
              <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                Resumo do Pedido
              </h2>
            </div>
            <OrderSummary editionId={editionId} expiresAt={reservationExpiresAt} />
          </div>

          <div className="w-full">
            <div className="flex items-center gap-2 mb-3 pb-2 border-b border-border">
              <Lock className="w-4 h-4 text-primary" />
              <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                Pagamento
              </h2>
            </div>

            {paymentUnavailable && (
              <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-4 text-yellow-800">
                <div className="flex items-center gap-2 mb-1">
                  <Clock className="h-4 w-4" />
                  <p className="text-sm font-bold uppercase tracking-wide">Pagamento indisponível</p>
                </div>
                <p className="text-xs leading-relaxed opacity-90">
                  Esta edição ainda não configurou uma conta para receber pagamentos.
                  Por favor, entre em contato com o organizador do evento.
                </p>
              </div>
            )}

            {!paymentUnavailable && (phase === "connecting" || phase === "awaiting_reservation") && (
              <p className="text-sm text-muted-foreground animate-pulse">Reservando itens…</p>
            )}

            {phase === "reservation_failed" && (
              <div className="space-y-3">
                <div className="flex items-start gap-2 text-destructive">
                  <XCircle className="w-4 h-4 mt-0.5 shrink-0" />
                  <p className="text-sm">{errorMessage ?? "Itens esgotados."}</p>
                </div>
                <button
                  onClick={() => void buyRequest(getCartItems())}
                  className="text-xs underline text-muted-foreground"
                >
                  Tentar novamente
                </button>
              </div>
            )}

            {/* ← name + reason */}
            {phase === "partial_reservation" && partialData && (
              <div className="space-y-4">
                <p className="text-sm font-medium">Alguns itens não estão disponíveis:</p>
                <ul className="space-y-1.5">
                  {partialData.unavailable.map((i) => (
                    <li key={i.product_id} className="text-xs text-muted-foreground">
                      <span className="font-medium text-foreground">{i.name}</span>: {i.reason}
                    </li>
                  ))}
                </ul>
                <div className="flex gap-2 pt-1">
                  <button
                    onClick={confirmPartial}
                    className="text-xs px-3 py-1.5 bg-primary text-primary-foreground rounded-md"
                  >
                    Continuar com disponíveis
                  </button>
                  <button
                    onClick={cancelReservation}
                    className="text-xs px-3 py-1.5 border border-border rounded-md text-muted-foreground"
                  >
                    Cancelar
                  </button>
                </div>
              </div>
            )}

            {(phase === "reservation_confirmed" || phase === "payment_failed") && (
              <div className="space-y-3">
                {phase === "payment_failed" && (
                  <div className="flex items-start gap-2 text-destructive">
                    <XCircle className="w-4 h-4 mt-0.5 shrink-0" />
                    <p className="text-sm">{errorMessage}</p>
                  </div>
                )}

                <PaymentProviderSelector
                  amount={totalCents}
                  handleSubmit={handleSubmitPayment}
                  seller_credential_id={edition?.trie_payments_credential_id ?? ""}
                  seller_public_key={edition?.trie_payments_provider_public_key ?? ""}
                />
              </div>
            )}

            {(phase === "awaiting_payment" || phase === "payment_processing") && (
              <p className="text-sm text-muted-foreground animate-pulse">
                {phase === "awaiting_payment" ? "Enviando pagamento…" : "Aguardando confirmação…"}
              </p>
            )}

            {/* ← pendingMessage*/}
            {phase === "payment_pending" && (
              <p className="text-sm text-muted-foreground">
                {pendingMessage ?? "Isso está demorando mais que o esperado."}
              </p>
            )}

            {phase === "order_confirmed" && (
              <div className="flex items-start gap-2 text-green-600">
                <CheckCircle className="w-4 h-4 mt-0.5 shrink-0" />
                <p className="text-sm font-medium">Pedido confirmado! #{orderId}</p>
              </div>
            )}

            {phase === "order_failed" && (
              <div className="space-y-2">
                <div className="flex items-start gap-2 text-destructive">
                  <XCircle className="w-4 h-4 mt-0.5 shrink-0" />
                  <p className="text-sm">{errorMessage}</p>
                </div>
                <button onClick={reset} className="text-xs underline text-muted-foreground">
                  Voltar
                </button>
              </div>
            )}

            {phase === "error" && (
              <div className="space-y-2">
                <div className="flex items-start gap-2 text-destructive">
                  <XCircle className="w-4 h-4 mt-0.5 shrink-0" />
                  <p className="text-sm">{errorMessage}</p>
                </div>
                <button
                  onClick={() => { void buyRequest(getCartItems()) }}
                  className="text-xs underline text-muted-foreground"
                >
                  Tentar novamente
                </button>
              </div>
            )}

          </div>
        </div>
      </div>
    </div>
  );
}