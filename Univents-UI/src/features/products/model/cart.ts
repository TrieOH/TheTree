import { Store } from "@tanstack/react-store";

export interface CartItem {
  id: string;
  name: string;
  price_cents: number;
  quantity: number;
  inventory_remaining?: number;
  has_inventory?: boolean;
}

export interface CartState {
  carts: Record<string, CartItem[]>;
}

const STORAGE_KEY = "univents-cart";

const getInitialState = (): CartState => {
  if (typeof window === "undefined") return { carts: {} };
  const saved = localStorage.getItem(STORAGE_KEY);
  if (saved) {
    try {
      const parsed = JSON.parse(saved) as Partial<CartState>;
      return { carts: parsed.carts ?? {} };
    } catch (e) {
      console.error("Error loading the cart from localStorage", e);
    }
  }
  return { carts: {} };
};

export const cartStore = new Store<CartState>(getInitialState());

cartStore.subscribe(() => {
  const state = cartStore.state;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
});

export const cartActions = {
  addItem: (editionId: string, product: {
    id: string;
    name: string;
    price_cents: number;
    inventory_remaining?: number;
    has_inventory?: boolean;
  }, quantity: number) => {
    cartStore.setState((prev) => {
      const currentItems = prev.carts[editionId] ?? [];
      const existing = currentItems.find((i) => i.id === product.id);

      let newItems;
      if (existing) {
        let newQuantity = existing.quantity + quantity;

        // Limit by inventory if applicable
        if (product.has_inventory && typeof product.inventory_remaining === 'number') {
          newQuantity = Math.min(newQuantity, product.inventory_remaining);
        }

        newItems = currentItems.map((i) =>
          i.id === product.id ? { ...i, quantity: newQuantity } : i
        );
      } else {
        let finalQuantity = quantity;
        if (product.has_inventory && typeof product.inventory_remaining === 'number') {
          finalQuantity = Math.min(quantity, product.inventory_remaining);
        }
        newItems = [...currentItems, { ...product, quantity: finalQuantity }];
      }

      return {
        ...prev,
        carts: {
          ...prev.carts,
          [editionId]: newItems,
        },
      };
    });
  },
  removeItem: (editionId: string, id: string) => {
    cartStore.setState((prev) => ({
      ...prev,
      carts: {
        ...prev.carts,
        [editionId]: (prev.carts[editionId] ?? []).filter((i) => i.id !== id),
      },
    }));
  },
  updateQuantity: (editionId: string, id: string, quantity: number) => {
    if (quantity <= 0) {
      cartActions.removeItem(editionId, id);
      return;
    }
    cartStore.setState((prev) => ({
      ...prev,
      carts: {
        ...prev.carts,
        [editionId]: (prev.carts[editionId] ?? []).map((i) => {
          if (i.id === id) {
            let newQuantity = quantity;
            if (i.has_inventory && typeof i.inventory_remaining === 'number') {
              newQuantity = Math.min(quantity, i.inventory_remaining);
            }
            return { ...i, quantity: newQuantity };
          }
          return i;
        }),
      },
    }));
  },
  clearCart: (editionId: string) => {
    cartStore.setState((prev) => ({
      ...prev,
      carts: {
        ...prev.carts,
        [editionId]: [],
      },
    }));
  },
};
