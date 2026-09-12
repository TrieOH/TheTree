import { createSignal } from "solid-js";
import { appLocalStorage } from "@/shared/lib/browser-storage";

export type CartItem = {
  id: string;
  type: "ticket" | "product";
  name: string;
  price_cents: number;
  quantity: number;
  stock: number | null;
};

export type Carts = Record<string, CartItem[]>;
export const CART_STORAGE_KEY = "univents-cart";

const isValidCartItem = (item: unknown): item is CartItem =>
  typeof item === "object" &&
  item !== null &&
  typeof (item as CartItem).id === "string" &&
  ((item as CartItem).type === "ticket" || (item as CartItem).type === "product") &&
  typeof (item as CartItem).name === "string" &&
  typeof (item as CartItem).price_cents === "number" &&
  typeof (item as CartItem).quantity === "number";

export function sanitizeCarts(raw: unknown): Carts {
  if (typeof raw !== "object" || raw === null) return {};
  const cartsObj = (raw as { carts?: Record<string, unknown[]> }).carts;
  if (typeof cartsObj !== "object" || cartsObj === null) return {};

  const sanitized: Carts = {};
  for (const [editionId, items] of Object.entries(cartsObj)) {
    if (Array.isArray(items)) {
      sanitized[editionId] = items.filter(isValidCartItem);
    }
  }
  return sanitized;
}

function initialCarts(): Carts {
  const raw = appLocalStorage.getJson<{ carts?: Carts } | null>(
    CART_STORAGE_KEY,
    null,
  );
  return sanitizeCarts(raw);
}

const [carts, setCarts] = createSignal<Carts>(initialCarts(), {
  ownedWrite: true,
});

function update(editionId: string, items: CartItem[]) {
  setCarts((current) => {
    const next = { ...current, [editionId]: items };
    appLocalStorage.setJson(CART_STORAGE_KEY, { carts: next });
    return next;
  });
}

export const cart = {
  items: (editionId: string) => carts()[editionId] ?? [],
  add(editionId: string, item: Omit<CartItem, "quantity">) {
    const items = cart.items(editionId);
    const existing = items.find(
      (candidate) => candidate.id === item.id && candidate.type === item.type,
    );
    if (
      item.type === "ticket" &&
      !existing &&
      items.some((i) => i.type === "ticket")
    )
      return;
    const max = item.type === "ticket" ? 1 : (item.stock ?? 999);
    update(
      editionId,
      existing
        ? items.map((candidate) =>
            candidate === existing
              ? {
                  ...candidate,
                  quantity: Math.min(candidate.quantity + 1, max),
                }
              : candidate,
          )
        : [...items, { ...item, quantity: 1 }],
    );
  },
  quantity(editionId: string, item: CartItem, quantity: number) {
    if (quantity <= 0) return cart.remove(editionId, item.id, item.type);
    const max = item.type === "ticket" ? 1 : (item.stock ?? 999);
    update(
      editionId,
      cart
        .items(editionId)
        .map((candidate) =>
          candidate.id === item.id && candidate.type === item.type
            ? { ...candidate, quantity: Math.min(quantity, max) }
            : candidate,
        ),
    );
  },
  remove(editionId: string, id: string, type: CartItem["type"]) {
    update(
      editionId,
      cart
        .items(editionId)
        .filter((item) => item.id !== id || item.type !== type),
    );
  },
  clear: (editionId: string) => update(editionId, []),
  stock(editionId: string, id: string, stock: number | null) {
    update(
      editionId,
      cart
        .items(editionId)
        .flatMap((item) =>
          item.id !== id
            ? [item]
            : stock === 0
              ? []
              : [
                  {
                    ...item,
                    stock,
                    quantity: Math.min(item.quantity, stock ?? 999),
                  },
                ],
        ),
    );
  },
};
