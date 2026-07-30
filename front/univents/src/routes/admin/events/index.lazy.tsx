import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { Calendar, Plus } from "lucide-react";
import { useState } from "react";
import {
  allJoinedEventsQueryOptions,
  allOwnEventsQueryOptions,
} from "@/features/events/api";
import {
  useCreateEventMutation,
  useDiscontinueEventMutation,
  usePublishEventMutation,
} from "@/features/events/api/mutations";
import type { EventI } from "@/features/events/model";
import AdminEventCard from "@/features/events/ui/AdminEventCard";
import { ManageEventModal } from "@/features/events/ui/ManageEventModal";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";

export const Route = createLazyFileRoute("/admin/events/")({
  component: RouteComponent,
});

const STATUS_SORT_ORDER: Record<EventI["status"], number> = {
  draft: 0,
  active: 1,
  discontinued: 2,
};

function RouteComponent() {
  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState<SortState<EventI>>({
    field: "created_at",
    direction: "desc",
  });
  const [modalOpen, setModalOpen] = useState(false);
  const [publishingEvent, setPublishingEvent] = useState<EventI | null>(null);
  const [discontinuingEvent, setDiscontinuingEvent] = useState<EventI | null>(
    null,
  );

  const { data: ownEvents = [] } = useQuery(allOwnEventsQueryOptions());
  const { data: joinedEvents = [] } = useQuery(allJoinedEventsQueryOptions());
  const events = [...ownEvents, ...joinedEvents].filter(
    (event, index, list) =>
      list.findIndex((candidate) => candidate.id === event.id) === index,
  );
  const createMutation = useCreateEventMutation();
  const publishEventMutation = usePublishEventMutation();
  const discontinueEventMutation = useDiscontinueEventMutation();

  const filteredEvents = [...events]
    .filter((event) => {
      const search = filter.trim().toLowerCase();
      if (!search) return true;

      return [
        event.full_name,
        event.slug,
        event.acronym ?? "",
        event.contact_email ?? "",
        event.status,
      ].some((value) => value.toLowerCase().includes(search));
    })
    .sort((a, b) => {
      const direction = sort.direction === "asc" ? 1 : -1;

      if (sort.field === "created_at") {
        return (
          (new Date(a.created_at).getTime() -
            new Date(b.created_at).getTime()) *
          direction
        );
      }

      if (sort.field === "status") {
        return (
          (STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status]) *
          direction
        );
      }

      return (
        String(a[sort.field]).localeCompare(String(b[sort.field])) * direction
      );
    });

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<EventI>
        items={filteredEvents}
        layout="grid"
        minItemWidth="16rem"
        pageSize={8}
        gap="6"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          {
            key: "created_at",
            label: "Criado em",
            comparator: (a, b) =>
              new Date(a.created_at).getTime() -
              new Date(b.created_at).getTime(),
          },
          { key: "full_name", label: "Nome" },
          { key: "slug", label: "Slug" },
          {
            key: "status",
            label: "Status",
            comparator: (a, b) =>
              STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status],
          },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome, slug, sigla ou e-mail..."
        itemLabel="eventos"
        headerActions={
          <Button
            type="button"
            onClick={() => setModalOpen(true)}
            size="sm"
            className="rounded-sm py-4 gap-2"
          >
            <Plus className="w-4 h-4" />
            Novo evento
          </Button>
        }
        emptyState={
          <EmptyState
            icon={Calendar}
            eyebrow="Eventos"
            title="Nenhum evento encontrado"
            description="Crie um evento para começar a organizar o dashboard do admin."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
            action={
              <Button
                type="button"
                onClick={() => setModalOpen(true)}
                size="sm"
                className="mt-0.5 h-9 rounded-sm gap-2 px-4 text-sm shadow-sm"
              >
                <Plus className="w-4 h-4" />
                Criar evento
              </Button>
            }
          />
        }
        renderItems={(slice) =>
          slice.map((event, idx) => (
            <AdminEventCard
              key={event.id}
              event={event}
              index={idx}
              onPublish={setPublishingEvent}
              onDiscontinue={setDiscontinuingEvent}
            />
          ))
        }
      />

      <ManageEventModal
        open={modalOpen}
        onOpenChange={setModalOpen}
        onCreate={(values) =>
          createMutation.mutateAsync(values).then(
            (res) => (res.success ? res.data : false),
            () => false,
          )
        }
      />

      <AlertModal
        open={!!publishingEvent}
        onOpenChange={() => setPublishingEvent(null)}
        title="Publicar evento?"
        description={`Ao publicar "${publishingEvent?.full_name}", ele ficará visível para o público.`}
        confirmLabel="Publicar"
        onConfirm={async () => {
          if (!publishingEvent) return;
          const response = await publishEventMutation.mutateAsync(
            publishingEvent.id,
          );
          if (response.success) setPublishingEvent(null);
        }}
        variant="success"
        loading={publishEventMutation.isPending}
      />

      <AlertModal
        open={!!discontinuingEvent}
        onOpenChange={() => setDiscontinuingEvent(null)}
        title="Descontinuar evento?"
        description={`Ao descontinuar "${discontinuingEvent?.full_name}", ele deixará de ser um evento ativo.`}
        confirmLabel="Descontinuar"
        onConfirm={() => {
          if (!discontinuingEvent) return;
          discontinueEventMutation.mutate(discontinuingEvent.id);
        }}
        variant="destructive"
        loading={discontinueEventMutation.isPending}
      />
    </div>
  );
}
