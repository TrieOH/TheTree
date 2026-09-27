import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";

import { CalendarEditor } from "@/features/calendar/ui/CalendarEditor";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/programs/calendar",
)({
  head: () => ({
    meta: [{ title: "Calendário da Programação - Admin Univents" }],
  }),
  component: AdminEditionCalendarRoute,
});

function AdminEditionCalendarRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;
  const editionId = () => params().editionId;

  return <CalendarEditor eventId={eventId()} editionId={editionId()} />;
}
