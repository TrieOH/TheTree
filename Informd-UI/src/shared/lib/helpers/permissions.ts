import { permission } from "@soramux/node-perm-sdk"
import { serverPerm } from "../api/server-auth"

export const isSuperAdmin = async (userId: string): Promise<boolean> => {
  const permB = permission()
    .resource("platform", "global")
    .permission("super_admin")
    .subject("user", userId)
    .build()

  const res = await serverPerm.check(permB)
  return res.success && res.data.permissionship === 'PERMISSIONSHIP_HAS_PERMISSION'
}