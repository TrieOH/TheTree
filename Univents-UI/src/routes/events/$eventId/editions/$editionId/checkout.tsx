import { createFileRoute } from "@tanstack/react-router";
import { CheckCircle, Clock, Copy, Lock, RefreshCw, Shield, XCircle, ShoppingBag } from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import type { SubmitPaymentPayloadI } from "@/features/payments/model";
import type { BuyRequestItem, ReservedItem } from "@/features/products/hooks/use-checkout-socket";
import { useCart } from "@/features/products/hooks/use-cart";
import { PaymentProviderSelector } from "@/features/payments/ui/PaymentProviderSelector";
import { useCheckoutSocket } from "@/features/products/hooks/use-checkout-socket";
import { env } from "@/env";
import { editionQueryOptions } from "@/features/editions/api";

export const Route = createFileRoute(
  "/events/$eventId/editions/$editionId/checkout",
)({
  beforeLoad: ({ context }) => {
    if (!context.auth?.isAuthenticated) {
      throw Route.redirect({ to: "/events/$eventId/editions/$editionId/products" });
    }
  },
  component: CheckoutPage,
});

const SESSION_EXPIRY_BUFFER_MS = 2 * 60 * 1000;

interface SavedSession {
  sessionId: string;
  expiresAt: string;
  items: BuyRequestItem[];
  reservedItems: ReservedItem[];
}

function saveSession(key: string, data: SavedSession) {
  sessionStorage.setItem(key, JSON.stringify(data));
}

function loadSession(key: string): SavedSession | null {
  try {
    const raw = sessionStorage.getItem(key);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as SavedSession;
    if (Date.now() >= new Date(parsed.expiresAt).getTime() - SESSION_EXPIRY_BUFFER_MS) return null;
    return parsed;
  } catch {
    return null;
  }
}

function itemsDiffer(a: BuyRequestItem[], b: BuyRequestItem[]): boolean {
  if (a.length !== b.length) return true;
  return a.some((ai) => {
    const bi = b.find((x) => x.product_id === ai.product_id);
    return bi?.quantity !== ai.quantity;
  });
}

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

function PixQRCode({ qrCode, qrCodeBase64 }: { qrCode: string; qrCodeBase64: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(qrCode);
    setCopied(true);
    setTimeout(() => { setCopied(false); }, 2000);
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-col items-center gap-3 p-4 rounded-lg border border-border bg-muted/30">
        <img
          src={`data:image/png;base64,${qrCodeBase64}`}
          alt="QR Code Pix"
          className="w-48 h-48 rounded-md"
        />
        <p className="text-xs text-muted-foreground text-center leading-relaxed">
          Escaneie o QR Code com o app do seu banco ou carteira digital
        </p>
      </div>

      <div className="space-y-1.5">
        <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          Ou copie o código Pix
        </p>
        <div className="flex gap-2">
          <code className="flex-1 text-[11px] font-mono bg-muted/50 border border-border rounded-md px-3 py-2 break-all leading-relaxed text-foreground">
            {qrCode}
          </code>
          <button
            onClick={handleCopy}
            className="shrink-0 flex items-center gap-1.5 text-xs px-3 py-2 rounded-md border border-border hover:bg-muted/50 transition-colors text-muted-foreground hover:text-foreground"
          >
            <Copy className="h-3.5 w-3.5" />
            {copied ? "Copiado!" : "Copiar"}
          </button>
        </div>
      </div>

      <p className="text-xs text-muted-foreground text-center">
        Aguardando confirmação do pagamento…
      </p>
    </div>
  );
}

function OrderSummary({
  editionId,
  expiresAt,
  expiresLabel = "Reserva expira em",
}: {
  editionId: string;
  expiresAt: string | null,
  expiresLabel?: string;
}) {
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
            <span className="text-xs font-medium uppercase tracking-wide">
              {expiresLabel}
            </span>
          </div>
          <span className="font-mono text-sm font-bold text-accent-foreground tabular-nums">
            {formatted}
          </span>
        </div>
      )}
    </div>
  );
}

function SessionConflictPrompt({
  savedSession,
  currentTotal,
  onResume,
  onStartNew,
}: {
  savedSession: SavedSession;
  currentTotal: number;
  onResume: () => void;
  onStartNew: () => void;
}) {
  const price = (cents: number) =>
    new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" }).format(cents / 100);

  const savedTotal = savedSession.reservedItems.reduce(
    (sum, i) => sum + i.price_cents * i.quantity, 0
  );

  return (
    <div className="rounded-lg border border-border bg-background p-4 space-y-4">
      <div className="flex items-start gap-2">
        <RefreshCw className="w-4 h-4 text-primary mt-0.5 shrink-0" />
        <p className="text-sm font-medium text-foreground">Você tem uma compra em andamento</p>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div className="rounded-md border border-border p-3 space-y-2">
          <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Reserva anterior</p>
          <ul className="space-y-1">
            {savedSession.reservedItems.map((i) => (
              <li key={i.product_id} className="text-xs text-foreground">
                {i.name} <span className="text-muted-foreground">× {i.quantity}</span>
              </li>
            ))}
          </ul>
          <p className="text-sm font-bold text-primary tabular-nums">{price(savedTotal)}</p>
        </div>

        <div className="rounded-md border border-border p-3 space-y-2">
          <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Carrinho atual</p>
          <ul className="space-y-1">
            {savedSession.items.map((i) => {
              const detail = savedSession.reservedItems.find((r) => r.product_id === i.product_id);
              return (
                <li key={i.product_id} className="text-xs text-foreground">
                  {detail?.name ?? i.product_id}{" "}
                  <span className="text-muted-foreground">× {i.quantity}</span>
                </li>
              );
            })}
          </ul>
          <p className="text-sm font-bold text-primary tabular-nums">{price(currentTotal)}</p>
        </div>
      </div>

      <div className="flex gap-2">
        <button
          onClick={onResume}
          className="text-xs px-3 py-1.5 bg-primary text-primary-foreground rounded-md"
        >
          Continuar reserva anterior
        </button>
        <button
          onClick={onStartNew}
          className="text-xs px-3 py-1.5 border border-border rounded-md text-muted-foreground"
        >
          Iniciar nova compra
        </button>
      </div>
    </div>
  );
}

function CheckoutPage() {
  const { eventId, editionId } = Route.useParams();
  const navigate = Route.useNavigate();
  const { items, totalCents, clearCart, replaceCart } = useCart(editionId);
  const { data: edition } = useQuery(editionQueryOptions(eventId, editionId));

  const wsUrl = `${env.VITE_API_URL.replace("http", "ws")}events/${eventId}/editions/${editionId}/products/purchase`;
  const sessionKey = `checkout_session_${editionId}`;

  const [pendingSession, setPendingSession] = useState<SavedSession | null>(null);

  const handlePartialReservation = useCallback(
    (reserved: ReservedItem[]) => {
      replaceCart(reserved.map((r) => ({
        id: r.product_id,
        name: r.name,
        price_cents: r.price_cents,
        quantity: r.quantity,
        has_inventory: false,
        inventory_remaining: 0,
      })));
    },
    [replaceCart],
  );

  const handlePixCreated = useCallback(() => {
    sessionStorage.removeItem(sessionKey);
    setPendingSession(null);
  }, [sessionKey]);

  const { state, buyRequest, resumeSession, confirmPartial, cancelReservation, submitPayment, reset } =
    useCheckoutSocket(wsUrl, handlePartialReservation, handlePixCreated);

  const {
    phase,
    errorMessage,
    partialData,
    reservationExpiresAt,
    sessionId,
    totalCents: socketTotal,
    reservedItems,
  } = state;

  const reservedItemsRef = useRef<ReservedItem[]>([]);
  useEffect(() => {
    reservedItemsRef.current = reservedItems;
  }, [reservedItems]);

  useEffect(() => {
    if (sessionId && reservationExpiresAt && reservedItems.length > 0) {
      saveSession(sessionKey, {
        sessionId,
        expiresAt: reservationExpiresAt,
        items: items.map((i) => ({ product_id: i.id, quantity: i.quantity })),
        reservedItems,
      });
    }
  }, [sessionId, reservationExpiresAt, reservedItems, items, sessionKey]);

  useEffect(() => {
    if (phase === "order_confirmed" || phase === "session_expired" || phase === "order_failed") {
      sessionStorage.removeItem(sessionKey);
      setPendingSession(null);
    }
    if (phase === "session_expired") void buyRequest(getCartItems());
    if (phase === "order_confirmed") clearCart();
  }, [phase, sessionKey, clearCart]); // eslint-disable-line react-hooks/exhaustive-deps

  const getCartItems = useCallback(
    (): BuyRequestItem[] => items.map((item) => ({ product_id: item.id, quantity: item.quantity })),
    [items],
  );

  const paymentReady = !!edition?.trie_payments_credential_id && !!edition?.trie_payments_provider_public_key;

  useEffect(() => {
    if (!paymentReady || phase !== "idle" || items.length === 0) return;

    const saved = loadSession(sessionKey);

    if (!saved) {
      sessionStorage.removeItem(sessionKey);
      void buyRequest(getCartItems());
      return;
    }

    if (itemsDiffer(getCartItems(), saved.items)) {
      setPendingSession(saved);
    } else {
      void resumeSession(saved.sessionId);
    }
  }, [buyRequest, resumeSession, getCartItems, phase, items.length, paymentReady, sessionKey]);

  const handleResume = useCallback(() => {
    if (!pendingSession) return;
    replaceCart(pendingSession.reservedItems.map((r) => ({
      id: r.product_id,
      name: r.name,
      price_cents: r.price_cents,
      quantity: r.quantity,
      has_inventory: false,
      inventory_remaining: 0,
    })));
    saveSession(sessionKey, {
      ...pendingSession,
      items: pendingSession.reservedItems.map((r) => ({ product_id: r.product_id, quantity: r.quantity })),
    });
    setPendingSession(null);
    void resumeSession(pendingSession.sessionId);
  }, [pendingSession, resumeSession, replaceCart, sessionKey]);

  const handleStartNew = useCallback(() => {
    sessionStorage.removeItem(sessionKey);
    setPendingSession(null);
    void buyRequest(getCartItems());
  }, [buyRequest, getCartItems, sessionKey]);

  useEffect(() => {
    const activePhases: typeof phase[] = [
      "reservation_confirmed",
      "awaiting_payment",
      "payment_processing",
      "payment_failed",
    ];
    if (!activePhases.includes(phase) || reservedItemsRef.current.length === 0) return;

    const current = getCartItems();
    const reservedAsItems = reservedItemsRef.current.map((r) => ({ product_id: r.product_id, quantity: r.quantity }));

    if (itemsDiffer(current, reservedAsItems)) {
      sessionStorage.removeItem(sessionKey);
      void buyRequest(current);
    }
  }, [items]); // eslint-disable-line react-hooks/exhaustive-deps

  const price = (cents: number) =>
    new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" }).format(cents / 100);

  const handleSubmitPayment = useCallback(
    (payload: SubmitPaymentPayloadI) => { submitPayment(payload); },
    [submitPayment],
  );

  const handleCancelReservation = useCallback(() => {
    cancelReservation();
    void navigate({ to: "/events/$eventId/editions/$editionId/products", params: { eventId, editionId } });
  }, [cancelReservation, navigate, eventId, editionId]);

  const handleReset = useCallback(() => {
    reset();
    void navigate({ to: "/events/$eventId/editions/$editionId/products", params: { eventId, editionId } });
  }, [reset, navigate, eventId, editionId]);

  const displayTotal = socketTotal || totalCents;

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
            <p className="text-sm font-bold text-primary tabular-nums">{price(displayTotal)}</p>
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
            <OrderSummary
              editionId={editionId}
              expiresAt={phase === "partial_reservation" && partialData
                ? partialData.confirmDeadline
                : reservationExpiresAt
              }
              expiresLabel={phase === "partial_reservation"
                ? "Confirmar em"
                : "Reserva expira em"
              }
            />
          </div>

          <div className="w-full">
            <div className="flex items-center gap-2 mb-3 pb-2 border-b border-border">
              <Lock className="w-4 h-4 text-primary" />
              <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                Pagamento
              </h2>
            </div>

            {pendingSession && (
              <SessionConflictPrompt
                savedSession={pendingSession}
                currentTotal={totalCents}
                onResume={handleResume}
                onStartNew={handleStartNew}
              />
            )}

            {!pendingSession && !paymentReady && (
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

            {!pendingSession && paymentReady && (
              phase === "connecting" ||
              phase === "awaiting_reservation" ||
              phase === "session_expired"
            ) && (
                <p className="text-sm text-muted-foreground animate-pulse">Reservando itens…</p>
              )}

            {!pendingSession && phase === "reservation_failed" && (
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

            {!pendingSession && phase === "partial_reservation" && partialData && (
              <div className="space-y-4">
                <p className="text-sm font-medium">Seu carrinho foi ajustado</p>
                <ul className="space-y-1">
                  {partialData.unavailable.map((i) => (
                    <li key={i.product_id} className="flex items-start gap-1.5 text-xs text-muted-foreground">
                      <XCircle className="w-3 h-3 text-destructive shrink-0 mt-0.5" />
                      <span>
                        <span className="font-medium text-foreground">{i.name}</span>
                        {" — "}
                        pedido {i.requested}, reservado {i.reserved}
                        {i.reason === "insufficient_inventory" ? " (estoque insuficiente)" : ` (${i.reason})`}
                      </span>
                    </li>
                  ))}
                </ul>
                <div className="border-t border-border pt-3">
                  <p className="text-xs text-muted-foreground mb-0.5">Novo total</p>
                  <p className="text-xl font-bold text-primary tabular-nums">{price(socketTotal)}</p>
                </div>
                <div className="flex gap-2 pt-1">
                  <button
                    onClick={confirmPartial}
                    className="text-xs px-3 py-1.5 bg-primary text-primary-foreground rounded-md"
                  >
                    Continuar
                  </button>
                  <button
                    onClick={handleCancelReservation}
                    className="text-xs px-3 py-1.5 border border-border rounded-md text-muted-foreground"
                  >
                    Voltar à loja
                  </button>
                </div>
              </div>
            )}

            {!pendingSession && (phase === "reservation_confirmed" || phase === "payment_failed") && (
              <div className="space-y-3">
                {phase === "payment_failed" && (
                  <div className="flex items-start gap-2 text-destructive">
                    <XCircle className="w-4 h-4 mt-0.5 shrink-0" />
                    <p className="text-sm">{errorMessage}</p>
                  </div>
                )}
                <PaymentProviderSelector
                  amount={socketTotal}
                  handleSubmit={handleSubmitPayment}
                  seller_public_key={edition?.trie_payments_provider_public_key ?? ""}
                />
              </div>
            )}

            {!pendingSession && (phase === "awaiting_payment" || phase === "payment_processing") && (
              <p className="text-sm text-muted-foreground animate-pulse">
                {phase === "awaiting_payment" ? "Enviando pagamento…" : "Aguardando confirmação…"}
              </p>
            )}

            {!pendingSession && phase === "payment_pending" && (
              state.pixData ? (
                <PixQRCode qrCode={state.pixData.qrCode} qrCodeBase64={state.pixData.qrCodeBase64} />
              ) : (
                <div className="space-y-4">
                  <div className="flex items-start gap-2 text-foreground">
                    <Clock className="w-4 h-4 mt-0.5 shrink-0 text-muted-foreground" />
                    <div className="space-y-1">
                      <p className="text-sm font-medium">Pagamento em processamento</p>
                      <p className="text-xs text-muted-foreground leading-relaxed">
                        Seu pagamento está sendo processado em segundo plano.
                        Você pode fechar esta página — avisaremos quando for confirmado.
                      </p>
                    </div>
                  </div>
                  <button
                    onClick={() => void navigate({ to: "/" })} // FIXME: Go to purchases
                    className="flex items-center gap-1.5 text-xs px-3 py-1.5 bg-primary text-primary-foreground rounded-md"
                  >
                    <ShoppingBag className="w-3.5 h-3.5" />
                    Ver meus pedidos
                  </button>
                </div>
              )
            )}

            {!pendingSession && phase === "order_confirmed" && (
              <div className="flex items-start gap-2 text-green-600">
                <CheckCircle className="w-4 h-4 mt-0.5 shrink-0" />
                <p className="text-sm font-medium">Pedido confirmado!</p>
              </div>
            )}

            {!pendingSession && phase === "order_failed" && (
              <div className="space-y-2">
                <div className="flex items-start gap-2 text-destructive">
                  <XCircle className="w-4 h-4 mt-0.5 shrink-0" />
                  <p className="text-sm">{errorMessage}</p>
                </div>
                <button onClick={handleReset} className="text-xs underline text-muted-foreground">
                  Voltar
                </button>
              </div>
            )}

            {!pendingSession && phase === "error" && (
              // errorMessage já foi normalizado pelo hook (i/o timeout → mensagem legível).
              // Dois caminhos: tentar novamente (nova reserva) ou desistir.
              <div className="space-y-3">
                <div className="flex items-start gap-2 text-destructive">
                  <XCircle className="w-4 h-4 mt-0.5 shrink-0" />
                  <p className="text-sm">{errorMessage}</p>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => void buyRequest(getCartItems())}
                    className="text-xs px-3 py-1.5 bg-primary text-primary-foreground rounded-md"
                  >
                    Tentar novamente
                  </button>
                  <button
                    onClick={handleReset}
                    className="text-xs px-3 py-1.5 border border-border rounded-md text-muted-foreground"
                  >
                    Desistir
                  </button>
                </div>
              </div>
            )}

          </div>
        </div>
      </div>
    </div>
  );
}