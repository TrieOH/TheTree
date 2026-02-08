import { cn } from "@/shared/lib/utils";
import { useNavigate } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";

interface PropsI {
  value: string;
  to?: string;
}

function goBack(
  navigate: ReturnType<typeof useNavigate>,
  to?: string
) {
  if (to) navigate({ to })
  else history.back()
}


export default function BackButton({value, to} : PropsI) {
  const navigate = useNavigate()

  return (
    <button
      type="button"
      onClick={() => goBack(navigate, to)}
      className={cn(
        "flex items-center group transition-all duration-200 cursor-pointer",
        "text-muted-foreground hover:text-foreground"
      )}
    >
      <div 
        className={cn(
          "flex items-center justify-center pr-1 mr-2",
          "border-r border-border/60 group-hover:border-border"
        )}
      >
        <ArrowLeft 
          size={18} 
          className="transition-transform duration-200 group-hover:-translate-x-px" 
        />
      </div>

      <span className="text-sm font-medium tracking-tight">
        {value}
      </span>
    </button>
  )
}