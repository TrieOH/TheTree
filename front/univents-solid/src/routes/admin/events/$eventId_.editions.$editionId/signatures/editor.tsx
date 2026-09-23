import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { SignatureEditor } from "@/features/signatures/ui";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/signatures/editor",
)({
  head: () => ({
    meta: [{ title: "Nova Assinatura - Admin Univents" }],
  }),
  component: AdminSignatureEditorRoute,
});

function AdminSignatureEditorRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;
  const editionId = () => params().editionId;

  return <SignatureEditor eventId={eventId()} editionId={editionId()} />;
}
