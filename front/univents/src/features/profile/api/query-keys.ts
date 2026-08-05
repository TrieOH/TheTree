export const profileKeys = {
  all: ["profiles"] as const,
  detail: (actorId: string) => [...profileKeys.all, "detail", actorId] as const,
  certificateName: (actorId: string) =>
    [...profileKeys.all, "certificate-name", actorId] as const,
};
