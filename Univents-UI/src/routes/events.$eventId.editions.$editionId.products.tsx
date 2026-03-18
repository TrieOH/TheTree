import { createFileRoute } from "@tanstack/react-router";
import { ShoppingBag } from "lucide-react";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { ProductList } from "@/features/products/ui/ProductList";
import { Cart } from "@/features/products/ui/Cart";
import { useCart } from "@/features/products/hooks/use-cart";
import { Button } from "@/shared/ui/shadcn/button";
import { cn } from "@/shared/lib/utils";
import { allProductsQueryOptions } from "@/features/products/api";

export const Route = createFileRoute("/events/$eventId/editions/$editionId/products")({
  component: ProductsPage,
});

function ProductsPage() {
  const { eventId, editionId } = Route.useParams();
  const { totalCents } = useCart(editionId);
  const [isCartOpen, setIsCartOpen] = useState(false);

  const totalFormatted = new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(totalCents / 100);

  const { data: products = [], isLoading } = useQuery(
    allProductsQueryOptions(eventId, editionId)
  );

  return (
    <div className="min-h-screen pb-20">
      <div className="max-w-7xl mx-auto px-4 pt-8 pb-4 flex items-center justify-end">
        <Button
          onClick={() => { setIsCartOpen(true); }}
          size="lg"
          className={cn(
            "rounded-full bg-accent hover:bg-accent/90 text-white shadow-xl",
            "shadow-accent/20 px-6 font-bold gap-3"
          )}
        >
          <ShoppingBag className="h-5 w-5" />
          <span>Carrinho</span>
          {totalCents > 0 && (
            <span className="border-l border-white/30 pl-3 font-mono text-sm">
              {totalFormatted}
            </span>
          )}
        </Button>
      </div>

      <ProductList products={products} isLoading={isLoading} />

      {/* Cart Drawer */}
      <Cart
        isOpen={isCartOpen}
        eventId={eventId}
        editionId={editionId}
        onClose={() => { setIsCartOpen(false); }}
      />
    </div>
  );
}