export const eventKeys = {
  all: ["events"] as const,
  lists: () => [...eventKeys.all, "list"] as const,
  publicLists: () => [...eventKeys.lists(), "public"] as const,
  details: () => [...eventKeys.all, "detail"] as const,
  detail: {
    publicBySlug: (slug: string) =>
      [...eventKeys.details(), "public", "slug", slug] as const,
  },
};
