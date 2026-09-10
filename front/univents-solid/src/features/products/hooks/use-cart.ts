import { createMemo } from "solid-js";
import { cart, type CartItem } from "../model/cart";

export function useCart(editionId: string) {
  const items = createMemo(() => cart.items(editionId));
  const count = createMemo(() =>
    items().reduce((total, item) => total + item.quantity, 0),
  );
  const total = createMemo(() =>
    items().reduce((sum, item) => sum + item.price_cents * item.quantity, 0),
  );
  return {
    items,
    count,
    total,
    add: (item: Omit<CartItem, "quantity">) => cart.add(editionId, item),
    quantity: (item: CartItem, quantity: number) =>
      cart.quantity(editionId, item, quantity),
    remove: (item: CartItem) => cart.remove(editionId, item.id, item.type),
    clear: () => cart.clear(editionId),
  };
}
