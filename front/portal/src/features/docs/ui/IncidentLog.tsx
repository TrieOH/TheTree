import { Link } from "@tanstack/solid-router";
import { Badge } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import { formatDate } from "../adr";
import { statusWord, type DocLink } from "../catalog";
import { DocStatus } from "../DocStatus";
import { SEVERITY_BADGE, severityLabel, durationOf, serviceName, severityOf } from "../incident";

/** Incident reports, newest first, each dated in the gutter. */
export function IncidentLog(props: { docs: readonly DocLink[] }) {
  return (
    <ol>
      <For each={props.docs}>{(doc) => <IncidentRow doc={doc} />}</For>
    </ol>
  );
}

function IncidentRow(props: { doc: DocLink }) {
  const meta = () => props.doc.meta;
  const severity = () => severityOf(meta());
  const open = () => statusWord(meta().status) !== "resolved";

  return (
    <li class="grid gap-x-8 border-t border-border py-5 sm:grid-cols-[120px_minmax(0,1fr)]">
      <time datetime={meta().date} class="text-[13px] tabular-nums text-muted-foreground sm:pt-0.5">
        {formatDate(meta().date)}
      </time>
      <div class="mt-1 min-w-0 sm:mt-0">
        <Link
          to="/docs/$category/$slug"
          params={{ category: "incident-reports", slug: props.doc.slug }}
          class="rounded-sm text-[16px] font-medium leading-snug hover:underline hover:underline-offset-[3px] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
        >
          {meta().title ?? props.doc.slug}
        </Link>
        <Show when={meta().summary}>
          <p class="mt-1.5 max-w-[62ch] text-[14px] leading-relaxed text-muted-foreground">{meta().summary}</p>
        </Show>
        <div class="mt-3 flex flex-wrap items-center gap-x-4 gap-y-1.5 text-[13px] text-muted-foreground">
          <Show when={severity()}>
            {(level) => (
              <Badge variant={SEVERITY_BADGE[level()]}>
                {severityLabel(level())}
              </Badge>
            )}
          </Show>
          <Show when={open()}>
            <DocStatus status={meta().status} />
          </Show>
          <Show when={durationOf(meta())}>{(span) => <span>Lasted {span()}</span>}</Show>
          <Show when={meta().services?.length}>
            <span>Affected {meta().services!.map(serviceName).join(", ")}</span>
          </Show>
        </div>
      </div>
    </li>
  );
}
