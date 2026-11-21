import { Accordion, EmptyState, PageSection, Sparkline, StatusDot, Table, ToggleGroup, cn } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal } from "solid-js";
import type { CiMetricStat, CiPayload, CiRun } from "../model";
import { isSlowLastRun, jobTimes, pipelineNumbers, type JobTimes, type StepRow } from "../pipeline";
import { formatDuration, formatMetric, formatStepTime } from "../stats";

const PIPELINES: Record<string, string> = { service: "APIs", front: "Fronts", sdk: "SDKs", pages: "Pages" };

const JOB: Record<string, string> = { parse: "Parse the tag", checks: "Checks", publish: "Publish" };

const note = "max-w-[68ch] text-[13px] text-muted-foreground";

function median(values: number[]): number | null {
  if (values.length === 0) return null;
  const sorted = [...values].sort((a, b) => a - b);
  const mid = Math.floor(sorted.length / 2);
  return sorted.length % 2 ? sorted[mid] : (sorted[mid - 1] + sorted[mid]) / 2;
}

const runMs = (run: CiRun) =>
  run.startedAt && run.completedAt ? Date.parse(run.completedAt) - Date.parse(run.startedAt) : null;

/** One release pipeline at a time: where its time goes, and what its checks measured. */
export function Pipelines(props: { ci: CiPayload }) {
  const types = createMemo(() => Object.keys(PIPELINES).filter((t) => props.ci.runs.some((r) => r.type === t)), {
    name: "delivery-pipeline-types",
  });
  const [picked, setPicked] = createSignal("service", { name: "delivery-pipeline" });
  const current = () => (types().includes(picked()) ? picked() : (types()[0] ?? "service"));
  const jobs = createMemo(() => jobTimes(props.ci.steps.filter((s) => s.type === current())), {
    name: "delivery-pipeline-jobs",
  });
  const numbers = createMemo(() => pipelineNumbers(props.ci.metrics.filter((m) => m.type === current())), {
    name: "delivery-pipeline-numbers",
  });
  const typical = createMemo(
    () =>
      median(
        props.ci.runs
          .filter((r) => r.type === current() && r.status === "success")
          .flatMap((r) => (runMs(r) === null ? [] : [runMs(r)!])),
      ),
    { name: "delivery-pipeline-typical" },
  );

  return (
    <Show
      when={types().length > 0}
      fallback={<EmptyState variant="inline" title="No release runs have reported step times yet." />}
    >
      <div class="flex flex-wrap items-center justify-between gap-3">
        <ToggleGroup
          ariaLabel="Pipeline"
          items={types().map((t) => ({ value: t, label: PIPELINES[t] }))}
          value={current()}
          onChange={setPicked}
        />
        <Show when={typical()}>
          {(ms) => <p class="text-[13px] text-muted-foreground">A release usually takes {formatDuration(ms())}.</p>}
        </Show>
      </div>

      <PageSection title="Where the time goes" aside="median of recent runs" class="mt-8">
        <div class="space-y-5">
          <For each={jobs()}>{(job) => <JobBar job={job} />}</For>
        </div>
      </PageSection>

      <Numbers numbers={numbers()} />

      <Accordion
        class="mt-10"
        multiple
        items={[
          {
            value: "steps",
            title: "Every step, with its recent runs",
            content: (
              <div class="space-y-6 pt-2">
                <For each={jobs()}>
                  {(job) => (
                    <div>
                      <h3 class="mb-1 text-[13px] font-medium">{JOB[job.job] ?? job.job}</h3>
                      <StepTable steps={job.steps} job={job.job} />
                      <Show when={job.hidden.count > 0}>
                        <p class="mt-2 text-[12px] text-muted-foreground">
                          {job.steps.length ? "Plus " : ""}
                          {job.hidden.count} quick {job.hidden.count === 1 ? "step" : "steps"} (setup, cleanup,
                          checkout), {formatStepTime(job.hidden.medianMs)} together.
                        </p>
                      </Show>
                    </div>
                  )}
                </For>
              </div>
            ),
          },
          ...(numbers().other.length
            ? [{ value: "numbers", title: "Every other recorded number", content: <OtherNumbers items={numbers().other} /> }]
            : []),
        ]}
      />
    </Show>
  );
}

function StepTable(props: { steps: readonly StepRow[]; job: string }) {
  return (
    <Show when={props.steps.length > 0}>
      <Table
        caption={`${props.job} steps`}
        class="text-[13.5px]"
        density="compact"
        items={props.steps}
        columns={[
          {
            id: "step",
            header: "Step",
            title: (s) => s.step,
            cell: (s) => s.label,
          },
          {
            id: "median",
            header: "Usually",
            align: "right",
            class: "whitespace-nowrap tabular-nums",
            cell: (s) => formatStepTime(s.medianMs),
          },
          {
            id: "p90",
            header: "Slow runs",
            align: "right",
            class: "whitespace-nowrap tabular-nums text-muted-foreground",
            title: () => "p90: nine runs in ten are this fast or faster",
            cell: (s) => formatStepTime(s.p90Ms),
          },
          {
            id: "last",
            header: "Last",
            align: "right",
            class: (s) => cn("whitespace-nowrap tabular-nums", isSlowLastRun(s) ? "font-medium" : "text-muted-foreground"),
            cell: (s) => (
              <>
                <Show when={isSlowLastRun(s)}>
                  <StatusDot tone="warning" label="Much slower than usual" class="mr-1.5" />
                </Show>
                {formatStepTime(s.lastMs)}
              </>
            ),
          },
          {
            id: "recent",
            header: "Recent runs",
            class: "w-[28%] min-w-32",
            cell: (s) => (
              <Sparkline
                type="bar"
                fill
                data={s.recentMs}
                label={`${s.label}, last ${s.recentMs.length} runs`}
                valueFormatter={formatStepTime}
                height={18}
                color="color-mix(in oklab, var(--muted-foreground) 55%, transparent)"
              />
            ),
          },
        ]}
      />
    </Show>
  );
}

function Numbers(props: { numbers: ReturnType<typeof pipelineNumbers> }) {
  const n = () => props.numbers;
  return (
    <>
      <Show when={n().findings || n().tests.length > 0}>
        <PageSection title="Checks" aside="from the newest run">
          <dl class="grid gap-x-8 gap-y-3 text-[13.5px] sm:grid-cols-[10rem_minmax(0,1fr)]">
            <Show when={n().findings}>
              {(findings) => (
                <>
                  <dt class="text-muted-foreground">Vulnerability scan</dt>
                  <dd class="flex flex-wrap gap-x-4 gap-y-1 tabular-nums">
                    <For each={findings()}>
                      {(f) => (
                        <span
                          class={cn(
                            f.count > 0 && (f.severity === "CRITICAL" || f.severity === "HIGH")
                              ? "font-medium text-destructive dark:text-destructive-foreground"
                              : f.count === 0 && "text-muted-foreground",
                          )}
                        >
                          {f.count} {f.severity.toLowerCase()}
                        </span>
                      )}
                    </For>
                    <span class="text-muted-foreground">
                      across {n().scannedPaths} {n().scannedPaths === 1 ? "path" : "paths"}
                    </span>
                  </dd>
                </>
              )}
            </Show>
            <For each={n().tests}>
              {(t) => (
                <>
                  <dt class="text-muted-foreground">{t.mode ? `Tests, ${t.mode}` : "Tests"}</dt>
                  <dd class="tabular-nums">
                    {t.pass} passed
                    <span class={t.fail > 0 ? "font-medium text-destructive dark:text-destructive-foreground" : "text-muted-foreground"}>
                      , {t.fail} failed
                    </span>
                    <Show when={t.skip > 0}>
                      <span class="text-muted-foreground">, {t.skip} skipped</span>
                    </Show>
                  </dd>
                </>
              )}
            </For>
          </dl>
        </PageSection>
      </Show>

      <Show when={n().slowTests.length > 0}>
        <PageSection title="Slowest tests" aside="a second or more, usually">
          <Table
            caption="Slowest tests"
            class="text-[13.5px]"
            density="compact"
            items={n().slowTests}
            columns={[
              {
                id: "test",
                header: "Test",
                class: "font-mono text-[12px]",
                cell: (t) => (
                  <>
                    <span class="text-muted-foreground">{t.pkg}.</span>
                    {t.test}
                  </>
                ),
              },
              { id: "mode", header: "Mode", class: "text-muted-foreground", cell: (t) => t.mode },
              { id: "median", header: "Usually", align: "right", class: "tabular-nums", cell: (t) => `${t.medianS.toFixed(1)} s` },
              {
                id: "last",
                header: "Last",
                align: "right",
                class: "tabular-nums text-muted-foreground",
                cell: (t) => `${t.lastS.toFixed(1)} s`,
              },
            ]}
          />
        </PageSection>
      </Show>

    </>
  );
}

function OtherNumbers(props: { items: readonly CiMetricStat[] }) {
  return (
    <div class="pt-2">
      <p class={cn(note, "mb-3")}>Recorded by .github/scripts/metric. Numbers that were zero in every run are left out.</p>
          <Table
            caption="Other numbers"
            class="text-[13.5px]"
            density="compact"
            items={props.items}
            columns={[
              {
                id: "name",
                header: "Number",
                cell: (m) => (
                  <>
                    <span class="text-muted-foreground">{m.source} · </span>
                    {m.name.replaceAll("_", " ")}
                    <span class="text-muted-foreground"> in {m.job}</span>
                    <For each={Object.entries(m.labels)}>
                      {([key, value]) => (
                        <span class="ml-2 rounded bg-muted px-1.5 py-0.5 font-mono text-[11px] text-muted-foreground">
                          {key}={value}
                        </span>
                      )}
                    </For>
                  </>
                ),
              },
              { id: "last", header: "Last", align: "right", class: "tabular-nums", cell: (m) => formatMetric(m.name, m.last) },
              {
                id: "median",
                header: "Usually",
                align: "right",
                class: "tabular-nums text-muted-foreground",
                cell: (m) => formatMetric(m.name, m.median),
              },
            ]}
          />
    </div>
  );
}

/**
 * One job as a single bar split by its steps, longest first, so the step that
 * eats the job stands out without reading a table.
 */
function JobBar(props: { job: JobTimes }) {
  const parts = createMemo(() => {
    const top = [...props.job.steps].filter((s) => s.medianMs >= 3_000).sort((a, b) => b.medianMs - a.medianMs).slice(0, 5);
    const rest = props.job.medianMs - top.reduce((sum, s) => sum + s.medianMs, 0);
    return { top, rest };
  });
  const share = (ms: number) => `${(ms / Math.max(1, props.job.medianMs)) * 100}%`;
  return (
    <div>
      <div class="mb-1.5 flex items-baseline justify-between gap-3 text-[13.5px]">
        <span class="font-medium">{JOB[props.job.job] ?? props.job.job}</span>
        <span class="tabular-nums text-muted-foreground">{formatStepTime(props.job.medianMs)}</span>
      </div>
      <div class="flex h-3 gap-px overflow-hidden rounded-full bg-muted" aria-hidden="true">
        <For each={parts().top}>
          {(step, i) => (
            <span
              title={`${step.label}: ${formatStepTime(step.medianMs)}`}
              class={i() === 0 ? "bg-primary" : "bg-muted-foreground/45"}
              style={{ width: share(step.medianMs) }}
            />
          )}
        </For>
      </div>
      <p class="mt-1.5 text-[12.5px] text-muted-foreground">
        <For each={parts().top}>
          {(step, i) => (
            <>
              <span class={i() === 0 ? "text-foreground" : undefined}>
                {step.label} {formatStepTime(step.medianMs)}
              </span>
              {", "}
            </>
          )}
        </For>
        everything else {formatStepTime(Math.max(0, parts().rest))}.
      </p>
    </div>
  );
}
