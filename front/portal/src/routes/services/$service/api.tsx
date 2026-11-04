import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { findService } from "@/features/catalog/services";
import { ApiReference } from "@/features/api-reference/ApiReference";

export const Route = createFileRoute("/services/$service/api")({
  loader: ({ params }) => {
    const service = findService(params.service);
    if (!service) throw notFound();
    return service;
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.name ?? ""} API · TrieOH dev` }] }),
  component: ApiPage,
});

function ApiPage() {
  const service = Route.useLoaderData();
  return (
    <>
      <nav class="mb-2 flex items-center gap-2 px-1 text-[13px] text-muted-foreground">
        <Link to="/services/$service" params={{ service: service().id }} class="hover:text-foreground">
          {service().name}
        </Link>
        <span>/</span>
        <span class="text-foreground">API reference</span>
      </nav>
      <ApiReference service={service().id} />
    </>
  );
}
