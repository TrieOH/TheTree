import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { Button, EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { For, createMemo, createSignal } from "solid-js";

import CalendarDaysIcon from "~icons/lucide/calendar-days";
import CalendarPlusIcon from "~icons/lucide/calendar-plus";

import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import {
  useCreateProgramMutation,
  useDeleteProgramMutation,
  useUpdateProgramMutation,
} from "@/features/programs/api/mutations";
import type {
  OccurrenceI,
  ProgramCreateInput,
  ProgramI,
} from "@/features/programs/model";
import { AdminCreateProgramCard } from "@/features/programs/ui/AdminCreateProgramCard";
import { AdminProgramCard } from "@/features/programs/ui/AdminProgramCard";
import { ManageProgramDialog } from "@/features/programs/ui/ManageProgramDialog";
import { toast } from "@/shared/ui/toast";
import { AlertModal } from "@/widgets/ui/AlertModal";

const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;
const CalendarPlus = CalendarPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/programs/",
)({
  head: () => ({
    meta: [{ title: "Programação - Univents Admin" }],
  }),
  component: AdminProgramsRoute,
});

function AdminProgramsRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;
  const editionId = () => params().editionId;
  const navigate = useNavigate();

  // Queries
  const programsQuery = useQuery(() => programsQueryOptions(editionId()));
  const occurrencesQuery = useQuery(() => occurrencesQueryOptions(editionId()));

  // Mutations
  const createMutation = useCreateProgramMutation(editionId);
  const updateMutation = useUpdateProgramMutation(editionId);
  const deleteMutation = useDeleteProgramMutation(editionId);

  const programs = createMemo(() => (programsQuery().data ?? []) as ProgramI[]);
  const occurrences = createMemo(() => (occurrencesQuery().data ?? []) as OccurrenceI[]);

  // States
  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<ProgramI>>({
    field: "name",
    direction: "asc",
  });

  const [creating, setCreating] = createSignal(false);
  const [editing, setEditing] = createSignal<ProgramI | null>(null);
  const [deleting, setDeleting] = createSignal<ProgramI | null>(null);

  const handleCreate = async (data: ProgramCreateInput) => {
    try {
      await createMutation.mutateAsync(data);
      toast.success("Programa criado com sucesso!");
      setCreating(false);
      return true;
    } catch {
      toast.error("Erro ao criar programa.");
      return false;
    }
  };

  const handleUpdate = async (id: string, data: ProgramCreateInput) => {
    try {
      await updateMutation.mutateAsync({ id, data });
      toast.success("Programa atualizado com sucesso!");
      setEditing(null);
      return true;
    } catch {
      toast.error("Erro ao atualizar programa.");
      return false;
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteMutation.mutateAsync(id);
      toast.success("Programa excluído com sucesso!");
      setDeleting(null);
      return true;
    } catch {
      toast.error("Erro ao excluir programa.");
      return false;
    }
  };

  const emptyStateAction = (
    <Button
      size="sm"
      class="gap-2 rounded-sm py-4"
      onClick={() => setCreating(true)}
    >
      <CalendarPlus class="size-4" />
      Novo programa
    </Button>
  );

  return (
    <div class="space-y-6">
      {/* Header with Navigation */}
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-foreground">
            Programação do Evento
          </h1>
          <p class="text-xs sm:text-sm text-muted-foreground">
            Gerencie as palestras, workshops, checkpoints e horários da edição.
          </p>
        </div>

        <div class="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() =>
              navigate({
                to: "/admin/events/$eventId/editions/$editionId/programs/calendar",
                params: { eventId: eventId(), editionId: editionId() },
              })
            }
            class="gap-2 cursor-pointer"
          >
            <CalendarDays class="size-4 text-primary" />
            <span>Abrir calendário</span>
          </Button>
        </div>
      </div>

      <PaginatedContainer<ProgramI>
        items={programs()}
        layout="grid"
        minItemWidth="16rem"
        maxRows={(columns) => (columns === 1 ? 8 : 4)}
        gap="2"
        sort={sort()}
        onSortChange={setSort}
        sortFields={[
          {
            key: "name",
            label: "Nome",
            ascLabel: "A → Z",
            descLabel: "Z → A",
          },
          {
            key: "kind",
            label: "Tipo",
            ascLabel: "Atividade primeiro",
            descLabel: "Checkpoint primeiro",
          },
        ]}
        filterValue={filter()}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome ou descrição..."
        filterFields={["name", "description"]}
        itemLabel="programas"
        emptyState={
          <EmptyState
            icon={<CalendarDays class="size-6 text-foreground/70" />}
            eyebrow="Programação"
            title="Nenhum programa cadastrado"
            description={
              filter()
                ? "Nenhum programa corresponde à busca informada."
                : "Crie o primeiro programa desta edição para definir a grade de horários."
            }
            class="border-0 bg-transparent px-0 py-4 shadow-none"
            action={emptyStateAction}
          />
        }
        renderItems={(slice, options) => (
          <>
            <AdminCreateProgramCard
              index={0}
              animate={options.animate}
              onCreate={() => setCreating(true)}
            />

            <For each={slice}>
              {(program, index) => (
                <AdminProgramCard
                  program={program}
                  index={index() + 1}
                  animate={options.animate}
                  occurrences={occurrences().filter((o) => o.program_id === program.id)}
                  occurrencesHref={`/admin/events/${eventId()}/editions/${editionId()}/programs/${program.id}/occurrences`}
                  onEdit={setEditing}
                  onDelete={(p) => setDeleting(p)}
                  onManageOccurrences={(p) =>
                    navigate({
                      to: "/admin/events/$eventId/editions/$editionId/programs/$programId/occurrences",
                      params: {
                        eventId: eventId(),
                        editionId: editionId(),
                        programId: p.id,
                      },
                    })
                  }
                  onOpenCalendar={(_p) =>
                    navigate({
                      to: "/admin/events/$eventId/editions/$editionId/programs/calendar",
                      params: { eventId: eventId(), editionId: editionId() },
                    })
                  }
                />
              )}
            </For>
          </>
        )}
      />

      {/* Dialog: Create or Edit */}
      <ManageProgramDialog
        open={creating() || editing() !== null}
        program={editing()}
        onOpenChange={(open) => {
          if (!open) {
            setCreating(false);
            setEditing(null);
          }
        }}
        onSubmit={(data: ProgramCreateInput) => {
          const current = editing();
          if (current) {
            return handleUpdate(current.id, data);
          }
          return handleCreate(data);
        }}
      />

      {/* Modal: Delete */}
      <AlertModal
        open={deleting() !== null}
        onOpenChange={(open) => {
          if (!open) setDeleting(null);
        }}
        title="Excluir programa"
        description={
          deleting()
            ? `Tem certeza que deseja excluir "${deleting()?.name}"? Esta ação removerá os horários associados.`
            : undefined
        }
        confirmLabel="Excluir programa"
        variant="destructive"
        loading={deleteMutation.result().isPending}
        onConfirm={() => {
          const p = deleting();
          if (p) void handleDelete(p.id);
        }}
      />
    </div>
  );
}
