import { createSignal, Switch, Match, createEffect } from "solid-js";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core-solid";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { checkoutQueryOptions } from "@/features/purchases/api";
import { usePurchaseSocket } from "@/features/purchases/hooks/use-purchase-socket";
import CheckoutPage from "@/features/purchases/ui/CheckoutPage";
import { toast } from "@/shared/ui/toast";
import type { Checkout } from "@trieoh/univents-api/schemas";

export const Route = createFileRoute("/checkouts/$purchaseId")({
  beforeLoad: requireAuth,
  component: CheckoutStatusPage,
});

function CheckoutStatusPage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [checkout, setCheckout] = createSignal<Checkout | null>(null);
  const [status, setStatus] = createSignal<"pending" | "success" | "error">("pending");

  const fetchCheckout = async (purchaseId: string, bypassCache = false) => {
    try {
      const options = checkoutQueryOptions(purchaseId);
      const data = await queryClient.fetchQuery(
        bypassCache ? { ...options, staleTime: 0 } : options
      );

      if (!data) throw new Error("Checkout não retornado");

      setCheckout(data);
      setStatus("success");
    } catch {
      setStatus("error");
    }
  };

  createEffect(
    () => params().purchaseId,
    (purchaseId) => {
      setStatus("pending");
      void fetchCheckout(purchaseId);
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
    (purchaseId) => {
      void fetchCheckout(purchaseId, true);
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