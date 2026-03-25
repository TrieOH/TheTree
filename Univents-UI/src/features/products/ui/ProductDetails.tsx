import type { ProductI } from "../model";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
  SheetClose,
} from "@/shared/ui/shadcn/sheet";
import { Button } from "@/shared/ui/shadcn/button";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Package, ShoppingCart } from "lucide-react";
import { useState } from "react";
import { cn } from "@/shared/lib/utils";
import { useCart } from "../hooks/use-cart";

interface ProductDetailsProps {
  product: ProductI | null;
  isOpen: boolean;
  onOpenChange: (isOpen: boolean) => void;
  inventoryRemaining: number;
}

export function ProductDetails({ product, isOpen, onOpenChange, inventoryRemaining }: ProductDetailsProps) {
  const [lastProduct, setLastProduct] = useState<ProductI | null>(product);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  if (product && product.id !== lastProduct?.id) {
    setLastProduct(product);
    setCurrentImageIndex(0);
  }

  const displayProduct = product || lastProduct;

  const {
    items,
    addItem,
    isLimitReached:
    checkLimitReached
  } = useCart(displayProduct?.edition_id ?? '');

  if (!displayProduct) return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-lg p-0 flex flex-col" />
    </Sheet>
  );

  const cartItem = items.find(i => i.id === displayProduct.id);
  const cartQuantity = cartItem?.quantity ?? 0;

  const isLimitReached = checkLimitReached(displayProduct, cartQuantity);
  const canAdd = displayProduct.status === "available" && !isLimitReached;

  const allImages = (displayProduct.gallery_urls ?? []).filter(Boolean) as string[];

  const priceFormatted = new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(displayProduct.price_cents / 100);

  const isAvailable = displayProduct.status === "available";
  const isOutOfStock = displayProduct.has_inventory && inventoryRemaining <= 0;

  const handleAddToCart = () => {
    if (!canAdd) return;
    addItem(
      {
        id: displayProduct.id,
        name: displayProduct.name,
        price_cents: displayProduct.price_cents,
        inventory_remaining: inventoryRemaining,
        has_inventory: displayProduct.has_inventory,
      },
      1
    );
  }

  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-lg p-0 flex flex-col">
        <SheetHeader className="p-4 border-b">
          <SheetTitle className="truncate">{displayProduct.name}</SheetTitle>
          <SheetDescription className="line-clamp-2">{displayProduct.description}</SheetDescription>
        </SheetHeader>
        <div className="flex-1 overflow-y-auto">
          {/* Image Gallery */}
          <div className="relative aspect-video bg-secondary/30">
            {allImages.length > 0 ? (
              <img
                src={allImages[currentImageIndex]}
                alt={displayProduct.name}
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="flex items-center justify-center h-full">
                <Package className="w-16 h-16 text-muted-foreground/40" />
              </div>
            )}
            {allImages.length > 1 && (
              <div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-2">
                {allImages.map((_, index) => (
                  <Button
                    key={index}
                    onClick={() => setCurrentImageIndex(index)}
                    className={cn(
                      "w-2 h-2 rounded-full",
                      currentImageIndex === index ? "bg-primary" : "bg-muted-foreground/50"
                    )}
                  />
                ))}
              </div>
            )}
          </div>

          <div className="p-4 space-y-4">
            {/* Price and Stock */}
            <div className="flex items-center justify-between">
              <span className="text-2xl font-bold text-accent">{priceFormatted}</span>
              <Badge variant={isAvailable && !isOutOfStock ? "default" : "destructive"}>
                {isAvailable && !isOutOfStock ? "Disponível" : isOutOfStock ? "Esgotado" : "Indisponível"}
              </Badge>
            </div>

            {displayProduct.has_inventory && isAvailable && (
              <p className="text-sm text-muted-foreground">
                {inventoryRemaining > 0
                  ? `${inventoryRemaining} unidades restantes`
                  : "Nenhuma unidade restante"}
              </p>
            )}

            {/* Product Details */}
            <div className="border-t pt-4">
              <h3 className="font-semibold mb-2">Detalhes</h3>
              <ul className="text-sm text-muted-foreground space-y-1">
                <li><span className="font-medium text-foreground">Tipo:</span> {displayProduct.type}</li>
                {displayProduct.available_from && <li><span className="font-medium text-foreground">Disponível de:</span> {new Date(displayProduct.available_from).toLocaleDateString()}</li>}
                {displayProduct.available_until && <li><span className="font-medium text-foreground">Disponível até:</span> {new Date(displayProduct.available_until).toLocaleDateString()}</li>}
              </ul>
            </div>
          </div>
        </div>

        <SheetFooter className="p-4 border-t bg-background flex-row gap-2">
          <Button
            onClick={handleAddToCart}
            disabled={!canAdd}
            className="flex-1"
          >
            <ShoppingCart className="w-4 h-4 mr-2" />
            {canAdd ? 'Adicionar ao carrinho' : 'Indisponível'}
          </Button>
          <SheetClose
            render={<Button variant="outline">Fechar</Button>}
          />
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
