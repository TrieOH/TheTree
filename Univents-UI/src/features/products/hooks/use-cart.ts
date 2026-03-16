import { useStore } from "@tanstack/react-store";
import { cartStore, cartActions } from "../model/cart";

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
    addItem: (product: {
      id: string;
      name: string;
      price_cents: number;
      inventory_remaining?: number;
      has_inventory?: boolean;
    }, quantity: number) => {
      cartActions.addItem(editionId, product, quantity);
    },
    removeItem: (id: string) => {
      cartActions.removeItem(editionId, id);
    },
    updateQuantity: (id: string, quantity: number) => {
      cartActions.updateQuantity(editionId, id, quantity);
    },
    clearCart: () => {
      cartActions.clearCart(editionId);
    },
  };
}
