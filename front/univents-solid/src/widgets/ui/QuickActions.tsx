import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { For, Match, Show, Switch } from "solid-js";
import CommandIcon from "~icons/lucide/command";
import CopyIcon from "~icons/lucide/copy";
import ExternalLinkIcon from "~icons/lucide/external-link";
import EyeIcon from "~icons/lucide/eye";
import PencilIcon from "~icons/lucide/pencil";
import { cn } from "@trieoh/ui-solid";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCommand = CommandIcon as unknown as IconComp;
const LucideCopy = CopyIcon as unknown as IconComp;
const LucideExternalLink = ExternalLinkIcon as unknown as IconComp;
const LucideEye = EyeIcon as unknown as IconComp;
const LucidePencil = PencilIcon as unknown as IconComp;

export interface QuickActionItem {
  label: string;
  shortcut: string;
  disabled?: boolean;
  variant?: "default" | "destructive";
  icon?: JSX.Element;
  onClick?: () => void;
  to?: "/events/$slug";
  params?: { slug: string };
}

function ActionIcon(props: { label: string }): JSX.Element {
  return (
    <Switch fallback={<LucideExternalLink class="size-4" />}>
      <Match when={props.label.includes("Editar")}>
        <LucidePencil class="size-4" />
      </Match>
      <Match when={props.label.includes("Copiar")}>
        <LucideCopy class="size-4" />
      </Match>
      <Match when={props.label.includes("Publicar")}>
        <LucideEye class="size-4" />
      </Match>
    </Switch>
  );
}

function compactLabel(label: string) {
  return label
    .replace(" evento", "")
    .replace(" edição", "")
    .replace(" público", "")
    .replace(" pública", "")
    .replace(" conta", "");
}

function Shortcut(props: { value: string }) {
  return (
    <kbd class="hidden rounded border border-border/70 bg-muted/70 px-1 py-0.5 font-mono text-[9px] text-muted-foreground sm:inline-block">
      {props.value.replace("Mod", "⌘/Ctrl")}
    </kbd>
  );
}

export interface QuickActionsProps {
  actions: QuickActionItem[];
  title?: string;
  subtitle?: string;
  ariaLabel?: string;
  class?: string;
}

export function QuickActions(props: QuickActionsProps): JSX.Element {
  const title = () => props.title ?? "Ações rápidas";
  const subtitle = () => props.subtitle ?? "Atalhos para as tarefas mais usadas";
  const ariaLabel = () => props.ariaLabel ?? "Atalhos rápidos";

  const getLinkTarget = (action: QuickActionItem) => {
    if (action.to && action.params) {
      return { to: action.to, params: action.params };
    }
    return undefined;
  };

  return (
    <div
      class={cn(
        "order-1 space-y-3 rounded-xl border border-border/60 bg-muted/20 p-3",
        props.class,
      )}
      role="toolbar"
      aria-label={ariaLabel()}
    >
      <div class="flex items-center justify-between gap-3 px-1">
        <div class="flex min-w-0 items-center gap-2">
          <div class="flex size-7 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
            <LucideCommand class="size-3.5" />
          </div>
          <div class="min-w-0">
            <p class="text-xs font-semibold text-foreground">
              {title()}
            </p>
            <p class="truncate text-[11px] text-muted-foreground">
              {subtitle()}
            </p>
          </div>
        </div>
        <span class="hidden shrink-0 text-[11px] text-muted-foreground sm:inline">
          Ctrl/⌘ + tecla
        </span>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <For each={props.actions}>
          {(action) => (
            <Show
              when={getLinkTarget(action)}
              fallback={
                <button
                  type="button"
                  disabled={action.disabled}
                  onClick={() => action.onClick?.()}
                  title={`${action.label} · ${action.shortcut}`}
                  aria-label={action.label}
                  class={cn(
                    "inline-flex h-9 shrink-0 flex-row items-center justify-center gap-1.5 rounded-md border bg-background px-2 text-xs font-medium text-foreground shadow-xs transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-50 sm:h-14! sm:min-w-28! sm:flex-col! sm:gap-1 sm:px-2 sm:py-1.5 sm:text-[11px] sm:leading-tight",
                    action.variant === "destructive"
                      ? "border-destructive/60 text-destructive hover:bg-destructive/10"
                      : "border-border",
                  )}
                >
                  <span class="flex items-center gap-1.5">
                    <Show when={action.icon} fallback={<ActionIcon label={action.label} />}>
                      {(icon) => icon()}
                    </Show>
                    <span>{compactLabel(action.label)}</span>
                  </span>
                  <Shortcut value={action.shortcut} />
                </button>
              }
            >
              {(target) => (
                <Link
                  to={target().to}
                  params={target().params}
                  aria-disabled={action.disabled ? "true" : undefined}
                  title={`${action.label} · ${action.shortcut}`}
                  aria-label={action.label}
                  class={cn(
                    "inline-flex h-9 shrink-0 flex-row items-center justify-center gap-1.5 rounded-md border border-border bg-background px-2 text-xs font-medium text-foreground shadow-xs transition-colors hover:bg-muted sm:h-14! sm:min-w-28! sm:flex-col! sm:gap-1 sm:px-2 sm:py-1.5 sm:text-[11px] sm:leading-tight",
                    action.disabled && "pointer-events-none opacity-50",
                  )}
                >
                  <span class="flex items-center gap-1.5">
                    <Show when={action.icon} fallback={<ActionIcon label={action.label} />}>
                      {(icon) => icon()}
                    </Show>
                    <span>{compactLabel(action.label)}</span>
                  </span>
                  <Shortcut value={action.shortcut} />
                </Link>
              )}
            </Show>
          )}
        </For>
      </div>
    </div>
  );
}
