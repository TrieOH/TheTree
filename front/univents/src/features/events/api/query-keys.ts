export const eventKeys = {
  all: ['events'] as const,

  lists: () => [...eventKeys.all, 'list'] as const,
  publicLists: () => [...eventKeys.lists(), 'public'] as const,
  ownLists: () => [...eventKeys.lists(), 'own'] as const,
  joinedLists: () => [...eventKeys.lists(), 'joined'] as const,
  members: (eventId: string) => [...eventKeys.all, eventId, 'members'] as const,
  // byEventId: (eventId: string) => [...eventKeys.all, eventId] as const,

  // editions: (eventId: string) => [...eventKeys.byEventId(eventId), 'editions'] as const,
  // detail: (eventId: string) => [...eventKeys.byEventId(eventId), 'detail'] as const,
}
