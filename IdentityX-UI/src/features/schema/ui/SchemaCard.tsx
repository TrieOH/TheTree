import { cn } from "@/shared/lib/utils"
import type { Schema } from "../model/types"
import { ArrowRight, Clock, ClockFading, Edit, Trash2 } from "lucide-react"
import { formatDate } from "@/shared/lib/date-utils"
import { useNavigate } from "@tanstack/react-router"
import { navigationActions } from "@/features/navigation"
import { schemaActions } from "../store"

interface PropsI {
  data: Schema
}

export default function SchemaCard({ data }: PropsI) {
  const navigate = useNavigate({ from: '/schemas' })

  const handleSchemaCardClick = () => {
    navigationActions.setCurrentProjectId(data.id);
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
      onClick={handleSchemaCardClick}
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
          
        </span>
        <span className="font-medium text-xl truncate max-w-50">
          {data.title}
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
            schemaActions.openEdit(data);
          }}
        />
        <Trash2 
          className="hover:text-card-foreground duration-300"
          onClick={(e) => {
            e.stopPropagation();
            schemaActions.openDelete(data);
          }}
        />
      </div>

      {/* Status */}
      {/* <div 
        className={cn(
          "absolute w-3.5 h-3.5 rounded-full top-2 right-2",
          data.is_active ? "bg-green-300" : "bg-destructive"
        )}
      /> */}


      <ArrowRight
        className={cn(
          "absolute bottom-2 right-2 text-muted-foreground",
          "group-hover:text-card-foreground duration-300"
        )}
      />
    </button>
  )
}