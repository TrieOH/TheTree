import { useQuery } from "@tanstack/react-query";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { getActorEmailsServerFn } from "@/features/events/api/actor-emails";
import { asUniventsProfile, profileDisplayName } from "../model/profile-data";
import { profileKeys } from "./query-keys";

export function useActorDisplayNames(actorIds: string[]) {
  const { auth } = useAuth();
  const ids = [...new Set(actorIds)];
  return useQuery({
    queryKey: profileKeys.displayNames(ids),
    queryFn: async () => {
      const [profiles, emails] = await Promise.all([
        Promise.all(
          ids.map(async (actorId) => {
            const response = await auth.getActorProfile(actorId);
            if (!response.success || !response.data) return null;
            return [
              actorId,
              profileDisplayName(
                asUniventsProfile(response.data.profile ?? {}),
              ),
            ] as const;
          }),
        ),
        getActorEmailsServerFn({ data: { actorIds: ids } }),
      ]);
      return Object.fromEntries(
        ids.map((id) => [
          id,
          profiles.find((entry) => entry?.[0] === id)?.[1] ?? emails[id] ?? id,
        ]),
      );
    },
    enabled: ids.length > 0,
  });
}
