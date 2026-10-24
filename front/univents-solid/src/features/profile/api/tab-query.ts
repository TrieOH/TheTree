import type { QueryClient } from "@trieoh/front-core-solid";
import { userBadgesQueryOptions } from "@/features/badges/api";
import { myCertificationsQueryOptions } from "@/features/certifications/api";
import { myPurchasesQueryOptions } from "@/features/purchases/api";
import type { ActorProfile } from "../model/profile-data";
import type { ProfileTab } from "../model/profile-tabs";
import { profileDetailQueryOptions, type ProfileApi } from "./queries";
import { profileKeys } from "./query-keys";

/** Prepare the same cache entries consumed by the collection components. */
export async function prefetchProfileTab(
  client: QueryClient,
  options: { tab: ProfileTab; identifier: string; viewerId?: string; auth: ProfileApi },
) {
  const { tab, identifier, viewerId, auth } = options;
  if (tab === "about") return;

  if (tab === "badges") {
    await client.prefetchQuery(profileDetailQueryOptions(identifier, auth));
    const profile = client.getQueryData<ActorProfile>(profileKeys.detail(identifier));
    if (profile?.actor_id) {
      await client.prefetchQuery(userBadgesQueryOptions(profile.actor_id));
    }
    return;
  }

  const viewerProfile = viewerId
    ? client.getQueryData<ActorProfile>(profileKeys.detail(viewerId))
    : undefined;
  const own = viewerId && (identifier === viewerId || identifier === viewerProfile?.handle);
  if (!own) return;

  if (tab === "certificates") {
    await client.prefetchQuery(myCertificationsQueryOptions());
  } else {
    await client.prefetchQuery(myPurchasesQueryOptions());
  }
}
