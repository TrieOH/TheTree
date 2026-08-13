import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { ShoppingBag } from "lucide-react";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { myPurchasesQueryOptions } from "@/features/purchases/api";
import { PurchaseSummary } from "@/features/purchases/ui/purchase-summary";

export const Route = createFileRoute("/profile/purchases")({
  beforeLoad: requireAuth,
  component: PurchasesPage,
});

function PurchasesPage() {
  const { data } = useSuspenseQuery(myPurchasesQueryOptions());

  return (
    <main className="mx-auto w-full max-w-3xl px-4 py-10 pb-28">
      <header className="mb-8">
        <h1 className="text-2xl font-bold">Minhas compras</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Acompanhe seus pagamentos e pedidos.
        </p>
      </header>

      {data.purchases.length === 0 ? (
        <div className="flex min-h-64 flex-col items-center justify-center border border-dashed border-border text-center">
          <ShoppingBag className="mb-3 size-10 text-muted-foreground" />
          <p className="font-medium">Você ainda não possui compras.</p>
        </div>
      ) : (
        <div className="space-y-4">
          {data.purchases.map((purchase) => (
            <Link
              key={purchase.purchase_id}
              to="/checkouts/$purchaseId"
              params={{ purchaseId: purchase.purchase_id }}
              className="block transition-transform hover:-translate-y-0.5"
            >
              <PurchaseSummary purchase={purchase} />
            </Link>
          ))}
        </div>
      )}
    </main>
  );
}
