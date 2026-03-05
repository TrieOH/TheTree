import type { User } from "../model/types";
import { useQuery, useMutation, useQueryClient, useQueries } from "@tanstack/react-query";
import { userRolesQueryOptions, userPermissionsQueryOptions, removeRoleOfUserFn, removePermissionOfUserFn } from "../api";
import { rolePermissionsQueryOptions } from "@/features/role/api";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Shield, Key, Trash2, Loader2, Info, ChevronDown, ChevronRight, Folder as FolderIcon } from "lucide-react";
import { toast } from "sonner";
import { cn } from "@/shared/lib/utils";

import type { Scope } from "@/features/scope/model/types";
import { useState, useMemo } from "react";
import type { Permission } from "@/features/permission/model/types";
import type { Role } from "@/features/role/model/types";

interface PropsI {
  user: User;
  project_id: string;
  onBack: () => void;
  allScopes: Scope[];
}

interface PermissionWithScope extends Permission {
  scope_id: string | null;
}

interface RoleWithPerms extends Role {
  permissions: Permission[];
  isLoadingPerms: boolean;
}

interface AccessTreeNode {
  scope: Scope | { id: string | null; name: string; isFolder?: boolean };
  roles: RoleWithPerms[];
  directPermissions: PermissionWithScope[];
  children: AccessTreeNode[];
}

export default function CurrentAccessList({ user, project_id, onBack, allScopes }: PropsI) {
  const queryClient = useQueryClient();

  const scopesWithGlobal = useMemo(() => [
    { id: 'global-scope', name: 'Root', project_id, created_at: '', type: 'global', external_id: 'global', parent_id: null },
    ...allScopes
  ], [allScopes, project_id]);

  const scopeIdMapping = (id: string | null) => id === 'global-scope' ? null : id;

  const { data: roles = [], isLoading: isLoadingRoles } = useQuery(
    userRolesQueryOptions(project_id, user.id)
  );

  const rolePermissionsQueries = useQueries({
    queries: roles.map((role) =>
      rolePermissionsQueryOptions(project_id, role.id)
    ),
  });

  const rolesWithPerms: RoleWithPerms[] = roles.map((role, index) => ({
    ...role,
    permissions: rolePermissionsQueries[index]?.data || [],
    isLoadingPerms: rolePermissionsQueries[index]?.isLoading || false,
  }));

  const permissionQueries = useQueries({
    queries: scopesWithGlobal.map((scope) =>
      userPermissionsQueryOptions(project_id, user.id, scopeIdMapping(scope.id))
    ),
  });

  const allPermissionsWithScope: PermissionWithScope[] = permissionQueries.flatMap((query, index) => {
    const scope = scopesWithGlobal[index];
    const currentScopeId = scopeIdMapping(scope.id);
    return (query.data || []).map(perm => ({ ...perm, scope_id: currentScopeId }));
  });

  const isLoadingPerms = permissionQueries.some((query) => query.isLoading);
  const isLoadingRolePerms = rolePermissionsQueries.some((query) => query.isLoading);

  const removeRoleMutation = useMutation({
    mutationFn: ({ roleId, scopeId }: { roleId: string; scopeId: string | null }) => 
      removeRoleOfUserFn(user, roleId, scopeId),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['userRoles', project_id, user.id] });
      queryClient.invalidateQueries({ queryKey: ['userPermissions', project_id, user.id, variables.scopeId] });
      toast.success(data?.message);
    },
    onError: () => toast.error("Failed to remove role")
  });

  const removePermissionMutation = useMutation({
    mutationFn: ({ permissionId, scopeId }: { permissionId: string; scopeId: string | null }) => 
      removePermissionOfUserFn(user, permissionId, scopeId),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['userPermissions', project_id, user.id, variables.scopeId] });
      toast.success(data?.message);
    },
    onError: () => toast.error("Failed to remove permission")
  });

  const getScopeDisplayName = (scope: Scope | { id: string | null; name: string; isFolder?: boolean }) => {
    if (!scope.id || scope.id === 'global-scope' || scope.id === 'root-node-id') return "Root";
    return scope.name;
  };

  const rolesByScope = useMemo(() => {
    const map = new Map<string, RoleWithPerms[]>();
    rolesWithPerms.forEach(role => {
      const scopeId = role.scope_id || 'root-node-id';
      if (!map.has(scopeId)) map.set(scopeId, []);
      map.get(scopeId)?.push(role);
    });
    return map;
  }, [rolesWithPerms]);

  const directPermsByScope = useMemo(() => {
    const map = new Map<string, PermissionWithScope[]>();
    allPermissionsWithScope.forEach(perm => {
      const scopeId = perm.scope_id || 'root-node-id';
      if (!map.has(scopeId)) map.set(scopeId, []);
      map.get(scopeId)?.push(perm);
    });
    return map;
  }, [allPermissionsWithScope]);

  const getInheritanceInfo = (perm: PermissionWithScope) => {
    const currentScopeId = perm.scope_id; // null for Root
    const normalizedCurrentScopeId = currentScopeId || 'root-node-id';

    // Helper to check if a permission is provided in a specific scope
    const checkInScope = (sid: string | null) => {
      const normalizedSid = sid || 'root-node-id';
      
      // Check roles
      const rolesInScope = rolesByScope.get(normalizedSid) || [];
      const roleProvidingPerm = rolesInScope.find(r => 
        r.permissions.some(p => p.object === perm.object && p.action === perm.action)
      );
      
      if (roleProvidingPerm) {
        const scopeObj = scopesWithGlobal.find(s => scopeIdMapping(s.id) === sid) || { id: null, name: 'Root' };
        const scopeName = getScopeDisplayName(scopeObj);
        return `role ${roleProvidingPerm.name} in ${scopeName}`;
      }

      // Check direct
      const permsInScope = directPermsByScope.get(normalizedSid) || [];
      const isDirectInScope = permsInScope.some(p => p.object === perm.object && p.action === perm.action);
      
      if (isDirectInScope) {
        const scopeObj = scopesWithGlobal.find(s => scopeIdMapping(s.id) === sid) || { id: null, name: 'Root' };
        const scopeName = getScopeDisplayName(scopeObj);
        return scopeName;
      }
      
      return null;
    };

    // 1. First check if it's in a Role in the SAME scope
    const rolesInSameScope = rolesByScope.get(normalizedCurrentScopeId) || [];
    const roleInSameScope = rolesInSameScope.find(r => 
      r.permissions.some(p => p.object === perm.object && p.action === perm.action)
    );
    if (roleInSameScope) {
      const scopeName = getScopeDisplayName(scopesWithGlobal.find(s => scopeIdMapping(s.id) === currentScopeId) || { id: null, name: 'Root' });
      return `role ${roleInSameScope.name} in ${scopeName}`;
    }

    // 2. Check parents recursively
    if (!currentScopeId) return null; // Root has no parents
    
    let parentId: string | null = allScopes.find(s => s.id === currentScopeId)?.parent_id || null;
    while (parentId) {
      const info = checkInScope(parentId);
      if (info) return info;
      parentId = allScopes.find(s => s.id === parentId)?.parent_id || null;
    }

    // 3. Finally check Root
    return checkInScope(null);
  };

  const accessTree = useMemo(() => {
    const nodesMap: Record<string, AccessTreeNode> = {};
    const folderNodesMap: Record<string, AccessTreeNode> = {};
    
    // Initialize Root Node
    const rootNode: AccessTreeNode = {
      scope: { id: 'root-node-id', name: 'Root' },
      roles: [],
      directPermissions: [],
      children: []
    };
    nodesMap['root-node-id'] = rootNode;

    // Map all scopes from the project
    allScopes.forEach(scope => {
      nodesMap[scope.id] = {
        scope,
        roles: [],
        directPermissions: [],
        children: []
      };
    });

    // Distribute Roles into the tree
    rolesWithPerms.forEach(role => {
      const targetId = role.scope_id || 'root-node-id';
      const targetNode = nodesMap[targetId] || rootNode;
      targetNode.roles.push(role);
    });

    // Distribute Direct Permissions into the tree
    allPermissionsWithScope.forEach(perm => {
      const targetId = perm.scope_id || 'root-node-id';
      const targetNode = nodesMap[targetId] || rootNode;
      targetNode.directPermissions.push(perm);
    });

    // Build the hierarchy with folders
    allScopes.forEach(scope => {
      const parentId = scope.parent_id || 'root-node-id';
      const parentNode = nodesMap[parentId] || rootNode;
      const folderName = scope.meta?.folder;

      if (folderName) {
        const folderId = `folder-${parentId}-${folderName}`;
        let folderNode = folderNodesMap[folderId];
        
        if (!folderNode) {
          folderNode = {
            scope: { id: folderId, name: folderName, isFolder: true },
            roles: [],
            directPermissions: [],
            children: []
          };
          folderNodesMap[folderId] = folderNode;
          parentNode.children.push(folderNode);
        }
        
        folderNode.children.push(nodesMap[scope.id]);
      } else {
        parentNode.children.push(nodesMap[scope.id]);
      }
    });

    return [rootNode];
  }, [allScopes, rolesWithPerms, allPermissionsWithScope]);

  const isLoading = isLoadingRoles || isLoadingPerms || isLoadingRolePerms;
  const hasAccess = roles.length > 0 || allPermissionsWithScope.length > 0;

  return (
    <div className="flex flex-col gap-4">
      <div className="text-center w-full">
        <span className="text-primary font-bold uppercase tracking-wider">Current Access</span>
        <p className="text-xs text-muted-foreground mt-1">
          Managing access for <span className="font-medium text-foreground">{user.email}</span>
        </p>
      </div>

      {isLoading ? (
        <div className="flex flex-col items-center justify-center py-10 gap-2">
          <Loader2 className="w-6 h-6 animate-spin text-primary" />
          <span className="text-xs text-muted-foreground">Loading access data...</span>
        </div>
      ) : !hasAccess ? (
        <div className="flex flex-col items-center justify-center py-10 px-4 bg-muted/30 rounded-lg border border-dashed border-muted gap-2">
          <Info className="w-8 h-8 text-muted-foreground/50" />
          <p className="text-sm text-muted-foreground text-center">
            This user currently has no direct roles or permissions assigned.
          </p>
        </div>
      ) : (
        <div className="space-y-2 max-h-112.5 overflow-auto pr-1 pb-2">
          {accessTree.map(node => (
            <AccessNode 
              key={node.scope.id || 'root'} 
              node={node} 
              level={0}
              onRemoveRole={(roleId, scopeId) => removeRoleMutation.mutate({ roleId, scopeId })}
              onRemovePermission={(permId, scopeId) => removePermissionMutation.mutate({ permissionId: permId, scopeId })}
              isRemovingRole={removeRoleMutation.isPending}
              isRemovingPermission={removePermissionMutation.isPending}
              getInheritanceInfo={getInheritanceInfo}
              getScopeDisplayName={getScopeDisplayName}
            />
          ))}
        </div>
      )}

      <hr className="w-full border-muted mt-2" />
      <div className="w-full">
        <ShadowButton 
          value="Back" 
          variant="ghost" 
          onClick={onBack} 
          className="w-full justify-center"
        />
      </div>
    </div>
  );
}

function AccessNode({ 
  node, 
  level,
  onRemoveRole,
  onRemovePermission,
  isRemovingRole,
  isRemovingPermission,
  getInheritanceInfo,
  getScopeDisplayName
}: { 
  node: AccessTreeNode; 
  level: number;
  onRemoveRole: (roleId: string, scopeId: string | null) => void;
  onRemovePermission: (permId: string, scopeId: string | null) => void;
  isRemovingRole: boolean;
  isRemovingPermission: boolean;
  getInheritanceInfo: (perm: PermissionWithScope) => string | null;
  getScopeDisplayName: (scope: Scope | { id: string | null; name: string; isFolder?: boolean }) => string;
}) {
  const [isExpanded, setIsExpanded] = useState(false);
  
  const hasDirectContent = node.roles.length > 0 || node.directPermissions.length > 0;
  
  const checkContent = (n: AccessTreeNode): boolean => {
    if (n.roles.length > 0 || n.directPermissions.length > 0) return true;
    return n.children.some(checkContent);
  };
  
  const hasVisibleChildren = node.children.some(checkContent);
  const isFolder = 'isFolder' in node.scope && node.scope.isFolder;

  if (!hasDirectContent && !hasVisibleChildren) return null;

  return (
    <div className={cn("flex flex-col", level > 0 && "ml-4 border-l border-muted/50 pl-4")}>
      <div 
        className="flex items-center gap-2 py-2 cursor-pointer group"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        {node.children.length > 0 ? (
          isExpanded ? <ChevronDown className="w-4 h-4 text-muted-foreground" /> : <ChevronRight className="w-4 h-4 text-muted-foreground" />
        ) : (
          <div className="w-4" />
        )}
        
        {isFolder && <FolderIcon className="w-3.5 h-3.5 text-indigo-400" />}
        
        <span className={cn(
          "text-[11px] font-bold uppercase tracking-wider transition-colors",
          level === 0 ? "text-primary" : (isFolder ? "text-indigo-500/80" : "text-muted-foreground group-hover:text-foreground")
        )}>
          {getScopeDisplayName(node.scope)}
        </span>
      </div>

      {isExpanded && (
        <div className="flex flex-col gap-2 mt-1 pb-2">
          {/* Roles */}
          {node.roles.length > 0 && (
            <div className="flex flex-col gap-2">
              {node.roles.map(role => (
                <RoleItem 
                  key={`${role.id}-${role.scope_id}`}
                  role={role}
                  onRemove={() => onRemoveRole(role.id, role.scope_id)}
                  isRemoving={isRemovingRole}
                />
              ))}
            </div>
          )}

          {/* Direct Permissions */}
          {node.directPermissions.length > 0 && (
            <div className="flex flex-col gap-2">
              {node.directPermissions.map(perm => {
                const inheritanceInfo = getInheritanceInfo(perm);
                return (
                  <PermissionItem 
                    key={`${perm.id}-${perm.scope_id}`}
                    permission={perm}
                    onRemove={() => onRemovePermission(perm.id, perm.scope_id)}
                    isRemoving={isRemovingPermission}
                    inheritedFrom={inheritanceInfo}
                  />
                );
              })}
            </div>
          )}

          {/* Children Scopes & Folders */}
          {node.children.map(child => (
            <AccessNode 
              key={child.scope.id || 'root'} 
              node={child} 
              level={level + 1}
              onRemoveRole={onRemoveRole}
              onRemovePermission={onRemovePermission}
              isRemovingRole={isRemovingRole}
              isRemovingPermission={isRemovingPermission}
              getInheritanceInfo={getInheritanceInfo}
              getScopeDisplayName={getScopeDisplayName}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function RoleItem({ role, onRemove, isRemoving }: { role: RoleWithPerms; onRemove: () => void; isRemoving: boolean }) {
  const [showPerms, setShowPerms] = useState(false);

  return (
    <div className="flex flex-col bg-muted/30 rounded-md border border-transparent hover:border-border transition-all overflow-hidden">
      <div className="flex items-center justify-between p-3">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-primary/10 rounded-full">
            <Shield className="w-4 h-4 text-primary" />
          </div>
          <div className="flex flex-col">
            <span className="text-sm font-medium">{role.name}</span>
            <button 
              type="button"
              onClick={(e) => { e.stopPropagation(); setShowPerms(!showPerms); }}
              className="text-[10px] text-muted-foreground hover:text-primary flex items-center gap-1"
            >
              {showPerms ? "Hide permissions" : "Show permissions"}
              {showPerms ? <ChevronDown className="w-2.5 h-2.5" /> : <ChevronRight className="w-2.5 h-2.5" />}
            </button>
          </div>
        </div>
        <button
          type="button"
          disabled={isRemoving}
          onClick={onRemove}
          className="p-2 text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded-md transition-colors disabled:opacity-50"
        >
          {isRemoving ? <Loader2 className="w-4 h-4 animate-spin" /> : <Trash2 className="w-4 h-4" />}
        </button>
      </div>
      
      {showPerms && (
        <div className="px-4 pb-3 flex flex-wrap gap-1.5 border-t border-muted/50 pt-2 bg-muted/20">
          {role.isLoadingPerms ? (
            <Loader2 className="w-3 h-3 animate-spin text-muted-foreground" />
          ) : role.permissions.length > 0 ? (
            role.permissions.map(p => (
              <span key={p.id} className="text-[10px] px-2 py-0.5 bg-background rounded-full border border-border text-muted-foreground">
                <span className="font-semibold text-foreground">{p.object}</span>:{p.action}
              </span>
            ))
          ) : (
            <span className="text-[10px] text-muted-foreground">No permissions in this role</span>
          )}
        </div>
      )}
    </div>
  );
}

function PermissionItem({ permission, onRemove, isRemoving, inheritedFrom }: { 
  permission: PermissionWithScope; 
  onRemove: () => void; 
  isRemoving: boolean;
  inheritedFrom: string | null;
}) {
  return (
    <div className="flex items-center justify-between p-3 bg-muted/50 rounded-md border border-transparent hover:border-border transition-all group">
      <div className="flex items-center gap-3">
        <div className="p-2 bg-accent/10 rounded-full">
          <Key className="w-4 h-4 text-accent" />
        </div>
        <div className="flex flex-col">
          <div className="flex items-center gap-1.5">
            <span className="text-sm font-medium">{permission.object}</span>
            <span className="text-[10px] text-muted-foreground">:</span>
            <span className="text-sm font-medium text-accent">{permission.action}</span>
          </div>
          {inheritedFrom && (
            <span className="text-[9px] text-muted-foreground italic">inherited from {inheritedFrom}</span>
          )}
        </div>
      </div>
      <button
        type="button"
        disabled={isRemoving || !!inheritedFrom}
        onClick={onRemove}
        className="p-2 text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded-md transition-colors disabled:opacity-50"
      >
        {isRemoving ? <Loader2 className="w-4 h-4 animate-spin" /> : <Trash2 className="w-4 h-4" />}
      </button>
    </div>
  );
}
