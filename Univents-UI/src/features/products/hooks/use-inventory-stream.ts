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

    void fetchEventSource(url, {
      credentials: "include",
      signal: controller.signal,

      onopen: async (res) => {
        if (res.ok && res.headers.get("content-type") === EventStreamContentType) {
          retryDelay.current = INITIAL_RETRY_MS;
          setStatus("open");
        } else if (res.status >= 400 && res.status < 500 && res.status !== 429) {
          throw new FatalError(`HTTP ${res.status}`);
        } else throw new RetriableError();
      },

      onmessage: (ev) => {
        if (ev.event !== "inventory_update") return;
        try {
          const updates = JSON.parse(ev.data) as InventoryUpdate[];
          setInventory((prev) => {
            const next = { ...prev };
            for (const item of updates) {
              next[item.product_id] = item.inventory_remaining;
            }
            return next;
          });
        } catch {
          console.error("Failed to parse SSE data:", ev.data);
        }
      },

      onclose: () => {
        throw new RetriableError();
      },

      onerror: (err) => {
        if (err instanceof FatalError) {
          setStatus("error");
          throw err;
        }
        setStatus("error");
        retryDelay.current = Math.min(retryDelay.current * 2, MAX_RETRY_MS);
        return retryDelay.current;
      },
    });

    return () => { controller.abort(); };
  }, [eventId, editionId]);

  return { inventory, status };
}