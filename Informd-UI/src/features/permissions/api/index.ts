import { createServerFn } from '@tanstack/react-start';
import { promoteToClientSchema  } from '../model';
import type { PromoteToClientI } from '../model';

import { serverPerm, serverRelationship } from '#/shared/lib/api/server-auth';
import { permission } from '@soramux/node-perm-sdk';

const isSuperAdmin = async (userId: string) => {
  const permB = permission().resource("platform", "global")
    .permission("super_admin").subject("user", userId).build()
  const hasThePerm = await serverPerm.check(permB)
  return hasThePerm.success && hasThePerm.data.permissionship === 'PERMISSIONSHIP_HAS_PERMISSION'
}
/**
 * Server function to promote a user to 'client'.
 * This allows the user to create namespaces and manage forms.
 */
export const promoteUserToClientFn = createServerFn({ method: 'POST' })
  .inputValidator((data: PromoteToClientI) => promoteToClientSchema.parse(data))
  .handler(async ({ data }) => {
    try {
      const hasThePerm = await isSuperAdmin(data.requesterId)
      if(!hasThePerm) {
        return {
          success: false,
          message: `You don't have permission to turn ${data.userId} to client`,
        }
      }
      const res = await serverRelationship.create({
        resourceType: "platform",
        resourceId: "global",
        relation: "client",
        subjectType: "user",
        subjectId: data.userId,
      })
      if(res.success) {
        return {
          success: true,
          message: `User ${data.userId} is now a client`,
        }
      } else {
        return {
          success: false,
          message: res.message || `Failed to promote user ${data.userId} to client`,
        }
      }
      
    } catch (error) {
      console.error('Failed to promote user:', error)
      throw new Error('Failed to process promotion on backend')
    }
  })

  /**
 * Server function to check if the current user can see admin actions.
 */
export const checkSuperAdminPrivilegesFn = createServerFn({ method: 'GET' })
  .inputValidator((userId: string) => userId)
  .handler(async ({ data: userId }) => isSuperAdmin(userId))