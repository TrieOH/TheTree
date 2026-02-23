import type { Scope } from "@/features/scope/model/types";
import { cn } from "@/shared/lib/utils"
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { ChevronRight, MapPin } from "lucide-react"

interface PropsI {
  setCurrentScopeID: (value: string) => void;
  setCurrentType: (value: null) => void;
  allScopes: Scope[];
  currentType: string;
}

export default function ScopeEditorSelector({ setCurrentScopeID, setCurrentType, allScopes, currentType }: PropsI) {
  return (
    <div className="flex flex-col items-center p-4 gap-3">
      <div className="text-center w-full">
        <span className="text-primary">SELECT SCOPE FOR {currentType.toUpperCase()}</span>
        <p className="text-xs text-muted-foreground">What type of access do you want to grant?</p>
      </div>
      <div className="w-full">
        {allScopes.map(scope => (
          <button
            onClick={() => setCurrentScopeID(scope.id)}
            type="button"
            className={cn(
              "w-full flex items-center xs:justify-between justify-center bg-muted rounded-sm p-4",
              "cursor-pointer group transition-colors duration-300 hover:bg-secondary/20"
            )}
          >
            <div className="space-x-1">
              <MapPin className="text-muted-foreground inline"/>
              <span>
                {scope.name.charAt(0).toUpperCase() + scope.name.substring(1)}
              </span>
            </div>
            <ChevronRight
              className={cn(
                "text-muted-foreground group-hover:opacity-100 opacity-0",
                "transition-opacity duration-300 xs:block hidden"
              )}
            />
          </button>
        ))}
      </div>
      <hr className="w-full"/>
      <div className="w-full">
        <ShadowButton value="Back" variant="ghost" onClick={() => setCurrentType(null)}/>
      </div>
    </div>
  )
}