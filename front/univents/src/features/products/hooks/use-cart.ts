import { useStore } from "@tanstack/react-store";
import { cartStore, cartActions, isLimitReached, getProductMaxQuantity } from "../model/cart";
import type { CartItem } from "../model/cart";

export function useCart(editionId: string) {
  const items = useStore(cartStore, (state) => state.carts[editionId] ?? []);

  const totalCents = items.reduce(
    (acc, item) => acc + item.price_cents * item.quantity,
    0
  );

  const itemCount = items.reduce((acc, item) => acc + item.quantity, 0);

  return {
    items,
    totalCents,
    itemCount,
    addItem: (product: Omit<CartItem, "quantity">, quantity: number) => {
      cartActions.addItem(editionId, product, quantity);
    },
    removeItem: (id: string) => {
      cartActions.removeItem(editionId, id);
    },
    updateQuantity: (id: string, quantity: number) => {
      cartActions.updateQuantity(editionId, id, quantity);
    },
    replaceCart: (newItems: CartItem[]) => {
      cartActions.replaceCart(editionId, newItems);
    },
    clearCart: () => {
      cartActions.clearCart(editionId);
    },
    isLimitReached: (product: Pick<CartItem, "has_inventory" | "inventory_remaining">, currentQuantity: number) => {
      return isLimitReached(product, currentQuantity);
    },
    getMaxQuantity: (product: Pick<CartItem, "has_inventory" | "inventory_remaining">) => {
      return getProductMaxQuantity(product);
    },
  };
}
