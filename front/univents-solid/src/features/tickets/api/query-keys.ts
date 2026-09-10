export const ticketKeys = {
  byEdition: (editionId: string) => ["tickets", "edition", editionId] as const,
  mine: (editionId: string) => ["tickets", "mine", editionId] as const,
};
