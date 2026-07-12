export const editionKeys = {
  all: ['editions'] as const,

  lists: () => [...editionKeys.all, 'list'] as const,
  publicLists: () => [...editionKeys.lists(), 'public'] as const,
  adminLists: () => [...editionKeys.lists(), 'admin'] as const,

  publicListByEvent: (eventId: string) => [...editionKeys.publicLists(), eventId] as const,
  adminListByEvent: (eventId: string) => [...editionKeys.adminLists(), eventId] as const,
}
