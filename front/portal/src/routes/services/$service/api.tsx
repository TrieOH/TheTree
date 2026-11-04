import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { For, Show, createEffect, createMemo, createSignal, onCleanup, onSettled } from "solid-js";
import { findService } from "@/features/catalog/services";
import { loadApi, schemaAnchor, tagAnchor, type ApiModel } from "@/features/api-reference/model";
import { ApiNav } from "@/features/api-reference/ui/ApiNav";
import { OperationSection } from "@/features/api-reference/ui/OperationSection";
import { SchemaView } from "@/features/api-reference/ui/SchemaView";

export const Route = createFileRoute("/services/$service/api")({
  loader: async ({ params }) => {
    const service = findService(params.service);
    if (!service) throw notFound();
    return { service, api: await loadApi(service.id) };
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.service.name ?? ""} API · TrieOH dev` }] }),
  component: ApiPage,
});

// The route list floats beside the reference: inset from the portal's
// sidebar, one fixed height centered on the screen, scrolling on its own.
const NAV_HEIGHT = "min(80dvh,760px)";
const floating =
  "sticky self-start overflow-y-auto rounded-2xl border border-border bg-card p-3 shadow-lg shadow-black/20 ring-1 ring-foreground/5 dark:shadow-black/50";

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
  // Nothing above the first endpoint.
  const activeOp = () => operations().find((op) => op.id === active()) ?? null;

  // Which endpoint is being read: the last one whose top has passed a line
  // a quarter down the screen. Read from the page on scroll (once a frame);
  // the URL is never rewritten.
  let nav: HTMLElement | undefined;
  onSettled(() => {
    let frame = 0;
    const update = () => {
      frame = 0;
      const line = window.innerHeight * 0.25;
      let current: string | null = null;
      for (const element of document.querySelectorAll<HTMLElement>("[data-operation]")) {
        if (element.getBoundingClientRect().top > line) break;
        current = element.dataset.operation ?? null;
      }
      setActive(current);
    };
    const onScroll = () => {
      if (!frame) frame = requestAnimationFrame(update);
    };
    update();
    window.addEventListener("scroll", onScroll, { passive: true });
    onCleanup(() => {
      window.removeEventListener("scroll", onScroll);
      cancelAnimationFrame(frame);
    });
  });

  // The nav follows: whatever marks the active endpoint (its route, or its
  // group while closed) is kept in the middle of the nav's own scroll.
  createEffect(
    () => activeOp()?.id,
    (id) => {
      if (!nav || !id) return;
      const tag = props.api.tags.find((t) => t.operations.some((op) => op.id === id))?.name;
      const target =
        nav.querySelector<HTMLElement>(`[data-nav-op="${CSS.escape(id)}"]`) ??
        (tag ? nav.querySelector<HTMLElement>(`[data-nav-tag="${CSS.escape(tag)}"]`) : null);
      if (!target) return;
      const offset = target.getBoundingClientRect().top - nav.getBoundingClientRect().top + nav.scrollTop;
      nav.scrollTo({ top: offset - nav.clientHeight / 2 + target.clientHeight / 2, behavior: "smooth" });
    },
  );

  return (
    // The reference sits centered like every page; the nav floats centered in
    // the space to its left.
    <div class="grid gap-8 lg:grid-cols-[minmax(250px,1fr)_minmax(0,860px)_minmax(0,1fr)]">
      <aside
        ref={(element) => (nav = element)}
        class={`hidden w-[250px] justify-self-center lg:block ${floating}`}
        style={{ height: NAV_HEIGHT, top: `calc((100dvh - ${NAV_HEIGHT}) / 2)` }}
      >
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

    </div>
  );
}
