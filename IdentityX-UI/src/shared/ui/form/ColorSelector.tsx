import { useFieldContext } from "@/shared/lib/forms";
import { cn } from "@/shared/lib/utils";

const COMMON_COLORS = [
  // Flat colors
  "#ef4444", "#f97316", "#f59e0b", "#10b981", "#06b6d4", "#3b82f6", "#6366f1", "#8b5cf6", "#d946ef", "#64748b",
  // Gradients
  "linear-gradient(135deg, #f43f5e 0%, #fb7185 100%)",
  "linear-gradient(135deg, #0ea5e9 0%, #38bdf8 100%)",
  "linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)",
  "linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)",
  "linear-gradient(135deg, #10b981 0%, #34d399 100%)",
  "linear-gradient(135deg, #6366f1 0%, #a855f7 100%)"
];

interface PropsI {
  label: string;
}

export default function ColorSelector({ label }: PropsI) {
  const field = useFieldContext<string>();
  const currentValue = field.state.value;

  return (
    <fieldset className="flex flex-col gap-2 mb-4 border-none p-0 m-0">
      <legend className="text-sm font-medium text-foreground">{label}</legend>
      <div className="grid grid-cols-5 sm:grid-cols-8 gap-2 p-1 border rounded-md bg-muted/20">
        {COMMON_COLORS.map((color) => {
          const isSelected = currentValue === color;
          const isGradient = color.includes("linear-gradient");
          
          return (
            <button
              key={color}
              type="button"
              onClick={() => field.handleChange(color)}
              className={cn(
                "h-8 w-8 rounded-full border-2 transition-all hover:scale-110",
                isSelected ? "border-primary shadow-sm ring-2 ring-primary/20" : "border-transparent"
              )}
              style={isGradient ? { backgroundImage: color } : { backgroundColor: color }}
            />
          );
        })}
      </div>
    </fieldset>
  );
}
