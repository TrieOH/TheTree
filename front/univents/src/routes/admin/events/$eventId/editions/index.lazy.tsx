import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { Calendar, Plus } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import {
  useCreateEditionMutation,
  usePatchEditionMutation,
  usePublishEditionMutation,
} from "@/features/editions/api/mutations";
import type { EditionI } from "@/features/editions/model";
import { editionRangesOverlap } from "@/features/editions/model";
import { AdminEditionCard } from "@/features/editions/ui/AdminEditionCard";
import { EditEditionModal } from "@/features/editions/ui/EditEditionModal";
import { ManageEditionModal } from "@/features/editions/ui/ManageEditionModal";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";

const STATUS_SORT_ORDER: Record<EditionI["status"], number> = {
  active: 0,
  future: 1,
  draft: 2,
  past: 3,
};

export const Route = createLazyFileRoute("/admin/events/$eventId/editions/")({
  component: EditionsRoute,
});

function EditionsRoute() {
  const { eventId } = Route.useParams();
  const { data: editions = [] } = useQuery(
    allAdminEditionsQueryOptions(eventId),
  );
  const createEditionMutation = useCreateEditionMutation();
  const patchEditionMutation = usePatchEditionMutation();
  const publishEditionMutation = usePublishEditionMutation();
  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState<SortState<EditionI>>({
    field: "starts_at",
    direction: "desc",
  });
  const [modalOpen, setModalOpen] = useState(false);
  const [editionToEdit, setEditionToEdit] = useState<EditionI | null>(null);
  const [editionToPublish, setEditionToPublish] = useState<EditionI | null>(
    null,
  );

  const filteredEditions = [...editions]
    .filter((edition) => {
      const search = filter.trim().toLowerCase();
      if (!search) return true;

      return [
        edition.name,
        edition.slug,
        edition.tagline ?? "",
        edition.location_name ?? "",
        edition.status,
      ].some((value) => value.toLowerCase().includes(search));
    })
    .sort((a, b) => {
      const direction = sort.direction === "asc" ? 1 : -1;

      if (sort.field === "starts_at") {
        return (
          (new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime()) *
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
        String(a[sort.field] ?? "").localeCompare(String(b[sort.field] ?? "")) *
        direction
      );
    });

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<EditionI>
        items={filteredEditions}
        layout="grid"
        minItemWidth="16rem"
        pageSize={4}
        gap="6"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          {
            key: "starts_at",
            label: "Início",
            comparator: (a, b) =>
              new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime(),
          },
          { key: "name", label: "Nome" },
          {
            key: "status",
            label: "Status",
            comparator: (a, b) =>
              STATUS_SORT_ORDER[a.status] - STATUS_SORT_ORDER[b.status],
          },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome, local ou status..."
        itemLabel="edições"
        headerActions={
          <Button
            type="button"
            onClick={() => setModalOpen(true)}
            className="h-9 gap-2"
          >
            <Plus className="size-4" />
            Nova edição
          </Button>
        }
        emptyState={
          <EmptyState
            icon={Calendar}
            eyebrow="Edições"
            title="Nenhuma edição encontrada"
            description="Crie a primeira edição para começar a organizar esse evento."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((edition, idx) => (
            <AdminEditionCard
              key={edition.id}
              edition={edition}
              eventId={eventId}
              index={idx}
              onEdit={setEditionToEdit}
              onPublish={edition.is_draft ? setEditionToPublish : undefined}
            />
          ))
        }
      />

      <ManageEditionModal
        open={modalOpen}
        onOpenChange={setModalOpen}
        onCreate={async (values) => {
          const overlaps = editions.some((edition) =>
            editionRangesOverlap(edition, values),
          );

          if (overlaps) {
            toast.error(
              "Já existe uma edição nesse período. Apenas uma edição pode estar ativa por vez.",
            );
            return false;
          }

          const edition = await createEditionMutation.mutateAsync({
            eventId,
            data: values,
          });
          return edition ? edition : false;
        }}
      />

      {editionToEdit ? (
        <EditEditionModal
          key={editionToEdit.id}
          open
          edition={editionToEdit}
          onOpenChange={(open) => {
            if (!open) setEditionToEdit(null);
          }}
          onUpdate={async (values) => {
            const overlaps = editions.some(
              (edition) =>
                edition.id !== editionToEdit.id &&
                editionRangesOverlap(edition, values),
            );
            if (overlaps) {
              toast.error(
                "Já existe uma edição nesse período. Apenas uma edição pode estar ativa por vez.",
              );
              return false;
            }

            const edition = await patchEditionMutation.mutateAsync({
              eventId,
              editionId: editionToEdit.id,
              data: values,
            });
            return edition ? edition : false;
          }}
        />
      ) : null}

      <AlertModal
        open={Boolean(editionToPublish)}
        onOpenChange={(open) => {
          if (!open && !publishEditionMutation.isPending) {
            setEditionToPublish(null);
          }
        }}
        title="Publicar edição?"
        description={
          editionToPublish
            ? `A edição "${editionToPublish.name}" ficará disponível publicamente.`
            : undefined
        }
        confirmLabel="Publicar edição"
        loading={publishEditionMutation.isPending}
        onConfirm={async () => {
          if (!editionToPublish) return;
          await publishEditionMutation.mutateAsync({
            eventId,
            editionId: editionToPublish.id,
          });
          setEditionToPublish(null);
        }}
      />
    </div>
  );
}
