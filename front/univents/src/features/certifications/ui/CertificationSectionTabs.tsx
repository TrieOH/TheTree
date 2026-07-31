import { AlertTriangle, Award, FileText } from "lucide-react";
import { cn } from "@/shared/lib/utils";

export type CertificationSection = "templates" | "certificates" | "errors";

export function CertificationSectionTabs({
  active,
  onChange,
}: {
  active: CertificationSection;
  onChange: (section: CertificationSection) => void;
}) {
  const items = [
    { id: "templates" as const, label: "Templates", icon: FileText },
    { id: "certificates" as const, label: "Certificados", icon: Award },
    { id: "errors" as const, label: "Erros de emissão", icon: AlertTriangle },
  ];

  return (
    <nav
      className="flex w-full min-w-0 gap-1 overflow-x-auto border-b border-border"
      aria-label="Certificações"
    >
      {items.map(({ id, label, icon: Icon }) => (
        <button
          key={id}
          type="button"
          onClick={() => onChange(id)}
          className={cn(
            "inline-flex shrink-0 items-center gap-2 border-b-2 px-3 py-2.5 text-sm font-medium transition-colors",
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
