export const badgeKeys = {
  all: ["badges"] as const,
  byEdition: (editionId: string) =>
    [...badgeKeys.all, "edition", editionId] as const,
  detail: (templateId: string) =>
    [...badgeKeys.all, "detail", templateId] as const,
};
