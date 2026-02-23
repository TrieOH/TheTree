import type { User } from "../model/types";
import { useState } from "react";
import { scopesQueryOptions } from "@/features/scope/api";
import { useQuery } from "@tanstack/react-query";
import PermEditorTypeSelector from "./PermEditorTypeSelector";
import ScopeEditorSelector from "./ScopeEditorSelector";
import AssignRoleEditor from "./AssignRoleEditor";
import { roleQueryOptions } from "@/features/role/api";
import type { Role } from "@/features/role/model/types";

interface UserPermEditorProps {
  project_id: string;
  user: User;
}

export default function UserPermEditor({
  project_id,
  // user,
}: UserPermEditorProps) {
  const [currentType, setCurrentType] = useState<null | "Roles" | "Permissions">(null);
  const [currentScopeID, setCurrentScopeID] = useState<null | string>(null);
  const [selectedRolesMap, setSelectedRolesMap] = useState<Map<string, Role>>(new Map());
  const [isReview, setIsReview] = useState(false);

  const { data: allScopes = [] } = useQuery(scopesQueryOptions(project_id));
  const { data: allRoles = [] } = useQuery(roleQueryOptions(project_id));

  const handleSelectRole = (role: Role) => {
    setSelectedRolesMap(prev => {
      const newState = new Map(prev);
      if (newState.has(role.id)) newState.delete(role.id);
      else newState.set(role.id, role);
      return newState;
    });
  };

  return (
    <>
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
      {currentType === "Roles" && isReview}
    </>
  )
}