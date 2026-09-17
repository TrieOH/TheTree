import type { JSX } from "@solidjs/web";
import { useQuery } from "@trieoh/front-core-solid";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { createMemo, createSignal, For, Show } from "solid-js";

import CalendarDaysIcon from "~icons/lucide/calendar-days";
import CalendarRangeIcon from "~icons/lucide/calendar-range";

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
  ProgramCreateOutput,
  ProgramI,
} from "@/features/programs/model";
import { AdminCreateProgramCard } from "@/features/programs/ui/AdminCreateProgramCard";
import { AdminProgramCard } from "@/features/programs/ui/AdminProgramCard";
import { ManageProgramDialog } from "@/features/programs/ui/ManageProgramDialog";
import { Button, EmptyState, PaginatedContainer, type SortState } from "@trieoh/ui-solid";
import { toast } from "@/shared/ui/toast";

const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;
const CalendarRange = CalendarRangeIcon as unknown as (props: { class?: string }) => JSX.Element;

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

  // Filter & dialog states
  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<ProgramI>>({
    field: "name",
    direction: "asc",
  });
  const [creating, setCreating] = createSignal(false);
  const [editing, setEditing] = createSignal<ProgramI | null>(null);

  const programs = createMemo(() => (programsQuery().data ?? []) as ProgramI[]);
  const occurrences = createMemo(() => (occurrencesQuery().data ?? []) as OccurrenceI[]);

  const handleCreate = async (
    data: ProgramCreateInput | ProgramCreateOutput,
  ): Promise<boolean> => {
    try {
      await createMutation.mutateAsync(data);
      toast.success("Programa cadastrado com sucesso!");
      setCreating(false);
      return true;
    } catch {
      toast.error("Erro ao cadastrar programa.");
      return false;
    }
  };

  const handleUpdate = async (
    id: string,
    data: ProgramCreateInput | ProgramCreateOutput,
  ): Promise<boolean> => {
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

  const handleDelete = async (id: string): Promise<boolean> => {
    try {
      await deleteMutation.mutateAsync(id);
      toast.success("Programa removido com sucesso!");
      return true;
    } catch {
      toast.error("Erro ao excluir programa.");
      return false;
    }
  };

  return (
    <div class="space-y-6">
      {/* Header */}
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-foreground">
            Programação
          </h1>
          <p class="text-sm text-muted-foreground mt-0.5">
            Gerencie as atividades, palestras, workshops e checkpoints desta edição.
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
            class="gap-1.5 cursor-pointer font-medium"
          >
            <CalendarDays class="size-4" />
            <span>Abrir calendário</span>
          </Button>
        </div>
      </div>

      {/* Grid / Paginated List */}
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
        filterPlaceholder="Buscar por nome ou tipo de programa..."
        filterFields={["name", "description", "kind"]}
        itemLabel="programas"
        emptyState={
          <EmptyState
            icon={<CalendarRange class="size-6 text-foreground/70" />}
            eyebrow="Programação"
            title="Nenhum programa encontrado"
            description="Crie o primeiro programa desta edição ou abra o calendário para gerenciar horários."
            class="border-0 bg-transparent px-0 py-4 shadow-none"
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
                  onEdit={setEditing}
                  onDelete={(p) => handleDelete(p.id)}
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
        open={creating()}
        program={null}
        onOpenChange={(open) => !open && setCreating(false)}
        onSubmit={handleCreate}
      />

      <Show when={editing()}>
        {(prog) => (
          <ManageProgramDialog
            open={Boolean(prog())}
            program={prog()}
            onOpenChange={(open) => !open && setEditing(null)}
            onSubmit={(data) => handleUpdate(prog().id, data)}
          />
        )}
      </Show>
    </div>
  );
}
