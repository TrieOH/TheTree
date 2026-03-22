import { Link } from "@tanstack/react-router";
import { ShoppingCart, X, Trash2, Plus, Minus, CreditCard } from "lucide-react";
import { useState, useEffect, useRef } from "react";
import { useCart } from "../hooks/use-cart";
import type { CartItem as CartItemType } from "../model/cart";
import { Button } from "@/shared/ui/shadcn/button";
import { cn } from "@/shared/lib/utils";

interface CartProps {
  isOpen: boolean;
  eventId: string;
  editionId: string;
  onClose: () => void;
}

interface CartItemProps {
  item: CartItemType;
  onRemove: (id: string) => void;
  onUpdateQuantity: (id: string, quantity: number) => void;
  priceFormatted: (cents: number) => string;
}

function CartItem({ item, onRemove, onUpdateQuantity, priceFormatted, getMaxQuantity }: CartItemProps & { getMaxQuantity: (p: Pick<CartItemType, "has_inventory" | "inventory_remaining">) => number }) {
  const max = getMaxQuantity(item);
  const maxReached = item.quantity >= max;
  const itemTotal = item.price_cents * item.quantity;
  const inputRef = useRef<HTMLInputElement>(null);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let value = parseInt(e.target.value, 10);
    if (!isNaN(value)) {
      if (value < 1) value = 1;
      if (value > max) value = max;
      onUpdateQuantity(item.id, value);
    }
  };

  const handleInputBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value, 10);
    if (isNaN(value) || value < 1) {
      onUpdateQuantity(item.id, 1);
    } else if (value > max) {
      onUpdateQuantity(item.id, max);
    }
  };

  return (
    <div className="group flex gap-3 p-3 bg-secondary/30 hover:bg-secondary/50 transition-colors border-b border-border/50 last:border-b-0">
      <div className="flex-1 min-w-0">
        <h4 className="font-medium text-sm text-foreground truncate leading-tight">
          {item.name}
        </h4>
        <p className="text-xs text-muted-foreground mt-0.5">
          {priceFormatted(item.price_cents)} un
        </p>

        <div className="flex items-center gap-1.5 mt-3">
          <button
            onClick={() => {
              if (item.quantity > 1) {
                onUpdateQuantity(item.id, item.quantity - 1);
              }
            }}
            className={cn(
              "h-8 w-8 flex items-center justify-center bg-background border border-border",
              "hover:bg-accent hover:text-accent-foreground hover:border-accent",
              "active:bg-accent/80 transition-colors select-none",
              "disabled:opacity-50 disabled:cursor-not-allowed"
            )}
            disabled={item.quantity <= 1}
            aria-label="Diminuir quantidade"
          >
            <Minus className="h-4 w-4" />
          </button>

          <input
            ref={inputRef}
            type="number"
            min="1"
            max={max}
            value={item.quantity}
            onChange={handleInputChange}
            onBlur={handleInputBlur}
            className={cn(
              "w-12 h-8 text-center text-sm font-semibold tabular-nums bg-background border border-border",
              "focus:outline-none focus:border-accent focus:ring-1 focus:ring-accent",
              "hover:border-accent/50 transition-colors",
              "[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
            )}
          />

          <button
            onClick={() => {
              if (!maxReached) {
                onUpdateQuantity(item.id, item.quantity + 1);
              }
            }}
            className={cn(
              "h-8 w-8 flex items-center justify-center bg-background border border-border",
              "hover:bg-accent hover:text-accent-foreground hover:border-accent",
              "active:bg-accent/80 transition-colors select-none",
              maxReached && "opacity-50 cursor-not-allowed hover:bg-background hover:border-border"
            )}
            disabled={maxReached}
            aria-label="Aumentar quantidade"
          >
            <Plus className="h-4 w-4" />
          </button>

          {maxReached && (
            <span className="text-[10px] text-destructive font-medium ml-2 uppercase tracking-wide">
              Máx
            </span>
          )}
        </div>
      </div>

      <div className="flex flex-col justify-between items-end gap-2">
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8 text-muted-foreground hover:text-destructive hover:bg-destructive/10 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity"
          onClick={() => { onRemove(item.id); }}
        >
          <Trash2 className="h-4 w-4" />
        </Button>
        <p className="text-sm font-bold text-foreground tabular-nums">
          {priceFormatted(itemTotal)}
        </p>
      </div>
    </div>
  );
}

export function Cart({ isOpen, eventId, editionId, onClose }: CartProps) {
  const [isAnimating, setIsAnimating] = useState(false);
  const [isVisible, setIsVisible] = useState(false);

  const { items, totalCents, removeItem, updateQuantity, clearCart, getMaxQuantity } = useCart(editionId);

  useEffect(() => {
    if (isOpen) {
      setIsVisible(true);
      setTimeout(() => {
        setIsAnimating(true);
      }, 10);
    } else {
      setIsAnimating(false);
      const timer = setTimeout(() => {
        setIsVisible(false);
      }, 300);
      return () => { clearTimeout(timer); };
    }
  }, [isOpen]);

  const handleClose = () => {
    setIsAnimating(false);
    setTimeout(() => {
      onClose();
    }, 300);
  };

  const priceFormatted = (cents: number) =>
    new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(cents / 100);

  if (!isVisible) return null;

  return (
    <div className="fixed inset-0 z-50 flex justify-end pointer-events-none">
      {/* Overlay */}
      <div
        className={cn(
          "fixed inset-0 bg-black/60 pointer-events-auto",
          "transition-opacity duration-300 ease-out",
          isAnimating ? "opacity-100" : "opacity-0"
        )}
        onClick={handleClose}
      />

      {/* Drawer */}
      <div className={cn(
        "relative min-w-[320px] max-w-100 w-full sm:w-90 bg-background h-full shadow-2xl flex flex-col pointer-events-auto",
        "transition-transform duration-300 ease-out",
        isAnimating ? "translate-x-0" : "translate-x-full"
      )}>
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b bg-primary text-primary-foreground">
          <div className="flex items-center gap-3">
            <div className="relative">
              <ShoppingCart className="h-5 w-5" />
              {items.length > 0 && (
                <span className="absolute -top-2 -right-2 flex h-5 w-5 items-center justify-center bg-accent text-accent-foreground text-[10px] font-bold border-2 border-primary">
                  {items.reduce((acc, item) => acc + item.quantity, 0)}
                </span>
              )}
            </div>
            <h2 className="font-semibold text-sm tracking-wide uppercase">Seu Carrinho</h2>
          </div>
          <Button
            variant="secondary"
            size="icon"
            onClick={handleClose}
            className="h-8 w-8 bg-primary-foreground/20 text-primary-foreground hover:bg-primary-foreground hover:text-primary border border-primary-foreground/30"
          >
            <X className="h-5 w-5" />
          </Button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto overscroll-contain bg-background">
          {items.length === 0 ? (
            <div className="h-full flex flex-col items-center justify-center text-muted-foreground px-6 py-12">
              <ShoppingCart className="h-16 w-16 opacity-20 mb-4" />
              <p className="font-medium text-foreground text-sm uppercase tracking-wide">Carrinho vazio</p>
              <p className="text-xs text-muted-foreground mt-1 mb-6">Adicione produtos para começar</p>
              <Button
                variant="outline"
                size="sm"
                onClick={handleClose}
                className="px-6 border-2 hover:bg-accent hover:text-accent-foreground hover:border-accent"
              >
                Explorar Produtos
              </Button>
            </div>
          ) : (
            <div className="divide-y divide-border/50">
              {items.map((item) => (
                <CartItem
                  key={item.id}
                  item={item}
                  onRemove={removeItem}
                  onUpdateQuantity={updateQuantity}
                  priceFormatted={priceFormatted}
                  getMaxQuantity={getMaxQuantity}
                />
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        {items.length > 0 && (
          <div className="border-t bg-secondary/20 p-4 space-y-4">
            <div className="space-y-2">
              <div className="flex justify-between items-center text-xs text-muted-foreground uppercase tracking-wide">
                <span>{items.reduce((acc, i) => acc + i.quantity, 0)} itens</span>
                <span>Subtotal {priceFormatted(totalCents)}</span>
              </div>
              <div className="flex justify-between items-center border-t border-border pt-2">
                <span className="text-sm font-semibold text-foreground uppercase tracking-wide">Total</span>
                <span className="text-2xl font-bold text-primary tabular-nums">
                  {priceFormatted(totalCents)}
                </span>
              </div>
            </div>

            <div className="space-y-2">
              <Link
                to="/events/$eventId/editions/$editionId/checkout"
                params={{ eventId, editionId }}
                onClick={handleClose}
                className={cn(
                  "flex items-center justify-center gap-2 w-full py-3 px-4",
                  "bg-primary text-primary-foreground! font-semibold",
                  "text-sm uppercase rounded-sm transition-colors duration-300",
                  "hover:text-accent-foreground hover:bg-accent"
                )}
              >
                <CreditCard className="h-4 w-4" />
                Finalizar Compra
              </Link>

              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  className={cn(
                    "flex-1 text-xs font-medium uppercase tracking-wide border-2",
                    "border-muted-foreground/30 hover:border-accent",
                    "hover:bg-accent hover:text-accent-foreground",
                    "transition-colors duration-300 rounded-sm"
                  )}
                  onClick={handleClose}
                >
                  Continuar
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  className={cn(
                    "flex-1 text-xs font-medium uppercase tracking-wide border-2",
                    "border-muted-foreground/30 hover:border-destructive",
                    "hover:bg-destructive hover:text-destructive-foreground",
                    "transition-colors duration-300 rounded-sm"
                  )}
                  onClick={clearCart}
                >
                  Limpar
                </Button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}