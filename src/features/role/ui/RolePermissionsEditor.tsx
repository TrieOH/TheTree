import { useState } from "react";
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
import { X, Loader2 } from "lucide-react";
import type { Role } from "../model/types";
import { toast } from "sonner";
import { cn } from "@/shared/lib/utils";
import { SearchInput } from "@/shared/ui/form/SearchInput";

interface RolePermissionsEditorProps {
  project_id: string;
  role: Role;
}

export default function RolePermissionsEditor({
  project_id,
  role,
}: RolePermissionsEditorProps) {
  const queryClient = useQueryClient();
  const [searchTerm, setSearchTerm] = useState("");

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

  const filteredPermissions = allPermissions?.filter((permission) =>
    `${permission.object}:${permission.action}`
      .toLowerCase()
      .includes(searchTerm.toLowerCase())
  );

  if (isLoadingAllPermissions || isLoadingRolePermissions) {
    return (
      <div className="flex flex-col items-center py-4 min-h-52 space-y-6">
        <div className="text-center space-y-2">
          <Skeleton className="h-6 w-48 mx-auto" />
          <Skeleton className="h-4 w-64 mx-auto" />
        </div>
        <div className="flex flex-wrap justify-center gap-2 max-w-2xl">
          {[...Array(12)].map(() => (
            <Skeleton key={crypto.randomUUID()} className="h-8 w-32" />
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

  return (
    <div className="flex flex-col items-center min-h-36 py-4 px-4">
      <div className="text-center space-y-1 mb-4">
        <h2 className="text-xl font-semibold flex items-center justify-center gap-2">
          Manage Permissions
        </h2>
        <p className="text-sm text-muted-foreground">
          Configure permissions for "{role.name}"
        </p>
      </div>

      <div className="w-full max-w-sm mb-4">
        <SearchInput
          placeholder="Search permissions (e.g., project:read)"
          value={searchTerm}
          onChange={(value) => setSearchTerm(value)}
        />
      </div>

      {filteredPermissions && filteredPermissions.length > 0 ? (
        <div className="flex flex-wrap justify-center gap-2 max-w-3xl max-h-80 overflow-y-auto">
          {filteredPermissions.map((permission) => {
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
                  "group flex items-center gap-2 px-3 py-1.5 rounded-md border text-sm transition-all duration-200",
                  "hover:border-primary/40 hover:bg-accent/40",
                  isAssigned 
                    ? "border-primary/30 bg-primary/5 text-foreground" 
                    : "border-border text-muted-foreground",
                  isPending && "opacity-60"
                )}
              >
                <div className="relative flex items-center justify-center">
                  <Checkbox
                    id={`permission-${role.id}-${permission.id}`}
                    checked={isAssigned}
                    onCheckedChange={(checked) =>
                      handlePermissionChange(permission, checked as boolean)
                    }
                    disabled={isMutating}
                    className="h-3.5 w-3.5"
                  />
                  {isPending && (
                    <Loader2 className="absolute h-3 w-3 animate-spin text-primary" />
                  )}
                </div>
                
                <label
                  htmlFor={`permission-${role.id}-${permission.id}`}
                  className="cursor-pointer select-none font-mono text-xs whitespace-nowrap"
                >
                  {permission.object}:{permission.action}
                </label>
              </div>
            );
          })}
        </div>
      ) : (
        <div className="text-center text-muted-foreground">
          <p>No permissions available</p>
        </div>
      )}
    </div>
  );
}