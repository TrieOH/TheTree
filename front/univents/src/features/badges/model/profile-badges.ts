interface ProfileBadgeGroups<T> {
  attendant: { current: T[]; past: T[] };
  staff: { current: T[]; past: T[] };
}

export const allProfileBadges = <T>(groups: ProfileBadgeGroups<T>): T[] => [
  ...groups.attendant.current,
  ...groups.attendant.past,
  ...groups.staff.current,
  ...groups.staff.past,
];
