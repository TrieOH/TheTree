import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { For, Show, createMemo, createSignal, onCleanup, onSettled } from "solid-js";
import { findService } from "@/features/catalog/services";
import { loadApi, schemaAnchor, tagAnchor, type ApiModel } from "@/features/api-reference/model";
import { ApiNav } from "@/features/api-reference/ui/ApiNav";
import { OperationSection } from "@/features/api-reference/ui/OperationSection";
import { SchemaView } from "@/features/api-reference/ui/SchemaView";
import { TryPanel } from "@/features/api-reference/ui/TryPanel";

export const Route = createFileRoute("/services/$service/api")({
  loader: async ({ params }) => {
    const service = findService(params.service);
    if (!service) throw notFound();
    return { service, api: await loadApi(service.id) };
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.service.name ?? ""} API · TrieOH dev` }] }),
  component: ApiPage,
});

// Floating cards beside the reference: inset from the portal's sidebar,
// sticky, scrolling on their own.
const floating =
  "sticky top-6 max-h-[calc(100dvh-3rem)] self-start overflow-y-auto rounded-2xl border border-border bg-card p-3 shadow-sm";

function ApiPage() {
  const data = Route.useLoaderData();
  return (
    <Show
      when={data().api}
      fallback={
        <p class="mx-auto max-w-[760px] text-[13px] text-muted-foreground">
          This build has no generated spec for {data().service.name}. Run <code>just gen</code>.
        </p>
      }
    >
      {(api) => <Reference api={api()} />}
    </Show>
  );
}

function Reference(props: { api: ApiModel }) {
  const data = Route.useLoaderData();
  const operations = createMemo(() => props.api.tags.flatMap((tag) => tag.operations));
  const [active, setActive] = createSignal<string | null>(null, { name: "api-active" });
  const activeOp = () => operations().find((op) => op.id === active()) ?? operations()[0];

  // Which endpoint is being read: the topmost one in the upper part of the
  // screen. Read from the page only; the URL is never rewritten on scroll.
  onSettled(() => {
    const visible = new Set<Element>();
    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) visible.add(entry.target);
          else visible.delete(entry.target);
        }
        const top = [...visible].sort((a, b) => a.getBoundingClientRect().top - b.getBoundingClientRect().top)[0];
        if (top) setActive((top as HTMLElement).dataset.operation ?? null);
      },
      { rootMargin: "-10% 0px -70% 0px" },
    );
    for (const element of document.querySelectorAll("[data-operation]")) observer.observe(element);
    onCleanup(() => observer.disconnect());
  });

  return (
    <div class="grid gap-6 lg:grid-cols-[220px_minmax(0,1fr)] xl:grid-cols-[230px_minmax(0,1fr)_360px]">
      <aside class={`hidden lg:block ${floating}`}>
        <ApiNav api={props.api} active={activeOp()?.id ?? null} />
      </aside>

      <div class="min-w-0 space-y-10">
        <header class="space-y-3">
          <nav class="flex items-center gap-2 text-[13px] text-muted-foreground">
            <Link to="/services/$service" params={{ service: data().service.id }} class="hover:text-foreground">
              {data().service.name}
            </Link>
            <span>/</span>
            <span class="text-foreground">API reference</span>
          </nav>
          <h1 class="font-heading text-[28px] font-semibold leading-tight tracking-tight">{props.api.title}</h1>
          <p class="font-mono text-[12px] text-muted-foreground">
            v{props.api.version} · base URL {props.api.baseUrl}
          </p>
          <details class="rounded-xl border border-border px-4 py-3">
            <summary class="cursor-pointer text-[13px] font-medium">Overview and authentication</summary>
            <div class="prose-doc mt-3 text-[13.5px]" innerHTML={props.api.description} />
            <ul class="mt-3 space-y-1 text-[13px]">
              <For each={props.api.auth}>
                {(auth) => (
                  <li>
                    <code class="font-mono">{auth.name}</code>
                    <span class="text-muted-foreground">: {auth.summary}</span>
                  </li>
                )}
              </For>
            </ul>
          </details>
        </header>

        <For each={props.api.tags}>
          {(tag) => (
            <section id={tagAnchor(tag.name)} class="scroll-mt-6 space-y-8">
              <div>
                <h2 class="font-heading text-xl font-semibold">{tag.name}</h2>
                <Show when={tag.description}>
                  <div class="prose-doc mt-1 text-[13.5px] text-muted-foreground" innerHTML={tag.description} />
                </Show>
              </div>
              <For each={tag.operations}>
                {(op) => <OperationSection op={op} service={data().service.id} auth={props.api.auth} />}
              </For>
            </section>
          )}
        </For>

        <section id="models" class="scroll-mt-6 space-y-2">
          <h2 class="mb-3 font-heading text-xl font-semibold">Models</h2>
          <For each={Object.entries(props.api.schemas)}>
            {([name, node]) => (
              <details id={schemaAnchor(name)} class="scroll-mt-6 rounded-xl border border-border">
                <summary class="cursor-pointer px-3 py-2.5 font-mono text-[13px] font-medium hover:bg-foreground/5">
                  {name}
                </summary>
                <div class="border-t border-border p-3">
                  <SchemaView node={node} />
                </div>
              </details>
            )}
          </For>
        </section>
      </div>

      <aside class={`hidden xl:block ${floating} p-4`}>
        <Show when={activeOp()} keyed>
          {(op) => <TryPanel op={op} auth={props.api.auth} />}
        </Show>
      </aside>
    </div>
  );
}
