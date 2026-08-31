import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { PaginatedContainer } from "@trieoh/ui-base";
import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import { ShoppingBag } from "lucide-react";
import { useMemo, useState } from "react";
import { editionPurchasesQueryOptions } from "@/features/purchases/api";
import { useRefundPurchaseMutation } from "@/features/purchases/api/mutations";
import { AdminPurchaseCard } from "@/features/purchases/ui/admin-purchase-card";
import { Card, CardContent } from "@/shared/ui/shadcn/card";
import { AlertModal } from "@/widgets/ui/alert-modal";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/purchases/",
)({ component: EditionPurchasesRoute });

function EditionPurchasesRoute() {
  const { editionId } = Route.useParams();
  const { data: purchases = [] } = useQuery(
    editionPurchasesQueryOptions(editionId),
  );
  const refundMutation = useRefundPurchaseMutation();
  const [search, setSearch] = useState("");
  const [refundId, setRefundId] = useState<string | null>(null);
  const [sort, setSort] = useState<SortState<EditionPurchase>>({
    field: "created_at",
    direction: "desc",
  });
  const filtered = useMemo(() => {
    const term = search.trim().toLowerCase();
    return purchases.filter(
      (purchase) =>
        !term ||
        [
          purchase.payer_email ?? "",
          purchase.purchase_id,
          purchase.status,
          purchase.payment_method ?? "",
          ...purchase.attendees.map((item) => item.email),
        ].some((value) => value.toLowerCase().includes(term)),
    );
  }, [purchases, search]);

  const refund = async () => {
    if (!refundId) return;
    try {
      await refundMutation.mutateAsync({ purchaseId: refundId, editionId });
      setRefundId(null);
    } catch {
      // Error feedback is centralized in the mutation.
    }
  };

  return (
    <div className="mx-auto max-w-7xl space-y-6 p-6 pb-28!">
      <PaginatedContainer
        items={filtered}
        layout="grid"
        pageSize={6}
        gap="4"
        minItemWidth="16rem"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          { key: "created_at", label: "Data" },
          {
            key: "total_cents",
            label: "Valor",
            comparator: (a, b) => a.total_cents - b.total_cents,
          },
        ]}
        filterValue={search}
        onFilterChange={setSearch}
        filterPlaceholder="Buscar por e-mail, status ou participante"
        itemLabel="compras"
        emptyState={
          <Card>
            <CardContent className="flex flex-col items-center gap-3 py-12 text-center text-muted-foreground">
              <ShoppingBag className="size-8" />
              <p>Nenhuma compra encontrada.</p>
            </CardContent>
          </Card>
        }
        renderItems={(slice) =>
          slice.map((purchase) => (
            <AdminPurchaseCard
              key={purchase.purchase_id}
              purchase={purchase}
              onRefund={(item) => setRefundId(item.purchase_id)}
            />
          ))
        }
      />
      <AlertModal
        open={refundId !== null}
        onOpenChange={(open) =>
          !open && !refundMutation.isPending && setRefundId(null)
        }
        title="Solicitar reembolso?"
        description="O valor pago será devolvido ao comprador. As taxas de processamento e de marketplace já cobradas não são devolvidas à organização."
        confirmLabel="Solicitar reembolso"
        variant="destructive"
        loading={refundMutation.isPending}
        onConfirm={refund}
      />
    </div>
  );
}
