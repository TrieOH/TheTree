import type { User } from "../model/types";
import { useQuery, useMutation, useQueryClient, useQueries } from "@tanstack/react-query";
import { userRolesQueryOptions, userPermissionsQueryOptions, removeRoleOfUserFn, removePermissionOfUserFn } from "../api";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/shared/ui/shadcn/select";
import { Shield, Key, Trash2, Loader2, Info } from "lucide-react";
import { toast } from "sonner";

import type { Scope } from "@/features/scope/model/types";
import { useState } from "react";

interface PropsI {
  user: User;
  project_id: string;
  onBack: () => void;
  allScopes: Scope[];
}

export default function CurrentAccessList({ user, project_id, onBack, allScopes }: PropsI) {
  const queryClient = useQueryClient();
  const [selectedScopeId, setSelectedScopeId] = useState<string | 'all'>('all');
  
  const { data: roles = [], isLoading: isLoadingRoles } = useQuery(
    userRolesQueryOptions(project_id, user.id)
  );

  const permissionQueries = useQueries({
    queries: allScopes.map((scope) =>
      userPermissionsQueryOptions(project_id, user.id, scope.id)
    ),
  });

  const allPermissions = permissionQueries.flatMap((query) => query.data || []);
  const isLoadingPerms = permissionQueries.some((query) => query.isLoading);

  const removeRoleMutation = useMutation({
    mutationFn: ({ roleId, scopeId }: { roleId: string; scopeId: string }) => 
      removeRoleOfUserFn(user, roleId, scopeId),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['userRoles', project_id, user.id] });
      queryClient.invalidateQueries({ queryKey: ['userPermissions', project_id, user.id, variables.scopeId] });
      toast.success(data.message);
    },
    onError: () => toast.error("Failed to remove role")
  });

  const removePermissionMutation = useMutation({
    mutationFn: ({ permissionId, scopeId }: { permissionId: string; scopeId: string }) => 
      removePermissionOfUserFn(user, permissionId, scopeId),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['userPermissions', project_id, user.id, variables.scopeId] });
      toast.success(data.message);
    },
    onError: () => toast.error("Failed to remove permission")
  });

  const isLoading = isLoadingRoles || isLoadingPerms;
  const hasAccess = roles.length > 0 || allPermissions.length > 0;

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
        <div className="space-y-6 max-h-112.5 overflow-y-auto pr-1">
          {/* Roles Section */}
          {roles.length > 0 && (
            <div className="space-y-2">
              <h3 className="text-xs font-bold text-muted-foreground uppercase px-1">Assigned Roles</h3>
              {roles.map((role) => (
                <div 
                  key={`${role.id}-${role.scope_id}`}
                  className="flex items-center justify-between p-3 bg-muted/50 rounded-md border border-transparent hover:border-border transition-all group"
                >
                  <div className="flex items-center gap-3">
                    <div className="p-2 bg-primary/10 rounded-full">
                      <Shield className="w-4 h-4 text-primary" />
                    </div>
                    <div className="flex flex-col">
                      <span className="text-sm font-medium">{role.name}</span>
                      <span className="text-[10px] text-muted-foreground">
                        Scope: <span className="text-foreground">{role.scope_name || role.scope_id}</span>
                      </span>
                    </div>
                  </div>
                  <button
                    disabled={removeRoleMutation.isPending}
                    onClick={() => removeRoleMutation.mutate({ roleId: role.id, scopeId: role.scope_id })}
                    className="p-2 text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded-md transition-colors disabled:opacity-50"
                  >
                    {removeRoleMutation.isPending ? (
                      <Loader2 className="w-4 h-4 animate-spin" />
                    ) : (
                      <Trash2 className="w-4 h-4" />
                    )}
                  </button>
                </div>
              ))}
            </div>
          )}

          {/* Permissions Section */}
          {allPermissions.length > 0 && (
            <div className="space-y-6">
              <div className="flex items-center justify-between py-1 gap-1 sm:flex-row flex-col">
                <h3 className="text-xs font-bold text-muted-foreground uppercase">
                  Direct Permissions by Scope
                </h3>
                <Select value={selectedScopeId} onValueChange={setSelectedScopeId}>
                  <SelectTrigger className="w-45">
                    <SelectValue placeholder="Filter by scope" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Scopes</SelectItem>
                    {allScopes.map((scope) => (
                      <SelectItem key={scope.id} value={scope.id}>
                        {scope.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              {allScopes.map((scope, index) => {
                if (selectedScopeId !== 'all' && selectedScopeId !== scope.id) return null;

                const scopePermissions = permissionQueries[index]?.data || [];

                if (scopePermissions.length === 0) return null;
                return (
                  <div key={scope.id} className="space-y-2">
                    <p className="text-sm font-semibold text-foreground px-1">{scope.name}</p>
                    {scopePermissions.map((perm) => (
                      <div 
                        key={`${scope.id}-${perm.id}`}
                        className="flex items-center justify-between p-3 bg-muted/50 rounded-md border border-transparent hover:border-border transition-all group"
                      >
                        <div className="flex items-center gap-3">
                          <div className="p-2 bg-accent/10 rounded-full">
                            <Key className="w-4 h-4 text-accent" />
                          </div>
                          <div className="flex flex-col">
                            <div className="flex items-center gap-1.5">
                              <span className="text-sm font-medium">{perm.object}</span>
                              <span className="text-[10px] text-muted-foreground">:</span>
                              <span className="text-sm font-medium text-accent">{perm.action}</span>
                            </div>
                          </div>
                        </div>
                        <button
                          disabled={removePermissionMutation.isPending}
                          onClick={() => removePermissionMutation.mutate({ 
                            permissionId: perm.id, 
                            scopeId: scope.id 
                          })}
                          className="p-2 text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded-md transition-colors disabled:opacity-50"
                        >
                          {removePermissionMutation.isPending ? (
                            <Loader2 className="w-4 h-4 animate-spin" />
                          ) : (
                            <Trash2 className="w-4 h-4" />
                          )}
                        </button>
                      </div>
                    ))}
                  </div>
                );
              })}
            </div>
          )}
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
