import {
  Check,
  ChevronLeft,
  ChevronRight,
  InfinityIcon,
  Package,
  ShoppingCart,
  UserCheck,
} from "lucide-react";
import { useState } from "react";
import { formatPrice } from "@/shared/lib/money";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/shared/ui/shadcn/tooltip";
import { useCart } from "../hooks/use-cart";
import type { ProductI, VariantI } from "../model";

interface ProductCardProps {
  product: ProductI;
  variants: VariantI[];
  maxVariants?: number;
  editionId?: string;
  onViewAllVariants?: (product: ProductI, variants: VariantI[]) => void;
}

export function ProductCard({
  product,
  variants,
  maxVariants = 3,
  editionId,
  onViewAllVariants,
}: ProductCardProps) {
  const { addItem, items } = useCart(editionId ?? "");
  const [imageIndexes, setImageIndexes] = useState<Record<string, number>>({});

  if (variants.length === 0) return null;

  const displayVariants = variants.slice(0, maxVariants);
  const remaining = variants.length - maxVariants;
  const hasMoreVariants = remaining > 0;

  const getInCart = (variantId: string) =>
    items.find((item) => item.id === variantId && item.type === "product");

  const nextImage = (variantId: string, totalImages: number) => {
    setImageIndexes((prev) => ({
      ...prev,
      [variantId]: ((prev[variantId] ?? 0) + 1) % totalImages,
    }));
  };

  const prevImage = (variantId: string, totalImages: number) => {
    setImageIndexes((prev) => ({
      ...prev,
      [variantId]: ((prev[variantId] ?? 0) - 1 + totalImages) % totalImages,
    }));
  };

  return (
    <div className="w-88 max-w-full space-y-2">
      <div className="relative flex h-10 items-center border-b border-border/50 px-1">
        <p className="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">
          {product.vendor_code}
        </p>

        {product.requires_registration && (
          <Tooltip>
            <TooltipTrigger
              aria-label="Necessita cadastro no evento"
              className="absolute right-0 top-1/2 z-10 -translate-y-1/2 cursor-help rounded-full bg-card p-2 text-primary transition-colors hover:bg-primary/20 focus:outline-none focus:ring-2 focus:ring-primary/30"
            >
              <UserCheck className="h-3 w-3" />
            </TooltipTrigger>
            <TooltipContent side="top" className="font-medium">
              Necessita cadastro no evento
            </TooltipContent>
          </Tooltip>
        )}
      </div>

      {/* Variants - each is its own card */}
      <div className="flex flex-col gap-2">
        {displayVariants.map((variant) => {
          const inCart = getInCart(variant.id);
          const currentImageIndex = imageIndexes[variant.id] ?? 0;
          const hasMultipleImages = (variant.gallery_urls?.length ?? 0) > 1;
          const isOutOfStock = variant.stock !== null && variant.stock <= 0;
          const isLowStock =
            variant.stock !== null && variant.stock > 0 && variant.stock <= 5;

          return (
            <div
              key={variant.id}
              className="group rounded-xl border border-border/60 bg-card p-4"
            >
              {/* Fixed-height top row: image + info + price/stock */}
              <div className="flex items-start gap-3 h-16">
                {/* Image - fixed 64x64 */}
                <div className="relative shrink-0">
                  {variant.gallery_urls?.[currentImageIndex] ? (
                    <div className="relative size-16 rounded-md overflow-hidden ring-1 ring-border/60 bg-muted">
                      <img
                        src={variant.gallery_urls[currentImageIndex]}
                        alt={variant.name}
                        className="size-full object-cover"
                      />
                      {hasMultipleImages && (
                        <>
                          <button
                            type="button"
                            onClick={(e) => {
                              e.stopPropagation();
                              prevImage(
                                variant.id,
                                variant.gallery_urls?.length ?? 0,
                              );
                            }}
                            className="absolute left-0.5 top-1/2 -translate-y-1/2 w-5 h-5 rounded-full bg-black/40 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/60"
                          >
                            <ChevronLeft className="w-3 h-3" />
                          </button>
                          <button
                            type="button"
                            onClick={(e) => {
                              e.stopPropagation();
                              nextImage(
                                variant.id,
                                variant.gallery_urls?.length ?? 0,
                              );
                            }}
                            className="absolute right-0.5 top-1/2 -translate-y-1/2 w-5 h-5 rounded-full bg-black/40 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/60"
                          >
                            <ChevronRight className="w-3 h-3" />
                          </button>
                          <div className="absolute bottom-1 left-1/2 -translate-x-1/2 flex gap-0.5">
                            {variant.gallery_urls.map((url, idx) => (
                              <div
                                key={url}
                                className={cn(
                                  "w-1 h-1 rounded-full transition-colors",
                                  idx === currentImageIndex
                                    ? "bg-white"
                                    : "bg-white/50",
                                )}
                              />
                            ))}
                          </div>
                        </>
                      )}
                    </div>
                  ) : (
                    <div className="size-16 rounded-md bg-muted flex items-center justify-center ring-1 ring-border/60">
                      <Package className="w-6 h-6 text-muted-foreground/50" />
                    </div>
                  )}
                </div>

                {/* Info - fixed 1 line name + 1 line desc */}
                <div className="min-w-0 flex-1 flex flex-col justify-center h-16">
                  <p className="text-sm font-semibold text-foreground leading-5 truncate">
                    {variant.name}
                  </p>
                  {variant.description ? (
                    <p className="text-xs text-muted-foreground leading-5 truncate">
                      {variant.description}
                    </p>
                  ) : (
                    <div className="h-5" />
                  )}
                </div>

                {/* Price + Stock */}
                <div className="shrink-0 flex flex-col justify-center h-16 text-right">
                  <span
                    className={cn(
                      "text-sm font-bold tabular-nums leading-5",
                      variant.price === 0
                        ? "text-emerald-600 dark:text-emerald-400"
                        : "text-foreground",
                    )}
                  >
                    {formatPrice(variant.price)}
                  </span>
                  <div className="mt-0.5 flex items-center justify-end gap-1 text-[11px] leading-5">
                    {variant.stock === null ? (
                      <InfinityIcon
                        className="w-3.5 h-3.5 text-muted-foreground"
                        aria-label="Estoque ilimitado"
                      />
                    ) : isOutOfStock ? (
                      <span className="text-destructive font-medium">
                        Esgotado
                      </span>
                    ) : (
                      <>
                        <Package
                          className={cn(
                            "w-3.5 h-3.5",
                            isLowStock
                              ? "text-amber-600 dark:text-amber-400"
                              : "text-muted-foreground",
                          )}
                          aria-label={`${variant.stock} em estoque`}
                        />
                        <span
                          className={cn(
                            isLowStock
                              ? "text-amber-600 dark:text-amber-400 font-medium"
                              : "text-muted-foreground",
                          )}
                        >
                          {variant.stock}
                        </span>
                      </>
                    )}
                  </div>
                </div>
              </div>

              {/* Cart button */}
              {editionId && !isOutOfStock && (
                <Button
                  size="sm"
                  variant={inCart ? "secondary" : "default"}
                  className={cn(
                    "mt-3 h-9 w-full gap-2 text-xs font-semibold transition-all",
                    inCart && "bg-emerald-600 hover:bg-emerald-700 text-white",
                  )}
                  onClick={() =>
                    addItem(
                      {
                        id: variant.id,
                        type: "product",
                        name: `${product.vendor_code} - ${variant.name}`,
                        price_cents: variant.price,
                        inventory_remaining: variant.stock ?? 999,
                        has_inventory: variant.stock !== null,
                      },
                      1,
                    )
                  }
                >
                  {inCart ? (
                    <>
                      <Check className="h-3.5 w-3.5" />
                      Adicionado ({inCart.quantity})
                    </>
                  ) : (
                    <>
                      <ShoppingCart className="h-3.5 w-3.5" />
                      Adicionar ao carrinho
                    </>
                  )}
                </Button>
              )}

              {isOutOfStock && (
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  disabled
                  className="mt-3 h-9 w-full text-xs font-semibold opacity-60 cursor-not-allowed"
                >
                  Indisponível
                </Button>
              )}
            </div>
          );
        })}

        {hasMoreVariants && (
          <button
            type="button"
            onClick={() => onViewAllVariants?.(product, variants)}
            className="py-2 text-center hover:bg-muted/30 rounded-xl transition-colors"
          >
            <span className="text-xs text-muted-foreground font-medium hover:text-foreground transition-colors">
              +{remaining} {remaining === 1 ? "variante" : "variantes"}{" "}
              disponível{remaining > 1 ? "is" : ""}
            </span>
          </button>
        )}
      </div>
    </div>
  );
}
