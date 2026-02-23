import type { Role } from "@/features/role/model/types";
import { cn } from "@/shared/lib/utils";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Checkbox } from "@/shared/ui/shadcn/checkbox";
import { ChevronRight } from "lucide-react";

interface PropsI {
  roles: Role[];
  setCurrentScopeID: (value: null) => void;
  selectedRolesMap: Map<string, Role>;
  handleSelectRole: (role: Role) => void;
  setIsReview: (value: boolean) => void;
}

export default function AssignRoleEditor({ 
  roles,
  setCurrentScopeID,
  selectedRolesMap,
  handleSelectRole,
  setIsReview
}: PropsI) {
  return (
    <div className="flex flex-col items-center p-4 gap-3">
      <div className="text-center w-full">
        <span className="text-primary">ASSIGN ROLES</span>
        <p className="text-xs text-muted-foreground">
          Select roles to assign. Click the arrow to inspect permissions.
        </p>
      </div>
      <div className="w-full space-y-2">
        {roles.map(role => (
          <div 
            className={cn(
              "w-full flex items-center xs:justify-between justify-center bg-muted rounded-sm p-4",
              "cursor-pointer transition-colors duration-300 hover:bg-secondary/20 border border-transparent",
              selectedRolesMap.get(role.id) && "bg-secondary/20 border-border"
            )}
            onClick={() => handleSelectRole(role)}
          >
            <div className="flex items-center gap-2">
              <Checkbox 
                className="rounded-sm w-5 h-5 cursor-pointer" 
                checked={!!selectedRolesMap.get(role.id)}
              />
              <span className="text-primary">{role.name}</span>
              <p className="text-muted-foreground text-xs">{role.description}</p>
            </div>
            <ChevronRight className="text-muted-foreground" />
          </div>
        ))}
      </div>
      <hr className="w-full"/>
      <div className="w-full flex justify-between items-center">
        <ShadowButton value="Back" variant="ghost" onClick={() => setCurrentScopeID(null)}/>
        <ShadowButton 
          value={`Assign Roles (${selectedRolesMap.size})`} 
          variant="solid"
          onClick={() => setIsReview(true)}
          disabled={selectedRolesMap.size <= 0}
        />
      </div>
    </div>
  )
}