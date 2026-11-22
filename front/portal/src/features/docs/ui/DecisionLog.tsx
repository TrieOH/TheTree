import { Link } from "@tanstack/solid-router";
import { SearchInput, ToggleGroup, cn } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal } from "solid-js";
import { AREAS, type Area, areasOf, byMonth, formatDate, isSuperseded, matches } from "../adr";
import { statusWord, type DocLink } from "../catalog";
import { DocStatus } from "../DocStatus";
import { AdrRelationLines } from "./AdrRelations";

type Filter = Area | "all";

/** The decisions, newest first, grouped by the month they were made. */
export function DecisionLog(props: { docs: readonly DocLink[] }) {
  const [query, setQuery] = createSignal("", { name: "decisions-query" });
  const [area, setArea] = createSignal<Filter>("all", { name: "decisions-area" });

  // Only offer the areas some decision actually touches.
  const areas = createMemo(
    () => AREAS.filter((a) => props.docs.some((doc) => areasOf(doc.meta).has(a.id))),
    { name: "decisions-areas" },
  );
  const shown = createMemo(
    () =>
      props.docs.filter(
        (doc) => matches(doc, query()) && (area() === "all" || areasOf(doc.meta).has(area() as Area)),
      ),
    { name: "decisions-shown" },
  );
  const months = createMemo(() => byMonth(shown()), { name: "decisions-months" });

  return (
    <>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
        <SearchInput
          value={query()}
          onInput={(event) => setQuery(event.currentTarget.value)}
          placeholder="Search decisions"
          aria-label="Search decisions"
          class="sm:max-w-[220px]"
        />
        <ToggleGroup
          ariaLabel="Area"
          size="sm"
          items={[{ value: "all" as Filter, label: "All" }, ...areas().map((a) => ({ value: a.id as Filter, label: a.label }))]}
          value={area()}
          onChange={setArea}
          class="flex-wrap"
        />
      </div>

      <Show
        when={months().length > 0}
        fallback={<p class="mt-10 text-[14px] text-muted-foreground">No decision matches. Clear the search or pick another area.</p>}
      >
        <div class="mt-8">
          <For each={months()}>
            {(group) => (
              <section class="grid gap-x-8 border-t border-border sm:grid-cols-[120px_minmax(0,1fr)]">
                <h2 class="pt-4 text-[13px] font-medium text-muted-foreground sm:sticky sm:top-4 sm:self-start sm:pb-4">
                  {group.month}
                </h2>
                <ol class="divide-y divide-border sm:pt-0">
                  <For each={group.docs}>{(doc) => <DecisionRow doc={doc} />}</For>
                </ol>
              </section>
            )}
          </For>
        </div>
      </Show>
    </>
  );
}

function DecisionRow(props: { doc: DocLink }) {
  const meta = () => props.doc.meta;
  const old = () => isSuperseded(meta());
  const notAccepted = () => statusWord(meta().status) !== "accepted";

  return (
    <li class={cn("py-4", old() && "opacity-60 hover:opacity-100 focus-within:opacity-100")}>
      <div class="flex items-baseline gap-3">
        <span class="w-[72px] shrink-0 font-mono text-[12.5px] text-muted-foreground">{meta().number}</span>
        <Link
          to="/docs/$category/$slug"
          params={{ category: "decisions", slug: props.doc.slug }}
          class="min-w-0 flex-1 rounded-sm font-medium leading-snug hover:underline hover:underline-offset-[3px] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
        >
          {meta().title ?? props.doc.slug}
        </Link>
        <time datetime={meta().date} class="shrink-0 text-[13px] tabular-nums text-muted-foreground">
          {formatDate(meta().date, false)}
        </time>
      </div>
      <div class="mt-1.5 sm:pl-[84px]">
        <Show when={meta().summary}>
          <p class="max-w-[62ch] text-[14px] leading-relaxed text-muted-foreground">{meta().summary}</p>
        </Show>
        <div class="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-[13px] text-muted-foreground empty:hidden">
          <Show when={notAccepted()}>
            <DocStatus status={meta().status} />
          </Show>
          <AdrRelationLines meta={meta()} />
          <Show when={meta().updated}>
            <span>Updated {formatDate(meta().updated, false)}</span>
          </Show>
        </div>
      </div>
    </li>
  );
}
