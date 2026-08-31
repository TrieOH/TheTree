import { useQueryClient } from "@tanstack/react-query";
import type { Checkout } from "@trieoh/univents-api/schemas";
import { useEffect, useRef, useState } from "react";
import { env } from "@/env";
import { checkoutQueryOptions, getWsTokenFn } from "../api";
import { invalidatePurchaseCaches } from "../api/cache";

export function usePurchaseSocket(purchaseId: string, pending: boolean) {
  const queryClient = useQueryClient();
  const [connected, setConnected] = useState(false);
  const retryRef = useRef<ReturnType<typeof setTimeout>>(undefined);

  useEffect(() => {
    if (!purchaseId || !pending) return;

    let stopped = false;
    let terminal = false;
    let socket: WebSocket | undefined;

    const connect = async () => {
      try {
        const storedToken = sessionStorage.getItem(`purchase-ws:${purchaseId}`);
        sessionStorage.removeItem(`purchase-ws:${purchaseId}`);
        const token = storedToken ?? (await getWsTokenFn(purchaseId)).token;
        if (stopped) return;

        const url = new URL("ws", env.VITE_API_URL);
        url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
        url.searchParams.set("token", token);
        socket = new WebSocket(url);

        socket.onopen = () => setConnected(true);
        socket.onmessage = (event) => {
          try {
            const frame = JSON.parse(event.data) as {
              type?: string;
              payload?: { status?: string };
            };
            terminal =
              frame.type === "purchase.confirmed" ||
              frame.type === "purchase.expired" ||
              frame.type === "purchase.cancelled" ||
              frame.type === "purchase.failed" ||
              frame.type === "purchase.rejected" ||
              (frame.type === "purchase.snapshot" &&
                frame.payload?.status !== "pending");
          } catch {
            // REST invalidation below remains the source of truth.
          }
          const editionId = queryClient.getQueryData<Checkout>(
            checkoutQueryOptions(purchaseId).queryKey,
          )?.edition_id;
          invalidatePurchaseCaches(queryClient, purchaseId, editionId);
        };
        socket.onclose = () => {
          setConnected(false);
          if (!stopped && !terminal)
            retryRef.current = setTimeout(connect, 1_000);
        };
      } catch {
        if (!stopped) retryRef.current = setTimeout(connect, 1_000);
      }
    };

    void connect();
    return () => {
      stopped = true;
      clearTimeout(retryRef.current);
      socket?.close();
      setConnected(false);
    };
  }, [pending, purchaseId, queryClient]);

  return connected;
}
