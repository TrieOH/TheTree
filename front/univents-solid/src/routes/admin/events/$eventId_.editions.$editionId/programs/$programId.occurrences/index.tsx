import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery, useQueryClient } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { For, createMemo, createSignal } from "solid-js";

import ArrowLeftIcon from "~icons/lucide/arrow-left";
import CalendarDaysIcon from "~icons/lucide/calendar-days";

import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import {
  useCreateOccurrenceMutation,
  useDeleteOccurrenceMutation,
  useUpdateOccurrenceMutation,
} from "@/features/programs/api/mutations";
import { programKeys } from "@/features/programs/api/query-keys";
import type {
  OccurrenceCreateOutput,
  OccurrenceI,
  ProgramI,
} from "@/features/programs/model";
import { ManageOccurrenceDialog } from "@/features/calendar/ui/ManageOccurrenceDialog";
import { AdminCreateOccurrenceCard } from "@/features/programs/ui/AdminCreateOccurrenceCard";
import { OccurrenceAdminCard } from "@/features/programs/ui/OccurrenceAdminCard";
import { OccurrenceAttendanceModal } from "@/features/programs/ui/OccurrenceAttendanceModal";
import { OccurrenceDrawModal } from "@/features/programs/ui/OccurrenceDrawModal";
import { toast } from "@/shared/ui/toast";
import { AlertModal } from "@/widgets/ui/AlertModal";

const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/programs/$programId/occurrences/",
)({
  head: () => ({
    meta: [{ title: "Ocorrências da Atividade - Univents Admin" }],
  }),
  component: AdminProgramOccurrencesRoute,
});

function AdminProgramOccurrencesRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;
  const editionId = () => params().editionId;
  const programId = () => params().programId;
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // Queries
  const programsQuery = useQuery(() => programsQueryOptions(editionId()));
  const occurrencesQuery = useQuery(() => occurrencesQueryOptions(editionId()));

  // Mutations
  const createMutation = useCreateOccurrenceMutation(editionId);
  const updateMutation = useUpdateOccurrenceMutation(editionId);
  const deleteMutation = useDeleteOccurrenceMutation(editionId);

  const programs = createMemo(() => (programsQuery().data ?? []) as ProgramI[]);
  const allOccurrences = createMemo(() => (occurrencesQuery().data ?? []) as OccurrenceI[]);

  const program = createMemo(() =>
    programs().find((p) => p.id === programId()),
  );

  const occurrences = createMemo(() =>
    allOccurrences().filter((occ) => occ.program_id === programId()),
  );

  // States
  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<OccurrenceI>>({
    field: "starts_at",
    direction: "asc",
  });

  // Modals state
  const [dialogOpen, setDialogOpen] = createSignal(false);
  const [selectedOccurrence, setSelectedOccurrence] = createSignal<OccurrenceI | null>(null);

  const [attendanceModalOpen, setAttendanceModalOpen] = createSignal(false);
  const [attendanceOccurrence, setAttendanceOccurrence] = createSignal<OccurrenceI | null>(null);

  const [drawModalOpen, setDrawModalOpen] = createSignal(false);
  const [drawOccurrence, setDrawOccurrence] = createSignal<OccurrenceI | null>(null);

  const [occurrenceToDelete, setOccurrenceToDelete] = createSignal<OccurrenceI | null>(null);
  const [deleting, setDeleting] = createSignal(false);

  const handleOpenAttendance = (occ: OccurrenceI) => {
    setAttendanceOccurrence(occ);
    setAttendanceModalOpen(true);
  };

  const handleOpenDraw = (occ: OccurrenceI) => {
    setDrawOccurrence(occ);
    setDrawModalOpen(true);
  };

  const handleEdit = (occ: OccurrenceI) => {
    setSelectedOccurrence(occ);
    setDialogOpen(true);
  };

  const handleNewOccurrence = () => {
    setSelectedOccurrence(null);
    setDialogOpen(true);
  };

  const handleSaveOccurrence = async (data: {
    programId: string;
    occurrenceData: OccurrenceCreateOutput;
    id?: string;
  }): Promise<boolean> => {
    try {
      if (data.id) {
        await updateMutation.mutateAsync({
          id: data.id,
          data: data.occurrenceData,
        });
        toast.success("Horário atualizado com sucesso!");
      } else {
        await createMutation.mutateAsync({
          programId: data.programId,
          data: data.occurrenceData,
        });
        toast.success("Horário criado com sucesso!");
      }
      setDialogOpen(false);
      return true;
    } catch {
      toast.error("Erro ao salvar horário.");
      return false;
    }
  };

  const handleDeleteOccurrence = async (id: string): Promise<boolean> => {
    try {
      await deleteMutation.mutateAsync(id);
      toast.success("Horário excluído com sucesso!");
      setDialogOpen(false);
      return true;
    } catch {
      toast.error("Erro ao excluir horário.");
      return false;
    }
  };

  const confirmDelete = async () => {
    const occ = occurrenceToDelete();
    if (!occ) return;
    setDeleting(true);
    try {
      await deleteMutation.mutateAsync(occ.id);
      toast.success("Horário excluído!");
      setOccurrenceToDelete(null);
      void queryClient.invalidateQueries({
        queryKey: programKeys.occurrences(editionId()),
      });
    } catch {
      toast.error("Não foi possível excluir o horário.");
    } finally {
      setDeleting(false);
    }
  };

  return (
    <div class="space-y-6">
      {/* Header with back button */}
      <div class="flex items-center gap-3">
        <button
          type="button"
          onClick={() =>
            navigate({
              to: "/admin/events/$eventId/editions/$editionId/programs",
              params: { eventId: eventId(), editionId: editionId() },
            })
          }
          class="inline-flex size-9 items-center justify-center rounded-full border border-border/60 transition-colors hover:bg-accent cursor-pointer"
          title="Voltar para a programação"
        >
          <ArrowLeft class="size-4" />
        </button>
        <div>
          <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-foreground">
            Ocorrências
          </h1>
          <p class="text-xs sm:text-sm text-muted-foreground">
            {program()?.name ?? "Programa"} · Horários e sessões cadastradas
          </p>
        </div>
      </div>

      <PaginatedContainer<OccurrenceI>
        items={occurrences()}
        layout="grid"
        minItemWidth="16rem"
        maxRows={(columns) => (columns === 1 ? 8 : 4)}
        gap="2"
        sort={sort()}
        onSortChange={setSort}
        sortFields={[
          {
            key: "starts_at",
            label: "Horário de início",
            ascLabel: "Mais cedo primeiro",
            descLabel: "Mais tarde primeiro",
          },
        ]}
        filterValue={filter()}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por data ou horário..."
        filterFields={["starts_at", "ends_at"]}
        itemLabel="ocorrências"
        emptyState={
          <EmptyState
            icon={<CalendarDays class="size-6 text-foreground/70" />}
            eyebrow="Ocorrências"
            title="Nenhuma ocorrência encontrada"
            description="Cadastre a primeira ocorrência desta atividade para permitir inscrições, presença e sorteio."
            class="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice, options) => (
          <>
            <AdminCreateOccurrenceCard
              index={0}
              animate={options.animate}
              onCreate={handleNewOccurrence}
            />
            <For each={slice}>
              {(occurrence, index) => (
                <OccurrenceAdminCard
                  occurrence={occurrence}
                  index={index() + 1}
                  animate={options.animate}
                  programKind={program()?.kind}
                  onAttendance={handleOpenAttendance}
                  onDraw={program()?.kind === "activity" ? handleOpenDraw : undefined}
                  onEdit={handleEdit}
                  onDelete={(occ) => setOccurrenceToDelete(occ)}
                />
              )}
            </For>
          </>
        )}
      />

      {/* Create / Edit Occurrence Dialog */}
      <ManageOccurrenceDialog
        open={dialogOpen()}
        onOpenChange={setDialogOpen}
        programs={program() ? [program()!] : []}
        initialProgramId={programId()}
        occurrence={selectedOccurrence()}
        onSave={handleSaveOccurrence}
        onDelete={handleDeleteOccurrence}
        onOpenAttendance={handleOpenAttendance}
        onOpenDraw={handleOpenDraw}
      />

      {/* Attendance Modal (Works on mobile & desktop, with camera QR scanner) */}
      <OccurrenceAttendanceModal
        open={attendanceModalOpen()}
        onOpenChange={setAttendanceModalOpen}
        occurrence={attendanceOccurrence()}
        programKind={program()?.kind}
        programName={program()?.name}
      />

      {/* Draw Modal with Projector Mode (Works on mobile & desktop) */}
      <OccurrenceDrawModal
        open={drawModalOpen()}
        onOpenChange={setDrawModalOpen}
        occurrence={drawOccurrence()}
        programName={program()?.name}
        programKind={program()?.kind}
      />

      {/* Confirm Delete Modal */}
      <AlertModal
        open={occurrenceToDelete() !== null}
        onOpenChange={(open) => !open && setOccurrenceToDelete(null)}
        title="Excluir ocorrência"
        description="Tem certeza que deseja excluir este horário? Participantes inscritos perderão o agendamento."
        confirmLabel="Excluir horário"
        variant="destructive"
        loading={deleting()}
        onConfirm={confirmDelete}
      />
    </div>
  );
}
