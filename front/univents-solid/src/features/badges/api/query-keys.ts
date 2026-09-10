export const badgeKeys = {
  user: (actorId: string) => ["badges", "user", actorId] as const,
};
