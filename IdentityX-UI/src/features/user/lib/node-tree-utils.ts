import type { Permission } from "@/features/permission/model/types";
import type { Scope } from "@/features/scope/model/types";
import type { Role } from "@/features/role/model/types";
import type { Node, RoleWithPermissions, NodeType, InheritedPermissionSource } from "../model/types";

export const buildRolePermissionsToNodeTree = (
  roleWithPermissions: RoleWithPermissions[],
  scope: string,
): Node => {
  const nodeTree: Node = {
    id: 'root',
    type: 'scope',
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
      type: 'role',
      roleId: role.id,
      scopeId: role.scope_id,
      icon: role.meta?.icon,
      color: role.meta?.color,
      children: []
    };

    const permissionsByObject: Record<string, Node> = {};

    permissions.forEach(permission => {
      if (!permissionsByObject[permission.object]) {
        permissionsByObject[permission.object] = {
          id: `${role.id}-${permission.object}`,
          name: permission.object,
          type: 'object',
          children: []
        };
      }

      permissionsByObject[permission.object].children?.push({
        id: `${role.id}-${permission.object}-${permission.action}`,
        name: permission.action,
        type: 'action',
        permissionId: permission.id,
        isFromRole: true,
        icon: permission.meta?.icon,
        color: permission.meta?.color,
        scopeId: role.scope_id,
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
): Node => {
  const nodeTree: Node = {
    id: 'root',
    type: 'scope',
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
        type: 'object',
        children: []
      };
    }

    permissionsByObject[permission.object].children?.push({
      id: `direct-${permission.id}`,
      name: permission.action,
      type: 'action',
      permissionId: permission.id,
      icon: permission.meta?.icon,
      color: permission.meta?.color,
    });
  });

  Object.values(permissionsByObject).forEach(objectNode => {
    nodeTree.children?.push(objectNode);
  });

  return nodeTree;
}

export const buildScopeHierarchyToNodeTree = (scopes: Scope[]): Node => {
  const rootNode: Node = {
    id: 'null',
    name: "Root",
    type: 'scope',
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
      type: 'scope',
      scopeId: scope.id,
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
          type: 'folder',
          children: []
        };
        parentNode.children.push(folderNode);
      }
      
      folderNode.children?.push(scopeMap[scope.id]);
    } else parentNode.children?.push(scopeMap[scope.id]);
  });

  return rootNode;
}

interface InternalAccessNode {
  id: string;
  name: string;
  type: NodeType;
  children: Node[];
  directRoles: (Role & { permissions: Permission[] })[];
  inheritedRoles: { role: Role & { permissions: Permission[] }; sourceScope: string }[];
  directPermissions: (Permission & { scopeId: string | null; isFromRole?: boolean })[];
  inheritedPermissions: InheritedPermissionSource[];
}

export const buildFullAccessControlNodeTree = (
  allScopes: Scope[],
  rolesWithPerms: (Role & { permissions: Permission[] })[],
  allPermissionsWithScope: (Permission & { scopeId: string | null })[],
): Node => {
  const getScopeDisplayName = (scope: Scope | { id: string | null; name: string }) => {
    if (!scope.id || scope.id === 'root-node-id') return "Root";
    return scope.name;
  };

  const getAncestors = (scopeId: string | null): (string | null)[] => {
    const ancestors: (string | null)[] = [];
    let currentId = scopeId;
    while (currentId) {
      const scope = allScopes.find(s => s.id === currentId);
      if (!scope?.parent_id) break;
      ancestors.push(scope.parent_id);
      currentId = scope.parent_id;
    }
    if (scopeId !== null) ancestors.push(null);
    return ancestors;
  };

  const nodesMap: Record<string, InternalAccessNode> = {};

  const rootNodeId = 'root-node-id';
  nodesMap[rootNodeId] = {
    id: rootNodeId,
    name: "Root",
    type: 'scope',
    children: [],
    directRoles: [],
    inheritedRoles: [],
    directPermissions: [],
    inheritedPermissions: []
  };

  allScopes.forEach(scope => {
    nodesMap[scope.id] = {
      id: scope.id,
      name: scope.name,
      type: 'scope',
      children: [],
      directRoles: [],
      inheritedRoles: [],
      directPermissions: [],
      inheritedPermissions: []
    };
  });

  rolesWithPerms.forEach(role => {
    const targetId = role.scope_id || rootNodeId;
    if (nodesMap[targetId]) nodesMap[targetId].directRoles.push(role);
  });

  allPermissionsWithScope.forEach(perm => {
    const targetId = perm.scopeId || rootNodeId;
    if (nodesMap[targetId]) nodesMap[targetId].directPermissions.push(perm);
  });

  Object.keys(nodesMap).forEach(nodeId => {
    const node = nodesMap[nodeId];
    const realScopeId: string | null = nodeId === rootNodeId ? null : nodeId;
    const ancestors = getAncestors(realScopeId);

    ancestors.forEach(ancestorId => {
      const ancestorScope = ancestorId ? allScopes.find(s => s.id === ancestorId) : { id: null, name: 'Root' };
      const sourceName = ancestorScope ? getScopeDisplayName(ancestorScope) : 'Root';

      rolesWithPerms.forEach(role => {
        if (role.scope_id === ancestorId) {
          node.inheritedRoles.push({ role, sourceScope: sourceName });
          role.permissions.forEach(p => {
            const exists = node.inheritedPermissions.some(ip => ip.object === p.object && ip.action === p.action && ip.sourceScope === sourceName && ip.sourceRole === role.name);
            if (!exists) {
              node.inheritedPermissions.push({
                object: p.object,
                action: p.action,
                sourceScope: sourceName,
                sourceRole: role.name,
                id: p.id,
                scopeId: ancestorId,
                icon: p.meta?.icon,
                color: p.meta?.color
              });
            }
          });
        }
      });

      allPermissionsWithScope.forEach(perm => {
        if (perm.scopeId === ancestorId) {
          const exists = node.inheritedPermissions.some(ip => ip.object === perm.object && ip.action === perm.action && ip.sourceScope === sourceName && !ip.sourceRole);
          if (!exists) {
            node.inheritedPermissions.push({
              object: perm.object,
              action: perm.action,
              sourceScope: sourceName,
              id: perm.id,
              scopeId: ancestorId,
              icon: perm.meta?.icon,
              color: perm.meta?.color
            });
          }
        }
      });
    });

    node.directPermissions = node.directPermissions.filter(dp => {
      return !node.inheritedPermissions.some(ip => ip.object === dp.object && ip.action === dp.action);
    });

    // Check if direct permissions are also provided by direct roles in the same scope
    node.directPermissions = node.directPermissions.map(dp => {
      const existsInRole = node.directRoles.some(role => 
        role.permissions.some((p: Permission) => p.object === dp.object && p.action === dp.action)
      );
      return { ...dp, isFromRole: existsInRole };
    });
  });

  // Now transform nodesMap to final Node structure with folders
  const finalNodes: Record<string, Node> = {};
  const folderNodesMap: Record<string, Node> = {};

  Object.keys(nodesMap).forEach(id => {
    const src = nodesMap[id];
    const node: Node = {
      id: src.id,
      name: src.name,
      type: src.type,
      children: []
    };

    // 1. Inherited
    if (src.inheritedRoles.length > 0 || src.inheritedPermissions.length > 0) {
      const inheritedNode: Node = {
        id: `${id}-inherited`,
        name: "Inherited",
        type: 'inherited',
        children: []
      };

      if (src.inheritedRoles.length > 0) {
        const rolesNode: Node = {
          id: `${id}-inherited-roles`,
          name: "Roles",
          type: 'folder',
          children: src.inheritedRoles.map((ir, idx) => ({
            id: `${id}-inherited-role-${ir.role.id}-${idx}`,
            name: ir.role.name,
            type: 'role',
            icon: ir.role.meta?.icon,
            color: ir.role.meta?.color,
            sourceScope: ir.sourceScope,
            roleId: ir.role.id,
            scopeId: ir.role.scope_id,
            isInherited: true,
            children: buildPermissionsTree(ir.role.permissions, `${id}-inherited-role-${ir.role.id}-${idx}`, true, ir.role.scope_id)
          }))
        };
        inheritedNode.children?.push(rolesNode);
      }

      if (src.inheritedPermissions.length > 0) {
        const permsNode: Node = {
          id: `${id}-inherited-permissions`,
          name: "Permissions",
          type: 'folder',
          children: buildInheritedPermissionsTree(src.inheritedPermissions, `${id}-inherited-permissions`)
        };
        inheritedNode.children?.push(permsNode);
      }

      node.children?.push(inheritedNode);
    }

    // 2. Direct Roles
    if (src.directRoles.length > 0) {
      const rolesNode: Node = {
        id: `${id}-direct-roles`,
        name: "Roles",
        type: 'role-folder',
        children: src.directRoles.map(role => ({
          id: `${id}-direct-role-${role.id}`,
          name: role.name,
          type: 'role',
          icon: role.meta?.icon,
          color: role.meta?.color,
          roleId: role.id,
          scopeId: role.scope_id,
          isInherited: false,
          children: buildPermissionsTree(role.permissions, `${id}-direct-role-${role.id}`, true, role.scope_id)
        }))
      };
      node.children?.push(rolesNode);
    }

    // 3. Direct Permissions
    if (src.directPermissions.length > 0) {
      const permsNode: Node = {
        id: `${id}-direct-permissions`,
        name: "Permissions",
        type: 'perm-folder',
        children: buildDirectPermissionsTree(src.directPermissions, `${id}-direct-permissions`)
      };
      node.children?.push(permsNode);
    }

    finalNodes[id] = node;
  });

  // Link children scopes and folders
  allScopes.forEach(scope => {
    const parentId = (scope.parent_id && finalNodes[scope.parent_id]) ? scope.parent_id : rootNodeId;
    const parentNode = finalNodes[parentId];
    const childNode = finalNodes[scope.id];
    
    if (!parentNode || !childNode) return;

    const folderName = scope.meta?.folder;

    if (folderName) {
      const folderId = `folder-${parentId}-${folderName}`;
      let folderNode = folderNodesMap[folderId];
      if (!folderNode) {
        folderNode = {
          id: folderId,
          name: folderName,
          type: 'folder',
          children: []
        };
        folderNodesMap[folderId] = folderNode;
        parentNode.children?.push(folderNode);
      }
      folderNode.children?.push(childNode);
    } else {
      parentNode.children?.push(childNode);
    }
  });

  return finalNodes[rootNodeId];
};

function buildPermissionsTree(permissions: Permission[], prefix: string, isFromRole = false, scopeId?: string | null): Node[] {
  const grouped: Record<string, Permission[]> = {};
  permissions.forEach(p => {
    if (!grouped[p.object]) grouped[p.object] = [];
    grouped[p.object].push(p);
  });

  return Object.entries(grouped).map(([object, perms]) => ({
    id: `${prefix}-obj-${object}`,
    name: object,
    type: 'object',
    children: perms.map(p => ({
      id: `${prefix}-obj-${object}-act-${p.action}`,
      name: p.action,
      type: 'action',
      icon: p.meta?.icon,
      color: p.meta?.color,
      isFromRole,
      scopeId: scopeId ?? null
    }))
  }));
}

function buildInheritedPermissionsTree(permissions: InheritedPermissionSource[], prefix: string): Node[] {
  const grouped: Record<string, Record<string, InheritedPermissionSource[]>> = {};
  permissions.forEach(p => {
    if (!grouped[p.object]) grouped[p.object] = {};
    if (!grouped[p.object][p.action]) grouped[p.object][p.action] = [];
    grouped[p.object][p.action].push(p);
  });

  return Object.entries(grouped).map(([object, actions]) => {
    return {
      id: `${prefix}-obj-${object}`,
      name: object,
      type: 'object',
      children: Object.entries(actions).map(([action, sources]) => ({
        id: `${prefix}-obj-${object}-act-${action}`,
        name: action,
        type: 'action',
        sources,
        scopeId: sources[0]?.scopeId ?? null,
        icon: sources[0]?.icon,
        color: sources[0]?.color
      }))
    };
  });
}

function buildDirectPermissionsTree(permissions: (Permission & { scopeId: string | null, isFromRole?: boolean })[], prefix: string): Node[] {
  const grouped: Record<string, (Permission & { scopeId: string | null, isFromRole?: boolean })[]> = {};
  permissions.forEach((p) => {
    if (!grouped[p.object]) grouped[p.object] = [];
    grouped[p.object].push(p);
  });

  return Object.entries(grouped).map(([object, perms]) => ({
    id: `${prefix}-obj-${object}`,
    name: object,
    type: 'object',
    children: perms.map(p => ({
      id: `${prefix}-obj-${object}-act-${p.action}`,
      name: p.action,
      type: 'action',
      icon: p.meta?.icon,
      color: p.meta?.color,
      permissionId: p.id,
      scopeId: p.scopeId,
      isFromRole: p.isFromRole ?? false 
    }))
  }));
}
