import { ShoppingCart, Plus, Minus, Package, AlertCircle } from "lucide-react";
import { useState } from "react";
import { useCart } from "../hooks/use-cart";
import type { ProductI } from "../model";
import { Button } from "@/shared/ui/shadcn/button";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Input } from "@/shared/ui/shadcn/input";

import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
  CardDescription
} from "@/shared/ui/shadcn/card";


interface ProductCardProps {
  product: ProductI;
}

export function ProductCard({ product }: ProductCardProps) {
  const { items, addItem } = useCart(product.edition_id);
  const [quantity, setQuantity] = useState(1);

  const cartItem = items.find(i => i.id === product.id);
  const cartQuantity = cartItem?.quantity ?? 0;
  const maxSelectable = product.has_inventory
    ? Math.max(0, product.inventory_remaining - cartQuantity)
    : 999;

  const priceFormatted = new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(product.price_cents / 100);

  const handleAddToCart = () => {
    if (maxSelectable <= 0) return;

    addItem(
      {
        id: product.id,
        name: product.name,
        price_cents: product.price_cents,
        inventory_remaining: product.inventory_remaining,
        has_inventory: product.has_inventory,
      },
      Math.min(quantity, maxSelectable)
    );
    setQuantity(1);
  };

  const isAvailable = product.status === "available";
  const isOutOfStock = product.has_inventory && product.inventory_remaining <= 0;
  const isLimitReached = product.has_inventory && cartQuantity >= product.inventory_remaining;

  const isLowStock = product.has_inventory && product.inventory_remaining <= 5 && product.inventory_remaining > 0;
  const showInventory = product.has_inventory || !isAvailable;

  return (
    <Card className={`group transition-all p-0! duration-200 ease-out hover:-translate-y-1 hover:shadow-lg ${isAvailable ? 'hover:border-primary/50' : 'opacity-75'}`} size="sm">
      {/* Image Section */}
      <div className="relative aspect-video bg-secondary/30 overflow-hidden shrink-0">
        <div className="absolute inset-0 flex items-center justify-center transition-transform duration-200 ease-out group-hover:scale-105">
          <Package className="w-10 h-10 text-muted-foreground/40" />
        </div>

        {!isAvailable || isOutOfStock ? (
          <div className="absolute inset-0 bg-background/80 backdrop-blur-[1px] flex items-center justify-center z-20">
            <div className="bg-destructive text-destructive-foreground px-3 py-1 rounded-md font-bold text-xs -rotate-12 shadow-sm border-2 border-destructive-foreground/20">
              {isOutOfStock ? "SEM ESTOQUE" : "INDISPONÍVEL"}
            </div>
          </div>
        ) : isLimitReached && (
          <div className="absolute inset-0 bg-background/60 backdrop-blur-[1px] flex items-center justify-center z-20">
            <div className="bg-accent text-accent-foreground px-3 py-1 rounded-md font-bold text-xs shadow-sm border-2 border-accent-foreground/20">
              NO CARRINHO (LIMITE)
            </div>
          </div>
        )}

        <Badge
          variant={isAvailable && !isOutOfStock ? "default" : "destructive"}
          className={`absolute top-2 right-2 text-[10px] px-2 py-0.5 ${isAvailable && !isOutOfStock ? "bg-primary text-primary-foreground" : ""}`}
        >
          {isAvailable && !isOutOfStock ? "Disponível" : isOutOfStock ? "Esgotado" : "Indisponível"}
        </Badge>

        {isLowStock && (
          <div className="absolute bottom-0 left-0 right-0 bg-accent px-2 py-0.5 flex items-center justify-center gap-1">
            <AlertCircle className="w-3 h-3 text-accent-foreground" />
            <p className="text-accent-foreground text-[10px] font-semibold">
              Apenas {product.inventory_remaining} restantes!
            </p>
          </div>
        )}
      </div>

      <CardHeader className="p-2 pb-0 gap-0.5">
        <CardTitle className="group-hover:text-primary transition-colors duration-200">
          {product.name}
        </CardTitle>
        <CardDescription className="line-clamp-2 text-xs">
          {product.description ?? "Nenhuma descrição disponível."}
        </CardDescription>
      </CardHeader>

      <CardContent className="p-2 pt-1 pb-0">
        <div className="flex items-center justify-between">
          <span className="text-lg font-bold text-accent">{priceFormatted}</span>
          {showInventory && isAvailable && product.has_inventory && (
            <span className="text-[10px] text-muted-foreground">{product.inventory_remaining} un.</span>
          )}
        </div>
      </CardContent>

      <CardFooter className="p-2 pt-1.5 flex-col gap-1.5 border-t-0 bg-transparent">
        {/* Quantity Selector */}
        <div className="flex items-center gap-1 w-full">
          <Button
            variant="outline"
            size="icon"
            className="h-8 w-8 shrink-0"
            onClick={() => { setQuantity(Math.max(1, quantity - 1)) }}
            disabled={quantity <= 1 || !isAvailable || isLimitReached}
          >
            <Minus className="h-3 w-3" />
          </Button>

          <Input
            type="number"
            min={1}
            max={maxSelectable}
            value={isLimitReached ? 0 : quantity}
            disabled={!isAvailable || isLimitReached}
            onChange={(e) => {
              const val = parseInt(e.target.value) || 1;
              setQuantity(Math.min(Math.max(1, val), maxSelectable));
            }}
            className="text-center h-8 text-sm [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
          />

          <Button
            variant="outline"
            size="icon"
            className="h-8 w-8 shrink-0"
            onClick={() => { setQuantity(Math.min(quantity + 1, maxSelectable)) }}
            disabled={!isAvailable || quantity >= maxSelectable || isLimitReached}
          >
            <Plus className="h-3 w-3" />
          </Button>
        </div>

        {/* Add to Cart Button */}
        <Button
          className="w-full h-10 font-bold text-xs"
          onClick={handleAddToCart}
          disabled={!isAvailable || isLimitReached || isOutOfStock}
        >
          <ShoppingCart className="h-3.5 w-3.5 -mt-1" />
          {isOutOfStock ? "Esgotado" : isLimitReached ? "No Carrinho" : "Adicionar"}
        </Button>
      </CardFooter>
    </Card>
  );
}