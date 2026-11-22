import { Link, createFileRoute } from "@tanstack/solid-router";
import { PageSection } from "@trieoh/ui-solid";
import { For, Show, createMemo } from "solid-js";
import { SERVICES, findService } from "@/features/catalog/services";
import { docsIn } from "@/features/docs/catalog";
import { loadDelivery } from "@/features/delivery/api";
import type { DeliveryPayload } from "@/features/delivery/model";
import { pullUrl } from "@/features/delivery/model";
import { groupReleases, releaseOutcome } from "@/features/delivery/releases";
import { timeAgo } from "@/features/delivery/stats";
import { loadIssues } from "@/features/issues/api";
import { isOpen, type IssuesPayload } from "@/features/issues/model";
import { loadMetrics } from "@/features/metrics/api";
import type { MetricsPayload } from "@/features/metrics/model";
import { StatusDot } from "@/widgets/shell/StatusDot";

export const Route = createFileRoute("/")({
  // Home renders whatever loads; a source that fails just leaves its part out.
  loader: async (): Promise<{
    issues: IssuesPayload | null;
    metrics: MetricsPayload | null;
    delivery: DeliveryPayload | null;
  }> => {
    const [issues, metrics, delivery] = await Promise.all([
      loadIssues().catch(() => null),
      loadMetrics().catch(() => null),
      loadDelivery().catch(() => null),
    ]);
    return { issues, metrics, delivery };
  },
  component: Home,
});

const external = { target: "_blank", rel: "noopener" } as const;

function Home() {
  const data = Route.useLoaderData();
  const onboarding = docsIn("onboarding").slice(0, 6);
  const decisions = docsIn("decisions").slice(0, 3);
  const status = (id: string) => data().metrics?.services.find((s) => s.id === id)?.status ?? "unknown";

  // Only what's actually wrong: an urgent issue or a degraded API. Nothing otherwise.
  const urgent = () => data().issues?.issues.filter((i) => isOpen(i) && i.priority === 1) ?? [];
  const degraded = () => data().metrics?.services.filter((s) => s.status === "degraded") ?? [];
  const shipped = createMemo(
    () =>
      data().delivery
        ? groupReleases(data().delivery!)
            .filter((r) => releaseOutcome(r) === "shipped")
            .slice(0, 4)
        : [],
    { name: "home-shipped" },
  );

  return (
    <>
      <section class="mb-12">
        <h1 class="font-heading text-[32px] font-semibold leading-tight tracking-tight">TheTree</h1>
        <p class="mt-3 max-w-[60ch] text-[15px] leading-relaxed text-muted-foreground">
          TrieOH's monorepo: four products, the code they share and the pipeline that ships them. This portal is its
          open, read-only side, for us and for anyone who wants to read along or contribute.
        </p>

        <Show when={urgent().length > 0 || degraded().length > 0}>
          <ul class="mt-6 space-y-1.5 text-[13.5px]">
            <For each={degraded()}>
              {(s) => (
                <li>
                  <Link to="/services/$service" params={{ service: s.id }} class="inline-flex items-center gap-2 hover:underline">
                    <StatusDot status="degraded" />
                    {findService(s.id)?.name} is degraded right now.
                  </Link>
                </li>
              )}
            </For>
            <For each={urgent()}>
              {(issue) => (
                <li>
                  <Link to="/issues" class="inline-flex items-baseline gap-2 hover:underline">
                    <span class="font-medium text-destructive dark:text-destructive-foreground">Urgent</span>
                    {issue.title}
                    <span class="font-mono text-[12px] text-muted-foreground">{issue.identifier}</span>
                  </Link>
                </li>
              )}
            </For>
          </ul>
        </Show>
      </section>

      <PageSection title="The products">
        <ul class="divide-y divide-border border-y border-border">
          <For each={SERVICES}>
            {(service) => (
              <li class="flex items-baseline gap-4 py-3.5">
                <div class="min-w-0 flex-1">
                  <Link
                    to="/services/$service"
                    params={{ service: service.id }}
                    class="inline-flex items-center gap-2 font-medium hover:underline"
                  >
                    <StatusDot status={status(service.id)} />
                    {service.name}
                  </Link>
                  <p class="mt-0.5 pl-3.5 text-[13.5px] text-muted-foreground">{service.description}</p>
                </div>
                <span class="flex shrink-0 gap-4 text-[13px] text-muted-foreground">
                  <Link to="/services/$service/api" params={{ service: service.id }} class="hover:text-foreground hover:underline">
                    API
                  </Link>
                  <a href={service.front} {...external} class="hover:text-foreground hover:underline">
                    App
                  </a>
                </span>
              </li>
            )}
          </For>
        </ul>
        <p class="mt-3 text-[13px] text-muted-foreground">
          Every app signs in through IdentityX, and Univents takes payments through Payssage.
        </p>
      </PageSection>

      <PageSection
        title="New here"
        aside={
          <Link to="/docs/$category" params={{ category: "onboarding" }} class="hover:underline">
            All onboarding
          </Link>
        }
      >
        <ol class="space-y-3">
          <For each={onboarding}>
            {(doc, i) => (
              <li class="flex gap-4">
                <span class="w-5 shrink-0 pt-px text-right font-mono text-[12.5px] text-muted-foreground">{i() + 1}</span>
                <span class="min-w-0">
                  <Link
                    to="/docs/$category/$slug"
                    params={{ category: "onboarding", slug: doc.slug }}
                    class="font-medium hover:underline"
                  >
                    {doc.meta.title}
                  </Link>
                  <Show when={doc.meta.summary}>
                    <span class="block text-[13.5px] text-muted-foreground">{doc.meta.summary}</span>
                  </Show>
                </span>
              </li>
            )}
          </For>
        </ol>
      </PageSection>

      <div class="grid gap-x-12 sm:grid-cols-2">
        <PageSection
          title="Recently shipped"
          aside={
            <Link to="/delivery" class="hover:underline">
              Delivery
            </Link>
          }
        >
          <Show when={shipped().length > 0} fallback={<p class="text-[13px] text-muted-foreground">Nothing recorded yet.</p>}>
            <ul class="space-y-3">
              <For each={shipped()}>
                {(release) => (
                  <li class="text-[13.5px]">
                    <Show when={release.pr} fallback={<span class="font-medium">{release.title}</span>}>
                      {(pr) => (
                        <a href={pullUrl(pr())} {...external} class="font-medium hover:underline">
                          {release.title}
                        </a>
                      )}
                    </Show>
                    <span class="block text-[12.5px] text-muted-foreground">
                      {release.runs.length === 1 ? release.runs[0].component : `${release.runs.length} components`}
                      <Show when={release.startedAt}>{(at) => <>, {timeAgo(at())}</>}</Show>
                    </span>
                  </li>
                )}
              </For>
            </ul>
          </Show>
        </PageSection>

        <PageSection
          title="Recently decided"
          aside={
            <Link to="/docs/$category" params={{ category: "decisions" }} class="hover:underline">
              All decisions
            </Link>
          }
        >
          <ul class="space-y-3">
            <For each={decisions}>
              {(doc) => (
                <li class="text-[13.5px]">
                  <Link
                    to="/docs/$category/$slug"
                    params={{ category: "decisions", slug: doc.slug }}
                    class="font-medium hover:underline"
                  >
                    {doc.meta.title}
                  </Link>
                  <span class="block text-[12.5px] text-muted-foreground">
                    {doc.meta.number}
                    <Show when={doc.meta.date}>, {doc.meta.date}</Show>
                  </span>
                </li>
              )}
            </For>
          </ul>
        </PageSection>
      </div>
    </>
  );
}

