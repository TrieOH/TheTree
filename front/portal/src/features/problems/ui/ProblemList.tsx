import { Link } from "@tanstack/solid-router";
import { Badge, cn } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal } from "solid-js";
import ChevronIcon from "~icons/lucide/chevron-right";
import { CopyButton } from "@/shared/ui/CopyButton";
import {
  byStatus,
  extraDescription,
  problemsOf,
  statusText,
  type ProblemDoc,
  type ProblemSet,
} from "../model";
import { ProblemDetails } from "./ProblemDetails";

export interface ProblemListProps {
  /** Whose problems: a service, or "shared" for the ones every service answers. */
  set: ProblemSet;
  /**
   * The shared problems on a service's list: "link" (default) adds one line
   * pointing at them, "include" lists them too, "hide" leaves them out.
   */
  shared?: "link" | "include" | "hide";
  /** Only these slugs, in the usual order. */
  only?: readonly string[];
  /** Group by status code (default) or one flat list. */
  group?: boolean;
  /** Keep every row's details open, with no toggle. */
  open?: boolean;
  /** Descriptions on the problem's own page (default) or in the list too. */
  description?: "page" | "inline";
}

interface Row {
  problem: ProblemDoc;
  /** Where the problem is documented: its own set, or "shared". */
  set: ProblemSet;
}

export function ProblemList(props: ProblemListProps) {
  const shared = () => (props.set === "shared" ? "hide" : (props.shared ?? "link"));
  const rows = createMemo<Row[]>(() => {
    const own = problemsOf(props.set).map((problem) => ({ problem, set: props.set }));
    const all =
      shared() === "include"
        ? [...own, ...problemsOf("shared").map((problem) => ({ problem, set: "shared" as const }))]
        : own;
    const only = props.only;
    return all
      .filter((row) => !only || only.includes(row.problem.slug))
      .sort((a, b) => a.problem.status - b.problem.status || a.problem.slug.localeCompare(b.problem.slug));
  });
  const groups = () =>
    props.group === false
      ? [{ status: 0, rows: rows() }]
      : byStatus(rows().map((row) => row.problem)).map((group) => ({
          status: group.status,
          rows: rows().filter((row) => row.problem.status === group.status),
        }));

  return (
    <div class="not-prose my-5 space-y-5">
      <For each={groups()}>
        {(group) => (
          <section>
            <Show when={group.status}>
              <h4 class="mb-1.5 flex items-baseline gap-2 text-[12px] font-medium text-muted-foreground">
                <StatusPill status={group.status} />
                {statusText(group.status)}
                <span class="tabular-nums">· {group.rows.length}</span>
              </h4>
            </Show>
            <ul class="divide-y divide-border rounded-xl border border-border">
              <For each={group.rows}>
                {(row) => (
                  <ProblemRow
                    row={row}
                    open={props.open ?? false}
                    inline={props.description === "inline"}
                    showStatus={props.group === false}
                    shared={row.set === "shared" && props.set !== "shared"}
                  />
                )}
              </For>
            </ul>
          </section>
        )}
      </For>
      <Show when={shared() === "link" && !props.only}>
        <p class="text-[13px] text-muted-foreground">
          Also answers the{" "}
          <Link to="/problems/$set" params={{ set: "shared" }} class="underline underline-offset-4">
            {problemsOf("shared").length} shared problems
          </Link>{" "}
          every service does.
        </p>
      </Show>
    </div>
  );
}

export function StatusPill(props: { status: number }) {
  return (
    <span
      class={cn(
        "inline-flex h-5 shrink-0 items-center rounded px-1.5 font-mono text-[11px] font-medium tabular-nums",
        props.status >= 500
          ? "bg-destructive/20 text-destructive dark:text-destructive-foreground"
          : "bg-warn/15 text-warn",
      )}
    >
      {props.status}
    </span>
  );
}

function ProblemRow(props: { row: Row; open: boolean; inline: boolean; showStatus: boolean; shared: boolean }) {
  const [expanded, setExpanded] = createSignal(false, { name: "problem-expanded" });
  const open = () => props.open || expanded();
  const description = () => (props.inline ? extraDescription(props.row.problem) : null);

  return (
    <li class="px-3 py-2.5">
      <div class="flex items-center gap-2">
        <Show when={!props.open}>
          <button
            type="button"
            aria-expanded={open() ? "true" : "false"}
            aria-label={open() ? "Hide details" : "Show details"}
            onClick={() => setExpanded(!expanded())}
            class="grid size-6 shrink-0 place-items-center rounded-md text-muted-foreground hover:bg-foreground/5 hover:text-foreground"
          >
            <ChevronIcon class={cn("size-3.5 transition-transform", open() && "rotate-90")} />
          </button>
        </Show>
        <Show when={props.showStatus}>
          <StatusPill status={props.row.problem.status} />
        </Show>
        <Link
          to="/problems/$set/$slug"
          params={{ set: props.row.set, slug: props.row.problem.slug }}
          class="font-mono text-[13.5px] font-medium hover:underline"
        >
          {props.row.problem.slug}
        </Link>
        <CopyButton text={() => props.row.problem.slug} label="Copy slug" />
        <Show when={props.shared}>
          <Badge variant="outline">shared</Badge>
        </Show>
        <span class="ml-auto hidden truncate text-[13px] text-muted-foreground sm:inline">{props.row.problem.title}</span>
      </div>
      <Show when={description()}>
        <p class="mt-1 pl-8 text-[13px] text-muted-foreground">{description()}</p>
      </Show>
      <Show when={open()}>
        <div class={cn("mt-3", !props.open && "pl-8")}>
          <ProblemDetails problem={props.row.problem} set={props.row.set} />
        </div>
      </Show>
    </li>
  );
}
