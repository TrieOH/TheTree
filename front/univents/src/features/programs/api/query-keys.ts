export const programKeys = {
  all: ["programs"] as const,
  byEdition: (editionId: string) => ["programs", editionId] as const,
  occurrences: (editionId: string) =>
    ["program-occurrences", editionId] as const,
};
