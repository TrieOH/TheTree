import { Link } from "@tanstack/react-router";
import type { LucideIcon } from "lucide-react";
import { cn } from "@/shared/lib/utils";

interface SectionTab {
  label: string;
  to: string;
  icon: LucideIcon;
  params: Record<string, string>;
}

interface SectionTabsProps {
  items: SectionTab[];
}

export function SectionTabs({ items }: SectionTabsProps) {
  return (
    <nav
      className="flex w-full gap-1 border-b border-border"
      aria-label="Seção"
    >
      {items.map(({ label, to, icon: Icon, params }) => (
        <Link
          key={to}
          to={to}
          params={params}
          activeOptions={{ exact: to.endsWith("signatures/") }}
          className={cn(
            "inline-flex items-center gap-2 border-b-2 border-transparent px-3 py-2.5",
            "text-sm font-medium text-muted-foreground transition-colors hover:text-foreground",
          )}
          activeProps={{ className: "border-primary text-primary" }}
        >
          <Icon className="size-4" />
          {label}
        </Link>
      ))}
    </nav>
  );
}
