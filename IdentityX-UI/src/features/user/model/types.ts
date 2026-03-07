import type { Permission } from "@/features/permission/model/types";
import type { Role } from "@/features/role/model/types";
import type { JsonValue } from "@/shared/model/types";


export interface User {
  id: string;
  email: string;
  is_active: boolean;
  is_verified: boolean;
  user_type: string;
  project_id: string;
  metadata: Record<string, JsonValue>
  last_login_at: string;
  verified_at: string;
  created_at: string;
  updated_at: string;
}

export type RoleWithPermissions = {
  role: Role;
  permissions: Permission[];
};

export interface NodeCustomName {
  receiverName: string;
  applicationName: string;
}

export type NodeType = 'scope' | 'role' | 'object' | 'action' | 'inherited' | 'direct' | 'folder'
  | 'perm-folder' | 'role-folder';

export interface InheritedPermissionSource {
  object: string;
  action: string;
  sourceScope: string;
  sourceRole?: string;
  id: string;
  scopeId: string | null;
  icon?: string;
  color?: string;
}

export interface Node {
  id: string
  name: string | NodeCustomName
  type: NodeType
  icon?: string
  color?: string
  children?: Node[]
  
  roleId?: string
  permissionId?: string
  scopeId?: string | null
  isInherited?: boolean
  sourceScope?: string
  isFromRole?: boolean
  sources?: InheritedPermissionSource[]
}