import { permission } from "@soramux/node-perm-sdk"

export const canCreateEvent = permission()
  .resource("platform", "global")
  .permission("create_events")

export const canPublishEvent = (eventId: string) => permission()
  .resource("event", eventId)
  .permission("publish")

export const canEditEvent = (eventId: string) => permission()
  .resource("event", eventId)
  .permission("edit")


// export const canPublishEvent = objEvent.permission("publish")
// export const canEditEvent = objEvent.permission("edit")

// authz.Subject("user", sub.ID),
// authz.Permission("create_events"),
// authz.Resource("platform", "global"),