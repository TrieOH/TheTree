export const profileKeys = {
  all: ["profiles"] as const,
  details: () => [...profileKeys.all, "detail"] as const,
  detail: (actorId: string) => [...profileKeys.details(), actorId] as const,
  certificateNames: () => [...profileKeys.all, "certificate-name"] as const,
  certificateName: (actorId: string) =>
    [...profileKeys.certificateNames(), actorId] as const,
  displayNameLists: () => [...profileKeys.all, "display-names"] as const,
  displayNames: (actorIds: string[]) =>
    [...profileKeys.displayNameLists(), actorIds] as const,
};
