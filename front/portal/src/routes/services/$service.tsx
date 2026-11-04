import { createFileRoute, notFound } from "@tanstack/solid-router";
import { For } from "solid-js";
import { REPO_URL, findService } from "@/features/catalog/services";
import { NotConnected, PageHeader, Section } from "@/widgets/shell/PageHeader";
import { StatusDot, statusLabel } from "@/widgets/shell/StatusDot";

export const Route = createFileRoute("/services/$service")({
  loader: ({ params }) => {
    const service = findService(params.service);
    if (!service) throw notFound();
    return service;
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.name ?? "Service"} · TrieOH dev` }] }),
  component: ServicePage,
});

const button =
  "inline-flex h-8 items-center rounded-lg border border-border px-3 text-[13px] transition-colors hover:bg-foreground/5";

function ServicePage() {
  const service = Route.useLoaderData();

  const links = () => [
    { label: "API reference", href: `${REPO_URL}/blob/main/${service().dir}/api-spec.yml` },
    { label: "Code", href: `${REPO_URL}/tree/main/${service().dir}` },
    { label: "Front", href: service().front },
  ];

  return (
    <>
      <PageHeader title={service().name} description={service().description}>
        <p class="font-mono text-[12px] text-muted-foreground">{service().apiUrl.replace("https://", "")}</p>
        <div class="flex flex-wrap gap-2 pt-3">
          <For each={links()}>
            {(link) => (
              <a class={button} href={link.href} target="_blank" rel="noopener">
                {link.label}
              </a>
            )}
          </For>
        </div>
      </PageHeader>

      <Section title="Health">
        <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
          <Kpi label="Status">
            <span class="inline-flex items-center gap-2">
              <StatusDot status="unknown" />
              {statusLabel("unknown")}
            </span>
          </Kpi>
          <Kpi label="Success rate 30d">—</Kpi>
          <Kpi label="p95">—</Kpi>
          <Kpi label="Open issues">—</Kpi>
        </div>
      </Section>

      <Section title="Issues">
        <NotConnected source="Linear" />
      </Section>
    </>
  );
}

function Kpi(props: { label: string; children: unknown }) {
  return (
    <div class="rounded-xl border border-border bg-card p-3.5">
      <span class="block text-[12px] text-muted-foreground">{props.label}</span>
      <b class="mt-1 block text-lg font-semibold tabular-nums">{props.children as never}</b>
    </div>
  );
}
