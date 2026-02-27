import type { Permission } from "@/features/permission/model/types";
import { cn } from "@/shared/lib/utils";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Box, ChevronRight, Zap } from "lucide-react";
import { useState, useMemo } from "react";
import { Checkbox } from "@/shared/ui/shadcn/checkbox";

interface PropsI {
  permissions: Permission[];
  setCurrentScopeID: (value: null) => void;
  selectedPermissionsMap: Map<string, Permission>;
  handleSelectPermission: (permission: Permission) => void;
  enableReview: () => void;
}

export default function AssignPermissionEditor({
  permissions,
  setCurrentScopeID,
  selectedPermissionsMap,
  handleSelectPermission,
  enableReview
}: PropsI) {
  const [selectedObject, setSelectedObject] = useState<string | null>(null);

  const objects = useMemo(() => {
    return Array.from(new Set(permissions.map(p => p.object))).sort();
  }, [permissions]);

  const actionsForSelectedObject = useMemo(() => {
    if (!selectedObject) return [];
    return permissions
      .filter(p => p.object === selectedObject)
      .sort((a, b) => a.action.localeCompare(b.action));
  }, [permissions, selectedObject]);

  if (selectedObject) {
    return (
      <div className="flex flex-col items-center gap-3">
        <div className="text-center w-full">
          <span className="text-primary font-bold uppercase tracking-wider">
            Actions for {selectedObject}
          </span>
          <p className="text-xs text-muted-foreground mt-1">
            Select actions you want to grant for this object.
          </p>
        </div>
        <div className="w-full space-y-2 max-h-100 overflow-y-auto pr-1">
           {actionsForSelectedObject.map(perm => (
            <div
              key={perm.id}
              className={cn(
                "w-full flex items-center p-4 gap-3 rounded-sm",
                "cursor-pointer transition-all duration-200 border border-transparent",
                "hover:bg-secondary/20",
                selectedPermissionsMap.has(perm.id) ? "bg-secondary/30 border-primary/20" : "bg-muted/50"
              )}
              onClick={() => handleSelectPermission(perm)}
            >
              <Checkbox 
                className="rounded-sm w-5 h-5 cursor-pointer" 
                checked={selectedPermissionsMap.has(perm.id)}
              />
              <div className="flex items-center gap-1.5">
                <Zap className="w-3.5 h-3.5 text-muted-foreground"/>
                <span className="text-primary text-sm">{perm.action}</span>
              </div>
            </div>
          ))}
        </div>
        <hr className="w-full border-muted"/>
        <div className="w-full flex justify-between items-center">
          <ShadowButton value="Back to Objects" variant="ghost" onClick={() => setSelectedObject(null)}/>
          <ShadowButton 
            value="Done" 
            variant="solid" 
            onClick={() => setSelectedObject(null)}
          />
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center gap-3">
      <div className="text-center w-full">
        <span className="text-primary font-bold uppercase tracking-wider">Select Object</span>
        <p className="text-xs text-muted-foreground mt-1">
          Select an object to choose actions. Build your permission set across multiple objects.
        </p>
      </div>
      <div className="w-full space-y-2 max-h-100 overflow-y-auto pr-1">
        {objects.map(obj => {
          const count = permissions.filter(p => p.object === obj && selectedPermissionsMap.has(p.id)).length;
          return (
            <div
              key={obj}
              className={cn(
                "w-full flex justify-between items-center bg-muted/50 rounded-sm p-4 group",
                "cursor-pointer transition-all duration-200 border border-transparent",
                "hover:bg-secondary/20",
                count > 0 && "bg-secondary/30 border-primary/20"
              )}
              onClick={() => setSelectedObject(obj)}
            >
              <div className="flex items-center gap-3">
                <Box 
                  className={cn(
                    "text-muted-foreground w-4 h-4 shrink-0 transition-colors duration-200",
                    "group-hover:text-primary",
                    count > 0 && "text-primary"
                  )}
                />
                <span className="text-primary text-sm font-medium">{obj}</span>
                {count > 0 && (
                  <span className="text-[10px] bg-primary text-primary-foreground px-2 py-0.5 rounded-full font-bold">
                    {count}
                  </span>
                )}
              </div>
              <ChevronRight
                className={cn(
                  "text-muted-foreground w-4 h-4 transition-all duration-200",
                  "group-hover:translate-x-0.5",
                  count > 0 ? "opacity-100" : "opacity-0 group-hover:opacity-100"
                )}
              />
            </div>
          );
        })}
      </div>
      <hr className="w-full border-muted"/>
      <div className="w-full flex justify-between items-center">
        <ShadowButton value="Back" variant="ghost" onClick={() => setCurrentScopeID(null)}/>
        <ShadowButton 
          value={`Review Permissions (${selectedPermissionsMap.size})`} 
          variant="solid"
          onClick={enableReview}
          disabled={selectedPermissionsMap.size <= 0}
        />
      </div>  
    </div>
  )
}
