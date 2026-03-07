import type { RoleWithPermissions, User } from "../model/types";
import { useState } from "react";
import { scopesQueryOptions } from "@/features/scope/api";
import {
  useQueries,
  useQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import {
  givePermissionToUserFn,
  giveRoleToUserFn,
  userPermissionsQueryOptions,
  userRolesQueryOptions,
} from "@/features/user/api";
import PermEditorTypeSelector from "./PermEditorTypeSelector";
import ScopeEditorSelector from "./ScopeEditorSelector";
import AssignRoleEditor from "./AssignRoleEditor";
import { roleQueryOptions, rolePermissionsQueryOptions } from "@/features/role/api";
import type { Role } from "@/features/role/model/types";
import UserPermTree from "./UserPermTree";
import { buildDirectPermissionsToNodeTree, buildRolePermissionsToNodeTree } from "../lib/node-tree-utils";
import AccessConfirmationPanel from "./AccessConfirmationPanel";
import { permissionsQueryOptions } from "@/features/permission/api";
import AssignPermissionEditor from "./AssignPermissionEditor";
import type { Permission } from "@/features/permission/model/types";
import CurrentAccessList from "./CurrentAccessList";


interface UserPermEditorProps {
  project_id: string;
  user: User;
}

export default function UserPermEditor({
  project_id,
  user,
}: UserPermEditorProps) {
  const [currentType, setCurrentType] = useState<null | "Roles" | "Permissions" | "Current">(null);
  const [currentScopeID, setCurrentScopeID] = useState<string | null | undefined>(undefined);
  const [selectedRolesMap, setSelectedRolesMap] = useState<Map<string, Role>>(new Map());
  const [selectedPermissionsMap, setSelectedPermissionsMap] = useState<Map<string, Permission>>(new Map());
  const [isReview, setIsReview] = useState(false);
  const [isTheEnd, setIsTheEnd] = useState(false);
  const [isError, setIsError] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  const queryClient = useQueryClient();

  const givePermissionMutation = useMutation({
    mutationFn: ({ permission_id, scope_id }: { permission_id: string, scope_id: string | null }) =>
      givePermissionToUserFn(user, permission_id, scope_id),
  });

  const giveRoleMutation = useMutation({
    mutationFn: ({ role_id, scope_id }: { role_id: string, scope_id: string | null }) =>
      giveRoleToUserFn(user, role_id, scope_id)
  });

  const { data: allScopes = [] } = useQuery(scopesQueryOptions(project_id));
  const { data: allRoles = [] } = useQuery(roleQueryOptions(project_id));
  const { data: allPermissions = [] } = useQuery(permissionsQueryOptions(project_id));

  const { data: userCurrentRoles = [] } = useQuery(userRolesQueryOptions(project_id, user.id));
  const { data: userCurrentPermissionsForScope = [] } = useQuery({
    ...userPermissionsQueryOptions(project_id, user.id, currentScopeID ?? null),
    enabled: currentScopeID !== undefined,
  });

  const assignedRolesInCurrentScope = new Set(
    userCurrentRoles
      .filter(role => role.scope_id === currentScopeID)
      .map(role => role.id)
  );
  const availableRoles = allRoles.filter(role => !assignedRolesInCurrentScope.has(role.id));

  const userCurrentPermissionIdsForScope = new Set(userCurrentPermissionsForScope.map(permission => permission.id));
  const availablePermissions = allPermissions.filter(permission => !userCurrentPermissionIdsForScope.has(permission.id));

  const rolePermissionsQueries = useQueries({
    queries: [...selectedRolesMap.values()].map((role) =>
      rolePermissionsQueryOptions(project_id, role.id)
    ),
  });

  const rolesWithPermissions: RoleWithPermissions[] =
    [...selectedRolesMap.values()].map((role, index) => {
      const permissionsForRole = rolePermissionsQueries[index]?.data || [];

      return {
        role,
        permissions: permissionsForRole,
      };
    }
  );


  const handleSelectRole = (role: Role) => {
    setSelectedRolesMap(prev => {
      const newState = new Map(prev);
      if (newState.has(role.id)) newState.delete(role.id);
      else newState.set(role.id, role);
      return newState;
    });
  };

  const handleSelectPermission = (permission: Permission) => {
    setSelectedPermissionsMap(prev => {
      const newState = new Map(prev);
      if (newState.has(permission.id)) newState.delete(permission.id);
      else newState.set(permission.id, permission);
      return newState;
    });
  };

  const resetAllStates = () => {
    setCurrentType(null);
    setCurrentScopeID(undefined);
    setSelectedRolesMap(new Map());
    setSelectedPermissionsMap(new Map());
    setIsReview(false);
    setIsTheEnd(false);
  }

  const handleGrantRoles = async () => {
    if (currentScopeID === undefined) return;
    const rolePromises = [...selectedRolesMap.values()].map((role) =>
      giveRoleMutation.mutateAsync({
        role_id: role.id,
        scope_id: currentScopeID,
      })
    );

    const results = await Promise.allSettled(rolePromises);
    await queryClient.invalidateQueries({
      queryKey: userRolesQueryOptions(project_id, user.id).queryKey
    });
    await queryClient.invalidateQueries({
      queryKey: userPermissionsQueryOptions(project_id, user.id, currentScopeID).queryKey
    });
    const hasError = results.some((result) => result.status === "rejected");
    setIsError(hasError);
    if (hasError) setErrorMessage("Some roles could not be assigned.");
    setIsTheEnd(true);
  };

  const handleGrantPermissions = async () => {
    if (currentScopeID === undefined) return;
    const permissionPromises = [...selectedPermissionsMap.values()].map(
      (permission) =>
        givePermissionMutation.mutateAsync({
          permission_id: permission.id,
          scope_id: currentScopeID,
        }
      )
    );

    const results = await Promise.allSettled(permissionPromises);
    await queryClient.invalidateQueries({
      queryKey: userPermissionsQueryOptions(project_id, user.id, currentScopeID).queryKey
    });
    const hasError = results.some((result) => result.status === "rejected");
    setIsError(hasError);
    if (hasError) setErrorMessage("Some permissions could not be assigned.");
    setIsTheEnd(true);
  };

  const getScopeName = () => {
    return currentScopeID === null ? "root" : 
      `${allScopes.find(item => item.id === currentScopeID)?.name}`
  }

  return (
    <>
      {currentType === null && <PermEditorTypeSelector setCurrentType={setCurrentType} />}

      {currentType === "Current" && (
        <CurrentAccessList 
          user={user} 
          project_id={project_id} 
          onBack={() => setCurrentType(null)} 
          allScopes={allScopes}
        />
      )}

      {currentType !== null && currentType !== "Current" && currentScopeID === undefined && 
        <ScopeEditorSelector 
          allScopes={allScopes} 
          currentType={currentType} 
          setCurrentScopeID={setCurrentScopeID} 
          setCurrentType={setCurrentType}
        />
      }

      {/* Roles Section */}
      {currentScopeID !== undefined && currentType === "Roles" && !isReview &&
        <AssignRoleEditor 
          roles={availableRoles} 
          setCurrentScopeID={setCurrentScopeID}
          selectedRolesMap={selectedRolesMap}
          handleSelectRole={handleSelectRole}
          setIsReview={setIsReview}
        />
      }
      {currentType === "Roles" && isReview && !isTheEnd &&
        <div className="w-full max-h-112.5 overflow-auto">
          <UserPermTree
            node={buildRolePermissionsToNodeTree(rolesWithPermissions, getScopeName())}
            goBack={() => setIsReview(false)}
            onSubmit={handleGrantRoles}
          />
        </div>
      }
      {currentType === "Roles" && isTheEnd && (
        <AccessConfirmationPanel
          title={isError ? "Error assigning roles" : `Access granted to ${user.email}`}
          subTitle={isError ? errorMessage : `on ${currentScopeID === null ? "root" : `scope ${allScopes.find(item => item.id === currentScopeID)?.name}`}`}
          state={isError ? "error" : "success"}
          onExit={resetAllStates}
        />
      )}

      {/* Permissions Sections */}
      {currentScopeID !== undefined && currentType === "Permissions" && !isReview &&
        <AssignPermissionEditor 
          permissions={availablePermissions} 
          setCurrentScopeID={setCurrentScopeID}
          selectedPermissionsMap={selectedPermissionsMap}
          handleSelectPermission={handleSelectPermission}
          enableReview={() => setIsReview(true)}
        />
      }
      {currentType === "Permissions" && isReview && !isTheEnd &&
        <div className="w-full max-h-112.5 overflow-auto">
          <UserPermTree
            node={buildDirectPermissionsToNodeTree([...selectedPermissionsMap.values()], getScopeName())}
            goBack={() => setIsReview(false)}
            onSubmit={handleGrantPermissions}
          />
        </div>
      }
      {currentType === "Permissions" && isTheEnd && (
        <AccessConfirmationPanel
          title={isError ? "Error assigning permissions" : `Access granted to ${user.email}`}
          subTitle={isError ? errorMessage : `on ${currentScopeID === null ? "root" : `scope ${allScopes.find(item => item.id === currentScopeID)?.name}`}`}
          state={isError ? "error" : "success"}
          onExit={resetAllStates}
        />
      )}
    </>
  )
}
