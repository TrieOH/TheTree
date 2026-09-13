export const purchaseKeys = {
  page: (slug: string, authenticated: boolean) =>
    ["purchases", "page", slug, authenticated] as const,
  mine: () => ["purchases", "mine"] as const,
  detail: (purchaseId: string) => ["purchases", "detail", purchaseId] as const,
  catalog: (purchaseId: string) =>
    ["purchases", "catalog", purchaseId] as const,
  edition: (editionId: string) =>
    ["purchases", "edition", editionId] as const,
};
