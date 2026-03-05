import { useState, useMemo } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  givePermissionToRoleFn,
  removePermissionOfRoleFn,
  rolePermissionsQueryOptions,
} from "../api";
import { permissionsQueryOptions } from "@/features/permission/api";
import type { Permission } from "@/features/permission/model/types";
import { Checkbox } from "@/shared/ui/shadcn/checkbox";
import { Skeleton } from "@/shared/ui/shadcn/skeleton";
import { X, Loader2, ChevronRight, ChevronDown } from "lucide-react";
import type { Role } from "../model/types";
import { toast } from "sonner";
import { cn } from "@/shared/lib/utils";
import { SearchInput } from "@/shared/ui/form/SearchInput";
import { Badge } from "@/shared/ui/shadcn/badge";

interface RolePermissionsEditorProps {
  project_id: string;
  role: Role;
}

interface GroupedPermissions {
  [object: string]: Permission[];
}

export default function RolePermissionsEditor({
  project_id,
  role,
}: RolePermissionsEditorProps) {
  const queryClient = useQueryClient();
  const [searchTerm, setSearchTerm] = useState("");
  const [expandedObjects, setExpandedObjects] = useState<Record<string, boolean>>({});

  const { data: allPermissions, isLoading: isLoadingAllPermissions } = useQuery(
    permissionsQueryOptions(project_id)
  );

  const {
    data: rolePermissions,
    isLoading: isLoadingRolePermissions,
    error,
  } = useQuery(rolePermissionsQueryOptions(project_id, role.id));

  const givePermissionMutation = useMutation({
    mutationFn: ({ role, permission_id }: { role: Role; permission_id: string }) =>
      givePermissionToRoleFn(role, permission_id),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: rolePermissionsQueryOptions(project_id, role.id).queryKey,
        });
        toast.success(response.message);
      } else {
        toast.error(`Failed to give permission: ${response.message}`);
      }
    },
  });

  const removePermissionMutation = useMutation({
    mutationFn: ({ role, permission_id }: { role: Role; permission_id: string }) =>
      removePermissionOfRoleFn(role, permission_id),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: rolePermissionsQueryOptions(project_id, role.id).queryKey,
        });
        toast.success(response.message);
      } else {
        toast.error(`Failed to remove permission: ${response.message}`);
      }
    },
  });

  const handlePermissionChange = (permission: Permission, isChecked: boolean) => {
    if (isChecked) {
      givePermissionMutation.mutate({ role, permission_id: permission.id });
    } else {
      removePermissionMutation.mutate({ role, permission_id: permission.id });
    }
  };

  const isPermissionAssigned = (permissionId: string) => {
    return rolePermissions?.some((p) => p.id === permissionId) ?? false;
  };

  const isMutating = givePermissionMutation.isPending || removePermissionMutation.isPending;

  const groupedPermissions = useMemo(() => {
    if (!allPermissions) return {};
    
    const filtered = allPermissions.filter((permission) =>
      `${permission.object}:${permission.action}`
        .toLowerCase()
        .includes(searchTerm.toLowerCase())
    );

    return filtered.reduce((acc: GroupedPermissions, p) => {
      if (!acc[p.object]) acc[p.object] = [];
      acc[p.object].push(p);
      return acc;
    }, {});
  }, [allPermissions, searchTerm]);

  const toggleObject = (object: string) => {
    setExpandedObjects(prev => ({ ...prev, [object]: !prev[object] }));
  };

  if (isLoadingAllPermissions || isLoadingRolePermissions) {
    return (
      <div className="flex flex-col items-center py-8 min-h-52 space-y-6">
        <div className="text-center space-y-2">
          <Skeleton className="h-6 w-48 mx-auto" />
          <Skeleton className="h-4 w-64 mx-auto" />
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 w-full max-w-4xl">
          {[...Array(6)].map(() => (
            <Skeleton key={crypto.randomUUID()} className="h-24 w-full rounded-lg" />
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center py-4 min-h-36 text-destructive text-center">
        <X className="h-8 w-8 mb-2" />
        <p>Failed to load permissions</p>
        <p className="text-sm opacity-70">{error.message}</p>
      </div>
    );
  }

  const objects = Object.keys(groupedPermissions).sort();

  return (
    <div className="flex flex-col items-center min-h-36 py-6 px-4 bg-muted/20 rounded-lg border border-border/50 my-2">
      <div className="text-center space-y-1 mb-6">
        <h2 className="text-lg font-semibold flex items-center justify-center gap-2">
          Assign Permissions to <span className="text-primary">{role.name}</span>
        </h2>
        <p className="text-xs text-muted-foreground">
          Select the specific actions this role is allowed to perform across different objects.
        </p>
      </div>

      <div className="w-full max-w-md mb-8">
        <SearchInput
          placeholder="Filter by object or action..."
          value={searchTerm}
          onChange={setSearchTerm}
        />
      </div>

      {objects.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 w-full max-w-6xl items-start">
          {objects.map((object) => {
            const permissions = groupedPermissions[object];
            const assignedCount = permissions.filter(p => isPermissionAssigned(p.id)).length;
            const isExpanded = expandedObjects[object] !== false; // Default expanded

            return (
              <div 
                key={object}
                className="flex flex-col border border-border bg-background rounded-lg overflow-hidden shadow-sm h-fit"
              >
                <div 
                  className="flex items-center justify-between px-3 py-2 bg-muted/40 border-b border-border cursor-pointer hover:bg-muted/60 transition-colors"
                  onClick={() => toggleObject(object)}
                >
                  <div className="flex items-center gap-2 overflow-hidden text-muted-foreground">
                    {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                    <span className="font-bold text-xs truncate uppercase tracking-wider">
                      {object}
                    </span>
                  </div>
                  <Badge variant={assignedCount > 0 ? "secondary" : "outline"} className="text-[10px] h-4.5 px-1.5 shrink-0">
                    {assignedCount}/{permissions.length}
                  </Badge>
                </div>

                {isExpanded && (
                  <div className="p-3 grid grid-cols-1 gap-1.5 max-h-48 overflow-y-auto">
                    {permissions.map((permission) => {
                      const isAssigned = isPermissionAssigned(permission.id);
                      const isPending = 
                        (givePermissionMutation.isPending && 
                         givePermissionMutation.variables?.permission_id === permission.id) ||
                        (removePermissionMutation.isPending && 
                         removePermissionMutation.variables?.permission_id === permission.id);

                      return (
                        <div
                          key={permission.id}
                          className={cn(
                            "flex items-center gap-2.5 px-2 py-1.5 rounded-md text-xs transition-all",
                            isAssigned 
                              ? "bg-primary/3 text-foreground font-medium" 
                              : "text-muted-foreground hover:bg-accent/50",
                            isPending && "opacity-50 pointer-events-none"
                          )}
                        >
                          <div className="relative flex items-center shrink-0">
                            <Checkbox
                              id={`perm-${role.id}-${permission.id}`}
                              checked={isAssigned}
                              onCheckedChange={(checked) =>
                                handlePermissionChange(permission, checked as boolean)
                              }
                              disabled={isMutating}
                              className="h-4 w-4 rounded-lg border-muted-foreground/30"
                            />
                            {isPending && (
                              <Loader2 className="absolute inset-0 h-4 w-4 animate-spin text-primary m-auto" />
                            )}
                          </div>
                          
                          <label
                            htmlFor={`perm-${role.id}-${permission.id}`}
                            className="flex-1 cursor-pointer select-none truncate font-mono"
                            title={permission.action}
                          >
                            {permission.action}
                          </label>
                        </div>
                      );
                    })}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      ) : (
        <div className="text-center py-10 text-muted-foreground bg-muted/10 rounded-lg border border-dashed border-border w-full">
          <p className="text-sm">No permissions match your search.</p>
        </div>
      )}
    </div>
  );
}
