import { FileText, Link2 } from "lucide-react";
import { cn } from "@/shared/lib/utils";

export type CertificationSection = "templates" | "links";

export function CertificationSectionTabs({
  active,
  onChange,
}: {
  active: CertificationSection;
  onChange: (section: CertificationSection) => void;
}) {
  const items = [
    { id: "templates" as const, label: "Templates", icon: FileText },
    { id: "links" as const, label: "Vínculos", icon: Link2 },
  ];

  return (
    <nav
      className="flex w-full gap-1 border-b border-border"
      aria-label="Certificações"
    >
      {items.map(({ id, label, icon: Icon }) => (
        <button
          key={id}
          type="button"
          onClick={() => onChange(id)}
          className={cn(
            "inline-flex items-center gap-2 border-b-2 px-3 py-2.5 text-sm font-medium transition-colors",
            active === id
              ? "border-primary text-primary"
              : "border-transparent text-muted-foreground hover:text-foreground",
          )}
          aria-current={active === id ? "page" : undefined}
        >
          <Icon className="size-4" />
          {label}
        </button>
      ))}
    </nav>
  );
}
