import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import { For, Show, createMemo, createSignal } from "solid-js";

import ShoppingBagIcon from "~icons/lucide/shopping-bag";

import { editionPurchasesQueryOptions } from "@/features/purchases/api";
import { useRefundPurchaseMutation } from "@/features/purchases/api/mutations";
import { AdminPurchaseCard } from "@/features/purchases/ui/AdminPurchaseCard";
import { AlertModal } from "@/widgets/ui/AlertModal";
import { toast } from "@/shared/ui/toast";

const ShoppingBag = ShoppingBagIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/purchases/",
)({
  head: () => ({
    meta: [{ title: "Compras - Admin Univents" }],
  }),
  component: AdminEditionPurchasesRoute,
});

function AdminEditionPurchasesRoute(): JSX.Element {
  const params = Route.useParams();
  const editionId = () => params().editionId;

  const purchasesQuery = useQuery(() => editionPurchasesQueryOptions(editionId()));
  const refundMutation = useRefundPurchaseMutation();

  const [refundTarget, setRefundTarget] = createSignal<EditionPurchase | null>(null);
  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<EditionPurchase>>({
    field: "created_at",
    direction: "desc",
  });

  const purchases = createMemo(() => (purchasesQuery().data ?? []) as EditionPurchase[]);

  const handleRefundConfirm = async () => {
    const target = refundTarget();
    if (!target) return;

    try {
      await refundMutation.mutateAsync({
        purchaseId: target.purchase_id,
        editionId: editionId(),
      });
      toast.success("Reembolso solicitado com sucesso!");
      setRefundTarget(null);
    } catch {
      toast.error("Não foi possível solicitar o reembolso.");
    }
  };

  return (
    <div class="mx-auto max-w-7xl space-y-6">
      <Show
        when={!purchasesQuery().isLoading}
        fallback={
          <div class="grid grid-cols-[repeat(auto-fill,minmax(16rem,1fr))] gap-2">
            <For each={[1, 2, 3, 4]}>
              {() => <div class="h-32 animate-pulse rounded-xl bg-muted" />}
            </For>
          </div>
        }
      >
        <AdminEditionPurchasesContent
          purchases={purchases()}
          filter={filter()}
          onFilterChange={setFilter}
          sort={sort()}
          onSortChange={setSort}
          onRefund={setRefundTarget}
        />
      </Show>

      <AlertModal
        open={refundTarget() !== null}
        onOpenChange={(open) => !open && setRefundTarget(null)}
        title="Solicitar reembolso?"
        description="O valor pago será devolvido ao comprador. As taxas de processamento e de marketplace já cobradas não são devolvidas à organização."
        confirmLabel="Solicitar reembolso"
        variant="destructive"
        loading={refundMutation.result().status === "pending"}
        onConfirm={handleRefundConfirm}
      />
    </div>
  );
}

function AdminEditionPurchasesContent(props: {
  purchases: EditionPurchase[];
  filter: string;
  onFilterChange: (value: string) => void;
  sort: SortState<EditionPurchase>;
  onSortChange: (sort: SortState<EditionPurchase>) => void;
  onRefund: (purchase: EditionPurchase) => void;
}): JSX.Element {
  const visible = createMemo(() => {
    const term = props.filter.trim().toLowerCase();
    if (!term) return props.purchases;

    return props.purchases.filter((purchase) =>
      [
        purchase.payer_email ?? "",
        purchase.purchase_id,
        purchase.status,
        purchase.payment_method ?? "",
        ...(purchase.attendees || []).map((att) => att.email),
        ...(purchase.attendees || []).map((att) => att.name || ""),
      ].some((value) => value.toLowerCase().includes(term)),
    );
  });

  return (
    <PaginatedContainer<EditionPurchase>
      items={visible()}
      layout="grid"
      minItemWidth="16rem"
      maxRows={(columns) => (columns === 1 ? 8 : 4)}
      gap="2"
      sort={props.sort}
      onSortChange={props.onSortChange}
      sortFields={[
        {
          key: "created_at",
          label: "Data",
          ascLabel: "Mais antigas primeiro",
          descLabel: "Mais recentes primeiro",
          comparator: (a, b) => {
            const dateA = new Date(a.created_at || 0).getTime();
            const dateB = new Date(b.created_at || 0).getTime();
            return dateA - dateB;
          },
        },
        {
          key: "total_cents",
          label: "Valor",
          ascLabel: "Menor valor primeiro",
          descLabel: "Maior valor primeiro",
          comparator: (a, b) => a.total_cents - b.total_cents,
        },
      ]}
      filterValue={props.filter}
      onFilterChange={props.onFilterChange}
      filterPlaceholder="Buscar por e-mail, status, participante ou ID..."
      itemLabel="compras"
      emptyState={
        <EmptyState
          class="border-0 bg-transparent px-0 py-4 shadow-none"
          icon={<ShoppingBag class="size-6 text-foreground/70" />}
          eyebrow="Compras da edição"
          title="Nenhuma compra encontrada"
          description={
            props.filter
              ? "Nenhuma compra corresponde à busca informada."
              : "As compras realizadas pelos participantes aparecerão aqui."
          }
        />
      }
      renderItems={(slice, options) => (
        <For each={slice}>
          {(purchase, index) => (
            <AdminPurchaseCard
              purchase={purchase}
              index={index()}
              animate={options.animate}
              onRefund={props.onRefund}
            />
          )}
        </For>
      )}
    />
  );
}
