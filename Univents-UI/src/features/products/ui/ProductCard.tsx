import { ShoppingCart, Package, AlertCircle, ChevronLeft, ChevronRight } from "lucide-react";
import { useState, useRef, useCallback } from "react";
import { useCart } from "../hooks/use-cart";
import type { ProductI } from "../model";
import { Button } from "@/shared/ui/shadcn/button";
import { Badge } from "@/shared/ui/shadcn/badge";
import {
  Card,
  CardContent,
} from "@/shared/ui/shadcn/card";
import { cn } from "@/shared/lib/utils";

interface ProductCardProps {
  product: ProductI;
  inventoryRemaining: number;
  onProductSelect: () => void;
}

interface ClickAnimation {
  id: number;
  x: number;
  y: number;
  value: string;
}

export function ProductCard({ product, inventoryRemaining, onProductSelect }: ProductCardProps) {
  const { items, addItem, isLimitReached: checkLimitReached } = useCart(product.edition_id);
  const [animations, setAnimations] = useState<ClickAnimation[]>([]);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const idCounter = useRef(0);

  const cartItem = items.find(i => i.id === product.id);
  const cartQuantity = cartItem?.quantity ?? 0;

  const isLimitReached = checkLimitReached(product, cartQuantity);
  const canAdd = product.status === "available" && !isLimitReached;

  const priceFormatted = new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(product.price_cents / 100);

  const images = product.thumbnail_url
    ? [product.thumbnail_url]
    : (product.gallery_urls ?? []).filter(Boolean) as string[];

  const hasMultipleImages = images.length > 1;

  const handleNextImage = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIndex((prev) => (prev + 1) % images.length);
  }, [images.length]);

  const handlePrevImage = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIndex((prev) => (prev - 1 + images.length) % images.length);
  }, [images.length]);

  const handleAddClick = useCallback((e: React.MouseEvent | React.TouchEvent) => {
    e.stopPropagation();
    if (!canAdd) return;

    let clientX, clientY;
    if ('touches' in e) {
      clientX = e.touches[0].clientX;
      clientY = e.touches[0].clientY;
    } else {
      clientX = (e).clientX;
      clientY = (e).clientY;
    }

    const newId = ++idCounter.current;
    const newAnim: ClickAnimation = {
      id: newId,
      x: clientX,
      y: clientY,
      value: "+1",
    };

    setAnimations(prev => [...prev, newAnim]);

    setTimeout(() => {
      setAnimations(prev => prev.filter(anim => anim.id !== newId));
    }, 800);

    addItem(
      {
        id: product.id,
        name: product.name,
        price_cents: product.price_cents,
        inventory_remaining: inventoryRemaining,
        has_inventory: product.has_inventory,
      },
      1
    );
  }, [canAdd, addItem, product]);

  const isAvailable = product.status === "available";
  const isOutOfStock = product.has_inventory && inventoryRemaining <= 0;
  const isLowStock = product.has_inventory && inventoryRemaining <= 5 && inventoryRemaining > 0;

  return (
    <>
      <Card
        className={cn(
          "group relative overflow-hidden transition-all cursor-pointer",
          "cursor-pointer duration-200 ease-out p-0!",
          isAvailable ? 'hover:-translate-y-1 hover:shadow-lg hover:border-primary/50' : 'opacity-75'
        )}
        onClick={onProductSelect}
      >
        <div className="absolute top-2 right-2 z-30">
          <Button
            size="icon"
            className={cn(
              "h-9 w-9 rounded-full shadow-md transition-all duration-200",
              cartQuantity > 0
                ? 'bg-primary text-primary-foreground hover:bg-primary/90'
                : 'bg-background/90 backdrop-blur-sm text-foreground hover:bg-primary hover:text-primary-foreground border border-border/50',
              !canAdd && 'opacity-50 cursor-not-allowed',
              canAdd && 'hover:scale-110 active:scale-95'
            )}
            onClick={handleAddClick}
            disabled={!canAdd}
          >
            <div className="relative">
              <ShoppingCart className="h-4 w-4" />
              {cartQuantity > 0 && (
                <span className={cn(
                  "absolute -top-2 -right-2 bg-destructive text-destructive-foreground",
                  "text-[10px] font-bold min-w-4 h-4 rounded-full flex items-center justify-center px-1 shadow-sm"
                )}>
                  {cartQuantity}
                </span>
              )}
            </div>
          </Button>
        </div>

        {/* Image */}
        <div className="relative aspect-video bg-secondary/30 overflow-hidden shrink-0">
          {images.length > 0 ? (
            <img
              src={images[currentImageIndex]}
              alt={product.name}
              className="absolute inset-0 w-full h-full object-cover transition-transform duration-200 ease-out group-hover:scale-105"
            />
          ) : (
            <div className="absolute inset-0 flex items-center justify-center transition-transform duration-200 ease-out group-hover:scale-105">
              <Package className="w-10 h-10 text-muted-foreground/40" />
            </div>
          )}

          {hasMultipleImages && (
            <div className="absolute inset-0 z-20 opacity-0 group-hover:opacity-100 transition-opacity">
              <Button
                variant="ghost"
                size="icon"
                className="absolute left-1 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-background/30 backdrop-blur-sm hover:bg-background/60"
                onClick={handlePrevImage}
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-background/30 backdrop-blur-sm hover:bg-background/60"
                onClick={handleNextImage}
              >
                <ChevronRight className="h-4 w-4" />
              </Button>
              <div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5">
                {images.map((_, i) => (
                  <div
                    key={i}
                    className={cn(
                      "w-1.5 h-1.5 rounded-full transition-all duration-300 shadow-sm",
                      i === currentImageIndex
                        ? "bg-primary scale-125"
                        : "bg-background/60"
                    )}
                  />
                ))}
              </div>
            </div>
          )}

          {/* Status Overlays */}
          {!isAvailable || isOutOfStock ? (
            <div className="absolute inset-0 bg-background/80 backdrop-blur-[1px] flex items-center justify-center z-10">
              <div className="bg-destructive text-destructive-foreground px-3 py-1 rounded-md font-bold text-xs -rotate-12 shadow-sm border-2 border-destructive-foreground/20">
                {isOutOfStock ? "SEM ESTOQUE" : "INDISPONÍVEL"}
              </div>
            </div>
          ) : isLimitReached && (
            <div className="absolute inset-0 bg-background/60 backdrop-blur-[1px] flex items-center justify-center z-10">
              <div className="bg-accent text-accent-foreground px-3 py-1 rounded-md font-bold text-xs shadow-sm border-2 border-accent-foreground/20">
                NO CARRINHO (LIMITE)
              </div>
            </div>
          )}

          <Badge
            variant={isAvailable && !isOutOfStock ? "default" : "destructive"}
            className={`absolute top-2 left-2 text-[10px] px-2 py-0.5 ${isAvailable && !isOutOfStock ? "bg-primary text-primary-foreground" : ""}`}
          >
            {isAvailable && !isOutOfStock ? "Disponível" : isOutOfStock ? "Esgotado" : "Indisponível"}
          </Badge>

          {isLowStock && (
            <div className="absolute bottom-0 left-0 right-0 bg-accent px-2 py-0.5 flex items-center justify-center gap-1">
              <AlertCircle className="w-3 h-3 text-accent-foreground" />
              <p className="text-accent-foreground text-[10px] font-semibold">
                Apenas {inventoryRemaining} restantes!
              </p>
            </div>
          )}
        </div>

        <CardContent className="px-3 pb-3 pt-0">
          <div className="flex items-center justify-between gap-2">
            <h3 className="font-semibold text-sm leading-tight line-clamp-1 flex-1 min-w-0 group-hover:text-primary transition-colors duration-200">
              {product.name}
            </h3>
            <span className="text-base font-bold text-accent whitespace-nowrap tabular-nums">
              {priceFormatted}
            </span>
          </div>

          <div className="flex items-start justify-between gap-2 mt-1">
            <p className="text-xs text-muted-foreground line-clamp-2 flex-1 min-w-0 leading-relaxed" title={product.description ?? undefined}>
              {product.description ?? "Nenhuma descrição disponível"}
            </p>

            {(product.has_inventory || !isAvailable) && isAvailable && (
              <span className={`text-[10px] whitespace-nowrap ml-2 mt-0.5 ${isLowStock ? 'text-amber-600 font-semibold' : 'text-muted-foreground'}`}>
                {inventoryRemaining} un.
              </span>
            )}
          </div>
        </CardContent>
      </Card>

      {animations.map((anim) => (
        <div
          key={anim.id}
          className="fixed pointer-events-none z-50 font-bold text-lg text-accent text-shadow-accent/50 select-none"
          style={{
            left: anim.x,
            top: anim.y,
            animation: 'clicker-float 0.8s ease-out forwards',
          }}
        >
          {anim.value}
        </div>
      ))}

      <style>{`
        @keyframes clicker-float {
          0% {
            opacity: 1;
            transform: translate(-50%, -50%) scale(0.5);
          }
          20% {
            transform: translate(-50%, -50%) scale(1.2);
          }
          100% {
            opacity: 0;
            transform: translate(-50%, -150%) scale(1);
          }
        }
      `}</style>
    </>
  );
}