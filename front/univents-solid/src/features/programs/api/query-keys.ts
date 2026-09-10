export const programKeys = {
  byEdition: (editionId: string) => ["programs", "edition", editionId] as const,
  occurrences: (editionId: string) =>
    ["programs", "occurrences", editionId] as const,
  occurrence: (occurrenceId: string) =>
    ["programs", "occurrence", occurrenceId] as const,
  mine: (editionId: string) => ["programs", "mine", editionId] as const,
};
