import { cn } from "@/shared/lib/utils";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { FolderPlus } from "lucide-react";

export default function CreateProjectButton() {
  return (
    <>
      <ShadowButton 
        value="New Project" 
        leftIcon={ <FolderPlus size={20}/> }
        className={cn(
          "xs:flex hidden text-primary-foreground bg-primary",
          "shadow-[1px_1px_0_0_var(--color-accent)] hover:shadow-[2px_2px_0_0_var(--color-accent)]"
        )}
        onClick={() => null}
      />
      <ShadowButton
        leftIcon={ <FolderPlus size={16}/> }
        className="xs:hidden flex text-primary-foreground bg-primary"
        onClick={() => null}
      />
    </>
  )
}