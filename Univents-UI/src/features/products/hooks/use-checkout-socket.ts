import { useCallback, useEffect, useRef, useState } from "react";
import { getWebsocketAuthToken } from "../api";
import type { SubmitPaymentPayloadI } from "@/features/payments/model";

type ServerMessage =
  | { type: "reservation_failed"; payload: { unavailable: { product_id: string; name: string; reason: string }[] } }
  | { type: "partial_reservation"; payload: { reserved: ReservedItem[]; unavailable: UnavailableItem[]; confirm_deadline: string } }
  | { type: "reservation_confirmed"; payload: { session_id: string; expires_at: string; items: ReservedItem[]; total: number } }
  | { type: "payment_processing" }
  | { type: "payment_failed"; payload: { reason: string } }
  | { type: "payment_pending"; payload: string }
  | { type: "order_confirmed"; payload: { order_id: string; receipt_url?: string } }
  | { type: "order_failed"; payload: { reason: string; order_id?: string } }
  | { type: "reservation_cancelled" }
  | { type: "error"; payload: string };

interface ReservedItem {
  product_id: string;
  name: string;
  quantity: number;
  price_cents: number;
}

interface UnavailableItem {
  product_id: string;
  name: string;
  reason: string;
}

export type CheckoutPhase =
  | "idle"
  | "connecting"
  | "awaiting_reservation"
  | "reservation_confirmed"
  | "partial_reservation"
  | "reservation_failed"
  | "awaiting_payment"
  | "payment_processing"
  | "payment_failed"
  | "payment_pending"
  | "order_confirmed"
  | "order_failed"
  | "error";

interface CheckoutState {
  phase: CheckoutPhase;
  errorMessage: string | null;
  sessionId: string | null;
  reservationExpiresAt: string | null;
  partialData: {
    reserved: ReservedItem[];
    unavailable: UnavailableItem[];
    confirmDeadline: string;
  } | null;
  reservedItems: ReservedItem[];
  total: number;
  orderId: string | null;
  pendingMessage: string | null;
}

const INITIAL: CheckoutState = {
  phase: "idle",
  errorMessage: null,
  sessionId: null,
  reservationExpiresAt: null,
  partialData: null,
  reservedItems: [],
  total: 0,
  orderId: null,
  pendingMessage: null,
};

export interface BuyRequestItem {
  product_id: string;
  quantity: number;
}

export function useCheckoutSocket(url: string) {
  const wsRef = useRef<WebSocket | null>(null);
  const intentionalClose = useRef(false);
  const hadServerError = useRef(false);
  const [state, setState] = useState<CheckoutState>(INITIAL);

  const patch = useCallback((updates: Partial<CheckoutState>) => {
    setState((prev) => ({ ...prev, ...updates }));
  }, []);

  const sendJSON = useCallback((payload: Record<string, unknown>) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(payload));
    }
  }, []);

  const handleMessage = useCallback(
    (raw: MessageEvent) => {
      let msg: ServerMessage;
      try {
        msg = JSON.parse(raw.data as string) as unknown as ServerMessage;
      } catch {
        return;
      }

      switch (msg.type) {
        case "reservation_confirmed":
          patch({
            phase: "reservation_confirmed",
            sessionId: msg.payload.session_id,
            reservationExpiresAt: msg.payload.expires_at,
            reservedItems: msg.payload.items,
            total: msg.payload.total,
          });
          break;

        case "partial_reservation":
          patch({
            phase: "partial_reservation",
            partialData: {
              reserved: msg.payload.reserved,
              unavailable: msg.payload.unavailable,
              confirmDeadline: msg.payload.confirm_deadline,
            },
          });
          break;

        case "reservation_failed":
          hadServerError.current = true;
          patch({
            phase: "reservation_failed",
            errorMessage: msg.payload.unavailable.map((i) => i.reason).join(", "),
          });
          break;

        case "reservation_cancelled":
          setState(INITIAL);
          break;

        case "payment_processing":
          patch({ phase: "payment_processing" });
          break;

        case "payment_failed":
          hadServerError.current = true;
          patch({ phase: "payment_failed", errorMessage: msg.payload.reason });
          break;

        case "payment_pending":
          patch({ phase: "payment_pending", pendingMessage: msg.payload });
          break;

        case "order_confirmed":
          patch({ phase: "order_confirmed", orderId: msg.payload.order_id });
          break;

        case "order_failed":
          hadServerError.current = true;
          patch({ phase: "order_failed", errorMessage: msg.payload.reason, orderId: msg.payload.order_id ?? null });
          break;

        case "error":
          hadServerError.current = true;
          patch({ phase: "error", errorMessage: msg.payload });
          break;
      }
    },
    [patch],
  );

  const buyRequest = useCallback(
    async (items: BuyRequestItem[]) => {
      if (wsRef.current) {
        intentionalClose.current = true;
        wsRef.current.close(1000, "restart");
        wsRef.current = null;
      }

      intentionalClose.current = false;
      hadServerError.current = false;
      setState({ ...INITIAL, phase: "connecting" });

      const res = await getWebsocketAuthToken();
      if (!res.success) {
        patch({ phase: "error", errorMessage: "Não foi possível autenticar. Tente novamente." });
        return;
      }

      const ws = new WebSocket(`${url}?token=${res.data.token}`);
      wsRef.current = ws;

      ws.onopen = () => {
        patch({ phase: "awaiting_reservation" });
        ws.send(JSON.stringify({ type: "buy_request", items }));
      };

      ws.onmessage = handleMessage;

      ws.onerror = () => {
        patch({ phase: "error", errorMessage: "Erro na conexão com o servidor." });
      };

      ws.onclose = () => {
        wsRef.current = null;
        if (!intentionalClose.current && !hadServerError.current) {
          patch({ phase: "error", errorMessage: "Conexão encerrada inesperadamente." });
        }
      };
    },
    [url, handleMessage, patch],
  );

  const confirmPartial = useCallback(() => {
    sendJSON({ type: "confirm_partial" });
    patch({ phase: "awaiting_reservation" });
  }, [sendJSON, patch]);

  const cancelReservation = useCallback(() => {
    sendJSON({ type: "cancel" });
    intentionalClose.current = true;
    wsRef.current?.close(1000, "cancelled");
    setState(INITIAL);
  }, [sendJSON]);

  const submitPayment = useCallback(
    (payload: SubmitPaymentPayloadI) => {
      if (wsRef.current?.readyState !== WebSocket.OPEN) {
        patch({ phase: "error", errorMessage: "Conexão perdida. Tente novamente." });
        return;
      }
      sendJSON({
        type: "submit_payment",
        payload: {
          card_token: payload.card_token,
          payment_method_id: payload.payment_method_id,
          payment_method_type: payload.payment_method_type,
          seller_credential_id: payload.seller_credential_id,
          payer_email: payload.payer_email,
          installments: payload.installments,
        }
      });
      patch({ phase: "awaiting_payment" });
    },
    [sendJSON, patch],
  );

  const reset = useCallback(() => {
    intentionalClose.current = true;
    wsRef.current?.close(1000, "reset");
    wsRef.current = null;
    setState(INITIAL);
  }, []);

  useEffect(() => {
    return () => {
      intentionalClose.current = true;
      wsRef.current?.close(1000, "unmount");
    };
  }, []);

  return { state, buyRequest, confirmPartial, cancelReservation, submitPayment, reset };
}