export const programKeys = {
  page: (slug: string, authenticated: boolean) =>
    ["programs", "page", slug, authenticated] as const,
  byEdition: (editionId: string) => ["programs", "edition", editionId] as const,
  occurrences: (editionId: string) =>
    ["programs", "occurrences", editionId] as const,
  occurrence: (occurrenceId: string) =>
    ["programs", "occurrence", occurrenceId] as const,
  participants: (occurrenceId: string) =>
    ["programs", "participants", occurrenceId] as const,
  mine: (editionId: string) => ["programs", "mine", editionId] as const,
};
