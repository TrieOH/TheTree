import { createFileRoute } from "@tanstack/solid-router";
import { Tabs, cn } from "@trieoh/ui-solid";
import { For, Match, Show, Switch, createMemo, createSignal } from "solid-js";
import { SERVICES, type ServiceId } from "@/features/catalog/services";
import { loadIssues } from "@/features/issues/api";
import { isOpen, type Project } from "@/features/issues/model";
import {
  cycleProgress,
  doneThisCycle,
  formatDay,
  groupByState,
  inCycle,
  medianDaysToClose,
} from "@/features/issues/stats";
import { IssueList, StateMark } from "@/features/issues/ui/IssueList";
import { Kpi, KpiRow } from "@/widgets/shell/Kpi";
import { PageHeader, Section } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/issues")({
  loader: () => loadIssues(),
  head: () => ({ meta: [{ title: "Issues · TrieOH dev" }] }),
  component: IssuesPage,
});

const TABS = [
  { id: "current", label: "Current cycle" },
  { id: "projects", label: "Projects" },
  { id: "bugs", label: "Bugs" },
  { id: "done", label: "Done" },
] as const;

type Tab = (typeof TABS)[number]["id"];
/** A service, "other" for issues outside the four, or everything. */
type Filter = ServiceId | "other" | "all";

const LINEAR_URL = "https://linear.app/trieoh";

function IssuesPage() {
  const data = Route.useLoaderData();
  const [tab, setTab] = createSignal<Tab>("current", { name: "issues-tab" });
  const [filter, setFilter] = createSignal<Filter>("all", { name: "issues-filter" });

  const issues = createMemo(() =>
    data().issues.filter((issue) =>
      filter() === "all" ? true : filter() === "other" ? issue.service === null : issue.service === filter(),
    ),
  );
  const open = createMemo(() => issues().filter(isOpen));
  const done = createMemo(() => issues().filter((issue) => !isOpen(issue)));
  const progress = () => cycleProgress(issues(), data().cycle);
  const median = () => medianDaysToClose(done());

  const currentList = () => {
    const cycle = data().cycle;
    return cycle ? issues().filter((issue) => inCycle(issue, cycle)) : open();
  };

  return (
    <>
      <PageHeader
        title="Issues"
        description={
          <>
            Synced from Linear as they change. Open any issue in{" "}
            <a class="underline underline-offset-4" href={LINEAR_URL} target="_blank" rel="noopener">
              Linear
            </a>{" "}
            to edit it.
          </>
        }
      />

      <KpiRow>
        <Kpi label="Open">{open().length}</Kpi>
        <Kpi label="Urgent">{open().filter((issue) => issue.priority === 1).length}</Kpi>
        <Kpi label="Done this cycle" note={data().cycle ? undefined : "last 14 days"}>
          {doneThisCycle(done(), data().cycle).length}
        </Kpi>
        <Kpi label="Median time to close" note="last 30 days">
          {median() === null ? "—" : `${median()!.toFixed(1)} d`}
        </Kpi>
      </KpiRow>

      <Show when={data().cycle}>
        {(cycle) => (
          <div class="mt-6">
            <div class="mb-1.5 flex justify-between text-[13px] text-muted-foreground">
              <span>
                Cycle {cycle().number} · {formatDay(cycle().startsAt)} – {formatDay(cycle().endsAt)}
              </span>
              <span class="tabular-nums">
                {progress().done} of {progress().total} done
              </span>
            </div>
            <div class="h-1.5 overflow-hidden rounded-full bg-muted">
              <div
                class="h-full rounded-full bg-accent"
                style={{ width: `${progress().total ? (100 * progress().done) / progress().total : 0}%` }}
              />
            </div>
          </div>
        )}
      </Show>

      <div class="mt-8">
        <Tabs ariaLabel="Issue views" items={TABS} value={tab()} onChange={setTab} />
      </div>

      <div role="group" aria-label="Filter by service" class="mt-4 flex flex-wrap gap-1.5">
        <FilterChip value="all" current={filter()} onSelect={setFilter}>All</FilterChip>
        <For each={SERVICES}>
          {(service) => (
            <FilterChip value={service.id} current={filter()} onSelect={setFilter}>
              {service.name}
            </FilterChip>
          )}
        </For>
        <FilterChip value="other" current={filter()} onSelect={setFilter}>Other</FilterChip>
      </div>

      <div class="mt-6">
        <Switch>
          <Match when={tab() === "current"}>
            <Show
              when={currentList().length > 0}
              fallback={<p class="text-[13px] text-muted-foreground">Nothing in this cycle.</p>}
            >
              <For each={groupByState(currentList())}>
                {(group) => (
                  <Section
                    title={group.state}
                    aside={<span class="tabular-nums">{group.issues.length}</span>}
                  >
                    <IssueList issues={group.issues} />
                  </Section>
                )}
              </For>
            </Show>
          </Match>
          <Match when={tab() === "projects"}>
            <ProjectList projects={data().projects} />
          </Match>
          <Match when={tab() === "bugs"}>
            <IssueList
              issues={open().filter((issue) => issue.labels.includes("Bug"))}
              empty="No open bugs."
            />
          </Match>
          <Match when={tab() === "done"}>
            <IssueList issues={done()} empty="Nothing completed in the last 30 days." />
          </Match>
        </Switch>
      </div>

      <Show when={data().syncedAt}>
        {(at) => (
          <p class="mt-10 text-[12px] text-muted-foreground">
            Last full sync {new Date(at()).toLocaleString("en", { dateStyle: "medium", timeStyle: "short" })}
          </p>
        )}
      </Show>
    </>
  );
}

function FilterChip(props: {
  value: Filter;
  current: Filter;
  onSelect: (value: Filter) => void;
  children: string;
}) {
  const active = () => props.value === props.current;
  return (
    <button
      type="button"
      aria-pressed={active() ? "true" : "false"}
      onClick={() => props.onSelect(props.value)}
      class={cn(
        "h-7 rounded-full border px-3 text-[12.5px] transition-colors",
        active()
          ? "border-foreground/30 bg-foreground/10 text-foreground"
          : "border-border text-muted-foreground hover:bg-foreground/5",
      )}
    >
      {props.children}
    </button>
  );
}

function ProjectList(props: { projects: readonly Project[] }) {
  return (
    <Show
      when={props.projects.length > 0}
      fallback={<p class="text-[13px] text-muted-foreground">No active projects.</p>}
    >
      <ul class="divide-y divide-border border-y border-border">
        <For each={props.projects}>
          {(project) => (
            <li class="flex items-center gap-4 py-3">
              <StateMark type={project.state === "started" ? "started" : "unstarted"} />
              <span class="min-w-0 flex-1">
                <span class="block truncate font-medium">{project.name}</span>
                <span class="text-[12px] text-muted-foreground">
                  {project.done} of {project.total} issues
                  {project.targetDate ? ` · target ${formatDay(project.targetDate)}` : ""}
                </span>
              </span>
              <span class="h-1.5 w-24 shrink-0 overflow-hidden rounded-full bg-muted">
                <span
                  class="block h-full rounded-full bg-accent"
                  style={{ width: `${project.total ? (100 * project.done) / project.total : 0}%` }}
                />
              </span>
            </li>
          )}
        </For>
      </ul>
    </Show>
  );
}
