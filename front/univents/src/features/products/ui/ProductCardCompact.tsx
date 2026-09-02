import {
  Check,
  ChevronLeft,
  ChevronRight,
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

interface ProductCardCompactProps {
  product: ProductI & { name?: string };
  variants: VariantI[];
  editionId?: string;
}

export function ProductCardCompact({
  product,
  variants,
  editionId,
}: ProductCardCompactProps) {
  const { addItem, items } = useCart(editionId ?? "");
  const [selectedId, setSelectedId] = useState<string>(variants[0]?.id ?? "");
  const [imageIndexes, setImageIndexes] = useState<Record<string, number>>({});

  if (variants.length === 0) return null;

  const selected = variants.find((v) => v.id === selectedId) ?? variants[0];
  const inCart = items.find(
    (item) => item.id === selected.id && item.type === "product",
  );
  const isOutOfStock = selected.stock !== null && selected.stock <= 0;
  const isLowStock =
    selected.stock !== null && selected.stock > 0 && selected.stock <= 5;

  const currentImageIndex = imageIndexes[selected.id] ?? 0;
  const hasMultipleImages = (selected.gallery_urls?.length ?? 0) > 1;
  const totalImages = selected.gallery_urls?.length ?? 0;

  const nextImage = () => {
    setImageIndexes((prev) => ({
      ...prev,
      [selected.id]: ((prev[selected.id] ?? 0) + 1) % totalImages,
    }));
  };

  const prevImage = () => {
    setImageIndexes((prev) => ({
      ...prev,
      [selected.id]: ((prev[selected.id] ?? 0) - 1 + totalImages) % totalImages,
    }));
  };

  return (
    <div className="flex w-88 max-w-full flex-col overflow-hidden rounded-xl border border-border bg-card">
      {/* Image */}
      <div className="group relative flex aspect-16/8 w-full shrink-0 items-center justify-center overflow-hidden bg-muted">
        {selected.gallery_urls?.[currentImageIndex] ? (
          <>
            <img
              src={selected.gallery_urls[currentImageIndex]}
              alt={selected.name}
              className="w-full h-full object-cover"
            />
            {hasMultipleImages && (
              <>
                <button
                  type="button"
                  onClick={(event) => {
                    event.stopPropagation();
                    prevImage();
                  }}
                  className="absolute left-2 top-1/2 flex h-7 w-7 -translate-y-1/2 items-center justify-center rounded-full bg-background/70 opacity-100 backdrop-blur-sm transition-opacity sm:opacity-0 sm:group-hover:opacity-100"
                >
                  <ChevronLeft className="w-3.5 h-3.5" />
                </button>
                <button
                  type="button"
                  onClick={(event) => {
                    event.stopPropagation();
                    nextImage();
                  }}
                  className="absolute right-2 top-1/2 flex h-7 w-7 -translate-y-1/2 items-center justify-center rounded-full bg-background/70 opacity-100 backdrop-blur-sm transition-opacity sm:opacity-0 sm:group-hover:opacity-100"
                >
                  <ChevronRight className="w-3.5 h-3.5" />
                </button>
                <div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5">
                  {selected.gallery_urls.map((url, idx) => (
                    <div
                      key={url}
                      className={cn(
                        "w-1.5 h-1.5 rounded-full transition-colors",
                        idx === currentImageIndex ? "bg-white" : "bg-white/40",
                      )}
                    />
                  ))}
                </div>
              </>
            )}
          </>
        ) : (
          <Package className="w-10 h-10 text-muted-foreground/40" />
        )}

        {product.requires_registration && (
          <Tooltip>
            <TooltipTrigger
              aria-label="Necessita cadastro no evento"
              className="absolute right-2 top-2 rounded-full bg-primary/10 p-2 text-primary cursor-help transition-colors hover:bg-primary/20 focus:outline-none focus:ring-2 focus:ring-primary/30"
            >
              <UserCheck className="h-4 w-4" />
            </TooltipTrigger>
            <TooltipContent side="bottom" align="end">
              Necessita cadastro no evento
            </TooltipContent>
          </Tooltip>
        )}
      </div>

      {/* Content */}
      <div className="flex flex-col px-3 pt-3">
        <h3 className="text-sm font-semibold text-foreground leading-snug">
          {product.name ?? product.vendor_code}: {selected.name}
        </h3>
        <p className="mt-0.5 min-h-8 text-[11px] leading-4 text-muted-foreground line-clamp-2">
          {selected.description}
        </p>

        {/* Price + Stock */}
        <div className="mt-2 pt-2 border-t border-border flex items-center justify-between">
          <span
            className={cn(
              "text-base font-bold tabular-nums",
              selected.price === 0
                ? "text-emerald-600 dark:text-emerald-400"
                : "text-foreground",
            )}
          >
            {formatPrice(selected.price)}
          </span>

          <span
            className={cn(
              "flex items-center gap-1 text-[11px]",
              isOutOfStock
                ? "text-destructive font-medium"
                : isLowStock
                  ? "text-amber-600 dark:text-amber-400 font-medium"
                  : "text-muted-foreground",
            )}
          >
            <Package className="w-3.5 h-3.5" />
            {selected.stock === null
              ? "Ilimitado"
              : isOutOfStock
                ? "Esgotado"
                : isLowStock
                  ? `${selected.stock} restantes`
                  : `${selected.stock} em estoque`}
          </span>
        </div>
      </div>

      {/* Variant chips */}
      <div className="min-h-16 px-3 pt-2">
        <p className="mb-1 text-[10px] font-medium uppercase tracking-wider text-muted-foreground">
          Opções
        </p>
        <div className="flex flex-wrap gap-1.5">
          {variants.map((variant) => {
            const vOut = variant.stock !== null && variant.stock <= 0;
            const isSelected = selected.id === variant.id;

            return (
              <button
                key={variant.id}
                type="button"
                onClick={() => !vOut && setSelectedId(variant.id)}
                disabled={vOut}
                className={cn(
                  "rounded-md border px-2.5 py-1 text-xs font-medium transition-all",
                  vOut
                    ? "cursor-not-allowed border-border bg-muted/30 text-muted-foreground opacity-35 line-through"
                    : isSelected
                      ? "border-foreground/30 bg-muted/30 text-foreground"
                      : "border-border bg-card text-foreground hover:bg-muted/20",
                )}
              >
                {variant.name}
              </button>
            );
          })}
        </div>
      </div>

      {/* Buy bar */}
      {editionId && !isOutOfStock && (
        <div className="flex items-center gap-2 px-3 py-3 mt-auto">
          <Button
            size="default"
            variant={inCart ? "secondary" : "default"}
            className={cn(
              "flex-1 h-9 gap-1.5 text-xs font-semibold transition-all",
              inCart && "bg-emerald-600 hover:bg-emerald-700 text-white",
            )}
            onClick={() =>
              addItem(
                {
                  id: selected.id,
                  type: "product",
                  name: `${product.vendor_code} - ${selected.name}`,
                  price_cents: selected.price,
                  inventory_remaining: selected.stock ?? 999,
                  has_inventory: selected.stock !== null,
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
                Adicionar
              </>
            )}
          </Button>
        </div>
      )}

      {isOutOfStock && (
        <div className="px-3 py-3 mt-auto">
          <Button
            size="default"
            variant="outline"
            disabled
            className="h-9 w-full text-xs font-semibold opacity-60 cursor-not-allowed"
          >
            Indisponível
          </Button>
        </div>
      )}
    </div>
  );
}
