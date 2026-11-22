import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { Badge, PageHeader } from "@trieoh/ui-solid";
import { Dynamic } from "@solidjs/web";
import { Show } from "solid-js";
import { findCategory, findDoc, neighboursOf } from "@/features/docs/catalog";
import { docComponents } from "@/features/docs/components";
import { DocPager } from "@/features/docs/DocPager";
import { DocStatus } from "@/features/docs/DocStatus";
import { formatDate } from "@/features/docs/adr";
import { AdrFacts, SupersededNotice } from "@/features/docs/ui/AdrFacts";
import { IncidentFacts } from "@/features/docs/ui/IncidentFacts";
import { SEVERITY_BADGE, severityLabel, severityOf } from "@/features/docs/incident";

export const Route = createFileRoute("/docs/$category/$slug")({
  loader: async ({ params }) => {
    const doc = findDoc(params.category, params.slug);
    if (!doc) throw notFound();
    const { default: Content } = await doc.load();
    return { meta: doc.meta, category: findCategory(params.category)!, Content, ...neighboursOf(doc) };
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.meta.title ?? "Doc"} · TrieOH dev` }] }),
  component: DocPage,
});

function DocPage() {
  const data = Route.useLoaderData();

  return (
    <>
      <PageHeader title={data().meta.title}>
        <div class="flex flex-wrap items-center gap-3 text-[13px] text-muted-foreground">
          <Link to="/docs/$category" params={{ category: data().category.id }}>
            {data().category.label}
          </Link>
          <Show when={data().meta.number}>
            <span class="font-mono text-[12px]">{data().meta.number}</span>
          </Show>
          <Show when={severityOf(data().meta)}>
            {(level) => (
              <Badge variant={SEVERITY_BADGE[level()]}>
                {severityLabel(level())}
              </Badge>
            )}
          </Show>
          <DocStatus status={data().meta.status} />
          <Show when={data().meta.date}>
            <span>
              {data().category.id === "decisions" ? "Decided " : ""}
              <time datetime={data().meta.date}>{formatDate(data().meta.date)}</time>
            </span>
          </Show>
          <Show when={data().meta.updated}>
            <span>
              Updated <time datetime={data().meta.updated}>{formatDate(data().meta.updated)}</time>
            </span>
          </Show>
        </div>
      </PageHeader>
      <Show when={data().category.id === "decisions"}>
        <div class="max-w-[68ch]">
          <SupersededNotice meta={data().meta} />
          <AdrFacts meta={data().meta} />
        </div>
      </Show>
      <Show when={data().category.id === "incident-reports"}>
        <IncidentFacts meta={data().meta} />
      </Show>
      <article class="prose-doc">
        {/* Dynamic, so moving to another doc swaps the body, not just the header. */}
        <Dynamic component={data().Content} components={docComponents} />
      </article>
      <DocPager previous={data().previous} next={data().next} />
    </>
  );
}
