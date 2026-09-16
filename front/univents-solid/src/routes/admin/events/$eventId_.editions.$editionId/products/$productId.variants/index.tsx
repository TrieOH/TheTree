import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { Button, EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal } from "solid-js";

import ArrowLeftIcon from "~icons/lucide/arrow-left";
import LayersIcon from "~icons/lucide/layers";
import PlusIcon from "~icons/lucide/plus";

import { productVariantsQueryOptions } from "@/features/products/api";
import {
  useCreateVariantMutation,
  useDeleteVariantMutation,
  useUpdateVariantMutation,
} from "@/features/products/api/mutations";
import type { VariantCreateOutputI, VariantI } from "@/features/products/model";
import { AdminCreateVariantCard } from "@/features/products/ui/AdminCreateVariantCard";
import { AdminVariantCard } from "@/features/products/ui/AdminVariantCard";
import { ManageVariantDialog } from "@/features/products/ui/ManageVariantDialog";
import { VariantGalleryDialog } from "@/features/products/ui/VariantGalleryDialog";
import { toast } from "@/shared/ui/toast";
import { AlertModal } from "@/widgets/ui/AlertModal";

const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const Layers = LayersIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/products/$productId/variants/",
)({
  head: () => ({
    meta: [{ title: "Variações do Produto - Admin Univents" }],
  }),
  component: AdminProductVariantsRoute,
});

function AdminProductVariantsRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;
  const editionId = () => params().editionId;
  const productId = () => params().productId;
  const navigate = useNavigate();

  const variantsQuery = useQuery(() => productVariantsQueryOptions(productId()));
  const createMutation = useCreateVariantMutation();
  const updateMutation = useUpdateVariantMutation();
  const deleteMutation = useDeleteVariantMutation();

  const [creating, setCreating] = createSignal(false);
  const [editing, setEditing] = createSignal<VariantI | null>(null);
  const [deleting, setDeleting] = createSignal<VariantI | null>(null);
  const [galleryVariant, setGalleryVariant] = createSignal<VariantI | null>(null);

  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<VariantI>>({
    field: "name",
    direction: "asc",
  });

  const variants = createMemo(() => (variantsQuery().data ?? []) as VariantI[]);

  const handleCreate = async (values: VariantCreateOutputI): Promise<boolean> => {
    try {
      await createMutation.mutateAsync({
        productId: productId(),
        data: values,
      });
      toast.success("Variação criada com sucesso!");
      setCreating(false);
      return true;
    } catch {
      toast.error("Erro ao criar variação.");
      return false;
    }
  };

  const handleUpdate = async (values: VariantCreateOutputI): Promise<boolean> => {
    const current = editing();
    if (!current) return false;
    try {
      await updateMutation.mutateAsync({
        variantId: current.id,
        productId: productId(),
        data: values,
      });
      toast.success("Variação atualizada com sucesso!");
      setEditing(null);
      return true;
    } catch {
      toast.error("Erro ao atualizar variação.");
      return false;
    }
  };

  const handleDelete = async () => {
    const current = deleting();
    if (!current) return;
    try {
      await deleteMutation.mutateAsync({
        variantId: current.id,
        productId: productId(),
      });
      toast.success("Variação excluída com sucesso!");
      setDeleting(null);
    } catch {
      toast.error("Erro ao excluir variação.");
    }
  };

  const handleBack = () => {
    void navigate({
      to: "/admin/events/$eventId/editions/$editionId/products",
      params: {
        eventId: eventId(),
        editionId: editionId(),
      },
    });
  };

  return (
    <div class="mx-auto max-w-7xl space-y-6">
      <div class="flex items-center gap-3">
        <Button
          type="button"
          variant="outline"
          size="icon"
          onClick={handleBack}
          aria-label="Voltar para produtos"
          class="size-9 rounded-full cursor-pointer hover:bg-muted"
        >
          <ArrowLeft class="size-4" />
        </Button>
        <div>
          <h1 class="text-xl font-bold tracking-tight text-foreground">
            Variações do Produto
          </h1>
          <p class="text-xs text-muted-foreground">
            Gerencie fotos, preços, estoque e SKUs das opções disponíveis para este produto.
          </p>
        </div>
      </div>

      <Show
        when={!variantsQuery().isLoading}
        fallback={
          <div class="grid grid-cols-[repeat(auto-fill,minmax(16rem,1fr))] gap-2">
            <For each={[1, 2, 3, 4]}>
              {() => <div class="h-44 animate-pulse rounded-xl bg-muted" />}
            </For>
          </div>
        }
      >
        <AdminProductVariantsContent
          variants={variants()}
          filter={filter()}
          onFilterChange={setFilter}
          sort={sort()}
          onSortChange={setSort}
          onCreate={() => setCreating(true)}
          onEdit={setEditing}
          onDelete={setDeleting}
          onManageGallery={setGalleryVariant}
        />
      </Show>

      {/* Modal de Criação / Edição */}
      <Show
        when={editing()}
        keyed
        fallback={
          <ManageVariantDialog
            open={creating()}
            variant={null}
            onOpenChange={(open) => !open && setCreating(false)}
            onSubmit={handleCreate}
          />
        }
      >
        {(current) => (
          <ManageVariantDialog
            open
            variant={current}
            onOpenChange={(open) => !open && setEditing(null)}
            onSubmit={handleUpdate}
          />
        )}
      </Show>

      {/* Modal de Galeria de Fotos */}
      <Show when={galleryVariant()} keyed>
        {(current) => (
          <VariantGalleryDialog
            open
            variant={current}
            onOpenChange={(open) => !open && setGalleryVariant(null)}
          />
        )}
      </Show>

      {/* Modal de Exclusão */}
      <AlertModal
        open={Boolean(deleting())}
        onOpenChange={(open) => !open && setDeleting(null)}
        title="Excluir variação?"
        description={
          deleting()
            ? `Tem certeza que deseja excluir a variação "${deleting()?.name}" (${deleting()?.vendor_code})?`
            : undefined
        }
        confirmLabel="Excluir variação"
        variant="destructive"
        loading={deleteMutation.result().isPending}
        onConfirm={handleDelete}
      />
    </div>
  );
}

function AdminProductVariantsContent(props: {
  variants: VariantI[];
  filter: string;
  onFilterChange: (value: string) => void;
  sort: SortState<VariantI>;
  onSortChange: (sort: SortState<VariantI>) => void;
  onCreate: () => void;
  onEdit: (variant: VariantI) => void;
  onDelete: (variant: VariantI) => void;
  onManageGallery: (variant: VariantI) => void;
}): JSX.Element {
  const visible = createMemo(() => {
    const search = props.filter.trim().toLowerCase();
    if (!search) return props.variants;

    return props.variants.filter((variant) =>
      [variant.name, variant.vendor_code, variant.description ?? ""].some((val) =>
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
      <Plus class="size-4" />
      Nova variação
    </Button>
  );

  return (
    <PaginatedContainer<VariantI>
      items={visible()}
      layout="grid"
      minItemWidth="16rem"
      maxRows={(columns) => (columns === 1 ? 8 : 4)}
      gap="2"
      sort={props.sort}
      onSortChange={props.onSortChange}
      sortFields={[
        {
          key: "name",
          label: "Nome",
          ascLabel: "A → Z",
          descLabel: "Z → A",
        },
        {
          key: "vendor_code",
          label: "Código",
          ascLabel: "A → Z",
          descLabel: "Z → A",
        },
        {
          key: "price",
          label: "Preço",
          ascLabel: "Menor preço primeiro",
          descLabel: "Maior preço primeiro",
          comparator: (a, b) => a.price - b.price,
        },
        {
          key: "stock",
          label: "Estoque",
          ascLabel: "Menor estoque primeiro",
          descLabel: "Maior estoque primeiro",
          comparator: (a, b) => (a.stock ?? 0) - (b.stock ?? 0),
        },
      ]}
      filterValue={props.filter}
      onFilterChange={props.onFilterChange}
      filterPlaceholder="Buscar por nome ou código..."
      itemLabel="variações"
      emptyState={
        <EmptyState
          class="border-0 bg-transparent px-0 py-4 shadow-none"
          icon={<Layers class="size-6 text-foreground/70" />}
          eyebrow="Variações do produto"
          title="Nenhuma variação cadastrada"
          description={
            props.filter
              ? "Nenhuma variação corresponde à busca informada."
              : "Cadastre variações para disponibilizar diferentes opções de tamanho ou modelo."
          }
          action={emptyStateAction}
        />
      }
      renderItems={(slice, options) => (
        <>
          <AdminCreateVariantCard
            index={0}
            animate={options.animate}
            onCreate={props.onCreate}
          />

          <For each={slice}>
            {(variant, index) => (
              <AdminVariantCard
                variant={variant}
                index={index() + 1}
                animate={options.animate}
                onEdit={props.onEdit}
                onDelete={props.onDelete}
                onManageGallery={props.onManageGallery}
              />
            )}
          </For>
        </>
      )}
    />
  );
}
