import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { Button, EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal, untrack } from "solid-js";

import CalendarDaysIcon from "~icons/lucide/calendar-days";
import CalendarPlusIcon from "~icons/lucide/calendar-plus";

import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import {
  useCreateEditionMutation,
  usePatchEditionMutation,
} from "@/features/editions/api/mutations";
import {
  editionRangesOverlap,
  type EditionI,
  type EditionStatus,
} from "@/features/editions/model";
import { AdminCreateEditionCard } from "@/features/editions/ui/AdminCreateEditionCard";
import { AdminEditionCard } from "@/features/editions/ui/AdminEditionCard";
import {
  ManageEditionDialog,
  type ManageEditionValues,
} from "@/features/editions/ui/ManageEditionDialog";
import { toast } from "@/shared/ui/toast";

const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;
const CalendarPlus = CalendarPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute("/admin/events/$eventId/editions/")({
  head: () => ({
    meta: [{ title: "Edições - Admin Univents" }],
  }),
  component: AdminEventEditionsRoute,
});

const STATUS_LABELS: Record<EditionStatus, string> = {
  draft: "Rascunho",
  future: "Futura",
  active: "Ativa",
  past: "Encerrada",
};

const STATUS_SORT_ORDER: Record<EditionStatus, number> = {
  active: 0,
  future: 1,
  draft: 2,
  past: 3,
};

function AdminEventEditionsRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;

  const editionsQuery = useQuery(() => allAdminEditionsQueryOptions(eventId()));
  const createMutation = useCreateEditionMutation();
  const patchMutation = usePatchEditionMutation();

  const [creating, setCreating] = createSignal(false);
  const [editing, setEditing] = createSignal<EditionI | null>(null);

  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<EditionI>>({
    field: "starts_at",
    direction: "desc",
  });

  const editions = createMemo(() => (editionsQuery().data ?? []) as EditionI[]);

  const saveEdition = async (values: ManageEditionValues): Promise<boolean> => {
    const currentTarget = untrack(editing);
    const currentEditions = untrack(editions);

    const overlaps = currentEditions.some((item) => {
      if (currentTarget && item.id === currentTarget.id) return false;
      return editionRangesOverlap(item, values);
    });

    if (overlaps) {
      toast.error("Já existe uma edição cadastrada nesse período.");
      return false;
    }

    try {
      if (currentTarget) {
        await patchMutation.mutateAsync({
          eventId: eventId(),
          editionId: currentTarget.id,
          data: {
            name: values.name,
            slug: values.slug,
            starts_at: values.starts_at,
            ends_at: values.ends_at,
            location_name: values.location_name || null,
          },
        });
        toast.success("Edição atualizada com sucesso!");
      } else {
        await createMutation.mutateAsync({
          eventId: eventId(),
          data: {
            name: values.name,
            slug: values.slug,
            starts_at: values.starts_at,
            ends_at: values.ends_at,
          },
        });
        toast.success("Edição criada com sucesso!");
      }
      setCreating(false);
      setEditing(null);
      return true;
    } catch {
      toast.error(currentTarget ? "Erro ao atualizar edição." : "Erro ao criar edição.");
      return false;
    }
  };

  return (
    <div class="space-y-6">
      <Show
        when={!editionsQuery().isLoading}
        fallback={
          <div class="grid grid-cols-[repeat(auto-fill,minmax(16rem,1fr))] gap-2">
            <For each={[1, 2, 3, 4]}>
              {() => <div class="h-24 animate-pulse rounded-xl bg-muted" />}
            </For>
          </div>
        }
      >
        <AdminEventEditionsContent
          editions={editions()}
          eventId={eventId()}
          filter={filter()}
          onFilterChange={setFilter}
          sort={sort()}
          onSortChange={setSort}
          onCreate={() => setCreating(true)}
          onEdit={setEditing}
        />
      </Show>

      <Show
        when={editing()}
        keyed
        fallback={
          <ManageEditionDialog
            open={creating()}
            edition={null}
            onOpenChange={(open) => !open && setCreating(false)}
            onSubmit={saveEdition}
          />
        }
      >
        {(current) => (
          <ManageEditionDialog
            open
            edition={current}
            onOpenChange={(open) => !open && setEditing(null)}
            onSubmit={saveEdition}
          />
        )}
      </Show>
    </div>
  );
}

function AdminEventEditionsContent(props: {
  editions: EditionI[];
  eventId: string;
  filter: string;
  onFilterChange: (value: string) => void;
  sort: SortState<EditionI>;
  onSortChange: (sort: SortState<EditionI>) => void;
  onCreate: () => void;
  onEdit: (edition: EditionI) => void;
}): JSX.Element {
  const visible = createMemo(() => {
    const search = props.filter.trim().toLowerCase();
    if (!search) return props.editions;

    return props.editions.filter((edition) =>
      [
        edition.name,
        edition.slug,
        edition.location_name,
        STATUS_LABELS[edition.status],
        edition.status,
      ].some((val) => (val ? val.toLowerCase().includes(search) : false)),
    );
  });

  const emptyStateAction = (
    <Button
      size="sm"
      class="gap-2 rounded-sm py-4"
      onClick={() => props.onCreate()}
    >
      <CalendarPlus class="size-4" />
      Nova edição
    </Button>
  );

  return (
    <PaginatedContainer<EditionI>
      items={visible()}
      layout="grid"
      minItemWidth="16rem"
      maxRows={(columns) => (columns === 1 ? 8 : 4)}
      gap="2"
      sort={props.sort}
      onSortChange={props.onSortChange}
      sortFields={[
        {
          key: "starts_at",
          label: "Início",
          ascLabel: "Mais antigas primeiro",
          descLabel: "Mais recentes primeiro",
          comparator: (a, b) =>
            new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime(),
        },
        {
          key: "name",
          label: "Nome",
          ascLabel: "A → Z",
          descLabel: "Z → A",
        },
        {
          key: "status",
          label: "Status",
          ascLabel: "Ativa primeiro",
          descLabel: "Encerrada primeiro",
          comparator: (a, b) =>
            STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status],
        },
        {
          key: "slug",
          label: "Slug",
          ascLabel: "A → Z",
          descLabel: "Z → A",
        },
      ]}
      filterValue={props.filter}
      onFilterChange={props.onFilterChange}
      filterPlaceholder="Buscar por nome, local, slug ou status..."
      itemLabel="edições"
      emptyState={
        <EmptyState
          class="border-0 bg-transparent px-0 py-4 shadow-none"
          icon={<CalendarDays class="size-6 text-foreground/70" />}
          eyebrow="Edições do evento"
          title="Nenhuma edição cadastrada"
          description={
            props.filter
              ? "Nenhuma edição corresponde à busca informada."
              : "Crie a primeira edição para começar a organizar este evento."
          }
          action={emptyStateAction}
        />
      }
      renderItems={(slice, options) => (
        <>
          <AdminCreateEditionCard
            index={0}
            animate={options.animate}
            onCreate={props.onCreate}
          />

          <For each={slice}>
            {(edition, index) => (
              <AdminEditionCard
                edition={edition}
                eventId={props.eventId}
                index={index() + 1}
                animate={options.animate}
                onEdit={props.onEdit}
              />
            )}
          </For>
        </>
      )}
    />
  );
}
