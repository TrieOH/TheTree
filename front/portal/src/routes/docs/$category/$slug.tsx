import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { Show } from "solid-js";
import { findCategory, findDoc } from "@/features/docs/catalog";
import { DocStatus } from "@/features/docs/DocStatus";
import { PageHeader } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/docs/$category/$slug")({
  loader: async ({ params }) => {
    const doc = findDoc(params.category, params.slug);
    if (!doc) throw notFound();
    const { default: Content } = await doc.load();
    return { meta: doc.meta, category: findCategory(params.category)!, Content };
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.meta.title ?? "Doc"} · TrieOH dev` }] }),
  component: DocPage,
});

function DocPage() {
  const data = Route.useLoaderData();
  const Content = data().Content;

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
          <DocStatus status={data().meta.status} />
          <Show when={data().meta.date}>
            <span class="tabular-nums">{data().meta.date}</span>
          </Show>
        </div>
      </PageHeader>
      <article class="prose-doc">
        <Content />
      </article>
    </>
  );
}
