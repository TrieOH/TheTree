import { InfinityIcon, Package, UserCheck } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import type { ProductI, VariantI } from "../model";

function formatPrice(price: number): string {
  if (price === 0) return "Gratuito";
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
    minimumFractionDigits: 0,
  }).format(price);
}

interface ProductCardProps {
  product: ProductI;
  variants: VariantI[];
  maxVariants?: number;
}

export function ProductCard({
  product,
  variants,
  maxVariants = 3,
}: ProductCardProps) {
  if (variants.length === 0) return null;

  const displayVariants = variants.slice(0, maxVariants);
  const remaining = variants.length - maxVariants;

  return (
    <div
      className={cn(
        "flex flex-col w-72 rounded-2xl border bg-card overflow-hidden",
        "transition-all duration-200",
        "hover:border-border hover:shadow-md",
        product.requires_registration
          ? "border-primary/30"
          : "border-border/60",
      )}
    >
      {/* Header */}
      <div
        className={cn(
          "flex items-center gap-3 px-5 py-4",
          product.requires_registration && "bg-primary/3",
        )}
      >
        <div
          className={cn(
            "shrink-0 w-9 h-9 rounded-xl flex items-center justify-center",
            product.requires_registration
              ? "bg-primary/15 text-primary"
              : "bg-muted text-muted-foreground",
          )}
        >
          <Package className="w-4 h-4" />
        </div>

        <div className="min-w-0 flex-1">
          <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider truncate">
            {product.vendor_code}
          </p>
        </div>

        {product.requires_registration && (
          <span
            className={cn(
              "shrink-0 inline-flex items-center gap-1",
              "px-2 py-0.5 rounded-full",
              "bg-primary/10 text-primary",
              "text-[10px] font-bold uppercase tracking-wide",
            )}
          >
            <UserCheck className="w-3 h-3" />
            Cadastro
          </span>
        )}
      </div>

      {/* Divider */}
      <div className="h-px bg-border/60 mx-5" />

      {/* Variants */}
      <div className="flex flex-col">
        {displayVariants.map((variant, index) => (
          <div
            key={variant.id}
            className={cn(
              "px-5 py-3.5",
              index !== displayVariants.length - 1 &&
                "border-b border-border/40",
            )}
          >
            <div className="flex items-start justify-between gap-2">
              <div className="min-w-0 flex-1">
                <p className="text-sm font-semibold text-foreground leading-snug">
                  {variant.name}
                </p>
                {variant.description && (
                  <p className="mt-1 text-xs text-muted-foreground line-clamp-2">
                    {variant.description}
                  </p>
                )}
              </div>
              <span
                className={cn(
                  "shrink-0 text-sm font-bold",
                  variant.price === 0
                    ? "text-emerald-600 dark:text-emerald-400"
                    : "text-foreground",
                )}
              >
                {formatPrice(variant.price)}
              </span>
            </div>

            {/* Stock */}
            <div className="mt-2 flex items-center gap-1.5">
              {variant.stock === null ? (
                <>
                  <InfinityIcon className="w-3 h-3 text-muted-foreground" />
                  <span className="text-[11px] font-medium text-muted-foreground">
                    Ilimitado
                  </span>
                </>
              ) : (
                <span
                  className={cn(
                    "text-[11px] font-medium",
                    variant.stock <= 5
                      ? "text-amber-600 dark:text-amber-400"
                      : "text-emerald-600 dark:text-emerald-400",
                  )}
                >
                  {variant.stock} em estoque
                </span>
              )}
            </div>
          </div>
        ))}

        {remaining > 0 && (
          <div className="px-5 py-3 text-center border-t border-border/40">
            <span className="text-xs text-muted-foreground font-medium">
              +{remaining} {remaining === 1 ? "variante" : "variantes"}
            </span>
          </div>
        )}
      </div>
    </div>
  );
}
