import { createFileRoute } from "@tanstack/react-router";
import { BadgeEditor } from "@/features/badges/editor/badge-editor";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/badges/editor",
)({ component: RouteComponent });

function RouteComponent() {
  const { eventId, editionId } = Route.useParams();
  return <BadgeEditor eventId={eventId} editionId={editionId} />;
}
