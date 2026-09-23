import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import MailIcon from "~icons/lucide/mail";
import PenLineIcon from "~icons/lucide/pen-line";
import { cn } from "@trieoh/ui-solid";

const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;

export type SignatureSection = "signatures" | "invites";

export interface SignatureSectionTabsProps {
  active: SignatureSection;
  onChange: (section: SignatureSection) => void;
  signaturesCount?: number;
  invitesCount?: number;
}

export function SignatureSectionTabs(props: SignatureSectionTabsProps): JSX.Element {
  const items = () => [
    {
      id: "signatures" as const,
      label: "Assinaturas",
      icon: PenLine,
      count: props.signaturesCount,
    },
    {
      id: "invites" as const,
      label: "Convites de assinatura",
      icon: Mail,
      count: props.invitesCount,
    },
  ];

  return (
    <nav
      class="flex w-full min-w-0 items-center gap-1 overflow-x-auto border-b border-border scrollbar-none"
      aria-label="Assinaturas da edição"
    >
      <For each={items()}>
        {({ id, label, icon: Icon, count }) => {
          const isActive = () => props.active === id;
          return (
            <button
              type="button"
              onClick={() => props.onChange(id)}
              class={cn(
                "group relative inline-flex shrink-0 items-center gap-2 border-b-2 px-3.5 py-2.5 text-sm font-medium transition-colors cursor-pointer",
                isActive()
                  ? "border-primary text-foreground font-semibold"
                  : "border-transparent text-muted-foreground hover:text-foreground",
              )}
              aria-current={isActive() ? "page" : undefined}
            >
              <Icon
                class={cn(
                  "size-4 shrink-0 transition-colors",
                  isActive()
                    ? "text-primary"
                    : "text-muted-foreground group-hover:text-foreground",
                )}
              />
              <span class="whitespace-nowrap">{label}</span>
              <Show when={typeof count === "number"}>
                <span
                  class={cn(
                    "ml-1 inline-flex shrink-0 items-center rounded-full px-2 py-0.5 text-[11px] font-mono leading-none transition-colors",
                    isActive()
                      ? "bg-primary/15 text-primary font-medium"
                      : "bg-muted text-muted-foreground",
                  )}
                >
                  {count}
                </span>
              </Show>
            </button>
          );
        }}
      </For>
    </nav>
  );
}
