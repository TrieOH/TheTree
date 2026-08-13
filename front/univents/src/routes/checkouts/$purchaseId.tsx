import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { ArrowLeft, Clock3 } from "lucide-react";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { checkoutQueryOptions } from "@/features/purchases/api";
import { usePurchaseSocket } from "@/features/purchases/hooks/use-purchase-socket";
import { PurchaseSummary } from "@/features/purchases/ui/purchase-summary";

export const Route = createFileRoute("/checkouts/$purchaseId")({
  beforeLoad: requireAuth,
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(
      checkoutQueryOptions(params.purchaseId),
    ),
  component: CheckoutStatusPage,
});

function CheckoutStatusPage() {
  const { purchaseId } = Route.useParams();
  const { data: checkout } = useSuspenseQuery(checkoutQueryOptions(purchaseId));
  const connected = usePurchaseSocket(
    purchaseId,
    checkout.status === "pending",
  );

  return (
    <main className="mx-auto w-full max-w-2xl px-4 py-10 pb-28">
      <Link
        to="/profile/purchases"
        className="mb-6 inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="size-4" />
        Minhas compras
      </Link>

      <h1 className="mb-6 text-2xl font-bold">Status da compra</h1>

      {checkout.status === "pending" && (
        <div className="mb-4 flex items-center gap-3 border border-amber-500/30 bg-amber-500/10 p-4 text-sm text-amber-800 dark:text-amber-300">
          <Clock3 className="size-5 shrink-0" />
          <p>
            Pagamento em processamento. Atualização em tempo real
            {connected ? " conectada." : " conectando…"}
          </p>
        </div>
      )}

      <PurchaseSummary purchase={checkout} />

      {checkout.intent_status && (
        <p className="mt-3 text-xs text-muted-foreground">
          Estado do pagamento: {checkout.intent_status}
        </p>
      )}
      {checkout.status_reason && (
        <p className="mt-1 text-xs text-muted-foreground">
          Motivo: {checkout.status_reason}
        </p>
      )}
    </main>
  );
}
