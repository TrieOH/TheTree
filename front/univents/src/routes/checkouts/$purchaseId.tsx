import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";
import { toast } from "sonner";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { checkoutQueryOptions } from "@/features/purchases/api";
import { usePurchaseSocket } from "@/features/purchases/hooks/use-purchase-socket";
import CheckoutPage from "@/features/purchases/ui/checkout-page";

export const Route = createFileRoute("/checkouts/$purchaseId")({
  beforeLoad: requireAuth,
  component: CheckoutStatusPage,
});

function CheckoutStatusPage() {
  const { purchaseId } = Route.useParams();
  const navigate = useNavigate();
  const checkoutQuery = useQuery(checkoutQueryOptions(purchaseId));
  const checkout = checkoutQuery.data;
  usePurchaseSocket(purchaseId, checkout?.status === "pending");

  useEffect(() => {
    if (checkoutQuery.isPending || (!checkoutQuery.isError && checkout)) return;
    toast.error("Não foi possível encontrar este checkout.");
    void navigate({ to: "/profile", search: { tab: "purchases" } });
  }, [checkout, checkoutQuery.isError, checkoutQuery.isPending, navigate]);

  if (checkoutQuery.isPending) {
    return (
      <main className="grid min-h-screen place-items-center text-sm text-muted-foreground">
        Carregando checkout…
      </main>
    );
  }

  if (checkoutQuery.isError || !checkout) return null;

  return <CheckoutPage purchase={checkout} />;
}
