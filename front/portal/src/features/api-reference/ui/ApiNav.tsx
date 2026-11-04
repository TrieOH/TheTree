import { cn } from "@trieoh/ui-solid";
import { For, Show, createEffect, createSignal } from "solid-js";
import ChevronIcon from "~icons/lucide/chevron-right";
import { operationAnchor, tagAnchor, type ApiModel } from "../model";
import { METHOD_TEXT } from "./MethodBadge";

/** Tag groups that open to their routes; the one being read is highlighted. */
export function ApiNav(props: { api: ApiModel; active: string | null }) {
  return (
    <nav aria-label="Endpoints" class="space-y-0.5 text-[13px]">
      <For each={props.api.tags}>{(tag) => <Group tag={tag} active={props.active} />}</For>
      <a href="#models" class="flex h-8 items-center rounded-lg px-2 text-foreground/75 hover:bg-foreground/5 hover:text-foreground">
        Models
      </a>
    </nav>
  );
}

function Group(props: { tag: ApiModel["tags"][number]; active: string | null }) {
  const holdsActive = () => props.tag.operations.some((op) => op.id === props.active);
  const [open, setOpen] = createSignal(false, { name: "nav-tag-open" });
  // Follows the reader: the group being read opens by itself.
  createEffect(holdsActive, (inside) => {
    if (inside) setOpen(true);
  });

  return (
    <div>
      <div class="flex items-center">
        <button
          type="button"
          aria-expanded={open() ? "true" : "false"}
          aria-label={open() ? `Close ${props.tag.name}` : `Open ${props.tag.name}`}
          onClick={() => setOpen(!open())}
          class="grid size-7 shrink-0 place-items-center rounded-md text-muted-foreground hover:bg-foreground/5"
        >
          <ChevronIcon class={cn("size-3.5 transition-transform", open() && "rotate-90")} />
        </button>
        <a
          href={`#${tagAnchor(props.tag.name)}`}
          class={cn(
            "flex h-8 flex-1 items-center justify-between rounded-lg px-1.5 hover:bg-foreground/5",
            holdsActive() ? "font-medium text-foreground" : "text-foreground/75",
          )}
        >
          {props.tag.name}
          <span class="text-[11px] tabular-nums opacity-70">{props.tag.operations.length}</span>
        </a>
      </div>
      <Show when={open()}>
        <ul class="mb-1 ml-3.5 border-l border-border pl-2">
          <For each={props.tag.operations}>
            {(op) => (
              <li>
                <a
                  href={`#${operationAnchor(op)}`}
                  aria-current={op.id === props.active ? "location" : undefined}
                  class={cn(
                    "flex min-h-7 items-center gap-2 rounded-md px-1.5 py-1 text-[12.5px] hover:bg-foreground/5",
                    op.id === props.active ? "bg-foreground/10 text-foreground" : "text-foreground/75",
                  )}
                >
                  <span class={cn("w-10 shrink-0 font-mono text-[10px] font-semibold uppercase", METHOD_TEXT[op.method])}>
                    {op.method}
                  </span>
                  <span class="min-w-0 truncate">{op.summary}</span>
                </a>
              </li>
            )}
          </For>
        </ul>
      </Show>
    </div>
  );
}
