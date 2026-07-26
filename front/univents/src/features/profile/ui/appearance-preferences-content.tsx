import { Monitor, Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { cn } from "@/shared/lib/utils";

export function AppearancePreferencesContent() {
  const { theme, setTheme } = useTheme();

  return (
    <div className="pt-4 pb-2">
      <div className="flex flex-col gap-4 rounded-lg border border-border/50 p-4 bg-card">
        <div className="space-y-1">
          <p className="text-sm font-medium leading-none">Tema da interface</p>
          <p className="text-sm text-muted-foreground">
            Personalize a aparência do sistema de acordo com a sua preferência.
          </p>
        </div>

        <div className="grid grid-cols-3 gap-3 pt-2">
          <button
            type="button"
            onClick={() => setTheme("light")}
            className={cn(
              "flex flex-col items-center justify-center rounded-md border-2 bg-popover p-4",
              "transition-colors hover:bg-accent hover:text-accent-foreground",
              theme === "light" ? "border-primary" : "border-muted",
            )}
          >
            <Sun className="mb-2 size-5" />
            <span className="text-xs font-semibold">Claro</span>
          </button>

          <button
            type="button"
            onClick={() => setTheme("dark")}
            className={cn(
              "flex flex-col items-center justify-center rounded-md border-2 bg-popover",
              "p-4 transition-colors hover:bg-accent hover:text-accent-foreground",
              theme === "dark" ? "border-primary" : "border-muted",
            )}
          >
            <Moon className="mb-2 size-5" />
            <span className="text-xs font-semibold">Escuro</span>
          </button>

          <button
            type="button"
            onClick={() => setTheme("system")}
            className={cn(
              "flex flex-col items-center justify-center rounded-md border-2 bg-popover p-4",
              "transition-colors hover:bg-accent hover:text-accent-foreground",
              theme === "system" ? "border-primary" : "border-muted",
            )}
          >
            <Monitor className="mb-2 size-5" />
            <span className="text-xs font-semibold">Sistema</span>
          </button>
        </div>
      </div>
    </div>
  );
}
