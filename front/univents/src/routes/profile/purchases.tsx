import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { ShoppingBag } from "lucide-react";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { myPurchasesQueryOptions } from "@/features/purchases/api";
import { PurchaseSummary } from "@/features/purchases/ui/purchase-summary";
import { isVerifiedEmailRequiredError } from "@/shared/lib/errors";

export const Route = createFileRoute("/profile/purchases")({
  beforeLoad: requireAuth,
  component: PurchasesPage,
});

function PurchasesPage() {
  const { data, error, isPending } = useQuery(myPurchasesQueryOptions());

  return (
    <main className="mx-auto w-full max-w-3xl px-4 py-10 pb-28">
      <header className="mb-8">
        <h1 className="text-2xl font-bold">Minhas compras</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Acompanhe seus pagamentos e pedidos.
        </p>
      </header>

      {isPending ? (
        <p className="text-sm text-muted-foreground">Carregando compras…</p>
      ) : error ? (
        <div className="flex min-h-64 flex-col items-center justify-center border border-dashed border-border px-6 text-center">
          <ShoppingBag className="mb-3 size-10 text-muted-foreground" />
          <p className="font-medium">
            {isVerifiedEmailRequiredError(error)
              ? "Verifique seu e-mail para acessar suas compras."
              : "Não foi possível carregar suas compras."}
          </p>
        </div>
      ) : data.purchases.length === 0 ? (
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
