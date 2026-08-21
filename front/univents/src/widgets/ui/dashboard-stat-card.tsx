import type { LucideIcon } from "lucide-react";
import type { ReactNode } from "react";

interface DashboardStatCardProps {
  label: string;
  value: ReactNode;
  hint: string;
  icon: LucideIcon;
}

export function DashboardStatCard({
  label,
  value,
  hint,
  icon: Icon,
}: DashboardStatCardProps) {
  return (
    <article className="rounded-lg bg-card p-4 ring-1 ring-foreground/10 transition-shadow hover:shadow-md">
      <div className="flex items-center justify-between text-muted-foreground">
        <p className="text-xs">{label}</p>
        <Icon className="size-4" />
      </div>
      <p className="mt-2 text-xl font-semibold tracking-tight">{value}</p>
      <p className="mt-2 truncate text-xs leading-4 text-muted-foreground">
        {hint}
      </p>
    </article>
  );
}
