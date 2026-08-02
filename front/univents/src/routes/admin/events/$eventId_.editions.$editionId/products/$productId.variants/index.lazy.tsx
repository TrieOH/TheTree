import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute, useRouter } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { ArrowLeft, Layers, Plus } from "lucide-react";
import { useState } from "react";
import { productVariantsQueryOptions } from "@/features/products/api";
import {
  useCreateVariantMutation,
  useDeleteVariantMutation,
  useUpdateVariantMutation,
} from "@/features/products/api/mutations";
import type { VariantI } from "@/features/products/model";
import { AdminVariantCard } from "@/features/products/ui/AdminVariantCard";
import { ManageVariantModal } from "@/features/products/ui/ManageVariantModal";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/products/$productId/variants/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { productId } = Route.useParams();
  const router = useRouter();
  const { data: variants = [] } = useQuery(
    productVariantsQueryOptions(productId),
  );
  const createVariantMutation = useCreateVariantMutation();
  const updateVariantMutation = useUpdateVariantMutation();
  const deleteVariantMutation = useDeleteVariantMutation();

  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState<SortState<VariantI>>({
    field: "name",
    direction: "asc",
  });
  const [modalState, setModalState] = useState<{
    open: boolean;
    variant?: VariantI;
  }>({ open: false });
  const [deletingVariant, setDeletingVariant] = useState<VariantI | null>(null);

  const filteredVariants = [...variants]
    .filter((variant) => {
      const search = filter.trim().toLowerCase();
      if (!search) return true;

      return [
        variant.name,
        variant.vendor_code,
        variant.description ?? "",
      ].some((value) => value.toLowerCase().includes(search));
    })
    .sort((a, b) => {
      const direction = sort.direction === "asc" ? 1 : -1;

      if (sort.field === "price") {
        return (a.price - b.price) * direction;
      }

      if (sort.field === "stock") {
        return ((a.stock ?? 0) - (b.stock ?? 0)) * direction;
      }

      return (
        String(a[sort.field]).localeCompare(String(b[sort.field])) * direction
      );
    });

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <div className="mb-4 flex w-full items-center gap-3">
        <button
          type="button"
          onClick={() => router.history.back()}
          className="inline-flex size-9 items-center justify-center rounded-full border border-border/60 text-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
        >
          <ArrowLeft className="size-4" />
        </button>
        <h2 className="text-lg font-semibold text-foreground">
          Variações do produto
        </h2>
      </div>

      <PaginatedContainer<VariantI>
        items={filteredVariants}
        layout="grid"
        minItemWidth="14rem"
        pageSize={10}
        gap="3"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          { key: "name", label: "Nome" },
          { key: "vendor_code", label: "Código" },
          {
            key: "price",
            label: "Preço",
            comparator: (a, b) => a.price - b.price,
          },
          {
            key: "stock",
            label: "Estoque",
            comparator: (a, b) => (a.stock ?? 0) - (b.stock ?? 0),
          },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome, código ou descrição..."
        itemLabel="variações"
        headerActions={
          <Button
            type="button"
            className="h-9 gap-2"
            onClick={() => setModalState({ open: true, variant: undefined })}
          >
            <Plus className="size-4" />
            Nova variação
          </Button>
        }
        emptyState={
          <EmptyState
            icon={Layers}
            eyebrow="Variações"
            title="Nenhuma variação encontrada"
            description="Crie a primeira variação para este produto."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((variant, index) => (
            <AdminVariantCard
              key={variant.id}
              variant={variant}
              index={index}
              onEdit={() => setModalState({ open: true, variant })}
              onDelete={() => setDeletingVariant(variant)}
            />
          ))
        }
      />

      <ManageVariantModal
        key={modalState.variant?.id ?? "variant-create"}
        open={modalState.open}
        productId={productId}
        variant={modalState.variant}
        onOpenChange={(open) => {
          if (open) {
            setModalState((prev) => ({ ...prev, open }));
            return;
          }
          setModalState({ open: false, variant: undefined });
        }}
        onCreate={async (values) => {
          const variant = await createVariantMutation.mutateAsync({
            productId,
            data: values,
          });
          return variant ? variant : false;
        }}
        onUpdate={async (variantId, values) => {
          const variant = await updateVariantMutation.mutateAsync({
            variantId,
            data: values,
          });
          return variant ? variant : false;
        }}
      />

      <AlertModal
        open={Boolean(deletingVariant)}
        onOpenChange={() => setDeletingVariant(null)}
        title="Excluir variação?"
        description={
          deletingVariant
            ? `Ao excluir a variação "${deletingVariant.name}".`
            : undefined
        }
        confirmLabel="Excluir variação"
        variant="destructive"
        loading={deleteVariantMutation.isPending}
        onConfirm={async () => {
          if (!deletingVariant) return;
          await deleteVariantMutation.mutateAsync({
            variantId: deletingVariant.id,
            productId,
          });
          setDeletingVariant(null);
        }}
      />
    </div>
  );
}
