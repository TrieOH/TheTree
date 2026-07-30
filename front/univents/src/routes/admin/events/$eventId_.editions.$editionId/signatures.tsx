import { createFileRoute, Outlet } from "@tanstack/react-router";
import { Mail, PenLine } from "lucide-react";
import { SectionTabs } from "@/widgets/ui/section-tabs";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/signatures",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId, editionId } = Route.useParams();
  return (
    <div className="flex flex-wrap gap-6 p-6 pb-28!">
      <SectionTabs
        items={[
          {
            label: "Assinaturas",
            to: "/admin/events/$eventId/editions/$editionId/signatures/",
            icon: PenLine,
            params: { eventId, editionId },
          },
          {
            label: "Convites",
            to: "/admin/events/$eventId/editions/$editionId/signatures/invites",
            icon: Mail,
            params: { eventId, editionId },
          },
        ]}
      />
      <Outlet />
    </div>
  );
}
