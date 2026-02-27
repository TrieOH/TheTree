import { cn } from "@/shared/lib/utils";
import { ChevronRight, Key, Shield } from "lucide-react";

interface PropsI {
  setCurrentType: (value: "Roles" | "Permissions") => void
}

export default function PermEditorTypeSelector({ setCurrentType }: PropsI) {
  return (
    <div className="flex flex-col items-center gap-3">
      <div className="text-center w-full">
        <span className="text-primary">GRANT ACCESS</span>
        <p className="text-xs text-muted-foreground">What type of access do you want to grant?</p>
      </div>
      <div
        onClick={() => setCurrentType("Roles")}
        className={cn(
          "w-full flex items-center xs:justify-between justify-center bg-muted rounded-md",
          "cursor-pointer group transition-colors duration-300 hover:bg-secondary/20 p-4"
        )}
      >
        <div className="flex xs:flex-row flex-col items-center gap-2">
          <Shield className="w-10 h-10 text-primary"/>
          <div className="flex flex-col items-start">
            <span className="w-full xs:text-start text-sm text-center font-medium">
              Roles
            </span>
            <p className="text-muted-foreground text-xs xs:text-start text-center">
              Assign predefined role blundes
            </p>
          </div>
        </div>
        <ChevronRight 
          className={cn(
            "text-muted-foreground group-hover:opacity-100 opacity-0",
            "transition-opacity duration-300 xs:block hidden"
          )}
        />
      </div>
      <div 
        onClick={() => setCurrentType("Permissions")}
        className={cn(
          "w-full flex items-center xs:justify-between justify-center bg-muted rounded-md p-4",
          "cursor-pointer group transition-colors duration-300 hover:bg-secondary/20"
        )}
      >
        <div className="flex xs:flex-row flex-col items-center gap-2">
          <Key className="w-10 h-10 text-accent" />
          <div className="flex flex-col items-start">
            <span className="w-full xs:text-start text-sm text-center font-medium">
              Permissions
            </span>
            <p className="text-muted-foreground text-xs xs:text-start text-center">
              Fine-grained object:action access
            </p>
          </div>
        </div>
        <ChevronRight 
          className={cn(
            "text-muted-foreground group-hover:opacity-100 opacity-0",
            "transition-opacity duration-300 xs:block hidden"
          )}
        />
      </div>
    </div>
  )
}