import { createSignal } from "solid-js";

export type CartItem = {
  id: string;
  type: "ticket" | "product";
  name: string;
  price_cents: number;
  quantity: number;
  stock: number | null;
};

type Carts = Record<string, CartItem[]>;
const storageKey = "univents-cart";

function initialCarts(): Carts {
  if (typeof window === "undefined") return {};
  try {
    return (
      (
        JSON.parse(localStorage.getItem(storageKey) ?? "{}") as {
          carts?: Carts;
        }
      ).carts ?? {}
    );
  } catch {
    return {};
  }
}

const [carts, setCarts] = createSignal<Carts>(initialCarts(), {
  ownedWrite: true,
});

function update(editionId: string, items: CartItem[]) {
  setCarts((current) => {
    const next = { ...current, [editionId]: items };
    localStorage.setItem(storageKey, JSON.stringify({ carts: next }));
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
