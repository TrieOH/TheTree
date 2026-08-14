import { Store } from "@tanstack/react-store";

export const GLOBAL_MAX_QUANTITY = 999;
export type CartItemType = "ticket" | "product" | "activity";

export interface CartItem {
  id: string;
  type: CartItemType;
  name: string;
  price_cents: number;
  quantity: number;
  inventory_remaining: number;
  has_inventory: boolean;
}

export const getProductMaxQuantity = (
  product: Pick<
    CartItem,
    "type" | "has_inventory" | "inventory_remaining"
  >,
) => {
  if (product.type === "ticket") return 1;
  return product.has_inventory
    ? product.inventory_remaining
    : GLOBAL_MAX_QUANTITY;
};

export const getValidQuantity = (
  product: Pick<
    CartItem,
    "type" | "has_inventory" | "inventory_remaining"
  >,
  quantity: number,
) => {
  const max = getProductMaxQuantity(product);
  return Math.max(0, Math.min(quantity, max));
};

export const isLimitReached = (
  product: Pick<
    CartItem,
    "type" | "has_inventory" | "inventory_remaining"
  >,
  currentQuantity: number,
) => {
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
  addItem: (
    editionId: string,
    product: Omit<CartItem, "quantity">,
    quantity: number,
  ) => {
    cartStore.setState((prev) => {
      const currentItems = prev.carts[editionId] ?? [];
      const existing = currentItems.find(
        (i) => i.id === product.id && i.type === product.type,
      );
      if (
        product.type === "ticket" &&
        !existing &&
        currentItems.some((item) => item.type === "ticket")
      ) {
        return prev;
      }

      let newItems: CartItem[];
      if (existing) {
        const newQuantity = getValidQuantity(
          product,
          existing.quantity + quantity,
        );

        newItems = currentItems.map((i) =>
          i.id === product.id && i.type === product.type
            ? { ...i, quantity: newQuantity }
            : i,
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
  removeItem: (editionId: string, id: string, type?: CartItemType) => {
    cartStore.setState((prev) => ({
      ...prev,
      carts: {
        ...prev.carts,
        [editionId]: (prev.carts[editionId] ?? []).filter(
          (i) => i.id !== id || (type !== undefined && i.type !== type),
        ),
      },
    }));
  },
  updateQuantity: (
    editionId: string,
    id: string,
    quantity: number,
    type?: CartItemType,
  ) => {
    if (quantity <= 0) {
      cartActions.removeItem(editionId, id, type);
      return;
    }
    cartStore.setState((prev) => ({
      ...prev,
      carts: {
        ...prev.carts,
        [editionId]: (prev.carts[editionId] ?? []).map((i) => {
          if (i.id === id && (type === undefined || i.type === type)) {
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
  updateStock: (editionId: string, id: string, stock: number | null) => {
    cartStore.setState((prev) => ({
      ...prev,
      carts: {
        ...prev.carts,
        [editionId]: (prev.carts[editionId] ?? []).flatMap((item) => {
          if (item.id !== id) return item;
          if (stock === 0) return [];
          return [
            {
              ...item,
              has_inventory: stock !== null,
              inventory_remaining: stock ?? GLOBAL_MAX_QUANTITY,
              quantity:
                stock === null ? item.quantity : Math.min(item.quantity, stock),
            },
          ];
        }),
      },
    }));
  },
};
