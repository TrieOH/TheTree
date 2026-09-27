import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { Button, EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal } from "solid-js";

import UserPlusIcon from "~icons/lucide/user-plus";
import UsersIcon from "~icons/lucide/users";

import {
  actorEmailsQueryOptions,
  allEventMembersQueryOptions,
} from "@/features/events/api/members";
import {
  useAddEventMemberMutation,
  useRemoveEventMemberMutation,
} from "@/features/events/api/mutations";
import type { EventMemberRole, EventMemberWithEmailI } from "@/features/events/model/member";
import { AdminAddMemberCard } from "@/features/events/ui/AdminAddMemberCard";
import { AdminEventMemberCard } from "@/features/events/ui/AdminEventMemberCard";
import {
  ManageEventMemberDialog,
  type ManageEventMemberValues,
} from "@/features/events/ui/ManageEventMemberDialog";
import { RemoveEventMemberDialog } from "@/features/events/ui/RemoveEventMemberDialog";
import { toast } from "@/shared/ui/toast";

const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;
const UserPlus = UserPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute("/admin/events/$eventId/members/")({
  head: () => ({
    meta: [{ title: "Membros - Admin Univents" }],
  }),
  component: AdminEventMembersRoute,
});

const ROLE_LABELS: Record<EventMemberRole, string> = {
  owner: "Proprietário",
  admin: "Administrador",
  staff: "Equipe",
};

const ROLE_SORT_ORDER: Record<EventMemberRole, number> = {
  owner: 0,
  admin: 1,
  staff: 2,
};

function AdminEventMembersRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;

  const membersQuery = useQuery(() => allEventMembersQueryOptions(eventId()));
  const addMutation = useAddEventMemberMutation();
  const removeMutation = useRemoveEventMemberMutation();

  const [addModalOpen, setAddModalOpen] = createSignal(false);
  const [memberToRemove, setMemberToRemove] = createSignal<EventMemberWithEmailI | null>(null);

  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<EventMemberWithEmailI>>({
    field: "created_at",
    direction: "desc",
  });

  const rawMembers = () => (membersQuery().data ?? []) as EventMemberWithEmailI[];
  const memberUserIds = createMemo(() => rawMembers().map((m) => m.user_id));
  const actorEmailsQuery = useQuery(() => actorEmailsQueryOptions(memberUserIds()));

  const members = createMemo(() => {
    const list = rawMembers();
    const details = actorEmailsQuery().data ?? {};
    return list.map((member) => ({
      ...member,
      email: member.email || details[member.user_id]?.email || undefined,
      pfp_url: member.pfp_url ?? details[member.user_id]?.pfp_url ?? null,
    }));
  });

  const handleAddMember = async (values: ManageEventMemberValues): Promise<boolean> => {
    try {
      await addMutation.mutateAsync({
        eventId: eventId(),
        email: values.email,
        role: values.role,
      });
      toast.success("Membro adicionado com sucesso!");
      return true;
    } catch {
      toast.error("Erro ao adicionar membro.");
      return false;
    }
  };

  const handleRemoveMember = async (userId: string, email: string): Promise<boolean> => {
    try {
      await removeMutation.mutateAsync({
        eventId: eventId(),
        userId,
        email,
      });
      toast.success("Membro removido com sucesso!");
      return true;
    } catch {
      toast.error("Erro ao remover membro.");
      return false;
    }
  };

  return (
    <div class="space-y-6">
      <Show
        when={!membersQuery().isLoading}
        fallback={
          <div class="grid grid-cols-[repeat(auto-fill,minmax(18rem,1fr))] gap-2">
            <For each={[1, 2, 3, 4]}>
              {() => <div class="h-16 animate-pulse rounded-xl bg-muted" />}
            </For>
          </div>
        }
      >
        <AdminEventMembersContent
          members={members()}
          filter={filter()}
          onFilterChange={setFilter}
          sort={sort()}
          onSortChange={setSort}
          onAdd={() => setAddModalOpen(true)}
          onRemove={setMemberToRemove}
        />
      </Show>

      <ManageEventMemberDialog
        open={addModalOpen()}
        onOpenChange={setAddModalOpen}
        onSubmit={handleAddMember}
      />

      <RemoveEventMemberDialog
        open={memberToRemove() !== null}
        onOpenChange={(open) => {
          if (!open) setMemberToRemove(null);
        }}
        member={memberToRemove()}
        onRemove={handleRemoveMember}
      />
    </div>
  );
}

function AdminEventMembersContent(props: {
  members: EventMemberWithEmailI[];
  filter: string;
  onFilterChange: (value: string) => void;
  sort: SortState<EventMemberWithEmailI>;
  onSortChange: (sort: SortState<EventMemberWithEmailI>) => void;
  onAdd: () => void;
  onRemove: (member: EventMemberWithEmailI) => void;
}): JSX.Element {
  const visible = createMemo(() => {
    const search = props.filter.trim().toLowerCase();
    if (!search) return props.members;

    return props.members.filter((member) =>
      [
        member.email,
        member.user_id,
        ROLE_LABELS[member.role],
        member.role,
      ].some((val) => (val ? val.toLowerCase().includes(search) : false)),
    );
  });

  const emptyStateAction = (
    <Button
      size="sm"
      class="gap-2 rounded-sm py-4"
      onClick={() => props.onAdd()}
    >
      <UserPlus class="size-4" />
      Adicionar membro
    </Button>
  );

  return (
    <PaginatedContainer<EventMemberWithEmailI>
      items={visible()}
      layout="grid"
      minItemWidth="16rem"
      maxRows={(columns) => (columns === 1 ? 8 : 4)}
      gap="2"
      sort={props.sort}
      onSortChange={props.onSortChange}
      sortFields={[
        {
          key: "created_at",
          label: "Adicionado em",
          ascLabel: "Mais antigos primeiro",
          descLabel: "Mais recentes primeiro",
          comparator: (a, b) =>
            new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
        },
        {
          key: "role",
          label: "Função",
          ascLabel: "Proprietário primeiro",
          descLabel: "Equipe primeiro",
          comparator: (a, b) => ROLE_SORT_ORDER[a.role] - ROLE_SORT_ORDER[b.role],
        },
        {
          key: "email",
          label: "E-mail",
          ascLabel: "A → Z",
          descLabel: "Z → A",
        },
      ]}
      filterValue={props.filter}
      onFilterChange={props.onFilterChange}
      filterPlaceholder="Buscar por e-mail, id ou função..."
      itemLabel="membros"
      emptyState={
        <EmptyState
          class="border-0 bg-transparent px-0 py-4 shadow-none"
          icon={<Users class="size-6 text-foreground/70" />}
          eyebrow="Equipe do evento"
          title="Nenhum membro encontrado"
          description={
            props.filter
              ? "Nenhum membro corresponde à busca informada."
              : "Adicione pessoas para colaborar na gestão deste evento."
          }
          action={emptyStateAction}
        />
      }
      renderItems={(slice, options) => (
        <>
          <AdminAddMemberCard
            index={0}
            animate={options.animate}
            onAdd={props.onAdd}
          />

          <For each={slice}>
            {(member, index) => (
              <AdminEventMemberCard
                member={member}
                index={index() + 1}
                animate={options.animate}
                onRemove={props.onRemove}
              />
            )}
          </For>
        </>
      )}
    />
  );
}
