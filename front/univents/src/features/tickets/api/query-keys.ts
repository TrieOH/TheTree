export const ticketKeys = {
  all: ['ticket'] as const,

  lists: () => [...ticketKeys.all, 'list'] as const,

  listByEdition: (editionId: string) =>
    [...ticketKeys.lists(), editionId] as const,

  detail: {
    byId: (id: string) =>
      [...ticketKeys.all, 'detail', 'id', id] as const,
  },
}