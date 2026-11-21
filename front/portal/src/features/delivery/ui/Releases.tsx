import { Badge, EmptyState, StatusDot, cn, type StatusDotTone } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import { commitUrl, pullUrl, runUrl, type CiRun } from "../model";
import { releaseOutcome, runSpan, versionOf, type Release, type ReleaseOutcome } from "../releases";
import { formatDuration, timeAgo } from "../stats";

const external = { target: "_blank", rel: "noopener" } as const;

const OUTCOME: Record<ReleaseOutcome, { tone: StatusDotTone; label: string }> = {
  shipped: { tone: "success", label: "Shipped" },
  failed: { tone: "destructive", label: "Failed" },
  cancelled: { tone: "muted", label: "Cancelled" },
  mixed: { tone: "warning", label: "Partly shipped" },
};

const BAR: Record<CiRun["status"], string> = {
  success: "bg-success",
  failure: "bg-destructive dark:bg-destructive-foreground",
  cancelled: "bg-muted-foreground/40",
};

function runState(run: CiRun): string {
  if (run.status === "cancelled") return "cancelled";
  if (run.status === "failure") return run.stage === "publish" ? "failed publishing" : `blocked at ${run.stage}`;
  return run.type === "sdk" ? "gate passed" : "shipped";
}

const ms = (from: string | null, to: string | null) => (from && to ? Date.parse(to) - Date.parse(from) : null);

/**
 * Every change released recently, newest first. A change ships one release run
 * per component it touched; the bars put those runs on the change's own clock,
 * so a batch that ran one after another reads as a staircase.
 */
export function Releases(props: { releases: readonly Release[] }) {
  return (
    <Show
      when={props.releases.length > 0}
      fallback={<EmptyState variant="inline" title="No release runs reported yet." />}
    >
      <p class="mb-6 max-w-[68ch] text-[13px] text-muted-foreground">
        Each change and the release runs it started, one per component. Bars share the change's clock, from the
        first run starting to the last one finishing.
      </p>
      <ol class="border-t border-border">
        <For each={props.releases}>{(release) => <ReleaseRow release={release} />}</For>
      </ol>
    </Show>
  );
}

function ReleaseRow(props: { release: Release }) {
  const outcome = () => OUTCOME[releaseOutcome(props.release)];
  const href = () =>
    props.release.pr ? pullUrl(props.release.pr) : props.release.sha ? commitUrl(props.release.sha) : null;

  return (
    <li class="border-b border-border py-4">
      <div class="flex items-baseline gap-3">
        <StatusDot tone={outcome().tone} label={outcome().label} size="md" class="relative top-[-1px] shrink-0" />
        <h3 class="min-w-0 flex-1 text-[14px] font-medium">
          <Show when={href()} fallback={props.release.title}>
            {(url) => (
              <a class="hover:underline" href={url()} {...external}>
                {props.release.title}
              </a>
            )}
          </Show>
          <Show when={props.release.pr}>
            {(pr) => <span class="ml-2 font-normal text-muted-foreground">#{pr()}</span>}
          </Show>
        </h3>
        <span class="shrink-0 text-[12px] tabular-nums text-muted-foreground">
          <Show when={props.release.startedAt}>{(at) => timeAgo(at())}</Show>
          <Show when={props.release.runs.length > 1 && ms(props.release.startedAt, props.release.completedAt)}>
            {(span) => <>, {formatDuration(span())} in all</>}
          </Show>
        </span>
      </div>

      <ul class="mt-3 space-y-1.5 pl-[22px]">
        <For each={props.release.runs}>
          {(run) => {
            const span = runSpan(props.release, run);
            const label = `${run.component} ${versionOf(run.tag)}: ${runState(run)}`;
            return (
              <li>
                <a
                  href={runUrl(run.runId)}
                  {...external}
                  aria-label={label}
                  class="group grid grid-cols-[minmax(0,9.5rem)_minmax(0,1fr)_auto] items-center gap-3 text-[12.5px] sm:grid-cols-[minmax(0,11rem)_minmax(0,1fr)_9rem]"
                >
                  <span class="truncate group-hover:underline">
                    {run.component}
                    <span class="ml-1.5 font-mono text-[11.5px] text-muted-foreground">{versionOf(run.tag)}</span>
                  </span>
                  <span class="relative h-2 rounded-full bg-muted" aria-hidden="true">
                    <Show when={span}>
                      {(s) => (
                        <span
                          class={cn("absolute inset-y-0 rounded-full", BAR[run.status])}
                          style={{
                            left: `${s().start * 100}%`,
                            width: `max(6px, ${s().width * 100}%)`,
                          }}
                        />
                      )}
                    </Show>
                  </span>
                  <span class="flex items-center justify-end gap-2 whitespace-nowrap text-[12px] tabular-nums text-muted-foreground">
                    <Show when={run.lane === "express"}>
                      <Badge shape="pill" variant="muted">
                        hotfix
                      </Badge>
                    </Show>
                    <Show when={run.status !== "success"} fallback={formatDuration(ms(run.startedAt, run.completedAt))}>
                      <span class={run.status === "failure" ? "text-destructive dark:text-destructive-foreground" : undefined}>
                        {runState(run)}
                      </span>
                    </Show>
                  </span>
                </a>
              </li>
            );
          }}
        </For>
      </ul>
    </li>
  );
}
