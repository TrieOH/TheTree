import { createFileRoute } from "@tanstack/react-router";
import { CertificateEditor } from "@/features/certifications/editor/ui/certificate-editor";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/certifications/editor",
)({
  validateSearch: (search: Record<string, unknown>) => ({
    templateId: typeof search.templateId === "string" ? search.templateId : "",
    duplicate: search.duplicate === true || search.duplicate === "true",
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId, editionId } = Route.useParams();
  const { templateId, duplicate } = Route.useSearch();
  return (
    <CertificateEditor
      eventId={eventId}
      editionId={editionId}
      templateId={templateId || undefined}
      duplicate={duplicate}
    />
  );
}
