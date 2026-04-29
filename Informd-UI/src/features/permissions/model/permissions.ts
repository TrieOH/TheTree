
import type { Relationship } from "@soramux/node-perm-sdk"

export const clientPermModel = (userId: string) => [
  {
    operation: 'OPERATION_TOUCH' as Relationship.RelationshipOperation,
    resourceType: 'user',
    resourceId: userId,
    relation: 'parent_platform',
    subjectType: 'platform',
    subjectId: 'global',
  },
  {
    operation: 'OPERATION_TOUCH' as Relationship.RelationshipOperation,
    resourceType: 'user',
    resourceId: userId,
    relation: 'self',
    subjectType: 'user',
    subjectId: userId,
  },
]
