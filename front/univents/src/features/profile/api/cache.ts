import type { QueryClient } from "@tanstack/react-query";
import type { ActorProfile, ProfileData } from "@trieoh/identityx-sdk-ts";
import { profileKeys } from "./query-keys";

export function syncActorProfileCache(
  queryClient: QueryClient,
  actorId: string,
  profile: ProfileData,
  handle?: string,
) {
  const { pfpUrl, ...profileData } = profile;
  queryClient.setQueriesData<ActorProfile>(
    { queryKey: profileKeys.details() },
    (old) =>
      old?.actor_id === actorId
        ? {
            ...old,
            handle: handle ?? old.handle,
            pfp_url: typeof pfpUrl === "string" ? pfpUrl : null,
            profile: profileData,
          }
        : old,
  );
  void queryClient.invalidateQueries({
    queryKey: profileKeys.certificateNames(),
  });
  void queryClient.invalidateQueries({
    queryKey: profileKeys.displayNameLists(),
  });
}

export function invalidateProfileCaches(queryClient: QueryClient) {
  return queryClient.invalidateQueries({ queryKey: profileKeys.all });
}
