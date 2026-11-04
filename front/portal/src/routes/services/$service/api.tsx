import { Link, createFileRoute, notFound, useLocation } from "@tanstack/solid-router";
import { For, Show } from "solid-js";
import { findService } from "@/features/catalog/services";
import { loadApi, operationAnchor, schemaAnchor, tagAnchor } from "@/features/api-reference/model";
import { OperationCard } from "@/features/api-reference/ui/OperationCard";
import { SchemaView } from "@/features/api-reference/ui/SchemaView";
import { PageHeader, Section } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/services/$service/api")({
  loader: async ({ params }) => {
    const service = findService(params.service);
    if (!service) throw notFound();
    return { service, api: await loadApi(service.id) };
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.service.name ?? ""} API · TrieOH dev` }] }),
  component: ApiPage,
});

const chip =
  "inline-flex h-7 items-center gap-1.5 rounded-full border border-border px-3 text-[12.5px] text-muted-foreground transition-colors hover:bg-foreground/5 hover:text-foreground";

function ApiPage() {
  const data = Route.useLoaderData();
  const location = useLocation();
  const hash = () => location().hash;

  return (
    <>
      <nav class="mb-4 flex items-center gap-2 text-[13px] text-muted-foreground">
        <Link to="/services/$service" params={{ service: data().service.id }} class="hover:text-foreground">
          {data().service.name}
        </Link>
        <span>/</span>
        <span class="text-foreground">API reference</span>
      </nav>

      <Show
        when={data().api}
        fallback={
          <p class="text-[13px] text-muted-foreground">
            This build has no generated spec for {data().service.name}. Run <code>just gen</code>.
          </p>
        }
      >
        {(api) => (
          <>
            <PageHeader title={`${api().title}`}>
              <p class="font-mono text-[12px] text-muted-foreground">
                v{api().version} · base URL {api().baseUrl}
              </p>
            </PageHeader>

            <div class="mb-8 flex flex-wrap gap-1.5">
              <For each={api().tags}>
                {(tag) => (
                  <a href={`#${tagAnchor(tag.name)}`} class={chip}>
                    {tag.name}
                    <span class="tabular-nums opacity-70">{tag.operations.length}</span>
                  </a>
                )}
              </For>
              <a href="#models" class={chip}>
                models
                <span class="tabular-nums opacity-70">{Object.keys(api().schemas).length}</span>
              </a>
            </div>

            <details class="mb-10 rounded-xl border border-border px-4 py-3">
              <summary class="cursor-pointer text-[13px] font-medium">Overview and authentication</summary>
              <div class="prose-doc mt-3 text-[13.5px]" innerHTML={api().description} />
              <ul class="mt-3 space-y-1 text-[13px]">
                <For each={api().auth}>
                  {(auth) => (
                    <li>
                      <code class="font-mono">{auth.name}</code>
                      <span class="text-muted-foreground">: {auth.summary}</span>
                    </li>
                  )}
                </For>
              </ul>
            </details>

            <For each={api().tags}>
              {(tag) => (
                <section id={tagAnchor(tag.name)} class="mb-10 scroll-mt-20">
                  <h2 class="mb-1 font-heading text-lg font-semibold">{tag.name}</h2>
                  <Show when={tag.description}>
                    <div class="prose-doc mb-3 text-[13.5px] text-muted-foreground" innerHTML={tag.description} />
                  </Show>
                  <div class="space-y-2">
                    <For each={tag.operations}>
                      {(op) => (
                        <OperationCard op={op} service={data().service.id} open={hash() === operationAnchor(op)} />
                      )}
                    </For>
                  </div>
                </section>
              )}
            </For>

            <Section title="Models">
              <div id="models" class="scroll-mt-20 space-y-2">
                <For each={Object.entries(api().schemas)}>
                  {([name, node]) => (
                    <details id={schemaAnchor(name)} open={hash() === schemaAnchor(name)} class="scroll-mt-20 rounded-xl border border-border">
                      <summary class="cursor-pointer px-3 py-2.5 font-mono text-[13px] font-medium hover:bg-foreground/5">
                        {name}
                      </summary>
                      <div class="border-t border-border p-3">
                        <Show when={!("ref" in node) && node.description}>
                          <p class="mb-2 text-[13px] text-muted-foreground">{!("ref" in node) && node.description}</p>
                        </Show>
                        <SchemaView node={node} />
                      </div>
                    </details>
                  )}
                </For>
              </div>
            </Section>
          </>
        )}
      </Show>
    </>
  );
}
