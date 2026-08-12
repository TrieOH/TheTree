import { useQueries, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { ArrowLeft, ShoppingCart } from "lucide-react";
import { useState } from "react";
import { activeEditionQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import {
  productsByEditionQueryOptions,
  productVariantsQueryOptions,
} from "@/features/products/api";
import { useCart } from "@/features/products/hooks/use-cart";
import { Cart } from "@/features/products/ui/Cart";
import { ProductCardCompact } from "@/features/products/ui/ProductCardCompact";
import { Logo } from "@/shared/ui/logo";
import { Button } from "@/shared/ui/shadcn/button";

export const Route = createFileRoute("/events/$slug/products")({
  loader: async ({ context, params }) => {
    const event = await context.queryClient.ensureQueryData(
      publicEventBySlugQueryOptions(params.slug),
    );
    if (!event) throw redirect({ to: "/events" });

    void context.queryClient.prefetchQuery(activeEditionQueryOptions(event.id));
    return event;
  },
  component: ProductsPage,
});

function ProductsPage() {
  const navigate = Route.useNavigate();
  const event = Route.useLoaderData();
  const { data: activeEdition } = useSuspenseQuery(
    activeEditionQueryOptions(event.id),
  );
  const { data: products = [] } = useSuspenseQuery(
    productsByEditionQueryOptions(activeEdition?.id ?? ""),
  );
  const variantsResults = useQueries({
    queries: products.map((product) => productVariantsQueryOptions(product.id)),
  });
  const productsWithVariants = products
    .map((product, index) => ({
      product,
      variants: variantsResults[index].data ?? [],
    }))
    .filter(({ variants }) => variants.length > 0);
  const [isCartOpen, setIsCartOpen] = useState(false);
  const cart = useCart(activeEdition?.id ?? "");

  return (
    <div className="min-h-screen bg-background pb-24">
      <header className="sticky top-0 z-30 border-b border-border bg-background/80 backdrop-blur-xl">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex h-14 items-center justify-between gap-3">
            <div className="flex items-center gap-3">
              <div className="h-8 w-8">
                <Logo variant="icon" imgClassName="object-left" />
              </div>
              <h1 className="border-l border-border pl-3 text-lg font-semibold text-foreground md:text-xl">
                Produtos
                <span className="ml-2 text-sm font-normal text-muted-foreground">
                  ({productsWithVariants.length})
                </span>
              </h1>
            </div>
            <Link
              to="/events/$slug"
              params={{ slug: event.slug }}
              className="inline-flex items-center gap-1.5 rounded-lg px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
            >
              <ArrowLeft className="h-4 w-4" />
              <span className="hidden sm:inline">Voltar</span>
            </Link>
          </div>
        </div>
      </header>

      <main className="mx-auto flex max-w-7xl flex-wrap items-stretch justify-start gap-6 px-4 py-8 md:gap-8 md:px-6 md:py-12 lg:px-8">
        {productsWithVariants.map(({ product, variants }) => (
          <ProductCardCompact
            key={product.id}
            product={product}
            variants={variants}
            editionId={activeEdition?.id}
          />
        ))}
      </main>

      {activeEdition && (
        <Cart
          isOpen={isCartOpen}
          eventId={event.id}
          editionId={activeEdition.id}
          onClose={() => setIsCartOpen(false)}
          onCheckout={() =>
            navigate({
              to: "/events/$slug/checkout",
              params: { slug: event.slug },
            })
          }
        />
      )}
      {activeEdition && (
        <Button
          type="button"
          onClick={() => setIsCartOpen(true)}
          className="fixed bottom-24 right-4 z-40 h-13 rounded-full px-5 shadow-md shadow-primary/10 transition-transform hover:scale-105 md:right-8"
          aria-label="Abrir carrinho"
        >
          <ShoppingCart className="mr-2 h-5 w-5" />
          <span className="hidden sm:inline">Carrinho</span>
          {cart.itemCount > 0 && (
            <span className="ml-2 rounded-full bg-background px-2 py-0.5 text-xs font-bold text-foreground">
              {cart.itemCount}
            </span>
          )}
        </Button>
      )}
    </div>
  );
}
