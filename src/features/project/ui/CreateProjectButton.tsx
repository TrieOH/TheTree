import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { FolderPlus } from "lucide-react";
import { projectActions } from "../store";

export default function CreateProjectButton() {
  return (
    <>
      <ShadowButton 
        value="New Project" 
        leftIcon={ <FolderPlus size={20}/> }
        className="xs:flex hidden"
        variant="accent-solid"
        onClick={projectActions.openCreate}
      />
      <ShadowButton
        leftIcon={ <FolderPlus size={16}/> }
        className="xs:hidden flex"
        variant="accent-solid"
        onClick={projectActions.openCreate}
      />
    </>
  )
}