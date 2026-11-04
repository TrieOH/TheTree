import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { For, Show } from "solid-js";
import { docsIn, findCategory } from "@/features/docs/catalog";
import { DocStatus } from "@/features/docs/DocStatus";
import { NotConnected, PageHeader } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/docs/$category/")({
  loader: ({ params }) => {
    const category = findCategory(params.category);
    if (!category) throw notFound();
    return category;
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.label ?? "Docs"} · TrieOH dev` }] }),
  component: CategoryPage,
});

function CategoryPage() {
  const category = Route.useLoaderData();
  const docs = () => docsIn(category().id);

  return (
    <>
      <PageHeader title={category().label} description={category().description} />
      <Show
        when={docs().length > 0}
        fallback={<NotConnected source="docs/">Nothing here yet. These pages arrive as the docs are converted.</NotConnected>}
      >
        <ul class="divide-y divide-border border-y border-border">
          <For each={docs()}>
            {(doc) => (
              <li>
                <Link
                  to="/docs/$category/$slug"
                  params={{ category: category().id, slug: doc.slug }}
                  class="flex items-baseline justify-between gap-4 py-3"
                >
                  <span class="min-w-0">
                    <span class="block font-medium">{doc.meta.title ?? doc.slug}</span>
                    <span class="mt-0.5 flex items-center gap-2 text-[13px] text-muted-foreground">
                      <span class="font-mono text-[12px]">{doc.meta.number}</span>
                      <DocStatus status={doc.meta.status} />
                    </span>
                  </span>
                  <span class="shrink-0 text-[13px] tabular-nums text-muted-foreground">{doc.meta.date}</span>
                </Link>
              </li>
            )}
          </For>
        </ul>
      </Show>
    </>
  );
}
