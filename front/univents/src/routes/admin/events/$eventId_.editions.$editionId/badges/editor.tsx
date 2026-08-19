import { createFileRoute } from "@tanstack/react-router";
import { BadgeEditor } from "@/features/badges/editor/badge-editor";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/badges/editor",
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
    <BadgeEditor
      eventId={eventId}
      editionId={editionId}
      templateId={templateId || undefined}
      duplicate={duplicate}
    />
  );
}
