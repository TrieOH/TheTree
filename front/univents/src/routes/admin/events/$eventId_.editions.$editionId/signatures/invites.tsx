import { createFileRoute } from "@tanstack/react-router";
import { SignatureRequestInvites } from "@/features/signatures/ui/SignatureRequestInvites";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/signatures/invites",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { editionId } = Route.useParams();

  return <SignatureRequestInvites editionId={editionId} />;
}
