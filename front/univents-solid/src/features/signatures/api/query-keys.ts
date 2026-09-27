export const signatureKeys = {
  all: ["signatures"] as const,

  lists: () => [...signatureKeys.all, "list"] as const,
  byEdition: (editionId: string) =>
    [...signatureKeys.lists(), editionId] as const,

  byId: (signatureId: string) =>
    [...signatureKeys.all, "detail", signatureId] as const,

  requests: () => [...signatureKeys.all, "requests"] as const,
  requestsByEdition: (editionId: string) =>
    [...signatureKeys.requests(), "list", editionId] as const,
  requestById: (requestId: string) =>
    [...signatureKeys.requests(), "detail", requestId] as const,
};
