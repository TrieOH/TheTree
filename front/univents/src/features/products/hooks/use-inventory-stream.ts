import { useEffect, useRef, useState } from "react";
import { fetchEventSource, EventStreamContentType } from "@microsoft/fetch-event-source";
import { env } from "@/env";

interface InventoryUpdate {
  product_id: string;
  inventory_remaining: number;
}

class RetriableError extends Error { }
class FatalError extends Error { }

const INITIAL_RETRY_MS = 1_000;
const MAX_RETRY_MS = 30_000;

export function useInventoryStream(eventId: string, editionId: string) {
  const [inventory, setInventory] = useState<Record<string, number>>({});
  const [status, setStatus] = useState<"connecting" | "open" | "error">("connecting");
  const retryDelay = useRef(INITIAL_RETRY_MS);

  useEffect(() => {
    if (!editionId) return;

    const controller = new AbortController();
    const url = `${env.VITE_API_URL}events/${eventId}/editions/${editionId}/products/inventory/stream`;

    console.log(`[InventoryStream] Connecting to ${url}`);
    setInventory({}); // Clear on new connection to avoid stale data

    void fetchEventSource(url, {
      credentials: "include",
      signal: controller.signal,

      // eslint-disable-next-line @typescript-eslint/require-await
      onopen: async (res) => {
        if (res.ok && res.headers.get("content-type") === EventStreamContentType) {
          console.log(`[InventoryStream] Connection established to ${url}`);
          retryDelay.current = INITIAL_RETRY_MS;
          setStatus("open");
        } else if (res.status >= 400 && res.status < 500 && res.status !== 429) {
          throw new FatalError(`HTTP ${res.status}`);
        } else throw new RetriableError();
      },

      onmessage: (ev) => {
        if (ev.event !== "inventory_update") return;
        try {
          console.log(`[InventoryStream] Message received:`, ev.data);
          const updates = JSON.parse(ev.data) as InventoryUpdate[];
          setInventory((prev) => {
            let hasChanges = false;
            const next = { ...prev };
            for (const item of updates) {
              if (next[item.product_id] !== item.inventory_remaining) {
                next[item.product_id] = item.inventory_remaining;
                hasChanges = true;
              }
            }
            return hasChanges ? next : prev;
          });
        } catch {
          console.error("Failed to parse SSE data:", ev.data);
        }
      },

      onclose: () => {
        console.warn(`[InventoryStream] Connection closed for ${url}`);
        throw new RetriableError();
      },

      onerror: (err) => {
        console.error(`[InventoryStream] Connection error for ${url}:`, err);
        if (err instanceof FatalError) {
          setStatus("error");
          throw err;
        }
        setStatus("error");
        retryDelay.current = Math.min(retryDelay.current * 2, MAX_RETRY_MS);
        return retryDelay.current;
      },
    });

    return () => {
      console.log(`[InventoryStream] Aborting connection to ${url}`);
      controller.abort();
    };
  }, [eventId, editionId]);

  return { inventory, status };
}