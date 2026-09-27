import type { QueryObserverOptions } from "@tanstack/query-core";
import type {
  BadgeProfileGroups,
  Certification,
} from "@trieoh/univents-api/schemas";
import { userBadgesQueryOptions } from "@/features/badges/api";
import { myCertificationsQueryOptions } from "@/features/certifications/api";
import { profileKeys } from "./query-keys";

export type ProfileTab = "badges" | "certificates" | "purchases";
export type ProfileTabData = BadgeProfileGroups | Certification[] | unknown[];

export function profileTabQueryOptions(
  tab: ProfileTab,
  actorId: string,
  ownProfile: boolean,
): QueryObserverOptions<ProfileTabData, Error> {
  return {
    queryKey: profileKeys.tab(tab, actorId, ownProfile),
    enabled: tab === "badges" || (ownProfile && tab === "certificates"),
    queryFn: () => {
      if (tab === "badges") {
        return userBadgesQueryOptions(actorId).queryFn();
      }
      if (!ownProfile) return [];
      if (tab === "certificates") {
        return myCertificationsQueryOptions().queryFn();
      }
      return [];
    },
  };
}
