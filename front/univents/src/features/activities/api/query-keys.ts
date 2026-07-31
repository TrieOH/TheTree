export const activityKeys = {
  all: ["activities"] as const,

  lists: () => [...activityKeys.all, "list"] as const,
  publicLists: () => [...activityKeys.lists(), "public"] as const,
  adminLists: () => [...activityKeys.lists(), "admin"] as const,

  publicListByEdition: (eventId: string, editionId: string) =>
    [...activityKeys.publicLists(), eventId, editionId] as const,

  adminListByEdition: (eventId: string, editionId: string) =>
    [...activityKeys.adminLists(), eventId, editionId] as const,

  attendanceRecords: (eventId: string, editionId: string, activityId: string) =>
    [
      ...activityKeys.all,
      "attendance-records",
      eventId,
      editionId,
      activityId,
    ] as const,
};
