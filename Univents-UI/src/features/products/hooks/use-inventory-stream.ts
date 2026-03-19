import { env } from "@/env";
import { useEffect, useRef, useState } from "react";

interface InventoryUpdate {
  product_id: string;
  inventory_remaining: number;
}

export function useInventoryStream(eventId: string, editionId: string) {
  const [inventory, setInventory] = useState<Record<string, number>>({});
  const [status, setStatus] = useState<"connecting" | "open" | "error">("connecting");
  const sourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    if (!editionId) return;

    const url = `${env.VITE_API_URL}events/${eventId}/editions/${editionId}/products/inventory/stream`;
    const source = new EventSource(url);
    sourceRef.current = source;

    source.onopen = () => setStatus("open");

    source.addEventListener("inventory_update", (e: MessageEvent) => {
      const updates: InventoryUpdate[] = JSON.parse(e.data);

      setInventory((prev) => {
        const next = { ...prev };
        for (const item of updates) {
          next[item.product_id] = item.inventory_remaining;
        }
        return next;
      });
    });

    source.onerror = () => {
      setStatus("error");
      source.close();
    };

    return () => {
      source.close();
    };
  }, [editionId]);

  return { inventory, status };
}