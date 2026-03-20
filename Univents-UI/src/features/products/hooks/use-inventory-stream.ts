import { useEffect, useState } from "react";
import { env } from "@/env";

interface InventoryUpdate {
  product_id: string;
  inventory_remaining: number;
}

export function useInventoryStream(eventId: string, editionId: string) {
  const [inventory, setInventory] = useState<Record<string, number>>({});
  const [status, setStatus] = useState<"connecting" | "open" | "error">("connecting");

  useEffect(() => {
    if (!editionId) return;

    const controller = new AbortController();
    const url = `${env.VITE_API_URL}events/${eventId}/editions/${editionId}/products/inventory/stream`;

    void (async () => {
      try {
        const res = await fetch(url, {
          credentials: "include",
          signal: controller.signal,
          headers: { Accept: "text/event-stream" },
        });

        if (!res.ok || !res.body) {
          setStatus("error");
          return;
        }
        setStatus("open");

        const reader = res.body.getReader();
        const decoder = new TextDecoder();
        let buffer = "";

        for (; ;) {
          const { done, value } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });

          const events = buffer.split("\n\n");

          buffer = events.pop() ?? "";

          for (const event of events) {
            const lines = event.split("\n");

            let eventType = "message";
            let data = "";

            for (const line of lines) {
              if (line.startsWith("event:")) {
                eventType = line.slice(6).trim();
              } else if (line.startsWith("data:")) {
                data = line.slice(5).trim();
              }
            }

            if (!data) continue;
            if (eventType !== "inventory_update") continue;

            try {
              const updates = JSON.parse(data) as InventoryUpdate[];
              setInventory((prev) => {
                const next = { ...prev };
                for (const item of updates) {
                  next[item.product_id] = item.inventory_remaining;
                }
                return next;
              });
            } catch {
              console.error("Failed to parse SSE data:", data);
            }
          }
        }
      } catch (err) {
        if ((err as Error).name === "AbortError") return;
        setStatus("error");
      }
    })();

    return () => {
      controller.abort();
    };
  }, [eventId, editionId]);

  return { inventory, status };
}
