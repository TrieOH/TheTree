import type { Permission } from "@/features/permission/model/types";
import type { Scope } from "@/features/scope/model/types";
import type { Node, RoleWithPermissions } from "../model/types";

export const buildRolePermissionsToNodeTree = (
  roleWithPermissions: RoleWithPermissions[],
  scope: string,
) => {
  const nodeTree: Node = {
    id: 'root',
    name: {
      receiverName: "User",
      applicationName: scope
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

export const buildDirectPermissionsToNodeTree = (
  permissions: Permission[], 
  scope: string
) => {
  const nodeTree: Node = {
    id: 'root',
    name: {
      receiverName: "User",
      applicationName: scope
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

export const buildScopeHierarchyToNodeTree = (scopes: Scope[]) => {
  const rootNode: Node = {
    id: 'null',
    name: "Root",
    children: []
  };

  const scopeMap: Record<string, Node> = {
    'null': rootNode
  };

  // Create all nodes
  scopes.forEach(scope => {
    scopeMap[scope.id] = {
      id: scope.id,
      name: scope.name,
      children: []
    };
  });

  // Link them
  scopes.forEach(scope => {
    const parentId = scope.parent_id && scopeMap[scope.parent_id] ? scope.parent_id : 'null';
    const parentNode = scopeMap[parentId];
    
    const folderName = scope.meta?.folder;
    
    if (folderName) {
      if (!parentNode.children) parentNode.children = [];
      
      let folderNode = parentNode.children.find(c => 
        c.name === folderName && c.id.startsWith(`folder-${parentId}-`)
      );
      
      if (!folderNode) {
        folderNode = {
          id: `folder-${parentId}-${folderName}`,
          name: folderName,
          children: []
        };
        parentNode.children.push(folderNode);
      }
      
      folderNode.children?.push(scopeMap[scope.id]);
    } else parentNode.children?.push(scopeMap[scope.id]);
  });

  return rootNode;
}
