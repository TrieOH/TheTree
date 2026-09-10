import { orvalData } from "@trieoh/api-client";
import { listUserBadges } from "@trieoh/univents-api";
import type { BadgeProfileGroups } from "@trieoh/univents-api/schemas";
import { badgeKeys } from "./query-keys";

export const userBadgesQueryOptions = (actorId: string) => ({
  queryKey: badgeKeys.user(actorId),
  queryFn: () =>
    listUserBadges(actorId, { public: true }).then(
      orvalData<BadgeProfileGroups>,
    ),
});
