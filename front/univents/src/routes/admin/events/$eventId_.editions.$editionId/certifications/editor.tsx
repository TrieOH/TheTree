import { createFileRoute } from "@tanstack/react-router";
import { CertificateEditor } from "@/features/certifications/editor/ui/certificate-editor";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/certifications/editor",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId, editionId } = Route.useParams();
  return <CertificateEditor eventId={eventId} editionId={editionId} />;
}
