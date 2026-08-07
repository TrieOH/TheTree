import { useQueries, useSuspenseQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";
import {
  productsByEditionQueryOptions,
  productVariantsQueryOptions,
} from "../api";
import { ProductCard } from "./ProductCard";

interface ProductsSectionProps {
  editionId?: string;
  eventSlug: string;
}

export function ProductsSection({ editionId, eventSlug }: ProductsSectionProps) {
  if (!editionId) return null;

  const { data: products } = useSuspenseQuery(
    productsByEditionQueryOptions(editionId),
  );

  const variantsResults = useQueries({
    queries: products.map((product) => productVariantsQueryOptions(product.id)),
  });

  const productsWithVariants = products
    .map((product, i) => ({
      product,
      variants: variantsResults[i].data ?? [],
    }))
    .filter(({ variants }) => variants.length > 0);

  if (productsWithVariants.length === 0) return null;

  return (
    <section className="w-full py-5">
      <div className="px-4">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-semibold text-foreground tracking-tight">
            Produtos
          </h2>
          <p className="mt-2 text-sm text-muted-foreground leading-relaxed max-w-xl mx-auto">
            Itens exclusivos disponíveis para esta edição do evento.
          </p>
        </div>

        <div className="flex flex-wrap justify-center gap-5 max-w-6xl mx-auto">
          {productsWithVariants.slice(0, 3).map(({ product, variants }) => (
            <ProductCard
              key={product.id}
              product={product}
              variants={variants}
              maxVariants={3}
              editionId={editionId}
            />
          ))}
        </div>

        <div className="mt-6 text-center">
          <Link
            to="/events/$slug/products"
            params={{ slug: eventSlug }}
            className="inline-flex items-center gap-1.5 text-sm font-semibold text-primary transition-all duration-200 hover:gap-2.5"
          >
            Ver todos os produtos
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
      </div>
    </section>
  );
}
