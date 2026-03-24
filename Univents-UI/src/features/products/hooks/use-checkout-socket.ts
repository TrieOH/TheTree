import { useCallback, useReducer, useRef } from "react";
import { getWebsocketAuthToken } from "../api";
import type { SubmitPaymentPayloadI } from "@/features/payments/model";
import type { BuyRequestItemI, ReservedItemI, UnavailableItemI } from "../model";

// ─── Server → Client Messages ─────────────────────────────────────────────────

type ServerMessage =
  | { type: "reservation_confirmed"; payload: { session_id: string; expires_at: string; reserved_items: ReservedItemI[]; total_cents: number } }
  | { type: "partial_reservation"; payload: { reserved: ReservedItemI[]; unavailable: UnavailableItemI[]; confirm_deadline: string } }
  | { type: "reservation_failed"; payload: { unavailable: UnavailableItemI[] } }
  | { type: "reservation_cancelled" }
  | { type: "session_expired" }
  | { type: "payment_processing" }
  | { type: "payment_failed"; payload: { payment_intent_id: string } }
  | { type: "payment_pending"; payload: string }
  | { type: "pix_created"; payload: { qr_code: string; qr_code_base64: string } }
  | { type: "purchase_failed"; payload: { reason: string; product_ids: string[] } | { invalid_products: UnavailableItemI[] } | string }
  | { type: "order_confirmed"; payload: { purchase_id: string } | { payment_intent_id: string } }
  | { type: "error"; payload: string };

// ─── Checkout Phase ───────────────────────────────────────────────────────────

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
  | "pix_pending"
  | "order_confirmed"
  | "session_expired"
  | "error";

// Phases where the backend closes the socket on its own after the last message.
// An onclose following one of these is expected — not a connection drop.
const TERMINAL_PHASES = new Set<CheckoutPhase>([
  "idle",
  "order_confirmed",
  "payment_pending",    // card: backend closes after this, webhook handles the rest
  "pix_pending",        // pix: backend deletes session before emitting, then closes
  "reservation_failed", // backend closes after sending this
  "session_expired",    // backend closes after sending this
  "error",              // backend closes after sending this
]);

// ─── State ────────────────────────────────────────────────────────────────────

export interface CheckoutState {
  phase: CheckoutPhase;
  errorMessage: string | null;
  paymentIntentId: string | null;
  sessionId: string | null;
  reservationExpiresAt: string | null;
  reservedItems: ReservedItemI[];
  totalCents: number;
  partialData: {
    reserved: ReservedItemI[];
    unavailable: UnavailableItemI[];
    confirmDeadline: string;
  } | null;
  pixData: { qrCode: string; qrCodeBase64: string } | null;
  pendingMessage: string | null;
}

const INITIAL_STATE: CheckoutState = {
  phase: "idle",
  errorMessage: null,
  paymentIntentId: null,
  sessionId: null,
  reservationExpiresAt: null,
  reservedItems: [],
  totalCents: 0,
  partialData: null,
  pixData: null,
  pendingMessage: null,
};

// ─── Reducer Actions ──────────────────────────────────────────────────────────

type Action =
  | { type: "RESET" }
  | { type: "CONNECTING" }
  | { type: "AWAITING_RESERVATION" }
  | { type: "AWAITING_PAYMENT" }
  | { type: "RESERVATION_CONFIRMED"; sessionId: string; expiresAt: string; items: ReservedItemI[]; totalCents: number }
  | { type: "PARTIAL_RESERVATION"; reserved: ReservedItemI[]; unavailable: UnavailableItemI[]; confirmDeadline: string }
  | { type: "RESERVATION_FAILED" }
  | { type: "SESSION_EXPIRED" }
  | { type: "PAYMENT_PROCESSING" }
  | { type: "PAYMENT_FAILED"; paymentIntentId: string }
  | { type: "PAYMENT_PENDING"; message: string }
  | { type: "PIX_CREATED"; qrCode: string; qrCodeBase64: string }
  | { type: "ORDER_CONFIRMED" }
  | { type: "ERROR"; message: string }
  | { type: "UNEXPECTED_CLOSE" };

function reducer(state: CheckoutState, action: Action): CheckoutState {
  switch (action.type) {
    case "RESET":
      return INITIAL_STATE;

    case "CONNECTING":
      return { ...INITIAL_STATE, phase: "connecting" };

    case "AWAITING_RESERVATION":
      return { ...state, phase: "awaiting_reservation" };

    case "AWAITING_PAYMENT":
      return { ...state, phase: "awaiting_payment" };

    case "RESERVATION_CONFIRMED":
      return {
        ...state,
        phase: "reservation_confirmed",
        sessionId: action.sessionId,
        reservationExpiresAt: action.expiresAt,
        reservedItems: action.items,
        totalCents: action.totalCents,
        errorMessage: null,
      };

    case "PARTIAL_RESERVATION": {
      const totalCents = action.reserved.reduce(
        (sum, item) => sum + item.price_cents * item.quantity,
        0,
      );
      return {
        ...state,
        phase: "partial_reservation",
        reservedItems: action.reserved,
        totalCents,
        partialData: {
          reserved: action.reserved,
          unavailable: action.unavailable,
          confirmDeadline: action.confirmDeadline,
        },
        errorMessage: null,
      };
    }

    case "RESERVATION_FAILED":
      return {
        ...state,
        phase: "reservation_failed",
        errorMessage: "Não foi possível reservar os itens selecionados.",
      };

    case "SESSION_EXPIRED":
      return {
        ...state,
        phase: "session_expired",
        errorMessage: "Sua reserva expirou.",
      };

    case "PAYMENT_PROCESSING":
      return { ...state, phase: "payment_processing", errorMessage: null };

    case "PAYMENT_FAILED":
      return {
        ...state,
        phase: "payment_failed",
        paymentIntentId: action.paymentIntentId,
        errorMessage: "Pagamento recusado. Tente novamente.",
      };

    case "PAYMENT_PENDING":
      return {
        ...state,
        phase: "payment_pending",
        pendingMessage: action.message,
        errorMessage: null,
      };

    case "PIX_CREATED":
      return {
        ...state,
        phase: "pix_pending",
        pixData: { qrCode: action.qrCode, qrCodeBase64: action.qrCodeBase64 },
        errorMessage: null,
      };

    case "ORDER_CONFIRMED":
      return { ...state, phase: "order_confirmed", errorMessage: null };

    case "ERROR":
      return { ...state, phase: "error", errorMessage: action.message };

    // The socket closed without us initiating it.
    // Only surface an error if we weren't already in a terminal phase —
    // those phases imply the backend intentionally closed the connection.
    case "UNEXPECTED_CLOSE":
      if (TERMINAL_PHASES.has(state.phase)) return state;
      return {
        ...state,
        phase: "error",
        errorMessage: "Conexão encerrada inesperadamente.",
      };

    default:
      return state;
  }
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

/**
 * Translates raw Go/network error strings into user-readable Portuguese messages.
 */
function parseErrorMessage(raw: string): string {
  if (raw.includes("i/o timeout")) return "Tempo esgotado para confirmar a reserva.";
  if (raw.includes("connection reset")) return "Conexão encerrada inesperadamente.";
  if (raw.includes("broken pipe")) return "Conexão encerrada inesperadamente.";
  return "Ocorreu um erro inesperado.";
}

function toAction(msg: ServerMessage): Action | null {
  switch (msg.type) {
    case "reservation_confirmed":
      return {
        type: "RESERVATION_CONFIRMED",
        sessionId: msg.payload.session_id,
        expiresAt: msg.payload.expires_at,
        items: msg.payload.reserved_items,
        totalCents: msg.payload.total_cents,
      };

    case "partial_reservation":
      return {
        type: "PARTIAL_RESERVATION",
        reserved: msg.payload.reserved,
        unavailable: msg.payload.unavailable,
        confirmDeadline: msg.payload.confirm_deadline,
      };

    case "reservation_failed":
      return { type: "RESERVATION_FAILED" };

    case "reservation_cancelled":
      // Backend confirms our cancel — reset to idle.
      return { type: "RESET" };

    case "session_expired":
      return { type: "SESSION_EXPIRED" };

    case "payment_processing":
      return { type: "PAYMENT_PROCESSING" };

    case "payment_failed":
      return { type: "PAYMENT_FAILED", paymentIntentId: msg.payload.payment_intent_id };

    case "payment_pending":
      return { type: "PAYMENT_PENDING", message: msg.payload };

    case "pix_created":
      return {
        type: "PIX_CREATED",
        qrCode: msg.payload.qr_code,
        qrCodeBase64: msg.payload.qr_code_base64,
      };

    case "purchase_failed":
      return { type: "ERROR", message: "Não foi possível processar a compra. Tente novamente." };

    case "order_confirmed":
      return { type: "ORDER_CONFIRMED" };

    case "error":
      return { type: "ERROR", message: parseErrorMessage(msg.payload) };

    default:
      return null;
  }
}

// ─── Hook ─────────────────────────────────────────────────────────────────────

export interface UseCheckoutSocketOptions {
  url: string;
  /**
   * Called when a partial_reservation arrives so the parent can
   * adjust its cart display before the user confirms or cancels.
   */
  onPartialReservation?: (reserved: ReservedItemI[]) => void;
  /**
   * Called when pix_created arrives. The backend already deleted the
   * session at this point — the caller is responsible for clearing sessionStorage.
   */
  onPixCreated?: () => void;
}

export interface UseCheckoutSocketReturn {
  state: CheckoutState;
  /** Start a new purchase flow for the given items. */
  buyRequest: (items: BuyRequestItemI[]) => Promise<void>;
  /** Resume an existing reservation session (e.g. after a page refresh). */
  resumeSession: (sessionId: string) => Promise<void>;
  /** Accept a partial reservation and proceed to payment. */
  confirmPartial: () => void;
  /** Cancel reservation and reset state. */
  cancelReservation: () => void;
  /** Submit payment details over the open socket. */
  submitPayment: (payload: SubmitPaymentPayloadI) => void;
  /** Hard reset — closes the socket and returns to idle. */
  reset: () => void;
}

export function useCheckoutSocket({
  url,
  onPartialReservation,
  onPixCreated,
}: UseCheckoutSocketOptions): UseCheckoutSocketReturn {
  const [state, dispatch] = useReducer(reducer, INITIAL_STATE);

  const wsRef = useRef<WebSocket | null>(null);

  // Guards against double-connecting while getWebsocketAuthToken is in-flight.
  // Set to true before the await, back to false in finally — the WebSocket
  const isConnectingRef = useRef(false);

  // When we close the socket intentionally (reset, cancel, restart), we flip
  // this to true so onclose doesn't treat it as an unexpected drop.
  // Cleared back to false inside onclose itself.
  const closingCleanlyRef = useRef(false);

  // Keep the latest callbacks in refs so openSocket never needs them as deps.
  const onPartialRef = useRef(onPartialReservation);
  const onPixRef = useRef(onPixCreated);
  onPartialRef.current = onPartialReservation;
  onPixRef.current = onPixCreated;

  // ── Internal helpers ─────────────────────────────────────────────────────────

  /**
   * Closes the current socket (if any) marking it as intentional so that
   * onclose won't dispatch UNEXPECTED_CLOSE.
   */
  const closeSocket = useCallback((reason = "done") => {
    const ws = wsRef.current;
    if (!ws) return;
    closingCleanlyRef.current = true;
    wsRef.current = null;
    if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
      ws.close(1000, reason);
    }
  }, []);

  /**
   * Sends a JSON frame over the socket.
   * Returns true if sent, false if the socket wasn't open.
   * By default, dispatches ERROR on failure; pass { silent: true } for
   * best-effort fire-and-forget sends (e.g. cancel on an already-dead socket).
   */
  const sendJSON = useCallback(
    (payload: Record<string, unknown>, { silent = false } = {}): boolean => {
      const ws = wsRef.current;
      if (ws?.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(payload));
        return true;
      }
      if (!silent) {
        dispatch({ type: "ERROR", message: "Conexão perdida. Tente novamente." });
      }
      return false;
    },
    [],
  );

  // ── Core socket opener ───────────────────────────────────────────────────────

  const openSocket = useCallback(
    async (firstMessage: Record<string, unknown>): Promise<void> => {
      // Debounce: ignore if already in the middle of authenticating.
      if (isConnectingRef.current) return;

      // Tear down any existing socket cleanly before re-opening.
      closeSocket("restart");

      isConnectingRef.current = true;
      closingCleanlyRef.current = false;

      try {
        const res = await getWebsocketAuthToken();
        if (!res.success) {
          dispatch({ type: "ERROR", message: "Não foi possível autenticar. Tente novamente." });
          return;
        }

        const ws = new WebSocket(`${url}?token=${res.data.token}`);
        // Assign synchronously — from this point onward any concurrent
        // openSocket call will hit the closeSocket("restart") guard above.
        wsRef.current = ws;

        ws.onopen = () => {
          ws.send(JSON.stringify(firstMessage));
          dispatch({ type: "AWAITING_RESERVATION" });
        };

        ws.onmessage = (event: MessageEvent<string>) => {
          // Discard frames from a stale socket that was already replaced.
          if (ws !== wsRef.current) return;

          let msg: ServerMessage;
          try {
            msg = JSON.parse(event.data) as ServerMessage;
          } catch {
            return; // malformed frame — ignore
          }

          // Side-effect callbacks (always fresh via ref, no dep needed).
          if (msg.type === "partial_reservation") {
            onPartialRef.current?.(msg.payload.reserved);
          }
          if (msg.type === "pix_created") {
            onPixRef.current?.();
          }

          const action = toAction(msg);
          if (action) dispatch(action);
        };

        ws.onerror = () => {
          // Ignore errors from a stale socket.
          if (ws !== wsRef.current) return;
          dispatch({ type: "ERROR", message: "Erro na conexão com o servidor." });
        };

        ws.onclose = () => {
          // Clear our ref only if it still points to this socket.
          if (wsRef.current === ws) wsRef.current = null;

          if (closingCleanlyRef.current) {
            closingCleanlyRef.current = false;
            return;
          }
          dispatch({ type: "UNEXPECTED_CLOSE" });
        };
      } finally {
        isConnectingRef.current = false;
      }
    },
    [url, closeSocket],
  );

  const buyRequest = useCallback(
    async (items: BuyRequestItemI[]) => {
      dispatch({ type: "CONNECTING" });
      await openSocket({ type: "buy_request", payload: { items } });
    },
    [openSocket],
  );

  const resumeSession = useCallback(
    async (sessionId: string) => {
      dispatch({ type: "CONNECTING" });
      await openSocket({ type: "resume_session", payload: { session_id: sessionId } });
    },
    [openSocket],
  );

  const confirmPartial = useCallback(() => {
    sendJSON({ type: "confirm_partial", payload: {} });
    dispatch({ type: "AWAITING_RESERVATION" });
  }, [sendJSON]);

  const cancelReservation = useCallback(() => {
    const sent = sendJSON({ type: "cancel", payload: {} }, { silent: true });
    if (!sent) {
      dispatch({ type: "RESET" });
      return;
    }
    // If sent, wait for reservation_cancelled → RESET via onmessage
    // The backend closes the socket next → onclose is cleared via TERMINAL_PHASES (idle)
  }, [sendJSON, closeSocket]);

  const submitPayment = useCallback(
    (payload: SubmitPaymentPayloadI) => {
      const sent = sendJSON({
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
      // Only advance the phase if the frame was actually sent.
      // sendJSON already dispatched ERROR if it wasn't.
      if (sent) dispatch({ type: "AWAITING_PAYMENT" });
    },
    [sendJSON],
  );

  const reset = useCallback(() => {
    closeSocket("reset");
    dispatch({ type: "RESET" });
  }, [closeSocket]);

  return {
    state,
    buyRequest,
    resumeSession,
    confirmPartial,
    cancelReservation,
    submitPayment,
    reset,
  };
}