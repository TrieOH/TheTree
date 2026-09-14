import type { JSX } from "@solidjs/web";
import { For, Match, Show, Switch } from "solid-js";
import { Link } from "@tanstack/solid-router";
import CommandIcon from "~icons/lucide/command";
import CopyIcon from "~icons/lucide/copy";
import ExternalLinkIcon from "~icons/lucide/external-link";
import EyeIcon from "~icons/lucide/eye";
import PencilIcon from "~icons/lucide/pencil";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCommand = CommandIcon as unknown as IconComp;
const LucideCopy = CopyIcon as unknown as IconComp;
const LucideExternalLink = ExternalLinkIcon as unknown as IconComp;
const LucideEye = EyeIcon as unknown as IconComp;
const LucidePencil = PencilIcon as unknown as IconComp;

export interface EditionQuickAction {
  label: string;
  shortcut: string;
  disabled?: boolean;
  variant: "default" | "destructive";
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
    .replace(" edição", "")
    .replace(" público", "")
    .replace(" pública", "");
}

function Shortcut(props: { value: string }) {
  return (
    <kbd class="hidden rounded border border-border/70 bg-muted/70 px-1 py-0.5 font-mono text-[9px] text-muted-foreground sm:inline-block">
      {props.value.replace("Mod", "⌘/Ctrl")}
    </kbd>
  );
}

export function EditionQuickActions(props: {
  actions: EditionQuickAction[];
}): JSX.Element {
  return (
    <div
      class="order-1 space-y-3 rounded-xl border border-border/60 bg-muted/20 p-3"
      role="toolbar"
      aria-label="Atalhos da edição"
    >
      <div class="flex items-center justify-between gap-3 px-1">
        <div class="flex min-w-0 items-center gap-2">
          <div class="flex size-7 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
            <LucideCommand class="size-3.5" />
          </div>
          <div class="min-w-0">
            <p class="text-xs font-semibold text-foreground">
              Ações rápidas
            </p>
            <p class="truncate text-[11px] text-muted-foreground">
              Atalhos para as tarefas mais usadas
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
              when={action.to && action.params}
              fallback={
                <button
                  type="button"
                  disabled={action.disabled}
                  onClick={() => action.onClick?.()}
                  title={`${action.label} · ${action.shortcut}`}
                  aria-label={action.label}
                  class="flex h-9 shrink-0 items-center gap-1.5 rounded-md border border-border bg-card px-2 text-xs font-medium transition-colors hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 sm:h-14 sm:min-w-28 sm:flex-col sm:justify-center sm:gap-1 sm:px-2 sm:py-1.5 sm:text-[11px] sm:leading-tight"
                >
                  <span class="flex items-center gap-1.5">
                    <ActionIcon label={action.label} />
                    {compactLabel(action.label)}
                  </span>
                  <Shortcut value={action.shortcut} />
                </button>
              }
            >
              <Link
                to={action.to!}
                params={action.params!}
                title={`${action.label} · ${action.shortcut}`}
                aria-label={action.label}
                class="flex h-9 shrink-0 items-center gap-1.5 rounded-md border border-border bg-card px-2 text-xs font-medium transition-colors hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:h-14 sm:min-w-28 sm:flex-col sm:justify-center sm:gap-1 sm:px-2 sm:py-1.5 sm:text-[11px] sm:leading-tight"
              >
                <span class="flex items-center gap-1.5">
                  <ActionIcon label={action.label} />
                  {compactLabel(action.label)}
                </span>
                <Shortcut value={action.shortcut} />
              </Link>
            </Show>
          )}
        </For>
      </div>
    </div>
  );
}
