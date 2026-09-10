export const storeStockQueryKey = (editionId: string) =>
  ["store", "stock", editionId] as const;

export const productKeys = {
  byEdition: (editionId: string) => ["products", "edition", editionId] as const,
  variants: (productId: string) => ["products", productId, "variants"] as const,
  stock: storeStockQueryKey,
};
