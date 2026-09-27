import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { Button, EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal } from "solid-js";

import TicketIcon from "~icons/lucide/ticket";
import TicketPlusIcon from "~icons/lucide/ticket-plus";

import { allTicketsQueryOptions } from "@/features/tickets/api";
import {
  useCreateTicketMutation,
  usePatchTicketMutation,
} from "@/features/tickets/api/mutations";
import type { TicketCreateOutputI, TicketI } from "@/features/tickets/model";
import { AdminCreateTicketCard } from "@/features/tickets/ui/AdminCreateTicketCard";
import { AdminTicketCard } from "@/features/tickets/ui/AdminTicketCard";
import { ManageTicketDialog } from "@/features/tickets/ui/ManageTicketDialog";
import { toast } from "@/shared/ui/toast";

const Ticket = TicketIcon as unknown as (props: { class?: string }) => JSX.Element;
const TicketPlus = TicketPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/tickets/",
)({
  head: () => ({
    meta: [{ title: "Tickets - Admin Univents" }],
  }),
  component: AdminEditionTicketsRoute,
});

function AdminEditionTicketsRoute(): JSX.Element {
  const params = Route.useParams();
  const editionId = () => params().editionId;

  const ticketsQuery = useQuery(() => allTicketsQueryOptions(editionId()));
  const createMutation = useCreateTicketMutation();
  const patchMutation = usePatchTicketMutation();

  const [creating, setCreating] = createSignal(false);
  const [editing, setEditing] = createSignal<TicketI | null>(null);

  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<TicketI>>({
    field: "name",
    direction: "asc",
  });

  const tickets = createMemo(() => (ticketsQuery().data ?? []) as TicketI[]);

  const saveTicket = async (values: TicketCreateOutputI): Promise<boolean> => {
    const current = editing();
    try {
      if (current) {
        await patchMutation.mutateAsync({
          ticketId: current.id,
          editionId: editionId(),
          data: {
            name: values.name,
            description: values.description ?? null,
            price_cents: values.price_cents,
            access_level: values.access_level,
            max_quantity: values.max_quantity ?? null,
          },
        });
        toast.success("Ticket atualizado com sucesso!");
      } else {
        await createMutation.mutateAsync({
          editionId: editionId(),
          data: {
            name: values.name,
            description: values.description ?? null,
            price_cents: values.price_cents,
            access_level: values.access_level,
            max_quantity: values.max_quantity ?? null,
          },
        });
        toast.success("Ticket criado com sucesso!");
      }
      setCreating(false);
      setEditing(null);
      return true;
    } catch {
      toast.error(current ? "Erro ao atualizar ticket." : "Erro ao criar ticket.");
      return false;
    }
  };

  return (
    <div class="mx-auto max-w-7xl space-y-6">
      <Show
        when={!ticketsQuery().isLoading}
        fallback={
          <div class="grid grid-cols-[repeat(auto-fill,minmax(16rem,1fr))] gap-2">
            <For each={[1, 2, 3, 4]}>
              {() => <div class="h-24 animate-pulse rounded-xl bg-muted" />}
            </For>
          </div>
        }
      >
        <AdminEditionTicketsContent
          tickets={tickets()}
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
          <ManageTicketDialog
            open={creating()}
            ticket={null}
            onOpenChange={(open) => !open && setCreating(false)}
            onSubmit={saveTicket}
          />
        }
      >
        {(current) => (
          <ManageTicketDialog
            open
            ticket={current}
            onOpenChange={(open) => !open && setEditing(null)}
            onSubmit={saveTicket}
          />
        )}
      </Show>
    </div>
  );
}

function AdminEditionTicketsContent(props: {
  tickets: TicketI[];
  filter: string;
  onFilterChange: (value: string) => void;
  sort: SortState<TicketI>;
  onSortChange: (sort: SortState<TicketI>) => void;
  onCreate: () => void;
  onEdit: (ticket: TicketI) => void;
}): JSX.Element {
  const visible = createMemo(() => {
    const search = props.filter.trim().toLowerCase();
    if (!search) return props.tickets;

    return props.tickets.filter((ticket) =>
      [
        ticket.name,
        ticket.description ?? "",
        String(ticket.access_level),
        String(ticket.price_cents),
      ].some((val) => val.toLowerCase().includes(search)),
    );
  });

  const emptyStateAction = (
    <Button
      size="sm"
      class="gap-2 rounded-sm py-4"
      onClick={() => props.onCreate()}
    >
      <TicketPlus class="size-4" />
      Novo ticket
    </Button>
  );

  return (
    <PaginatedContainer<TicketI>
      items={visible()}
      layout="grid"
      minItemWidth="16rem"
      maxRows={(columns) => (columns === 1 ? 8 : 4)}
      gap="2"
      sort={props.sort}
      onSortChange={props.onSortChange}
      sortFields={[
        {
          key: "name",
          label: "Nome",
          ascLabel: "A → Z",
          descLabel: "Z → A",
        },
        {
          key: "price_cents",
          label: "Preço",
          ascLabel: "Menor preço primeiro",
          descLabel: "Maior preço primeiro",
          comparator: (a, b) => a.price_cents - b.price_cents,
        },
        {
          key: "access_level",
          label: "Nível de acesso",
          ascLabel: "Menor nível primeiro",
          descLabel: "Maior nível primeiro",
          comparator: (a, b) => a.access_level - b.access_level,
        },
      ]}
      filterValue={props.filter}
      onFilterChange={props.onFilterChange}
      filterPlaceholder="Buscar por nome, descrição ou nível..."
      itemLabel="tickets"
      emptyState={
        <EmptyState
          class="border-0 bg-transparent px-0 py-4 shadow-none"
          icon={<Ticket class="size-6 text-foreground/70" />}
          eyebrow="Ingressos da edição"
          title="Nenhum ticket encontrado"
          description={
            props.filter
              ? "Nenhum ticket corresponde à busca informada."
              : "Crie o primeiro ticket para começar a vender nesta edição."
          }
          action={emptyStateAction}
        />
      }
      renderItems={(slice, options) => (
        <>
          <AdminCreateTicketCard
            index={0}
            animate={options.animate}
            onCreate={props.onCreate}
          />

          <For each={slice}>
            {(ticket, index) => (
              <AdminTicketCard
                ticket={ticket}
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
