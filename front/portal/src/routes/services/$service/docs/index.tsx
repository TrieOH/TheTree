import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { PageHeader, buttonVariants } from "@trieoh/ui-solid";
import { findService } from "@/features/catalog/services";
import { ServiceDocIndex, ServiceRelated } from "@/features/docs/ui/ServiceDocIndex";

export const Route = createFileRoute("/services/$service/docs/")({
  loader: ({ params }) => {
    const service = findService(params.service);
    if (!service) throw notFound();
    return service;
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.name ?? ""} docs · TrieOH dev` }] }),
  component: ServiceDocsPage,
});

const action = buttonVariants({ variant: "outline", size: "sm" });

function ServiceDocsPage() {
  const service = Route.useLoaderData();

  return (
    <>
      <PageHeader title={`${service().name} docs`} description={service().description}>
        <div class="flex flex-wrap gap-2 pt-3">
          <Link to="/services/$service" params={{ service: service().id }} class={action}>
            Service
          </Link>
          <Link to="/services/$service/api" params={{ service: service().id }} class={action}>
            API reference
          </Link>
          <Link to="/problems/$set" params={{ set: service().id }} class={action}>
            Problems
          </Link>
        </div>
      </PageHeader>
      <ServiceDocIndex service={service().id} />
      <ServiceRelated service={service().id} />
    </>
  );
}
