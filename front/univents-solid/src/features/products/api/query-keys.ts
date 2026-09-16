export const storeStockQueryKey = (editionId: string) =>
  ["store", "stock", editionId] as const;

export const productKeys = {
  all: ["products"] as const,
  storePage: (slug: string, authenticated: boolean) =>
    ["products", "store-page", slug, authenticated] as const,
  byEdition: (editionId: string) => ["products", "edition", editionId] as const,
  detail: (productId: string) => ["products", productId] as const,
  variants: (productId: string) => ["products", productId, "variants"] as const,
  stock: storeStockQueryKey,
};
