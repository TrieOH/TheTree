import * as Icons from "lucide-react";
import { Badge } from "@/shared/ui/shadcn/badge";

export interface VisualMetadata {
  icon?: string;
  color?: string;
  folder?: string;
  description?: string;
  status?: "active" | "restricted" | "beta" | "deprecated";
}

interface PropsI {
  name: string;
  meta?: VisualMetadata;
}

export function MetadataVisualizer({ name, meta: incomingMeta }: PropsI) {
  const defaultMeta: VisualMetadata = {
    icon: "Shield",
    color: "linear-gradient(135deg, #6366f1 0%, #a855f7 100%)",
    description: "No description provided.",
    status: "active"
  };

  const meta: VisualMetadata = { ...defaultMeta, ...incomingMeta };
  
  const IconComponent = (meta.icon && (Icons as unknown as Record<string, Icons.LucideIcon>)[meta.icon]) || Icons.Globe;
  
  const isGradient = meta.color?.includes("linear-gradient");
  const iconBgStyle = isGradient 
    ? { backgroundImage: meta.color } 
    : { backgroundColor: meta.color || "#6366f1" };

  return (
    <div className="flex items-center gap-4 group">
      <div 
        className="relative flex items-center justify-center w-10 h-10 rounded-sm shadow-sm transition-transform group-hover:scale-105"
        style={iconBgStyle}
      >
        <IconComponent size={20} className="text-white drop-shadow-md" />
        
        {meta.status === "restricted" && (
          <div className="absolute -top-1 -right-1 w-3 h-3 bg-destructive border-2 border-white rounded-full" />
        )}
      </div>

      <div className="flex flex-col gap-0.5">
        <div className="flex items-center gap-2">
          <span className="font-semibold text-sm text-foreground tracking-tight">
            {name}
          </span>
          
          {meta.status && meta.status !== "active" && (
            <Badge 
              variant={meta.status === "restricted" ? "destructive" : "outline"}
              className="text-[10px] px-1.5 py-0 h-4 uppercase font-bold"
            >
              {meta.status}
            </Badge>
          )}
        </div>

        {meta.description && (
          <p className="text-xs text-muted-foreground line-clamp-1 max-w-50">
            {meta.description}
          </p>
        )}
      </div>
    </div>
  );
}
