export function upsertById<T extends { id: string }>(
  list: T[] | undefined,
  item: T,
) {
  if (!list) return list;
  const index = list.findIndex((candidate) => candidate.id === item.id);
  if (index === -1) return [...list, item];
  const next = [...list];
  next[index] = item;
  return next;
}
