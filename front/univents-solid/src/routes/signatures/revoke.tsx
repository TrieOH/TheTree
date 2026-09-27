import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { signatureTokenSearchSchema } from "@/features/signatures/model";
import { RevokeSignatureView } from "@/features/signatures/ui/RevokeSignatureView";

export const Route = createFileRoute("/signatures/revoke")({
  validateSearch: (search) => signatureTokenSearchSchema.parse(search),
  head: () => ({
    meta: [{ title: "Revogar Assinatura - Univents" }],
  }),
  component: RevokeSignatureRoute,
});

function RevokeSignatureRoute(): JSX.Element {
  const search = Route.useSearch();
  return <RevokeSignatureView token={search().token} />;
}
