import { permission } from "@soramux/node-auth-sdk"
import { env } from "@/env"

const objActivities = permission().object("activities")
  .project(env.VITE_TRIEOH_AUTH_PROJECT_ID)

export const canAnnounceActivity = objActivities.action("announce")
export const canReadActivity = objActivities.action("read")
export const canCreateActivity = objActivities.action("create")
export const canManageActivity = objActivities.action("manage")
export const canPublishActivity = objActivities.action("publish") 