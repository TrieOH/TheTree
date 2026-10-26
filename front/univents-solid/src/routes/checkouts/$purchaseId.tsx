import { createEffect, createMemo } from "solid-js";
import { createFileRoute } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { checkoutQueryOptions } from "@/features/purchases/api";
import { usePurchaseSocket } from "@/features/purchases/hooks/use-purchase-socket";
import CheckoutPage from "@/features/purchases/ui/CheckoutPage";
import { classifyError } from "@/shared/lib/service-error";
import { QueryBoundary } from "@/shared/ui/QueryBoundary";
import { toast } from "@/shared/ui/toast";

export const Route = createFileRoute("/checkouts/$purchaseId")({
  beforeLoad: requireAuth,
  component: CheckoutStatusPage,
});

function CheckoutStatusPage() {
  const params = Route.useParams();
  const navigate = Route.useNavigate();
  const checkoutQuery = useQuery(() => checkoutQueryOptions(params().purchaseId));

  // Only a definitive "no such checkout" sends the user away. A slow or failing
  // API (or a failed background refresh while the order is on screen) must not
  // eject someone from their own payment.
  const notFound = createMemo(
    () => checkoutQuery().isError && classifyError(checkoutQuery().error).kind === "not-found",
    { name: "CheckoutStatusPage.notFound" },
  );

  createEffect(
    () => notFound(),
    (missing) => {
      if (!missing) return;
      toast.error("Não foi possível encontrar este checkout.");
      void navigate({ to: "/profile", search: { tab: "purchases" } });
    },
  );

  usePurchaseSocket(
    () => params().purchaseId,
    (_purchaseId) => {
      void checkoutQuery().refetch();
    },
  );

  return (
    <QueryBoundary
      query={checkoutQuery}
      what="o pedido"
      skeleton={<CheckoutStatusSkeleton />}
      class="mx-auto my-16 max-w-xl"
    >
      {(checkout) => <CheckoutPage purchase={checkout()} />}
    </QueryBoundary>
  );
}

function CheckoutStatusSkeleton() {
  return (
    <div class="mx-auto min-h-screen w-full max-w-3xl space-y-4 px-4 py-10">
      <div class="h-8 w-48 animate-pulse rounded bg-muted" />
      <div class="h-40 animate-pulse rounded-md bg-muted" />
      <div class="h-56 animate-pulse rounded-md bg-muted" />
    </div>
  );
}
