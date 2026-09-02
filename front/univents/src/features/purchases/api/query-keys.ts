export const purchaseQueryKeys = {
  all: ["purchases"] as const,
  details: () => [...purchaseQueryKeys.all, "detail"] as const,
  detail: (purchaseId: string) =>
    [...purchaseQueryKeys.details(), purchaseId] as const,
  lists: () => [...purchaseQueryKeys.all, "list"] as const,
  mine: () => [...purchaseQueryKeys.lists(), "mine"] as const,
  edition: (editionId: string) =>
    [...purchaseQueryKeys.lists(), "edition", editionId] as const,
};
