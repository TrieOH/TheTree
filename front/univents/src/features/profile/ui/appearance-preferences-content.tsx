import { Monitor, Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { cn } from "@/shared/lib/utils";

export function AppearancePreferencesContent() {
  const { theme, setTheme } = useTheme();

  return (
    <div className="flex flex-col gap-3 rounded-sm border border-border bg-background p-3 text-left">
      <div className="min-w-0">
        <p className="w-full text-left text-xs text-muted-foreground sm:text-center">
          Personalize a aparência do sistema de acordo com a sua preferência.
        </p>
      </div>

      <div className="grid gap-2 sm:grid-cols-3 sm:gap-3">
        <button
          type="button"
          onClick={() => setTheme("light")}
          className={cn(
            "flex items-center gap-3 rounded-md border-2 bg-popover p-3 text-left",
            "sm:flex-col sm:justify-center sm:gap-0 sm:p-4 sm:text-center",
            "transition-colors hover:bg-accent hover:text-accent-foreground",
            theme === "light" ? "border-primary" : "border-muted",
          )}
        >
          <Sun className="size-5 shrink-0 sm:mb-2" />
          <span className="text-xs font-semibold">Claro</span>
        </button>

        <button
          type="button"
          onClick={() => setTheme("dark")}
          className={cn(
            "flex items-center gap-3 rounded-md border-2 bg-popover p-3 text-left",
            "sm:flex-col sm:justify-center sm:gap-0 sm:p-4 sm:text-center",
            "transition-colors hover:bg-accent hover:text-accent-foreground",
            theme === "dark" ? "border-primary" : "border-muted",
          )}
        >
          <Moon className="size-5 shrink-0 sm:mb-2" />
          <span className="text-xs font-semibold">Escuro</span>
        </button>

        <button
          type="button"
          onClick={() => setTheme("system")}
          className={cn(
            "flex items-center gap-3 rounded-md border-2 bg-popover p-3 text-left",
            "sm:flex-col sm:justify-center sm:gap-0 sm:p-4 sm:text-center",
            "transition-colors hover:bg-accent hover:text-accent-foreground",
            theme === "system" ? "border-primary" : "border-muted",
          )}
        >
          <Monitor className="size-5 shrink-0 sm:mb-2" />
          <span className="text-xs font-semibold">Sistema</span>
        </button>
      </div>
    </div>
  );
}
