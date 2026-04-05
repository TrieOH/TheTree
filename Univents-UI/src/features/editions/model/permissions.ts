import { permission } from "@soramux/node-auth-sdk"
import { env } from "@/env"

const objEditions = permission().object("editions")
  .project(env.VITE_TRIEOH_AUTH_PROJECT_ID)

export const canAnnounceEdition = objEditions.action("announce")
export const canReadEdition = objEditions.action("read")
export const canCreateEdition = objEditions.action("create")