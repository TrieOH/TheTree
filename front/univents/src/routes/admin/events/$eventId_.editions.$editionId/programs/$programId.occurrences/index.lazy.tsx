import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute, useRouter } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { ArrowLeft, CalendarDays, Plus } from "lucide-react";
import { useMemo, useState } from "react";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import {
  useDeleteOccurrenceMutation,
  useOccurrenceMutation,
} from "@/features/programs/api/mutations";
import type { OccurrenceI } from "@/features/programs/model";
import { ManageOccurrenceModal } from "@/features/programs/ui/ManageOccurrenceModal";
import { OccurrenceAdminCard } from "@/features/programs/ui/OccurrenceAdminCard";
import { OccurrenceAttendanceDialog } from "@/features/programs/ui/OccurrenceAttendanceDialog";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/programs/$programId/occurrences/",
)({ component: OccurrencesRoute });

function OccurrencesRoute() {
  const { eventId, editionId, programId } = Route.useParams();
  const router = useRouter();
  const { data: programs = [] } = useQuery(programsQueryOptions(editionId));
  const { data: allOccurrences = [] } = useQuery(
    occurrencesQueryOptions(editionId),
  );
  const mutation = useOccurrenceMutation(editionId);
  const deleteMutation = useDeleteOccurrenceMutation(editionId);
  const program = programs.find((item) => item.id === programId);
  const occurrences = useMemo(
    () => allOccurrences.filter((item) => item.program_id === programId),
    [allOccurrences, programId],
  );
  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState<SortState<OccurrenceI>>({
    field: "starts_at",
    direction: "asc",
  });
  const [modal, setModal] = useState<{
    open: boolean;
    occurrence?: OccurrenceI;
  }>({ open: false });
  const [occurrenceToDelete, setOccurrenceToDelete] =
    useState<OccurrenceI | null>(null);
  const [attendanceOccurrence, setAttendanceOccurrence] =
    useState<OccurrenceI | null>(null);
  const filtered = occurrences.filter(
    (item) =>
      !filter ||
      new Date(item.starts_at).toLocaleString("pt-BR").includes(filter),
  );

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <div className="mb-4 flex w-full items-center gap-3">
        <button
          type="button"
          onClick={() => router.history.back()}
          className="inline-flex size-9 items-center justify-center rounded-full border border-border/60 transition-colors hover:bg-accent"
        >
          <ArrowLeft className="size-4" />
        </button>
        <div>
          <h2 className="text-lg font-semibold">Ocorrências</h2>
          <p className="text-sm text-muted-foreground">
            {program?.name ?? "Programa"}
          </p>
        </div>
      </div>
      <PaginatedContainer<OccurrenceI>
        items={filtered}
        layout="list"
        pageSize={10}
        gap="3"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          {
            key: "starts_at",
            label: "Início",
            comparator: (a, b) =>
              new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime(),
          },
          { key: "ends_at", label: "Término" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por data ou horário..."
        itemLabel="ocorrências"
        headerActions={
          <Button
            type="button"
            className="h-9 gap-2"
            onClick={() => setModal({ open: true, occurrence: undefined })}
          >
            <Plus className="size-4" />
            Nova ocorrência
          </Button>
        }
        emptyState={
          <EmptyState
            icon={CalendarDays}
            eyebrow="Ocorrências"
            title="Nenhuma ocorrência encontrada"
            description="Crie a primeira ocorrência deste programa."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((occurrence) => (
            <OccurrenceAdminCard
              key={occurrence.id}
              occurrence={occurrence}
              onEdit={() => setModal({ open: true, occurrence })}
              onAttendance={() => setAttendanceOccurrence(occurrence)}
              onDraw={
                program?.kind === "activity"
                  ? () =>
                      router.navigate({
                        to: "/admin/events/$eventId/editions/$editionId/programs/$programId/occurrences/$occurrenceId/draw",
                        params: {
                          eventId,
                          editionId,
                          programId,
                          occurrenceId: occurrence.id,
                        },
                      })
                  : undefined
              }
              onDelete={() => {
                setOccurrenceToDelete(occurrence);
              }}
            />
          ))
        }
      />
      <AlertModal
        open={Boolean(occurrenceToDelete)}
        onOpenChange={(open) => !open && setOccurrenceToDelete(null)}
        title="Excluir ocorrência?"
        description="Esta ocorrência será removida permanentemente. Essa ação não pode ser desfeita."
        confirmLabel="Excluir ocorrência"
        variant="destructive"
        loading={deleteMutation.isPending}
        onConfirm={async () => {
          if (!occurrenceToDelete) return;
          await deleteMutation.mutateAsync(occurrenceToDelete.id);
          setOccurrenceToDelete(null);
        }}
      />
      {attendanceOccurrence ? (
        <OccurrenceAttendanceDialog
          occurrenceId={attendanceOccurrence.id}
          programKind={program?.kind ?? "activity"}
          open
          onOpenChange={(open) => !open && setAttendanceOccurrence(null)}
        />
      ) : null}
      <ManageOccurrenceModal
        open={modal.open}
        occurrence={modal.occurrence}
        onOpenChange={(open) =>
          setModal(
            open ? { ...modal, open } : { open: false, occurrence: undefined },
          )
        }
        onSave={async (data) => {
          const occurrence = await mutation.mutateAsync({
            id: modal.occurrence?.id,
            programId,
            data,
          });
          return Boolean(occurrence);
        }}
      />
    </div>
  );
}
