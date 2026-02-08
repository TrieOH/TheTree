import { cn } from "@/shared/lib/utils"
import type { Project } from "../model/types"
import { ArrowRight, Clock, ClockFading, Edit, Trash2 } from "lucide-react"
import { formatDate } from "@/shared/lib/date-utils"
import { projectActions } from "../store"
import { useNavigate } from "@tanstack/react-router"

interface PropsI {
  data: Project
}

export default function ProjectCard({ data }: PropsI) {
  const navigate = useNavigate({ from: '/projects' })

  const handleProjectCardClick = () => {
    projectActions.setCurrentProjectId(data.id);
    navigate({ to: '/schemas' });
  };

  return (
    <button 
      type="button"
      className={cn(
        "relative min-w-78 bg-card p-5 text-card-foreground",
        "cursor-pointer transition-all duration-300 ease-out group",
        "border-2 border-border rounded-lg",
        "shadow-[1px_1px_0_0_var(--color-border)] hover:shadow-[2px_2px_0_0_var(--color-border)]"
      )}
      onClick={handleProjectCardClick}
    >
      {/* Top */}
      <div className="flex items-center gap-2.5 mb-5">
        <span 
          className={cn(
            "flex justify-center items-center w-11 h-11 rounded-sm",
            "bg-primary-foreground border border-primary font-bold text-2xl",
            "shadow-[1px_1px_0_0_var(--color-primary)] hover:shadow-[2px_2px_0_0_var(--color-primary)]",
            "transition-all duration-300 ease-out",
          )}
        >
          {data.project_name.charAt(0).toUpperCase()}
        </span>
        <span className="font-medium text-xl truncate max-w-50">
          {data.project_name}
        </span>
      </div>

      {/* Middle */}
      <div className="space-y-1.5">
        <div className="flex items-center gap-1.5 text-sm">
          <Clock size={16}/>
          <span>Created at {formatDate(data.created_at)}</span>
        </div>
        <div className="flex items-center gap-1.5 text-sm">
          <ClockFading size={16}/>
          <span>Updated at {formatDate(data.updated_at)}</span>
        </div>
      </div>

      {/* Bottom */}
      <div 
        className={cn(
          "flex opacity-0 justify-center gap-4 mt-4 text-muted-foreground",
          "group-hover:opacity-100 duration-300"
        )}
      >
        <Edit 
          className="hover:text-card-foreground duration-300"
          onClick={(e) => {
            e.stopPropagation();
            projectActions.openEdit(data);
          }}
        />
        <Trash2 
          className="hover:text-card-foreground duration-300"
          onClick={(e) => {
            e.stopPropagation();
            projectActions.openDelete(data);
          }}
        />
      </div>

      {/* Active State */}
      <div 
        className={cn(
          "absolute w-3.5 h-3.5 rounded-full top-2 right-2",
          data.is_active ? "bg-green-300" : "bg-destructive"
        )}
      />


      <ArrowRight
        className={cn(
          "absolute bottom-2 right-2 text-muted-foreground",
          "group-hover:text-card-foreground duration-300"
        )}
      />
    </button>
  )
}