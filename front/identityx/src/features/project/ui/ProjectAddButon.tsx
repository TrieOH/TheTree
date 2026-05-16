import { Plus } from "lucide-react"
import { cn } from "@/shared/lib/utils"

interface PropsI {
  onCreate?: () => void
}

export function ProjectAddButton({ onCreate }: PropsI) {
  return (
    <button 
      type="button"
      className={cn(
        "relative min-w-78 bg-card p-5 text-card-foreground",
        "cursor-pointer transition-all duration-300 ease-out group",
        "border-2 border-border rounded-lg border-dashed",
        "flex flex-col items-center justify-center text-center space-y-4"
      )}
      onClick={onCreate}
    >
      <div 
        className={cn(
          "flex justify-center items-center w-11 h-11 rounded-sm",
          "bg-primary-foreground border border-primary font-bold text-2xl",
          "shadow-[1px_1px_0_0_var(--color-primary)] hover:shadow-[2px_2px_0_0_var(--color-primary)]",
          "transition-all duration-300 ease-out",
        )}
      >
        <Plus size={24} />
      </div>
      <h3 className="font-medium text-xl">Create new project</h3>
    </button>
  )
}