import { cn } from "@trieoh/ui-solid";
import { For, Show, createSignal } from "solid-js";
import ChevronIcon from "~icons/lucide/chevron-right";
import { operationAnchor, tagAnchor, type ApiModel } from "../model";
import { METHOD_TEXT } from "./MethodBadge";

/**
 * Tag groups that open to their routes. The group and route being read are
 * highlighted; groups only open when the reader opens them.
 */
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
          // Where the nav scrolls to while the group is closed.
          data-nav-tag={props.tag.name}
          aria-current={holdsActive() && !open() ? "location" : undefined}
          class={cn(
            "relative flex h-8 flex-1 items-center justify-between rounded-lg px-1.5 hover:bg-foreground/5",
            holdsActive() ? "bg-foreground/10 font-medium text-foreground" : "text-foreground/75",
          )}
        >
          <Show when={holdsActive()}>
            <span aria-hidden="true" class="absolute left-0 top-1/2 h-4 w-0.5 -translate-y-1/2 rounded-full bg-accent" />
          </Show>
          {props.tag.name}
          <span class="text-[11px] tabular-nums opacity-70">{props.tag.operations.length}</span>
        </a>
      </div>
      <Show when={open()}>
        <ul class="mb-1 ml-3 border-l border-border pl-1">
          <For each={props.tag.operations}>
            {(op) => (
              <li>
                <a
                  href={`#${operationAnchor(op)}`}
                  data-nav-op={op.id}
                  aria-current={op.id === props.active ? "location" : undefined}
                  title={`${op.method.toUpperCase()} ${op.path}\n${op.summary}`}
                  class={cn(
                    "flex min-h-7 items-center gap-1.5 rounded-md px-1.5 py-1 text-[12.5px]",
                    op.id === props.active
                      ? "bg-accent/10 font-medium text-foreground shadow-[inset_2px_0_0_var(--accent)]"
                      : "text-foreground/75 hover:bg-foreground/5",
                  )}
                >
                  <span class={cn("w-8 shrink-0 font-mono text-[10px] font-semibold uppercase", METHOD_TEXT[op.method])}>
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
