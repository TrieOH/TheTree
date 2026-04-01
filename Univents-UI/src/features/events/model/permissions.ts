import { permission } from "@soramux/node-auth-sdk"
import { env } from "@/env"

const objEvent = permission().object("events")
  .project(env.VITE_TRIEOH_AUTH_PROJECT_ID)

export const canPublishEvent = objEvent.action("publish")
export const canCreateEvent = objEvent.action("create")
export const canEditEvent = objEvent.action("edit")