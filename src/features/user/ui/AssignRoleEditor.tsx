import type { Role } from "@/features/role/model/types";
import { cn } from "@/shared/lib/utils";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Checkbox } from "@/shared/ui/shadcn/checkbox";

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
    <div className="flex flex-col items-center gap-3">
      <div className="text-center w-full">
        <span className="text-primary font-bold">ASSIGN ROLES</span>
        <p className="text-xs text-muted-foreground">
          Select roles to assign. Click the arrow to inspect permissions.
        </p>
      </div>
      <div className="w-full space-y-2">
        {roles.map(role => (
          <label
            key={role.id}
            htmlFor={`check-${role.id}`}
            className={cn(
              "w-full flex items-start p-4 gap-3 text-left",
              "cursor-pointer transition-colors duration-300 hover:bg-secondary/20 border border-transparent",
              selectedRolesMap.get(role.id) && "bg-secondary/20 border-border"
            )}
          >
            <Checkbox 
              id={`check-${role.id}`}
              className="rounded-sm w-5 h-5 cursor-pointer" 
              checked={!!selectedRolesMap.get(role.id)}
              onCheckedChange={() => handleSelectRole(role)}
            />
            <div className="flex gap-x-2 items-baseline flex-wrap">
              <span className="text-primary text-sm font-medium">{role.name}</span>
              <span className="text-muted-foreground text-xs">
                {role.description}
              </span>
            </div>
          </label>
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