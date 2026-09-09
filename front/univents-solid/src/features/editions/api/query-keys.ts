export const editionKeys = {
  all: ["editions"] as const,
  publicLists: () => [...editionKeys.all, "list", "public"] as const,
  publicListByEvent: (eventId: string) =>
    [...editionKeys.publicLists(), eventId] as const,
};
