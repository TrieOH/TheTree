import type { QueryObserverOptions } from "@tanstack/query-core";
import { queryError } from "@trieoh/front-core-solid";
import type { ActorProfile } from "../model/profile-data";
import { profileKeys } from "./query-keys";

export type CachedProfile = {
  success: boolean;
  data?: ActorProfile;
  message?: string;
  code?: number;
};

export type LoadProfile = (identifier: string) => Promise<{
  success: boolean;
  data?: ActorProfile;
  message?: string;
  code?: number;
}>;

export function profileDetailQueryOptions(
  identifier: string | undefined,
  auth: {
    getActorProfile: (id: string) => Promise<{ success: boolean; data?: ActorProfile; message?: string; code?: number }>;
    getProfileByHandle: (handle: string) => Promise<{ success: boolean; data?: ActorProfile; message?: string; code?: number }>;
  },
): QueryObserverOptions<ActorProfile | null, Error> {
  return {
    queryKey: profileKeys.detail(identifier),
    enabled: Boolean(identifier),
    queryFn: async () => {
      if (!identifier) return null;
      const isActorId =
        /^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(identifier);
      const response = isActorId
        ? await auth.getActorProfile(identifier)
        : await auth.getProfileByHandle(identifier);
      if (!response.success || !response.data) {
        throw queryError(
          response.message || "Não foi possível carregar este perfil.",
          response.code,
        );
      }
      return response.data;
    },
  };
}

export function profileQueryOptions<T extends LoadProfile>(
  identifier: string | undefined,
  loadProfile: T,
): QueryObserverOptions<ActorProfile | null, Error> {
  return {
    queryKey: profileKeys.detail(identifier),
    enabled: Boolean(identifier),
    queryFn: async () => {
      if (!identifier) return null;
      const response = await loadProfile(identifier);
      if (!response.success || !response.data) {
        throw queryError(
          response.message || "Não foi possível carregar este perfil.",
          response.code,
        );
      }
      return response.data;
    },
  };
}

export * from "./cache";
export * from "./query-keys";
export * from "./tab-query";
