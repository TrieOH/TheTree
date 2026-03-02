import { useFieldContext } from "@/shared/lib/forms";
import * as Icons from "lucide-react";
import { cn } from "@/shared/lib/utils";

const COMMON_ICONS = [
  "Shield", "Lock", "Key", "Users", "Eye", "Edit3", "Settings", "Zap", 
  "Globe", "Database", "Bell", "Mail", "Cloud", "HardDrive", "Layout", 
  "Search", "Star", "Heart", "Trash2", "FileText", "Folder", "Image"
];

interface PropsI {
  label: string;
}

export default function IconSelector({ label }: PropsI) {
  const field = useFieldContext<string>();
  const currentValue = field.state.value;

  return (
    <fieldset className="flex flex-col gap-2 mb-4 border-none p-0 m-0">
      <legend className="text-sm font-medium text-foreground">{label}</legend>
      <div className="grid grid-cols-6 sm:grid-cols-8 gap-2 max-h-40 overflow-y-auto p-1">
        {COMMON_ICONS.map((iconName) => {
          const Icon = (Icons as unknown as Record<string, Icons.LucideIcon>)[iconName];
          const isSelected = currentValue === iconName;
          
          return (
            <button
              key={iconName}
              type="button"
              title={iconName}
              onClick={() => field.handleChange(iconName)}
              className={cn(
                "p-2 rounded-md flex items-center justify-center transition-all border",
                isSelected 
                  ? "bg-primary text-primary-foreground border-primary shadow-sm" 
                  : "bg-white border-transparent hover:border-slate-300 text-slate-600 shadow-xs"
              )}
            >
              {Icon && <Icon size={18} />}
            </button>
          );
        })}
      </div>
    </fieldset>
  );
}
