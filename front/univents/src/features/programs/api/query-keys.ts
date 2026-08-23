export const programKeys = {
  all: ["programs"] as const,
  byEdition: (editionId: string) => ["programs", editionId] as const,
  occurrences: (editionId: string) =>
    ["program-occurrences", editionId] as const,
  myParticipations: (editionId: string) =>
    ["program-participations", "mine", editionId] as const,
  participants: (occurrenceId: string) =>
    ["program-participants", occurrenceId] as const,
};
