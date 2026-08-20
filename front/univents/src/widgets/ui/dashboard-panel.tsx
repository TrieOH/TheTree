import type { LucideIcon } from "lucide-react";
import type { ReactNode } from "react";

interface DashboardPanelProps {
  title: string;
  description?: string;
  icon?: LucideIcon;
  children: ReactNode;
  className?: string;
}

export function DashboardPanel({
  title,
  description,
  icon: Icon,
  children,
  className = "",
}: DashboardPanelProps) {
  return (
    <section className={`min-w-0 max-w-full space-y-3 ${className}`}>
      <div className="flex min-w-0 items-center gap-3 px-1">
        {Icon ? (
          <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
            <Icon className="size-5" />
          </div>
        ) : null}
        <div className="min-w-0">
          <h2 className="text-base font-semibold tracking-tight">{title}</h2>
          {description ? (
            <p className="mt-0.5 truncate text-xs text-muted-foreground">
              {description}
            </p>
          ) : null}
        </div>
      </div>
      {children}
    </section>
  );
}
