import { cn } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import { findService } from "@/features/catalog/services";
import { PRIORITY_LABEL, type Issue, type StateType } from "../model";

/** Linear-style state ring: hollow for todo, half for in progress, full for done. */
export function StateMark(props: { type: StateType | string; class?: string }) {
  return (
    <span
      aria-hidden="true"
      class={cn(
        "inline-block size-2.5 shrink-0 rounded-full border-[1.5px] border-muted-foreground",
        props.type === "backlog" && "border-dashed",
        props.type === "started" && "border-accent bg-[conic-gradient(var(--accent)_0_50%,transparent_50%)]",
        props.type === "completed" && "border-ok bg-ok",
        props.class,
      )}
    />
  );
}

function PriorityMark(props: { priority: Issue["priority"] }) {
  return (
    <span
      title={PRIORITY_LABEL[props.priority]}
      aria-label={PRIORITY_LABEL[props.priority]}
      class={cn(
        "inline-grid w-4 shrink-0 place-items-center text-[11px] font-semibold",
        props.priority === 1 ? "text-destructive dark:text-destructive-foreground" : "text-muted-foreground",
      )}
    >
      {props.priority === 0 ? "–" : props.priority === 1 ? "!" : "·".repeat(5 - props.priority)}
    </span>
  );
}

function initials(name: string): string {
  return name
    .split(/\s+/)
    .map((part) => part[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();
}

export function IssueRow(props: { issue: Issue }) {
  const product = () => findService(props.issue.service ?? "")?.name ?? props.issue.project;
  return (
    <li>
      <a
        href={props.issue.url ?? undefined}
        target="_blank"
        rel="noopener"
        class="flex min-h-10 items-center gap-3 px-1 py-2 text-[13.5px] transition-colors hover:bg-foreground/5"
      >
        <PriorityMark priority={props.issue.priority} />
        <span class="w-16 shrink-0 font-mono text-[12px] text-muted-foreground">{props.issue.identifier}</span>
        <span class="min-w-0 flex-1 truncate">{props.issue.title}</span>
        <Show when={product()}>
          <span class="hidden shrink-0 rounded-full border border-border px-2 py-0.5 text-[11px] text-muted-foreground sm:inline">
            {product()}
          </span>
        </Show>
        <Show when={props.issue.estimate !== null}>
          <span class="hidden w-5 shrink-0 text-right text-[12px] tabular-nums text-muted-foreground sm:inline">
            {props.issue.estimate}
          </span>
        </Show>
        <span
          title={props.issue.assignee ?? "Unassigned"}
          class="grid size-6 shrink-0 place-items-center rounded-full bg-muted text-[10px] font-medium text-muted-foreground"
        >
          {props.issue.assignee ? initials(props.issue.assignee) : ""}
        </span>
      </a>
    </li>
  );
}

export function IssueList(props: { issues: readonly Issue[]; empty?: string }) {
  return (
    <Show
      when={props.issues.length > 0}
      fallback={<p class="py-3 text-[13px] text-muted-foreground">{props.empty ?? "No issues."}</p>}
    >
      <ul class="divide-y divide-border border-y border-border">
        <For each={props.issues}>{(issue) => <IssueRow issue={issue} />}</For>
      </ul>
    </Show>
  );
}
