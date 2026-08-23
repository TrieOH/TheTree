import { useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import { env } from "@/env";
import { productKeys } from "../api/query-keys";
import { cartActions } from "../model/cart";

type StockItem = {
  id: string;
  item_type: "ticket" | "product" | "program_occurrence";
  stock: number | null;
};

export function useInventoryStream(editionId: string) {
  const queryClient = useQueryClient();
  useEffect(() => {
    if (!editionId) return;

    const url = new URL(`editions/${editionId}/store/stream`, env.VITE_API_URL);
    const stream = new EventSource(url);
    const apply = (item: StockItem) => {
      cartActions.updateStock(editionId, item.id, item.stock);
      queryClient.setQueryData<StockItem[] | undefined>(
        productKeys.storeStock(editionId),
        (current) => {
          const items = current ?? [];
          const index = items.findIndex((entry) => entry.id === item.id);
          if (index < 0) return [...items, item];
          return items.map((entry, entryIndex) =>
            entryIndex === index ? item : entry,
          );
        },
      );
    };
    const parse = (data: string) => JSON.parse(data) as StockItem;

    const onSnapshot = (event: MessageEvent<string>) => {
      try {
        for (const item of JSON.parse(event.data) as StockItem[]) apply(item);
      } catch {
        // A future valid event can restore the authoritative stock.
      }
    };
    const onStock = (event: MessageEvent<string>) => {
      try {
        apply(parse(event.data));
      } catch {
        // Ignore malformed network data.
      }
    };

    stream.addEventListener("snapshot", onSnapshot);
    stream.addEventListener("stock", onStock);
    return () => stream.close();
  }, [editionId, queryClient]);
}
