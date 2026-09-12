export const profileKeys = {
  detail: (actorId?: string) => ["profile", "detail", actorId] as const,
  tab: (tab: string, actorId: string, ownProfile: boolean) =>
    ["profile", "tab", tab, actorId, ownProfile] as const,
};
