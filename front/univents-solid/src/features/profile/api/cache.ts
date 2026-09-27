import type { QueryClient } from "@trieoh/front-core-solid";
import type { ActorProfile, ProfileData } from "../model/profile-data";
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
    (old: ActorProfile | undefined) =>
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

export function cacheProfile(
  queryClient: QueryClient,
  profile: { success: boolean; data?: ActorProfile },
  identifier?: string,
) {
  if (!profile.success || !profile.data) return;

  const data = profile.data;
  if (identifier) {
    queryClient.setQueryData(profileKeys.detail(identifier), data);
  }
  if (data.actor_id) {
    queryClient.setQueryData(profileKeys.detail(data.actor_id), data);
  }
  const handle = data.handle;
  if (handle && handle !== identifier) {
    queryClient.setQueryData(profileKeys.detail(handle), data);
  }
}

export function getCachedOwnProfile(
  queryClient: QueryClient,
  actorId?: string,
): ActorProfile | undefined {
  if (!actorId) return undefined;

  return queryClient
    .getQueryCache()
    .findAll({ queryKey: profileKeys.details() })
    .map((query) => query.state.data as ActorProfile | { data?: ActorProfile } | undefined)
    .map((d) => (d && "actor_id" in d ? d : d?.data))
    .find((p) => p?.actor_id === actorId);
}

export function invalidateProfileCaches(queryClient: QueryClient) {
  return queryClient.invalidateQueries({ queryKey: profileKeys.all });
}
