export const storeStockQueryKey = (editionId: string) =>
  ["store", "stock", editionId] as const;

export const productKeys = {
  storePage: (slug: string, authenticated: boolean) =>
    ["products", "store-page", slug, authenticated] as const,
  byEdition: (editionId: string) => ["products", "edition", editionId] as const,
  variants: (productId: string) => ["products", productId, "variants"] as const,
  stock: storeStockQueryKey,
};
