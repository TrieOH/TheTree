export const eventKeys = {
  all: ["events"] as const,
  lists: () => [...eventKeys.all, "list"] as const,
  publicLists: () => [...eventKeys.lists(), "public"] as const,
  /** Events the current actor owns or is a member of (admin list). */
  ownedLists: () => [...eventKeys.lists(), "owned"] as const,
  joinedLists: () => [...eventKeys.lists(), "joined"] as const,
  details: () => [...eventKeys.all, "detail"] as const,
  detail: {
    publicBySlug: (slug: string) =>
      [...eventKeys.details(), "public", "slug", slug] as const,
  },
};
