import { createFileRoute } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core/solid";
import { For, Loading, Show, createMemo, createSignal } from "solid-js";
import type { JSX } from "@solidjs/web";
import SearchIcon from "~icons/lucide/search";
import SlidersIcon from "~icons/lucide/sliders-horizontal";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import { EventCard } from "@/widgets/landing/ui/EventCard";
import Logo from "@/shared/ui/Logo";
import { Drawer } from "@/shared/ui/Drawer";

const Search = SearchIcon as unknown as () => JSX.Element;
const Sliders = SlidersIcon as unknown as () => JSX.Element;

export const Route = createFileRoute("/events/")({
  head: () => ({ meta: [{ title: "Eventos - Univents" }] }),
  component: EventsPage,
});

function EventsPage() {
  const queryClient = useQueryClient();
  const [filter, setFilter] = createSignal<"all" | "series">("all");
  const [drawerOpen, setDrawerOpen] = createSignal(false);

  const events = createMemo(() =>
    queryClient.fetchQuery(allPublicEventsQueryOptions()),
  );

  return (
    <Loading
      fallback={
        <EventsShell count="...">
          <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 md:gap-8">
            <For each={[1, 2, 3, 4]}>
              {() => (
                <div class="aspect-4/3 animate-pulse rounded-2xl bg-muted" />
              )}
            </For>
          </div>
        </EventsShell>
      }
    >
      <EventsContent
        events={events()}
        filter={filter}
        setFilter={setFilter}
        drawerOpen={drawerOpen}
        setDrawerOpen={setDrawerOpen}
      />
    </Loading>
  );
}

function EventsContent(props: {
  events: EventI[];
  filter: () => "all" | "series";
  setFilter: (value: "all" | "series") => void;
  drawerOpen: () => boolean;
  setDrawerOpen: (value: boolean) => void;
}) {
  const filteredEvents = createMemo(() =>
    [...props.events]
      .sort(
        (a, b) =>
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
      )
      .filter(() => props.filter() === "all"),
  );

  return (
    <EventsShell
      count={String(filteredEvents().length)}
      actions={
        <>
          <div class="hidden rounded-lg bg-muted p-1 sm:flex">
            <FilterButton
              active={props.filter() === "all"}
              onClick={() => props.setFilter("all")}
            >
              Todos
            </FilterButton>
            <FilterButton
              active={props.filter() === "series"}
              onClick={() => props.setFilter("series")}
            >
              Séries
            </FilterButton>
          </div>
          <button
            type="button"
            aria-label="Filtrar eventos"
            class={`flex size-9 items-center justify-center rounded-lg sm:hidden ${props.drawerOpen() ? "bg-muted" : "hover:bg-muted"}`}
            onClick={() => props.setDrawerOpen(true)}
          >
            <Sliders />
          </button>
          <Drawer
            open={props.drawerOpen()}
            title="Filtrar eventos"
            onClose={() => props.setDrawerOpen(false)}
          >
            <div class="space-y-1">
              <FilterButton
                active={props.filter() === "all"}
                onClick={() => {
                  props.setFilter("all");
                  props.setDrawerOpen(false);
                }}
              >
                Todos os eventos
              </FilterButton>
              <FilterButton
                active={props.filter() === "series"}
                onClick={() => {
                  props.setFilter("series");
                  props.setDrawerOpen(false);
                }}
              >
                Apenas séries
              </FilterButton>
            </div>
          </Drawer>
        </>
      }
    >
      <Show
        when={filteredEvents().length > 0}
        fallback={
          <div class="flex flex-col items-center justify-center space-y-6 py-12 md:py-16">
            <div class="flex size-16 items-center justify-center rounded-full bg-muted">
              <Search />
            </div>
            <div class="space-y-1 text-center">
              <h2 class="text-lg font-medium">Nenhum evento encontrado</h2>
              <p class="text-sm text-muted-foreground">
                Tente ajustar os filtros ou volte mais tarde.
              </p>
            </div>
          </div>
        }
      >
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 md:gap-8">
          <For each={filteredEvents()}>
            {(event, index) => <EventCard event={event} index={index()} />}
          </For>
        </div>
      </Show>
    </EventsShell>
  );
}

function FilterButton(props: {
  active: boolean;
  onClick: () => void;
  children: string;
}) {
  return (
    <button
      type="button"
      class={`rounded-md px-3 py-1.5 text-sm transition-all ${props.active ? "bg-background text-foreground shadow-sm" : "text-muted-foreground hover:text-foreground"}`}
      onClick={() => props.onClick()}
    >
      {props.children}
    </button>
  );
}

function EventsShell(props: {
  count: string;
  actions?: JSX.Element;
  children: JSX.Element;
}) {
  return (
    <div class="min-h-screen bg-background pb-24">
      <header class="sticky top-0 z-30 border-b border-border bg-background/80 backdrop-blur-xl">
        <div class="mx-auto flex h-14 max-w-7xl items-center justify-between gap-2 px-4">
          <div class="flex items-center gap-3">
            <div class="size-8">
              <Logo variant="icon" imgClassName="object-left" />
            </div>
            <h1 class="border-l border-border pl-3 text-lg font-semibold md:text-xl">
              Eventos{" "}
              <span class="ml-2 text-sm font-normal text-muted-foreground">
                ({props.count})
              </span>
            </h1>
          </div>
          <div class="ml-auto flex items-center gap-2">{props.actions}</div>
        </div>
      </header>
      <main class="mx-auto max-w-7xl px-4 py-8 md:py-12">{props.children}</main>
    </div>
  );
}
