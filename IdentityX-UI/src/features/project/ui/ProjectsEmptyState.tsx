import { ShadowButton } from "@/shared/ui/buttons/ShadowButton"
import { EmptyState } from "@/shared/ui/placeholders/EmptyState"
import { FolderOpen } from "lucide-react"

interface PropsI {
  onCreate?: () => void
}

export function ProjectsEmptyState({ onCreate }: PropsI) {
  return (
    <EmptyState
      icon={FolderOpen}
      title="No projects yet"
      description="Get started by creating your first project to organize your work."
      action={
        onCreate && (
          <ShadowButton onClick={onCreate} value="Create Project"/>
        )
      }
      className="w-full border-2 border-dashed border-border py-10 max-w-xl"
    />
  )
}