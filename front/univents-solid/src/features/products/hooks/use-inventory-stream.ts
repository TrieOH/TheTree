import { useQueryClient } from "@trieoh/front-core/solid";
import type { StoreStockItem } from "@trieoh/univents-api/schemas";
import { createSignal, onSettled } from "solid-js";
import { cart } from "../model/cart";

export function useInventoryStream(
  editionId: string,
  initialStock: StoreStockItem[],
) {
  const queryClient = useQueryClient();
  const [stockItems, setStockItems] = createSignal(initialStock);

  onSettled(() => {
    if (!editionId) return;
    const baseUrl = import.meta.env.VITE_API_URL ?? "http://localhost:8081";
    const stream = new EventSource(
      new URL(`/editions/${editionId}/store/stream`, baseUrl),
    );

    const apply = (item: StoreStockItem) => {
      cart.stock(editionId, item.id, item.stock);
      setStockItems((current) => {
        const exists = current.some((entry) => entry.id === item.id);
        return exists
          ? current.map((entry) => (entry.id === item.id ? item : entry))
          : [...current, item];
      });
      queryClient.setQueryData<StoreStockItem[]>(
        ["store", "stock", editionId],
        (current = []) => {
          const exists = current.some((entry) => entry.id === item.id);
          return exists
            ? current.map((entry) => (entry.id === item.id ? item : entry))
            : [...current, item];
        },
      );
    };

    const snapshot = (event: MessageEvent<string>) => {
      try {
        for (const item of JSON.parse(event.data) as StoreStockItem[])
          apply(item);
      } catch {
        return;
      }
    };
    const stock = (event: MessageEvent<string>) => {
      try {
        apply(JSON.parse(event.data) as StoreStockItem);
      } catch {
        return;
      }
    };

    stream.addEventListener("snapshot", snapshot);
    stream.addEventListener("stock", stock);
    return () => stream.close();
  });

  return stockItems;
}
