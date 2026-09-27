import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { signatureTokenSearchSchema } from "@/features/signatures/model";
import { FulfillSignatureRequestView } from "@/features/signatures/ui/FulfillSignatureRequestView";

export const Route = createFileRoute("/signature-requests/fulfill")({
  validateSearch: (search) => signatureTokenSearchSchema.parse(search),
  head: () => ({
    meta: [{ title: "Completar Solicitação de Assinatura - Univents" }],
  }),
  component: FulfillSignatureRequestRoute,
});

function FulfillSignatureRequestRoute(): JSX.Element {
  const search = Route.useSearch();
  return <FulfillSignatureRequestView token={search().token} />;
}
