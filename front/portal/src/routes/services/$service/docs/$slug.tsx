import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { PageHeader } from "@trieoh/ui-solid";
import { Dynamic } from "@solidjs/web";
import { Show } from "solid-js";
import { findService } from "@/features/catalog/services";
import { docComponents } from "@/features/docs/components";
import { DocPager } from "@/features/docs/DocPager";
import { AUDIENCES, findServiceDoc, serviceNeighbours, type ServiceDocLink } from "@/features/docs/service-docs";

export const Route = createFileRoute("/services/$service/docs/$slug")({
  loader: async ({ params }) => {
    const service = findService(params.service);
    const doc = service && findServiceDoc(service.id, params.slug);
    if (!service || !doc) throw notFound();
    const { default: Content } = await doc.load();
    return {
      service,
      meta: doc.meta,
      audience: AUDIENCES.find((entry) => entry.id === doc.audience)!,
      Content,
      ...serviceNeighbours(doc),
    };
  },
  head: ({ loaderData }) => ({
    meta: [{ title: `${loaderData?.meta.title ?? "Doc"} · ${loaderData?.service.name ?? ""} · TrieOH dev` }],
  }),
  component: ServiceDocPage,
});

const hrefOf = (doc: ServiceDocLink) => `/services/${doc.service}/docs/${doc.slug}`;

function ServiceDocPage() {
  const data = Route.useLoaderData();

  return (
    <>
      <PageHeader title={data().meta.title}>
        <nav class="flex flex-wrap items-center gap-2 text-[13px] text-muted-foreground">
          <Link to="/services/$service/docs" params={{ service: data().service.id }} class="hover:text-foreground">
            {data().service.name} docs
          </Link>
          <span>/</span>
          <span>{data().audience.label}</span>
          <Show when={data().meta.group}>
            <span>/</span>
            <span>{data().meta.group}</span>
          </Show>
        </nav>
      </PageHeader>
      <article class="prose-doc">
        {/* Dynamic, so moving to another doc swaps the body, not just the header. */}
        <Dynamic component={data().Content} components={docComponents} />
      </article>
      <DocPager previous={data().previous} next={data().next} href={hrefOf} />
    </>
  );
}
