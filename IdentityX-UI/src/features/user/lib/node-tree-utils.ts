import type { Permission } from "@/features/permission/model/types";
import type { Node, RoleWithPermissions } from "../model/types";

export const buildRolePermissionsToNodeTree = (roleWithPermissions: RoleWithPermissions[]) => {
  const nodeTree: Node = {
    id: 'root',
    name: {
      receiverName: "User",
      applicationName: "Scope"
    },
    children: []
  };

  roleWithPermissions.forEach(({ role, permissions }) => {
    const roleNode: Node = {
      id: role.id,
      name: role.name,
      children: []
    };

    const permissionsByObject: Record<string, Node> = {};

    permissions.forEach(permission => {
      if (!permissionsByObject[permission.object]) {
        permissionsByObject[permission.object] = {
          id: `${role.id}-${permission.object}`,
          name: permission.object,
          children: []
        };
      }

      permissionsByObject[permission.object].children?.push({
        id: `${role.id}-${permission.object}-${permission.action}`,
        name: permission.action,
      });
    });

    Object.values(permissionsByObject).forEach(objectNode => {
      roleNode.children?.push(objectNode);
    });

    nodeTree.children?.push(roleNode);
  });

  return nodeTree;
}

export const buildDirectPermissionsToNodeTree = (permissions: Permission[]) => {
  const nodeTree: Node = {
    id: 'root',
    name: {
      receiverName: "User",
      applicationName: "Scope"
    },
    children: []
  };

  const permissionsByObject: Record<string, Node> = {};

  permissions.forEach(permission => {
    if (!permissionsByObject[permission.object]) {
      permissionsByObject[permission.object] = {
        id: `direct-${permission.object}`,
        name: permission.object,
        children: []
      };
    }

    permissionsByObject[permission.object].children?.push({
      id: `direct-${permission.id}`,
      name: permission.action,
    });
  });

  Object.values(permissionsByObject).forEach(objectNode => {
    nodeTree.children?.push(objectNode);
  });

  return nodeTree;
}
