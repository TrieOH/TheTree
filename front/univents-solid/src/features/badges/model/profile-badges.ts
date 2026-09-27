import type {
  BadgeProfileBadge,
  BadgeProfileGroups,
} from "@trieoh/univents-api/schemas";

export interface ProfileBadgeGroups<T = BadgeProfileBadge> {
  attendant: { current: T[]; past: T[] };
  staff: { current: T[]; past: T[] };
}

export const allProfileBadges = <T = BadgeProfileBadge>(
  groups?: ProfileBadgeGroups<T> | BadgeProfileGroups | null,
): T[] => {
  if (!groups) return [];
  return [
    ...((groups.attendant?.current ?? []) as unknown as T[]),
    ...((groups.attendant?.past ?? []) as unknown as T[]),
    ...((groups.staff?.current ?? []) as unknown as T[]),
    ...((groups.staff?.past ?? []) as unknown as T[]),
  ];
};
