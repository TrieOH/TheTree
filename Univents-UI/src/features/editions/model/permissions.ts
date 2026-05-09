import { permission } from "@soramux/node-perm-sdk"

export const canReadAdminEdition = (eventId: string) => permission()
  .resource("event", eventId)
  .permission("view_editions")

export const canCreateEdition = (eventId: string) => permission()
  .resource("event", eventId)
  .permission("create_editions")

export const canAnnounceEdition = (editionId: string) => permission()
  .resource("edition", editionId)
  .permission("announce")

export const canConnectPayment = (editionId: string) => permission()
  .resource("edition", editionId)
  .permission("connect_payments")

export const canDisconnectPayment = (editionId: string) => permission()
  .resource("edition", editionId)
  .permission("disconnect_payments")