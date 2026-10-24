import { queryError } from "@trieoh/front-core-solid";
import type { ActorProfile } from "../model/profile-data";
import { profileKeys } from "./query-keys";

export type CachedProfile = {
  success: boolean;
  data?: ActorProfile;
  message?: string;
  code?: number;
};

export type LoadProfile = (identifier: string) => Promise<CachedProfile>;

export interface ProfileApi {
  getActorProfile: LoadProfile;
  getProfileByHandle: LoadProfile;
}

export function loadProfile(identifier: string, auth: ProfileApi) {
  const isActorId = /^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(identifier);
  return isActorId ? auth.getActorProfile(identifier) : auth.getProfileByHandle(identifier);
}

export function profileQueryOptions(
  identifier: string | undefined,
  load: LoadProfile,
  { allowMissing = false }: { allowMissing?: boolean } = {},
) {
  return {
    queryKey: profileKeys.detail(identifier),
    enabled: Boolean(identifier),
    queryFn: async (): Promise<ActorProfile | null> => {
      if (!identifier) return null;
      const response = await load(identifier);
      if (response.success && response.data) return response.data;
      if (allowMissing && (response.success || response.code === 404)) return null;
      throw queryError(response.message || "Não foi possível carregar este perfil.", response.code);
    },
  };
}

export function profileDetailQueryOptions(identifier: string | undefined, auth: ProfileApi) {
  return profileQueryOptions(identifier, (id) => loadProfile(id, auth));
}

/** Guards accept a missing profile so they can redirect to setup. */
export function ownProfileQueryOptions(actorId: string, auth: Pick<ProfileApi, "getActorProfile">) {
  return profileQueryOptions(actorId, (id) => auth.getActorProfile(id), { allowMissing: true });
}
