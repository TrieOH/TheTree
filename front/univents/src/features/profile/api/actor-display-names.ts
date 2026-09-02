import { useQuery } from "@tanstack/react-query";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { getActorEmailsServerFn } from "@/features/events/server";
import { asUniventsProfile, profileDisplayName } from "../model/profile-data";
import { profileKeys } from "./query-keys";

export function useActorDisplayNames(actorIds: string[]) {
  const { auth } = useAuth();
  const ids = [...new Set(actorIds)];
  return useQuery({
    queryKey: profileKeys.displayNames(ids),
    queryFn: async () => {
      const profiles = [] as Array<readonly [string, string] | null>;
      for (const actorId of ids) {
        const response = await auth.getActorProfile(actorId);
        profiles.push(
          response.success && response.data
            ? [
                actorId,
                profileDisplayName(
                  asUniventsProfile(response.data.profile ?? {}),
                ),
              ]
            : null,
        );
      }
      const emails = await getActorEmailsServerFn({ data: { actorIds: ids } });
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
