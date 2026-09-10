export const purchaseKeys = {
  mine: () => ["purchases", "mine"] as const,
  detail: (purchaseId: string) => ["purchases", "detail", purchaseId] as const,
  catalog: (purchaseId: string) =>
    ["purchases", "catalog", purchaseId] as const,
};
