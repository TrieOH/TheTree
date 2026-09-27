import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { Button, EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal } from "solid-js";

import PackageIcon from "~icons/lucide/package";
import PackagePlusIcon from "~icons/lucide/package-plus";

import { productsByEditionQueryOptions } from "@/features/products/api";
import {
  useCreateInitialProductMutation,
  useDeleteProductMutation,
  useUpdateProductMutation,
} from "@/features/products/api/mutations";
import type {
  CreateInitialProductOutputI,
  ProductI,
  ProductPatchOutputI,
} from "@/features/products/model";
import { AdminCreateProductCard } from "@/features/products/ui/AdminCreateProductCard";
import { AdminProductCard } from "@/features/products/ui/AdminProductCard";
import { ManageProductDialog } from "@/features/products/ui/ManageProductDialog";
import { toast } from "@/shared/ui/toast";
import { AlertModal } from "@/widgets/ui/AlertModal";

const Package = PackageIcon as unknown as (props: { class?: string }) => JSX.Element;
const PackagePlus = PackagePlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/products/",
)({
  head: () => ({
    meta: [{ title: "Produtos - Admin Univents" }],
  }),
  component: AdminEditionProductsRoute,
});

function AdminEditionProductsRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;
  const editionId = () => params().editionId;
  const navigate = useNavigate();

  const productsQuery = useQuery(() => productsByEditionQueryOptions(editionId()));
  const createMutation = useCreateInitialProductMutation();
  const updateMutation = useUpdateProductMutation();
  const deleteMutation = useDeleteProductMutation();

  const [creating, setCreating] = createSignal(false);
  const [editing, setEditing] = createSignal<ProductI | null>(null);
  const [deleting, setDeleting] = createSignal<ProductI | null>(null);

  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<ProductI>>({
    field: "vendor_code",
    direction: "asc",
  });

  const products = createMemo(() => (productsQuery().data ?? []) as ProductI[]);

  const handleCreate = async (values: CreateInitialProductOutputI): Promise<boolean> => {
    try {
      await createMutation.mutateAsync({
        editionId: editionId(),
        data: values,
      });
      toast.success("Produto criado com sucesso!");
      setCreating(false);
      return true;
    } catch {
      toast.error("Erro ao criar produto.");
      return false;
    }
  };

  const handleUpdate = async (values: ProductPatchOutputI): Promise<boolean> => {
    const current = editing();
    if (!current) return false;
    try {
      await updateMutation.mutateAsync({
        productId: current.id,
        editionId: editionId(),
        data: values,
      });
      toast.success("Produto atualizado com sucesso!");
      setEditing(null);
      return true;
    } catch {
      toast.error("Erro ao atualizar produto.");
      return false;
    }
  };

  const handleDelete = async () => {
    const current = deleting();
    if (!current) return;
    try {
      await deleteMutation.mutateAsync({
        productId: current.id,
        editionId: editionId(),
      });
      toast.success("Produto excluído com sucesso!");
      setDeleting(null);
    } catch {
      toast.error("Erro ao excluir produto.");
    }
  };

  const handleManageVariants = (product: ProductI) => {
    void navigate({
      to: "/admin/events/$eventId/editions/$editionId/products/$productId/variants",
      params: {
        eventId: eventId(),
        editionId: editionId(),
        productId: product.id,
      },
    });
  };

  return (
    <div class="mx-auto max-w-7xl space-y-6">
      <Show
        when={!productsQuery().isLoading}
        fallback={
          <div class="grid grid-cols-[repeat(auto-fill,minmax(16rem,1fr))] gap-2">
            <For each={[1, 2, 3, 4]}>
              {() => <div class="h-28 animate-pulse rounded-xl bg-muted" />}
            </For>
          </div>
        }
      >
        <AdminEditionProductsContent
          products={products()}
          filter={filter()}
          onFilterChange={setFilter}
          sort={sort()}
          onSortChange={setSort}
          onCreate={() => setCreating(true)}
          onEdit={setEditing}
          onDelete={setDeleting}
          onManageVariants={handleManageVariants}
        />
      </Show>

      {/* Modal de Criação / Edição */}
      <Show
        when={editing()}
        keyed
        fallback={
          <ManageProductDialog
            open={creating()}
            product={null}
            onOpenChange={(open) => !open && setCreating(false)}
            onCreate={handleCreate}
          />
        }
      >
        {(current) => (
          <ManageProductDialog
            open
            product={current}
            onOpenChange={(open) => !open && setEditing(null)}
            onUpdate={handleUpdate}
          />
        )}
      </Show>

      {/* Modal de Exclusão */}
      <AlertModal
        open={Boolean(deleting())}
        onOpenChange={(open) => !open && setDeleting(null)}
        title="Excluir produto?"
        description={
          deleting()
            ? `Tem certeza que deseja excluir o produto "${deleting()?.vendor_code}"? Suas variações e estoque também serão afetados.`
            : undefined
        }
        confirmLabel="Excluir produto"
        variant="destructive"
        loading={deleteMutation.result().isPending}
        onConfirm={handleDelete}
      />
    </div>
  );
}

function AdminEditionProductsContent(props: {
  products: ProductI[];
  filter: string;
  onFilterChange: (value: string) => void;
  sort: SortState<ProductI>;
  onSortChange: (sort: SortState<ProductI>) => void;
  onCreate: () => void;
  onEdit: (product: ProductI) => void;
  onDelete: (product: ProductI) => void;
  onManageVariants: (product: ProductI) => void;
}): JSX.Element {
  const visible = createMemo(() => {
    const search = props.filter.trim().toLowerCase();
    if (!search) return props.products;

    return props.products.filter((product) =>
      [product.vendor_code, product.requires_registration ? "cadastro" : "livre"].some((val) =>
        val.toLowerCase().includes(search),
      ),
    );
  });

  const emptyStateAction = (
    <Button
      size="sm"
      class="gap-2 rounded-sm py-4"
      onClick={() => props.onCreate()}
    >
      <PackagePlus class="size-4" />
      Novo produto
    </Button>
  );

  return (
    <PaginatedContainer<ProductI>
      items={visible()}
      layout="grid"
      minItemWidth="16rem"
      maxRows={(columns) => (columns === 1 ? 8 : 4)}
      gap="2"
      sort={props.sort}
      onSortChange={props.onSortChange}
      sortFields={[
        {
          key: "vendor_code",
          label: "Código",
          ascLabel: "A → Z",
          descLabel: "Z → A",
        },
        {
          key: "created_at",
          label: "Data de criação",
          ascLabel: "Mais antigos primeiro",
          descLabel: "Mais recentes primeiro",
          comparator: (a, b) =>
            new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
        },
        {
          key: "requires_registration",
          label: "Exige cadastro",
          comparator: (a, b) =>
            Number(a.requires_registration) - Number(b.requires_registration),
        },
      ]}
      filterValue={props.filter}
      onFilterChange={props.onFilterChange}
      filterPlaceholder="Buscar por código..."
      itemLabel="produtos"
      emptyState={
        <EmptyState
          class="border-0 bg-transparent px-0 py-4 shadow-none"
          icon={<Package class="size-6 text-foreground/70" />}
          eyebrow="Produtos da edição"
          title="Nenhum produto cadastrado"
          description={
            props.filter
              ? "Nenhum produto corresponde à busca informada."
              : "Cadastre produtos para vender itens físicos ou digitais nesta edição."
          }
          action={emptyStateAction}
        />
      }
      renderItems={(slice, options) => (
        <>
          <AdminCreateProductCard
            index={0}
            animate={options.animate}
            onCreate={props.onCreate}
          />

          <For each={slice}>
            {(product, index) => (
              <AdminProductCard
                product={product}
                index={index() + 1}
                animate={options.animate}
                onEdit={props.onEdit}
                onDelete={props.onDelete}
                onManageVariants={props.onManageVariants}
              />
            )}
          </For>
        </>
      )}
    />
  );
}
