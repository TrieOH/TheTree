import type { JSX } from "@solidjs/web";
import { For } from "solid-js";
import AlertTriangleIcon from "~icons/lucide/alert-triangle";
import AwardIcon from "~icons/lucide/award";
import FileTextIcon from "~icons/lucide/file-text";
import { cn } from "@trieoh/ui-solid";

const AlertTriangle = AlertTriangleIcon as unknown as (props: { class?: string }) => JSX.Element;
const Award = AwardIcon as unknown as (props: { class?: string }) => JSX.Element;
const FileText = FileTextIcon as unknown as (props: { class?: string }) => JSX.Element;

export type CertificationSection = "templates" | "certificates" | "errors";

export function CertificationSectionTabs(props: {
  active: CertificationSection;
  onChange: (section: CertificationSection) => void;
}): JSX.Element {
  const items = [
    { id: "templates" as const, label: "Templates", icon: FileText },
    { id: "certificates" as const, label: "Certificados", icon: Award },
    { id: "errors" as const, label: "Erros de emissão", icon: AlertTriangle },
  ];

  return (
    <nav
      class="flex w-full min-w-0 gap-1 overflow-x-auto border-b border-border"
      aria-label="Certificações"
    >
      <For each={items}>
        {({ id, label, icon: Icon }) => (
          <button
            type="button"
            onClick={() => props.onChange(id)}
            class={cn(
              "inline-flex shrink-0 items-center gap-2 border-b-2 px-3 py-2.5 text-sm font-medium transition-colors cursor-pointer",
              props.active === id
                ? "border-primary text-primary"
                : "border-transparent text-muted-foreground hover:text-foreground",
            )}
            aria-current={props.active === id ? "page" : undefined}
          >
            <Icon class="size-4" />
            {label}
          </button>
        )}
      </For>
    </nav>
  );
}
