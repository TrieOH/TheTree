import { createLazyFileRoute } from "@tanstack/react-router";
import { CalendarEditor } from "@/features/calendar/ui/CalendarEditor";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/programs/calendar",
)({ component: CalendarRoute });

function CalendarRoute() {
  const { eventId, editionId } = Route.useParams();
  return <CalendarEditor eventId={eventId} editionId={editionId} />;
}
