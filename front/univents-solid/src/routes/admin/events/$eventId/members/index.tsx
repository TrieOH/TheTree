import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import UsersIcon from "~icons/lucide/users";
import { EmptyState } from "@trieoh/ui-solid";

const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute("/admin/events/$eventId/members/")({
  head: () => ({
    meta: [{ title: "Membros - Admin Univents" }],
  }),
  component: AdminEventMembersRoute,
});

function AdminEventMembersRoute(): JSX.Element {
  return (
    <div class="space-y-6">
      <div class="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-foreground sm:text-3xl">
            Membros
          </h1>
          <p class="text-sm text-muted-foreground">
            Gerencie a equipe, permissões e administradores com acesso a este evento.
          </p>
        </div>
      </div>

      <EmptyState
        icon={<Users class="h-6 w-6 text-foreground/70" />}
        eyebrow="Equipe do Evento"
        title="Membros do evento"
        description="Gerencie os membros com permissão de proprietário, administrador ou equipe neste evento."
      />
    </div>
  );
}
