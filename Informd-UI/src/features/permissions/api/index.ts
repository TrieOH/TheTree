import { createServerFn } from '@tanstack/react-start';
import { promoteToClientSchema  } from '../model';
import type { PromoteToClientI } from '../model';

import { serverRelationship } from '#/shared/lib/api/server-auth';
import { isSuperAdmin } from '#/shared/lib/helpers/permissions';
import { clientPermModel } from '../model/permissions';


/**
 * Server function to promote a user to 'client'.
 * This allows the user to create namespaces and manage forms.
 */
export const promoteUserToClientFn = createServerFn({ method: 'POST' })
  .inputValidator((data: PromoteToClientI) => promoteToClientSchema.parse(data))
  .handler(async ({ data }) => {
    if(!(await isSuperAdmin(data.requesterId))) {
      return {
        success: false,
        message: `You don't have permission to turn ${data.userId} to client`,
      }
    }
    const res = await serverRelationship.batchWrite(clientPermModel(data.userId))
    if (!res.success) {
      return {
        success: false,
        message: `Failed to promote user ${data.userId} to client`,
      }
    }
    return {
      success: true,
      message: `User ${data.userId} is now a client`,
    }
  })

  /**
 * Server function to check if the current user can see admin actions.
 */
export const checkSuperAdminPrivilegesFn = createServerFn({ method: 'GET' })
  .inputValidator((userId: string) => userId)
  .handler(async ({ data: userId }) => isSuperAdmin(userId))