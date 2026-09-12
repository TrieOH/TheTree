import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core-solid";
import { Button, EmptyState, PaginatedContainer, type SortState } from "@trieoh/ui-solid";
import { For, Loading, Show, createMemo, createSignal } from "solid-js";

import CalendarIcon from "~icons/lucide/calendar";
import PlusIcon from "~icons/lucide/plus";

import {
  allJoinedEventsQueryOptions,
  allOwnEventsQueryOptions,
  createEventFn,
  patchEventFn,
} from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import { AdminEventCard } from "@/features/events/ui/AdminEventCard";
import {
  ManageEventDialog,
  type ManageEventValues,
} from "@/features/events/ui/ManageEventDialog";

// The icon casts are the app's existing convention: `unplugin-icons` resolves
// the icon at build time and this package only ever sees an element. The props
// are declared because casting to `() => JSX.Element` drops them.
const Calendar = CalendarIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute("/admin/events/")({
  head: () => ({ meta: [{ title: "Eventos - Admin - Univents" }] }),
  component: AdminEventsPage,
});

const STATUS_SORT_ORDER: Record<EventI["status"], number> = {
  draft: 0,
  active: 1,
  discontinued: 2,
};

function AdminEventsPage(): JSX.Element {
  const queryClient = useQueryClient();
  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<EventI>>({
    field: "created_at",
    direction: "desc",
  });
  const [editing, setEditing] = createSignal<EventI | null>(null);
  const [creating, setCreating] = createSignal(false);

  const events = createMemo(async () => {
    const [owned, joined] = await Promise.all([
      queryClient.fetchQuery(allOwnEventsQueryOptions()),
      queryClient.fetchQuery(allJoinedEventsQueryOptions()),
    ]);

    // Owned first, then the ones the actor only takes part in, without repeats.
    const seen = new Set<string>();
    return [...owned, ...joined].filter((event) => {
      if (seen.has(event.id)) return false;
      seen.add(event.id);
      return true;
    });
  });

  // Named instead of inline: an `async` arrow inside JSX trips the
  // "tracked scope should not be async" diagnostic. The dialog already trims.
  const saveEvent = async (values: ManageEventValues): Promise<boolean> => {
    const current = editing();
    const payload = {
      full_name: values.full_name,
      slug: values.slug,
      acronym: values.acronym || null,
      description: values.description || null,
      contact_email: values.contact_email || null,
    };

    try {
      if (current) await patchEventFn(current.id, payload);
      else await createEventFn(payload);
    } catch {
      return false;
    }

    await Promise.all([
      queryClient.invalidateQueries({ queryKey: allOwnEventsQueryOptions().queryKey }),
      queryClient.invalidateQueries({ queryKey: allJoinedEventsQueryOptions().queryKey }),
    ]);

    setCreating(false);
    setEditing(null);
    return true;
  };

  return (
    <Loading
      fallback={
        <div class="grid grid-cols-[repeat(auto-fill,minmax(16rem,1fr))] gap-6">
          <For each={[1, 2, 3, 4]}>
            {() => <div class="h-40 animate-pulse rounded-lg bg-muted" />}
          </For>
        </div>
      }
    >
      <AdminEventsContent
        events={events()}
        filter={filter()}
        onFilterChange={setFilter}
        sort={sort()}
        onSortChange={setSort}
        onEdit={setEditing}
        onCreate={() => setCreating(true)}
      />

      {/* Keyed so switching events remounts the form instead of needing a
          re-seed effect inside it. */}
      <Show when={editing()} keyed fallback={
        <ManageEventDialog
          open={creating()}
          event={null}
          onOpenChange={(open) => !open && setCreating(false)}
          onSubmit={saveEvent}
        />
      }>
        {(current) => (
          <ManageEventDialog
            open
            event={current}
            onOpenChange={(open) => !open && setEditing(null)}
            onSubmit={saveEvent}
          />
        )}
      </Show>
    </Loading>
  );
}

function AdminEventsContent(props: {
  events: EventI[];
  filter: string;
  onFilterChange: (value: string) => void;
  sort: SortState<EventI>;
  onSortChange: (sort: SortState<EventI>) => void;
  onEdit: (event: EventI) => void;
  onCreate: () => void;
}): JSX.Element {
  // Filtering only: sorting and pagination belong to the container, same split
  // as the React admin screen.
  const visible = createMemo(() => {
    const search = props.filter.trim().toLowerCase();
    if (!search) return props.events;

    return props.events.filter((event) =>
      [
        event.full_name,
        event.slug,
        event.acronym ?? "",
        event.contact_email ?? "",
        event.status,
      ].some((value) => value.toLowerCase().includes(search)),
    );
  });

  const createButton = (
    <Button size="sm" class="gap-2 rounded-sm py-4" onClick={() => props.onCreate()}>
      <Plus class="size-4" />
      Novo evento
    </Button>
  );

  return (
    <PaginatedContainer<EventI>
      items={visible()}
      layout="grid"
      minItemWidth="16rem"
      maxRows={2}
      gap="6"
      sort={props.sort}
      onSortChange={props.onSortChange}
      sortFields={[
        {
          key: "created_at",
          label: "Criado em",
          comparator: (a, b) =>
            new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
        },
        { key: "full_name", label: "Nome" },
        { key: "slug", label: "Slug" },
        {
          key: "status",
          label: "Status",
          comparator: (a, b) => STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status],
        },
      ]}
      filterValue={props.filter}
      onFilterChange={props.onFilterChange}
      filterPlaceholder="Buscar por nome, slug, sigla ou e-mail..."
      itemLabel="eventos"
      headerActions={createButton}
      renderItems={(slice, options) => (
        <For each={slice}>
          {(event, index) => (
            <AdminEventCard
              event={event}
              index={index()}
              animate={options.animate}
              onEdit={props.onEdit}
            />
          )}
        </For>
      )}
      emptyState={
        <EmptyState
          class="border-0 bg-transparent px-0 py-4 shadow-none"
          icon={<Calendar class="size-6" />}
          eyebrow="Eventos"
          title="Nenhum evento encontrado"
          description="Crie um evento para começar a organizar o dashboard do admin."
          action={createButton}
        />
      }
    />
  );
}

