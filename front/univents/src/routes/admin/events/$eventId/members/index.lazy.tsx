import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { UserPlus, Users } from "lucide-react";
import { useState } from "react";
import {
  allEventMembersQueryOptions,
  type EventMemberI,
} from "@/features/events/api/members";
import {
  useAddEventMemberMutation,
  useRemoveEventMemberMutation,
} from "@/features/events/api/mutations";
import type {
  EventMemberCreateOutput,
  EventMemberRole,
} from "@/features/events/model/member";
import { AdminEventMemberCard } from "@/features/events/ui/AdminEventMemberCard";
import { ManageEventMemberModal } from "@/features/events/ui/ManageEventMemberModal";
import { RemoveEventMemberModal } from "@/features/events/ui/RemoveEventMemberModal";
import { Button } from "@/shared/ui/shadcn/button";

export const Route = createLazyFileRoute("/admin/events/$eventId/members/")({
  component: EventMembersRoute,
});

const roleLabels: Record<EventMemberRole, string> = {
  owner: "Proprietário",
  admin: "Administrador",
  staff: "Equipe",
};

const roleSortOrder: Record<EventMemberRole, number> = {
  owner: 0,
  admin: 1,
  staff: 2,
};

function EventMembersRoute() {
  const { eventId } = Route.useParams();
  const { data: members = [] } = useQuery(allEventMembersQueryOptions(eventId));
  const addMutation = useAddEventMemberMutation();
  const removeMutation = useRemoveEventMemberMutation();

  const [addModalOpen, setAddModalOpen] = useState(false);
  const [memberToRemove, setMemberToRemove] = useState<EventMemberI | null>(
    null,
  );
  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState<SortState<EventMemberI>>({
    field: "created_at",
    direction: "desc",
  });

  const visibleMembers = [...members]
    .filter((member) => {
      const search = filter.trim().toLowerCase();
      if (!search) return true;

      return [member.user_id, roleLabels[member.role], member.role].some(
        (value) => value.toLowerCase().includes(search),
      );
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

      if (sort.field === "role") {
        return (roleSortOrder[a.role] - roleSortOrder[b.role]) * direction;
      }

      return (
        String(a[sort.field] ?? "").localeCompare(String(b[sort.field] ?? "")) *
        direction
      );
    });

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<EventMemberI>
        items={visibleMembers}
        layout="grid"
        minItemWidth="16rem"
        pageSize={8}
        gap="4"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          {
            key: "created_at",
            label: "Adicionado em",
            comparator: (a, b) =>
              new Date(a.created_at).getTime() -
              new Date(b.created_at).getTime(),
          },
          {
            key: "role",
            label: "Função",
            comparator: (a, b) => roleSortOrder[a.role] - roleSortOrder[b.role],
          },
          { key: "user_id", label: "Usuário" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por usuário ou função..."
        itemLabel="membros"
        headerActions={
          <Button
            type="button"
            size="sm"
            className="h-9 gap-2"
            onClick={() => setAddModalOpen(true)}
          >
            <UserPlus className="size-4" />
            Adicionar membro
          </Button>
        }
        emptyState={
          <EmptyState
            icon={Users}
            eyebrow="Equipe do evento"
            title="Nenhum membro encontrado"
            description={
              filter
                ? "Nenhum membro corresponde à busca informada."
                : "Adicione pessoas para colaborar na gestão deste evento."
            }
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((member, index) => (
            <AdminEventMemberCard
              key={member.id}
              member={member}
              index={index}
              onRemove={setMemberToRemove}
            />
          ))
        }
      />

      <ManageEventMemberModal
        open={addModalOpen}
        onOpenChange={setAddModalOpen}
        onCreate={(values: EventMemberCreateOutput) =>
          addMutation
            .mutateAsync({
              eventId,
              email: values.email.trim().toLowerCase(),
              role: values.role,
            })
            .then(
              (res) => res.success,
              () => false,
            )
        }
      />

      <RemoveEventMemberModal
        open={memberToRemove !== null}
        onOpenChange={(open) => {
          if (!open) setMemberToRemove(null);
        }}
        member={memberToRemove}
        onRemove={(userId, email) =>
          removeMutation.mutateAsync({ eventId, userId, email }).then(
            (res) => res.success,
            () => false,
          )
        }
      />
    </div>
  );
}
