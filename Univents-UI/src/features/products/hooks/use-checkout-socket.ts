import { useCallback, useEffect, useRef, useState } from "react";
import { getWebsocketAuthToken } from "../api";
import type { SubmitPaymentPayloadI } from "@/features/payments/model";
import type { ProductType } from "../model";

type ServerMessage =
  | { type: "reservation_failed"; payload: { unavailable: UnavailableItem[] } }
  | { type: "partial_reservation"; payload: { reserved: ReservedItem[]; unavailable: UnavailableItem[]; confirm_deadline: string } }
  | { type: "reservation_confirmed"; payload: { session_id: string; expires_at: string; reserved_items: ReservedItem[]; total_cents: number } }
  | { type: "payment_processing" }
  | { type: "pix_created"; payload: { qr_code: string; qr_code_base64: string } }
  | { type: "payment_failed"; payload: { payment_intent_id: string } }
  | { type: "purchase_failed"; payload: { reason: string; product_ids: string[] } }
  | { type: "purchase_failed"; payload: { invalid_products: UnavailableItem[] } }
  | { type: "purchase_failed"; payload: string }
  | { type: "payment_pending"; payload: string }
  | { type: "order_confirmed"; payload: { purchase_id: string } | { payment_intent_id: string } }
  | { type: "reservation_cancelled" }
  | { type: "session_expired" }
  | { type: "error"; payload: string };

export interface ReservedItem {
  product_id: string;
  name: string;
  quantity: number;
  price_cents: number;
  product_type: ProductType;
  ticket_id?: string;
}

export interface UnavailableItem {
  product_id: string;
  name: string;
  reason: string;
  requested: number;
  reserved: number;
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
  | "payment_pending"   // cartão: processando em background, socket fechado
  | "pix_pending"       // pix: aguardando scan do QR code
  | "order_confirmed"
  | "session_expired"
  | "error";

interface CheckoutState {
  phase: CheckoutPhase;
  errorMessage: string | null;
  payment_intent_id: string | null;
  sessionId: string | null;
  reservationExpiresAt: string | null;
  partialData: {
    reserved: ReservedItem[];
    unavailable: UnavailableItem[];
    confirmDeadline: string;
  } | null;
  reservedItems: ReservedItem[];
  pixData: { qrCode: string; qrCodeBase64: string } | null;
  totalCents: number;
  pendingMessage: string | null;
}

const INITIAL: CheckoutState = {
  phase: "idle",
  errorMessage: null,
  payment_intent_id: null,
  sessionId: null,
  reservationExpiresAt: null,
  partialData: null,
  reservedItems: [],
  totalCents: 0,
  pendingMessage: null,
  pixData: null,
};

export interface BuyRequestItem {
  product_id: string;
  quantity: number;
}

// Normaliza mensagens de erro cruas do Go em strings legíveis para o usuário.
function normalizeErrorMessage(raw: string): string {
  if (raw.includes("i/o timeout")) return "Tempo esgotado para confirmar a reserva.";
  if (raw.includes("connection reset")) return "Conexão encerrada inesperadamente.";
  if (raw.includes("broken pipe")) return "Conexão encerrada inesperadamente.";
  return "Ocorreu um erro inesperado.";
}

export function useCheckoutSocket(
  url: string,
  onPartialReservation?: (reserved: ReservedItem[]) => void,
  // Chamado quando pix_created chega: o backend já deletou a sessão nesse
  // momento, então o CheckoutPage precisa limpar o sessionStorage.
  onPixCreated?: () => void,
) {
  const wsRef = useRef<WebSocket | null>(null);
  const connectingRef = useRef(false);
  const intentionalClose = useRef(false);
  const hadServerError = useRef(false);
  // Fases terminais onde o backend fecha o socket após a última mensagem.
  // Sem essa ref o onclose interpretaria o fechamento como inesperado.
  const terminalPhaseRef = useRef(false);
  const handleMessageRef = useRef<((raw: MessageEvent) => void) | null>(null);
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
        msg = JSON.parse(raw.data as string) as ServerMessage;
      } catch {
        return;
      }

      switch (msg.type) {
        case "reservation_confirmed":
          patch({
            phase: "reservation_confirmed",
            sessionId: msg.payload.session_id,
            reservationExpiresAt: msg.payload.expires_at,
            reservedItems: msg.payload.reserved_items,
            totalCents: msg.payload.total_cents,
          });
          break;

        case "partial_reservation":
          patch({
            phase: "partial_reservation",
            reservedItems: msg.payload.reserved,
            totalCents: msg.payload.reserved.reduce((sum, i) => sum + i.price_cents * i.quantity, 0),
            partialData: {
              reserved: msg.payload.reserved,
              unavailable: msg.payload.unavailable,
              confirmDeadline: msg.payload.confirm_deadline,
            },
          });
          onPartialReservation?.(msg.payload.reserved);
          break;

        case "reservation_failed":
          hadServerError.current = true;
          patch({
            phase: "reservation_failed",
            errorMessage: "Não foi possível reservar os itens selecionados.",
          });
          break;

        case "reservation_cancelled":
          setState(INITIAL);
          break;

        case "session_expired":
          hadServerError.current = true;
          patch({
            phase: "session_expired",
            errorMessage: "Sua reserva expirou.",
          });
          break;

        case "payment_processing":
          patch({
            phase: "payment_processing",
            errorMessage: null,
          });
          break;

        case "payment_failed":
          hadServerError.current = true;
          patch({
            phase: "payment_failed",
            errorMessage: "Pagamento recusado. Tente novamente.",
            payment_intent_id: msg.payload.payment_intent_id,
          });
          break;

        case "purchase_failed":
          hadServerError.current = true;
          patch({
            phase: "error",
            errorMessage: "Não foi possível processar a compra. Tente novamente.",
          });
          break;

        case "payment_pending":
          // Terminal: socket fechado, pagamento em background via webhook.
          terminalPhaseRef.current = true;
          patch({
            phase: "payment_pending",
            pendingMessage: msg.payload,
            errorMessage: null,
          });
          break;

        case "pix_created":
          // Terminal do lado da sessão: backend já deletou a sessão antes de
          // emitir. onPixCreated limpa o sessionStorage no CheckoutPage.
          onPixCreated?.();
          patch({
            phase: "pix_pending",
            pixData: {
              qrCode: msg.payload.qr_code,
              qrCodeBase64: msg.payload.qr_code_base64,
            },
            errorMessage: null,
          });
          break;

        case "order_confirmed":
          // Terminal: backend fecha o socket imediatamente após esta mensagem.
          terminalPhaseRef.current = true;
          patch({
            phase: "order_confirmed",
            errorMessage: null,
          });
          break;

        case "error":
          hadServerError.current = true;
          patch({
            phase: "error",
            errorMessage: normalizeErrorMessage(msg.payload),
          });
          break;
      }
    },
    [patch, onPartialReservation, onPixCreated],
  );

  useEffect(() => {
    handleMessageRef.current = handleMessage;
  }, [handleMessage]);

  const openSocket = useCallback(
    async (firstMessage: Record<string, unknown>) => {
      if (connectingRef.current) return;

      if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
        intentionalClose.current = true;
        wsRef.current.close(1000, "restart");
        wsRef.current = null;
      }

      connectingRef.current = true;
      intentionalClose.current = false;
      hadServerError.current = false;
      terminalPhaseRef.current = false;

      try {
        const res = await getWebsocketAuthToken();
        if (!res.success) {
          patch({ phase: "error", errorMessage: "Não foi possível autenticar. Tente novamente." });
          return;
        }

        const ws = new WebSocket(`${url}?token=${res.data.token}`);
        wsRef.current = ws;

        ws.onopen = () => { ws.send(JSON.stringify(firstMessage)); };
        ws.onmessage = (raw) => handleMessageRef.current?.(raw);
        ws.onerror = () => { patch({ phase: "error", errorMessage: "Erro na conexão com o servidor." }); };
        ws.onclose = () => {
          wsRef.current = null;
          if (!intentionalClose.current && !hadServerError.current && !terminalPhaseRef.current) {
            patch({ phase: "error", errorMessage: "Conexão encerrada inesperadamente." });
          }
        };
      } finally {
        connectingRef.current = false;
      }
    },
    [url, patch],
  );

  const buyRequest = useCallback(
    async (items: BuyRequestItem[]) => {
      setState({ ...INITIAL, phase: "connecting" });
      await openSocket({ type: "buy_request", payload: { items } });
      patch({ phase: "awaiting_reservation" });
    },
    [openSocket, patch],
  );

  const resumeSession = useCallback(
    async (sessionId: string) => {
      setState({ ...INITIAL, phase: "connecting" });
      await openSocket({ type: "resume_session", payload: { session_id: sessionId } });
      patch({ phase: "awaiting_reservation" });
    },
    [openSocket, patch],
  );

  const confirmPartial = useCallback(() => {
    sendJSON({ type: "confirm_partial", payload: {} });
    patch({ phase: "awaiting_reservation" });
  }, [sendJSON, patch]);

  const cancelReservation = useCallback(() => {
    sendJSON({ type: "cancel", payload: {} });
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
          payer_email: payload.payer_email,
          installments: payload.installments,
          identification_type: payload.identification_type,
          identification_number: payload.identification_number,
        },
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

  return { state, buyRequest, resumeSession, confirmPartial, cancelReservation, submitPayment, reset };
}