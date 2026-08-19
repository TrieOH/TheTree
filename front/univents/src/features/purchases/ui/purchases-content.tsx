import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { ShoppingBag } from "lucide-react";
import { myPurchasesQueryOptions } from "../api";
import { resolvePurchaseCatalog } from "../api/purchase-catalog";
import { type Order, OrderCard, orderIcons } from "./order-card";

const money = (cents: number, currency: string) =>
  new Intl.NumberFormat("pt-BR", { style: "currency", currency }).format(
    cents / 100,
  );

export function PurchasesContent() {
  const { data, error, isPending } = useQuery(myPurchasesQueryOptions());
  const purchases = data?.purchases ?? [];
  const catalogItems = purchases.flatMap((purchase) =>
    purchase.items.map((item) => ({
      ...item,
      edition_id: purchase.edition_id,
    })),
  );
  const { data: catalog = {} } = useQuery({
    queryKey: [
      "purchase-catalog",
      catalogItems.map((item) => `${item.item_type}:${item.item_id}`).join(","),
    ],
    queryFn: () => resolvePurchaseCatalog(catalogItems),
    enabled: catalogItems.some((item) => item.item_id.includes("-")),
  });
  const orders: Order[] = purchases.map((purchase) => ({
    purchaseId: purchase.purchase_id,
    status:
      purchase.status === "approved" || purchase.status === "pending"
        ? purchase.status
        : purchase.status === "refunded"
          ? "refunded"
          : "cancelled",
    paymentIcon:
      purchase.payment_method === "pix"
        ? orderIcons.Zap
        : orderIcons.CreditCard,
    paymentText:
      purchase.payment_method === "pix" ? "Pix" : "Cartão de crédito",
    date: new Intl.DateTimeFormat("pt-BR", { dateStyle: "medium" }).format(
      new Date(purchase.created_at ?? purchase.expires_at),
    ),
    total: money(purchase.total_cents, purchase.currency),
    items: purchase.items.map((item) => {
      const details = catalog[`${item.item_type}:${item.item_id}`];
      return {
        icon:
          item.item_type === "ticket"
            ? orderIcons.Ticket
            : item.item_type === "product"
              ? orderIcons.Wallet
              : orderIcons.FileText,
        image: details?.image ?? undefined,
        title: details?.name ?? "Item não identificado",
        description:
          details?.description ??
          `${item.quantity} unidade${item.quantity > 1 ? "s" : ""}`,
        price: money(item.unit_price_cents * item.quantity, purchase.currency),
      };
    }),
  }));

  if (isPending) {
    return <p className="text-sm text-muted-foreground">Carregando compras…</p>;
  }

  if (error || orders.length === 0) {
    return (
      <div className="flex min-h-64 flex-col items-center justify-center border border-dashed border-border text-center">
        <ShoppingBag className="mb-3 size-10 text-muted-foreground" />
        <p className="font-medium">
          {error
            ? "Não foi possível carregar suas compras."
            : "Você ainda não possui compras."}
        </p>
      </div>
    );
  }
  return (
    <div className="grid gap-4 lg:grid-cols-2">
      {orders.map((order) => (
        <Link
          key={order.purchaseId}
          to="/checkouts/$purchaseId"
          params={{ purchaseId: order.purchaseId }}
          className="block h-full"
        >
          <OrderCard order={order} />
        </Link>
      ))}
    </div>
  );
}
