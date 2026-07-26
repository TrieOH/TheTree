export const signatureKeys = {
  all: ["signatures"] as const,

  lists: () => [...signatureKeys.all, "list"] as const,
  byEdition: (eventId: string, editionId: string) =>
    [...signatureKeys.lists(), eventId, editionId] as const,

  byId: (eventId: string, editionId: string, signatureId: string) =>
    [...signatureKeys.all, eventId, editionId, signatureId] as const,
};
