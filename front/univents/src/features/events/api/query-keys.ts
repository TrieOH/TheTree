export const eventKeys = {
  all: ["events"] as const,

  lists: () => [...eventKeys.all, "list"] as const,
  publicLists: () => [...eventKeys.lists(), "public"] as const,
  ownLists: () => [...eventKeys.lists(), "own"] as const,
  joinedLists: () => [...eventKeys.lists(), "joined"] as const,
  members: (eventId: string) => [...eventKeys.all, eventId, "members"] as const,

  details: () => [...eventKeys.all, "detail"] as const,
  detail: {
    publicBySlug: (slug: string) =>
      [...eventKeys.details(), "public", "slug", slug] as const,
  },
};
