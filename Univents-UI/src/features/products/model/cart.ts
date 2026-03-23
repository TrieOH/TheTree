import { Store } from "@tanstack/react-store";

export const GLOBAL_MAX_QUANTITY = 999;

export interface CartItem {
  id: string;
  name: string;
  price_cents: number;
  quantity: number;
  inventory_remaining: number;
  has_inventory: boolean;
}

export const getProductMaxQuantity = (product: Pick<CartItem, "has_inventory" | "inventory_remaining">) => {
  return product.has_inventory ? product.inventory_remaining : GLOBAL_MAX_QUANTITY;
};

export const getValidQuantity = (product: Pick<CartItem, "has_inventory" | "inventory_remaining">, quantity: number) => {
  const max = getProductMaxQuantity(product);
  return Math.max(0, Math.min(quantity, max));
};

export const isLimitReached = (product: Pick<CartItem, "has_inventory" | "inventory_remaining">, currentQuantity: number) => {
  return currentQuantity >= getProductMaxQuantity(product);
};

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
  addItem: (editionId: string, product: Omit<CartItem, "quantity">, quantity: number) => {
    cartStore.setState((prev) => {
      const currentItems = prev.carts[editionId] ?? [];
      const existing = currentItems.find((i) => i.id === product.id);

      let newItems;
      if (existing) {
        const newQuantity = getValidQuantity(product, existing.quantity + quantity);

        newItems = currentItems.map((i) =>
          i.id === product.id ? { ...i, quantity: newQuantity } : i
        );
      } else {
        const finalQuantity = getValidQuantity(product, quantity);
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
            return { ...i, quantity: getValidQuantity(i, quantity) };
          }
          return i;
        }),
      },
    }));
  },
  replaceCart: (editionId: string, items: CartItem[]) => {
    cartStore.setState((prev) => ({
      ...prev,
      carts: {
        ...prev.carts,
        [editionId]: items,
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
