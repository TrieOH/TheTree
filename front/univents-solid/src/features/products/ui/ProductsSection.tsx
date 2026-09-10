import type {
  Product,
  ProductVariant,
  StoreStockItem,
} from "@trieoh/univents-api/schemas";
import { For, Show, createMemo } from "solid-js";
import { Link } from "@tanstack/solid-router";
import { ProductCard } from "./ProductCard";

export function ProductsSection(props: {
  products: { product: Product; variants: ProductVariant[] }[];
  stock: StoreStockItem[];
  eventSlug: string;
  editionId: string;
}) {
  const products = createMemo(() =>
    props.products.filter(({ variants }) => variants.length > 0),
  );
  const stock = createMemo(
    () => new Map(props.stock.map((item) => [item.id, item.stock])),
  );

  return (
    <Show when={products().length > 0}>
      <section class="w-full py-5">
        <div class="mb-8 text-center">
          <h2 class="text-3xl font-semibold tracking-tight text-foreground">
            Produtos
          </h2>
          <p class="mx-auto mt-2 max-w-xl text-sm leading-relaxed text-muted-foreground">
            Itens exclusivos disponíveis para esta edição do evento.
          </p>
        </div>
        <div class="mx-auto flex max-w-6xl flex-wrap justify-start gap-5">
          <For each={products().slice(0, 3)}>
            {({ product, variants }) => (
              <ProductCard
                product={product}
                variants={variants}
                stock={stock()}
                editionId={props.editionId}
              />
            )}
          </For>
        </div>
        <div class="mt-6 text-center">
          <Link
            class="text-sm font-semibold text-primary"
            to="/events/$slug/store"
            params={{ slug: props.eventSlug }}
            search={{ tab: "products" }}
          >
            Ver todos os produtos →
          </Link>
        </div>
      </section>
    </Show>
  );
}
