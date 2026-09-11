import { createSignal, createEffect } from "solid-js";
import { purchaseWsToken } from "../api";

export function usePurchaseSocket(
  purchaseId: () => string,
  onChange: (purchaseId: string) => void,
) {
  const [connected, setConnected] = createSignal(false);

  createEffect(
    () => purchaseId(),
    (id) => {
      let stopped = false;
      let terminal = false;
      let socket: WebSocket | undefined;
      let retry: ReturnType<typeof setTimeout> | undefined;

      const connect = async () => {
        try {
          const stored = sessionStorage.getItem(`purchase-ws:${id}`);
          sessionStorage.removeItem(`purchase-ws:${id}`);
          const token = stored ?? (await purchaseWsToken(id)).token;
          if (stopped) return;
          const url = new URL("ws", import.meta.env.VITE_API_URL ?? "http://localhost:8081");
          url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
          url.searchParams.set("token", token);
          socket = new WebSocket(url);
          socket.onopen = () => setConnected(true);
          socket.onmessage = (event) => {
            try {
              const frame = JSON.parse(String(event.data)) as { type?: string; payload?: { status?: string } };
              terminal = ["purchase.confirmed", "purchase.expired", "purchase.cancelled", "purchase.failed", "purchase.rejected"].includes(frame.type ?? "") ||
                (frame.type === "purchase.snapshot" && frame.payload?.status !== "pending");
            } finally {
              onChange(id);
            }
          };
          socket.onclose = () => {
            setConnected(false);
            if (!stopped && !terminal) retry = setTimeout(() => void connect(), 1000);
          };
        } catch {
          if (!stopped) retry = setTimeout(() => void connect(), 1000);
        }
      };

      void connect();

      return () => {
        stopped = true;
        if (retry) clearTimeout(retry);
        socket?.close();
        setConnected(false);
      };
    },
  );

  return connected;
}
