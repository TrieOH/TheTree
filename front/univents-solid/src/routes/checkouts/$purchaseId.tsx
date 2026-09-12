import { Switch, Match, createEffect } from "solid-js";
import { createFileRoute } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { checkoutQueryOptions } from "@/features/purchases/api";
import { usePurchaseSocket } from "@/features/purchases/hooks/use-purchase-socket";
import CheckoutPage from "@/features/purchases/ui/CheckoutPage";
import { toast } from "@/shared/ui/toast";

export const Route = createFileRoute("/checkouts/$purchaseId")({
  beforeLoad: requireAuth,
  component: CheckoutStatusPage,
});

function CheckoutStatusPage() {
  const params = Route.useParams();
  const navigate = Route.useNavigate();
  const checkoutQuery = useQuery(() => checkoutQueryOptions(params().purchaseId));
  const checkout = () => checkoutQuery().data;
  const status = () => checkoutQuery().isError ? "error" : checkoutQuery().isPending ? "pending" : "success";

  createEffect(
    () => params().purchaseId,
    () => {
    }
  );

  createEffect(
    () => status(),
    (currentStatus) => {
      if (currentStatus === "error") {
        toast.error("Não foi possível encontrar este checkout.");
        void navigate({ to: "/profile", search: { tab: "purchases" } });
      }
    }
  );

  usePurchaseSocket(
    () => params().purchaseId,
    (_purchaseId) => {
      void checkoutQuery().refetch();
    },
  );

  return (
    <Switch>
      <Match when={status() === "pending"}>
        <main class="grid min-h-screen place-items-center text-sm text-muted-foreground">
          Carregando checkout…
        </main>
      </Match>

      <Match when={status() === "success" && checkout()}>
        {(loadedCheckout) => <CheckoutPage purchase={loadedCheckout()} />}
      </Match>
    </Switch>
  );
}
