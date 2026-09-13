import type { QueryClient } from "@trieoh/front-core-solid";
import type { QueryObserverOptions } from "@tanstack/query-core";
import { profileKeys } from "./query-keys";

export type CachedProfile = {
  success: boolean;
  data?: {
    actor_id?: string;
    handle?: string | null;
    pfp_url?: string | null;
    profile?: Record<string, unknown>;
  };
  message?: string;
};

export type LoadProfile = (identifier: string) => Promise<CachedProfile>;

export function profileQueryOptions<T extends LoadProfile>(
  identifier: string | undefined,
  loadProfile: T,
): QueryObserverOptions<CachedProfile | null, Error> {
  return {
    queryKey: profileKeys.detail(identifier),
    enabled: Boolean(identifier),
    queryFn: async () => {
      if (!identifier) return null;
      const response = await loadProfile(identifier);
      return response.success ? response : null;
    },
  };
}

export function getCachedOwnProfile(
  queryClient: QueryClient,
  actorId?: string,
): CachedProfile | undefined {
  if (!actorId) return undefined;

  return queryClient
    .getQueryCache()
    .findAll({ queryKey: profileKeys.details() })
    .map((query) => query.state.data as CachedProfile | undefined)
    .find((profile) => profile?.data?.actor_id === actorId);
}

export function cacheProfile(
  queryClient: QueryClient,
  profile: CachedProfile,
  identifier?: string,
) {
  if (!profile.success || !profile.data) return;

  queryClient.setQueryData(profileKeys.detail(identifier), profile);
  const handle = profile.data.handle;
  if (handle && handle !== identifier) {
    queryClient.setQueryData(profileKeys.detail(handle), profile);
  }
}

export * from "./tab-query";
