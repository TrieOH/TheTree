import type { User, Node, InheritedPermissionSource } from "../model/types";
import { useQuery, useMutation, useQueryClient, useQueries } from "@tanstack/react-query";
import { userRolesQueryOptions, userPermissionsQueryOptions, removeRoleOfUserFn, removePermissionOfUserFn } from "../api";
import { rolePermissionsQueryOptions } from "@/features/role/api";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Trash2, Loader2, Info } from "lucide-react";
import { toast } from "sonner";

import type { Scope } from "@/features/scope/model/types";
import { useMemo } from "react";
import UserPermTree from "./UserPermTree";
import { buildFullAccessControlNodeTree } from "../lib/node-tree-utils";

interface PropsI {
  user: User;
  project_id: string;
  onBack: () => void;
  allScopes: Scope[];
}

export default function CurrentAccessList({ user, project_id, onBack, allScopes }: PropsI) {
  const queryClient = useQueryClient();

  const scopesWithGlobal = useMemo(() => [
    { id: 'global-scope', name: 'Root', project_id, created_at: '', type: 'global', external_id: 'global', parent_id: null } as Scope,
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

  const rolesWithPerms = roles.map((role, index) => ({
    ...role,
    permissions: rolePermissionsQueries[index]?.data || [],
    isLoadingPerms: rolePermissionsQueries[index]?.isLoading || false,
  }));

  const permissionQueries = useQueries({
    queries: scopesWithGlobal.map((scope) =>
      userPermissionsQueryOptions(project_id, user.id, scopeIdMapping(scope.id))
    ),
  });

  const allPermissionsWithScope = permissionQueries.flatMap((query, index) => {
    const scope = scopesWithGlobal[index];
    const currentScopeId = scopeIdMapping(scope.id);
    return (query.data || []).map(perm => ({ ...perm, scopeId: currentScopeId }));
  });

  const accessTree = useMemo(() => {
    if (isLoadingRoles || permissionQueries.some(q => q.isLoading) || rolePermissionsQueries.some(q => q.isLoading)) {
        return null;
    }
    return buildFullAccessControlNodeTree(allScopes, rolesWithPerms, allPermissionsWithScope);
  }, [allScopes, rolesWithPerms, allPermissionsWithScope, isLoadingRoles, permissionQueries, rolePermissionsQueries]);

  const isLoading = isLoadingRoles || permissionQueries.some(q => q.isLoading) || rolePermissionsQueries.some(q => q.isLoading);
  const hasAccess = roles.length > 0 || allPermissionsWithScope.length > 0;

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

  const renderExtra = (node: Node) => {
    if (node.type === 'role') {
        const isInherited = node.isInherited;
        const roleId = node.roleId;
        const scopeId = node.scopeId;
        if (isInherited) {
            return (
                <span className="text-[8px] text-muted-foreground/60 italic">
                    from {node.sourceScope}
                </span>
            );
        }
        return (
            <button 
                type="button"
                onClick={() => roleId && removeRoleMutation.mutate({ roleId, scopeId: scopeId ?? null })} 
                disabled={removeRoleMutation.isPending}
                className="p-1 text-muted-foreground/30 hover:text-destructive transition-all"
            >
                {removeRoleMutation.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <Trash2 className="w-3 h-3" />}
            </button>
        );
    }

    if (node.type === 'action') {
        if (node.sources) {
            const sources: InheritedPermissionSource[] = node.sources;
            return (
                <div className="flex flex-wrap gap-1 justify-end max-w-37.5">
                    {sources.map((s) => (
                        <span key={`${s.sourceRole}-${s.sourceScope}`} className="text-[7px] text-muted-foreground italic bg-muted/20 px-1 rounded border border-muted/10 whitespace-nowrap">
                            from {s.sourceRole ? `${s.sourceRole} (${s.sourceScope})` : s.sourceScope}
                        </span>
                    ))}
                </div>
            );
        }
        if (node.permissionId) {
            const permissionId = node.permissionId;
            const scopeId = node.scopeId;
            const isFromRole = node.isFromRole === true;

            if (isFromRole) {
                return (
                    <span className="text-[7px] text-muted-foreground italic bg-muted/20 px-1 rounded border border-muted/10 whitespace-nowrap">
                        from roles
                    </span>
                );
            }

            return (
                <button 
                    type="button"
                    onClick={() => removePermissionMutation.mutate({ permissionId, scopeId: scopeId ?? null })} 
                    disabled={removePermissionMutation.isPending}
                    className="p-1 text-muted-foreground/30 hover:text-destructive transition-all"
                >
                    {removePermissionMutation.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <Trash2 className="w-3 h-3" />}
                </button>
            );
        }
    }

    return null;
  };

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
      ) : !hasAccess || !accessTree ? (
        <div className="flex flex-col items-center justify-center py-10 px-4 bg-muted/30 rounded-lg border border-dashed border-muted gap-2">
          <Info className="w-8 h-8 text-muted-foreground/50" />
          <p className="text-sm text-muted-foreground text-center">
            This user currently has no direct roles or permissions assigned.
          </p>
        </div>
      ) : (
        <div className="space-y-2 max-h-112.5 overflow-auto pr-1 pb-2">
          <UserPermTree 
            node={accessTree}
            goBack={onBack}
            showFooter={false}
            renderExtra={renderExtra}
            defaultExpanded={false}
          />
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