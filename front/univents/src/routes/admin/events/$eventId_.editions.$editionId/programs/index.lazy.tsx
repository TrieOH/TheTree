import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute, useNavigate } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { CalendarDays, CalendarRange, Plus } from "lucide-react";
import { useState } from "react";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import { useProgramMutation } from "@/features/programs/api/mutations";
import type { ProgramI } from "@/features/programs/model";
import { ManageProgramModal } from "@/features/programs/ui/ManageProgramModal";
import { ProgramAdminCard } from "@/features/programs/ui/ProgramAdminCard";
import { Button } from "@/shared/ui/shadcn/button";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/programs/",
)({ component: ProgramsRoute });

function ProgramsRoute() {
  const { eventId, editionId } = Route.useParams();
  const navigate = useNavigate();
  const { data: programs = [] } = useQuery(programsQueryOptions(editionId));
  const { data: occurrences = [] } = useQuery(
    occurrencesQueryOptions(editionId),
  );
  const mutation = useProgramMutation(editionId);
  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState<SortState<ProgramI>>({
    field: "name",
    direction: "asc",
  });
  const [editing, setEditing] = useState<ProgramI>();
  const [modalOpen, setModalOpen] = useState(false);

  const filtered = programs
    .filter((program) =>
      program.name.toLowerCase().includes(filter.trim().toLowerCase()),
    )
    .sort((a, b) =>
      sort.direction === "asc"
        ? a.name.localeCompare(b.name)
        : b.name.localeCompare(a.name),
    );

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<ProgramI>
        items={filtered}
        layout="grid"
        minItemWidth="16rem"
        pageSize={9}
        gap="6"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          { key: "name", label: "Nome" },
          { key: "kind", label: "Tipo" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar programa..."
        itemLabel="programas"
        headerActions={
          <>
            <Button
              type="button"
              className="h-9 gap-2"
              onClick={() => {
                setEditing(undefined);
                setModalOpen(true);
              }}
            >
              <Plus className="mr-2 size-4" />
              Novo programa
            </Button>
            <Button
              type="button"
              variant="outline"
              className="h-9 gap-2"
              onClick={() =>
                navigate({
                  to: "/admin/events/$eventId/editions/$editionId/programs/calendar",
                  params: { eventId, editionId },
                })
              }
            >
              <CalendarDays className="size-4" />
              Abrir calendário
            </Button>
          </>
        }
        emptyState={
          <EmptyState
            icon={CalendarRange}
            eyebrow="Programação"
            title="Nenhum programa encontrado"
            description="Crie o primeiro programa desta edição."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((program, index) => (
            <ProgramAdminCard
              key={program.id}
              program={program}
              index={index}
              occurrences={occurrences.filter(
                (item) => item.program_id === program.id,
              )}
              onEdit={() => {
                setEditing(program);
                setModalOpen(true);
              }}
              onManageOccurrences={() =>
                navigate({
                  to: "/admin/events/$eventId/editions/$editionId/programs/$programId/occurrences",
                  params: {
                    eventId: program.edition_id,
                    editionId,
                    programId: program.id,
                  },
                })
              }
            />
          ))
        }
      />
      <ManageProgramModal
        open={modalOpen}
        program={editing}
        onOpenChange={setModalOpen}
        onSave={async (data) => {
          const response = await mutation.mutateAsync({
            id: editing?.id,
            data,
          });
          return response.success;
        }}
      />
    </div>
  );
}
