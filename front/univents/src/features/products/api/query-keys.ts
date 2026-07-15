export const productKeys = {
  all: ['products'] as const,

  lists: () => [...productKeys.all, 'list'] as const,
  publicLists: () => [...productKeys.lists(), 'public'] as const,
  adminLists: () => [...productKeys.lists(), 'admin'] as const,

  publicListByEdition: (eventId: string, editionId: string) =>
    [...productKeys.publicLists(), eventId, editionId] as const,

  adminListByEdition: (eventId: string, editionId: string) =>
    [...productKeys.adminLists(), eventId, editionId] as const,
}
