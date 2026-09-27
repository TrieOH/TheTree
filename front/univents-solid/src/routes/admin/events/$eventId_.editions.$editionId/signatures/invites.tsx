import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { createEffect } from "solid-js";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/signatures/invites",
)({
  head: () => ({
    meta: [{ title: "Convites de Assinatura - Admin Univents" }],
  }),
  component: AdminSignatureInvitesRoute,
});

function AdminSignatureInvitesRoute(): JSX.Element {
  const params = Route.useParams();
  const navigate = useNavigate();

  createEffect(
    () => [params().eventId, params().editionId] as const,
    ([eventId, editionId]) => {
      if (!eventId || !editionId) return;
      navigate({
        to: "/admin/events/$eventId/editions/$editionId/signatures",
        params: { eventId, editionId },
      });
    },
  );

  return <div class="p-4 text-xs text-muted-foreground">Redirecionando para convites...</div>;
}
