export const badgeKeys = {
  all: ["badges"] as const,
  byEdition: (editionId: string) =>
    [...badgeKeys.all, "edition", editionId] as const,
  detail: (templateId: string) =>
    [...badgeKeys.all, "detail", templateId] as const,
  user: (userId: string) => [...badgeKeys.all, "user", userId] as const,
  emissions: (editionId: string) =>
    [...badgeKeys.all, "emissions", editionId] as const,
  print: (editionId: string, emissionIds?: string[]) =>
    [...badgeKeys.all, "print", editionId, emissionIds ?? []] as const,
  printActorEmails: (editionId: string, actorIds: string[]) =>
    [...badgeKeys.all, "print-actor-emails", editionId, actorIds] as const,
  printActorNames: (editionId: string, actorIds: string[]) =>
    [...badgeKeys.all, "print-actor-names", editionId, actorIds] as const,
};
