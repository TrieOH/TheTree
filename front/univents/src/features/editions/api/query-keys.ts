export const editionKeys = {
  all: ["editions"] as const,

  lists: () => [...editionKeys.all, "list"] as const,
  publicLists: () => [...editionKeys.lists(), "public"] as const,
  adminLists: () => [...editionKeys.lists(), "admin"] as const,

  publicListByEvent: (eventId: string) =>
    [...editionKeys.publicLists(), eventId] as const,
  adminListByEvent: (eventId: string) =>
    [...editionKeys.adminLists(), eventId] as const,
  activeByEvent: (eventId: string) =>
    [...editionKeys.publicLists(), eventId, "active"] as const,
  pastByEvent: (eventId: string) =>
    [...editionKeys.publicLists(), eventId, "past"] as const,
  upcomingByEvent: (eventId: string) =>
    [...editionKeys.publicLists(), eventId, "upcoming"] as const,
};
