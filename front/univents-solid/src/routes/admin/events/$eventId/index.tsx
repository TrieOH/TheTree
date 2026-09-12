import type { JSX } from "@solidjs/web";
import { Link, createFileRoute } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core-solid";
import { EmptyState } from "@trieoh/ui-solid";
import { createMemo } from "solid-js";

import { allOwnEventsQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";

export const Route = createFileRoute("/admin/events/$eventId/")({
  head: ({ params }) => ({
    meta: [{ title: `Evento ${params.eventId} - Admin - Univents` }],
  }),
  component: AdminEventPanel,
});

function AdminEventPanel(): JSX.Element {
  const params = Route.useParams();
  const queryClient = useQueryClient();

  const event = createMemo(() => {
    const owned = queryClient.getQueryData<EventI[]>(allOwnEventsQueryOptions().queryKey);
    return owned?.find((candidate) => candidate.id === params().eventId) ?? null;
  });

  return (
    <div class="flex flex-col gap-6">
      <Link
        to="/admin/events"
        class="w-fit text-sm text-muted-foreground transition-colors hover:text-foreground"
      >
        ← Eventos
      </Link>

      <EmptyState
        eyebrow="Painel do evento"
        title={event()?.full_name ?? params().eventId}
        description="Edições, produtos, programação e equipe do evento."
      />
    </div>
  );
}
