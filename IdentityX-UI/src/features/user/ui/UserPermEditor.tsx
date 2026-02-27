import type { RoleWithPermissions, User } from "../model/types";
import { useState } from "react";
import { scopesQueryOptions } from "@/features/scope/api";
import { useQueries, useQuery } from "@tanstack/react-query";
import PermEditorTypeSelector from "./PermEditorTypeSelector";
import ScopeEditorSelector from "./ScopeEditorSelector";
import AssignRoleEditor from "./AssignRoleEditor";
import { roleQueryOptions, rolePermissionsQueryOptions } from "@/features/role/api";
import type { Role } from "@/features/role/model/types";
import UserPermTree from "./UserPermTree";
import { buildRolePermissionsToNodeTree } from "../lib/node-tree-utils";
import AccessConfirmationPanel from "./AccessConfirmationPanel";


interface UserPermEditorProps {
  project_id: string;
  user: User;
}


export default function UserPermEditor({
  project_id,
  user,
}: UserPermEditorProps) {
  const [currentType, setCurrentType] = useState<null | "Roles" | "Permissions">(null);
  const [currentScopeID, setCurrentScopeID] = useState<null | string>(null);
  const [selectedRolesMap, setSelectedRolesMap] = useState<Map<string, Role>>(new Map());
  const [isReview, setIsReview] = useState(false);
  const [isTheEnd, setIsTheEnd] = useState(false);

  const { data: allScopes = [] } = useQuery(scopesQueryOptions(project_id));
  const { data: allRoles = [] } = useQuery(roleQueryOptions(project_id));

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

  const resetAllStates = () => {
    setCurrentType(null);
    setCurrentScopeID(null);
    setSelectedRolesMap(new Map());
    setIsReview(false);
    setIsTheEnd(false);
  }

  return (
    <div className="cursor-default">
      {currentType === null && <PermEditorTypeSelector setCurrentType={setCurrentType} />}
      {currentType !== null && currentScopeID === null && 
        <ScopeEditorSelector 
          allScopes={allScopes} 
          currentType={currentType} 
          setCurrentScopeID={setCurrentScopeID} 
          setCurrentType={setCurrentType}
        />
      }
      {currentScopeID !== null && currentType === "Roles" && !isReview &&
        <AssignRoleEditor 
          roles={allRoles} 
          setCurrentScopeID={setCurrentScopeID}
          selectedRolesMap={selectedRolesMap}
          handleSelectRole={handleSelectRole}
          setIsReview={setIsReview}
        />
      }
      {currentType === "Roles" && isReview && !isTheEnd &&
        <UserPermTree
          node={buildRolePermissionsToNodeTree(rolesWithPermissions)}
          goBack={() => setIsReview(false)}
          onSubmit={() => {setIsTheEnd(true); }}
        />
      }
      {currentType === "Roles" && isTheEnd && 
        <AccessConfirmationPanel
          title={`Access granted to ${user.email}`}
          subTitle={`on scope ${allScopes.find(item => item.id === currentScopeID)?.name}`}
          state="success"
          onExit={resetAllStates}
        />
      }
      {/* {currentType === "Permissions" && isReview &&
        <UserPermTree
          node={buildRolePermissionsToNodeTree(rolesWithPermissions)}
        />
      } */}
    </div>
  )
}
