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
    <div className="flex flex-col items-center gap-3 text-foreground">
      <div className="text-center w-full">
        <span className="text-primary">SELECT SCOPE FOR {currentType.toUpperCase()}</span>
        <p className="text-xs text-muted-foreground">What type of access do you want to grant?</p>
      </div>
      <div className="w-full">
        {allScopes.map(scope => (
          <div
            key={scope.id}
            onClick={() => setCurrentScopeID(scope.id)}
            className={cn(
              "w-full flex justify-between bg-muted rounded-sm p-4 group",
              "cursor-pointer group transition-colors duration-300 hover:bg-secondary/20"
            )}
          >
            <div className="flex items-center gap-2">
              <MapPin 
                className={cn(
                  "text-muted-foreground inline w-4 h-4 shrink-0",
                  "group-hover:text-primary transition-colors duration-300"
                )}
              />
              <span className="font-medium text-sm">
                {scope.name.charAt(0).toUpperCase() + scope.name.substring(1)}
              </span>
            </div>
            <ChevronRight
              className={cn(
                "text-muted-foreground group-hover:opacity-100 opacity-0",
                "transition-opacity duration-300"
              )}
            />
          </div>
        ))}
      </div>
      <hr className="w-full"/>
      <div className="w-full">
        <ShadowButton value="Back" variant="ghost" onClick={() => setCurrentType(null)}/>
      </div>
    </div>
  )
}