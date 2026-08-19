export const purchaseQueryKeys = {
  all: ["purchases"] as const,
  detail: (purchaseId: string) => ["purchases", purchaseId] as const,
  edition: (editionId: string) => ["purchases", "edition", editionId] as const,
};
