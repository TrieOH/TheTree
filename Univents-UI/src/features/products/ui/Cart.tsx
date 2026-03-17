import { Link } from "@tanstack/react-router";
import { ShoppingCart, X, Trash2, Plus, Minus, CreditCard } from "lucide-react";
import { useCart } from "../hooks/use-cart";
import { Button } from "@/shared/ui/shadcn/button";
import { CardHeader, CardTitle, CardContent } from "@/shared/ui/shadcn/card";
import { cn } from "@/shared/lib/utils";

interface CartProps {
  isOpen: boolean;
  eventId: string;
  editionId: string;
  onClose: () => void;
}

export function Cart({ isOpen, eventId, editionId, onClose }: CartProps) {
  const { items, totalCents, removeItem, updateQuantity, clearCart } = useCart(editionId);

  const priceFormatted = (cents: number) =>
    new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(cents / 100);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      {/* Overlay */}
      <div
        className="fixed inset-0 bg-black/40 backdrop-blur-sm animate-in fade-in duration-300"
        onClick={onClose}
      />

      {/* Drawer */}
      <div className={cn(
        "relative w-full max-w-md bg-background h-full shadow-2xl flex flex-col animate-in slide-in-from-right duration-300"
      )}>
        <CardHeader className="border-b flex flex-row items-center justify-between py-4 px-6">
          <div className="flex items-center gap-2">
            <ShoppingCart className="h-5 w-5 text-primary" />
            <CardTitle className="text-lg">Seu Carrinho</CardTitle>
          </div>
          <Button variant="ghost" size="icon" onClick={onClose} className="rounded-full">
            <X className="h-5 w-5" />
          </Button>
        </CardHeader>

        <CardContent className="grow overflow-y-auto p-6 space-y-4">
          {items.length === 0 ? (
            <div className="h-full flex flex-col items-center justify-center text-muted-foreground space-y-4">
              <div className="bg-muted/50 p-6 rounded-full">
                <ShoppingCart className="h-12 w-12 opacity-20" />
              </div>
              <div className="text-center">
                <p className="font-medium text-foreground">Seu carrinho está vazio</p>
                <p className="text-sm">Explore nossos produtos e adicione-os aqui.</p>
              </div>
              <Button variant="outline" onClick={onClose}>Continuar Comprando</Button>
            </div>
          ) : (
            <div className="space-y-4">
              {items.map((item) => {
                const maxReached = item.has_inventory && typeof item.inventory_remaining === 'number' && item.quantity >= item.inventory_remaining;

                return (
                  <div key={item.id} className="flex gap-4 p-4 rounded-xl ring-1 ring-foreground/10 bg-card/50 transition-colors hover:bg-card">
                    <div className="grow">
                      <h4 className="font-bold text-primary">{item.name}</h4>
                      <p className="text-sm text-accent font-semibold">{priceFormatted(item.price_cents)}</p>

                      <div className="flex items-center gap-3 mt-3">
                        <Button
                          variant="outline"
                          size="icon"
                          className="rounded-full"
                          onClick={() => { updateQuantity(item.id, item.quantity - 1) }}
                        >
                          <Minus className="h-3.5 w-3.5" />
                        </Button>
                        <span className="text-sm font-bold w-4 text-center">{item.quantity}</span>
                        <Button
                          variant="outline"
                          size="icon"
                          className="rounded-full"
                          onClick={() => { updateQuantity(item.id, item.quantity + 1) }}
                          disabled={maxReached}
                        >
                          <Plus className="h-3.5 w-3.5" />
                        </Button>
                      </div>
                      {maxReached && (
                        <p className="text-[10px] text-destructive mt-1 font-medium">
                          Limite de estoque atingido
                        </p>
                      )}
                    </div>
                    <div className="flex flex-col justify-between items-end">
                      <Button
                        variant="ghost"
                        size="icon"
                        className="text-muted-foreground hover:text-destructive rounded-full"
                        onClick={() => { removeItem(item.id) }}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                      <p className="text-sm font-black text-primary">
                        {priceFormatted(item.price_cents * item.quantity)}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>

        {items.length > 0 && (
          <div className="border-t p-6 space-y-4 bg-muted/30">
            <div className="space-y-1.5">
              <div className="flex justify-between items-center text-sm text-muted-foreground">
                <span>Subtotal</span>
                <span>{priceFormatted(totalCents)}</span>
              </div>
              <div className="flex justify-between items-center text-xl font-black">
                <span>Total</span>
                <span className="text-accent">{priceFormatted(totalCents)}</span>
              </div>
            </div>

            <div className="grid grid-cols-1 gap-3">
              <Link
                to="/events/$eventId/editions/$editionId/checkout"
                params={{
                  eventId,
                  editionId,
                }}
                onClick={onClose}
                className={cn(
                  "flex items-center justify-center font-bold rounded-sm bg-primary",
                  "text-primary-foreground! py-2"
                )}
              >
                <CreditCard className="mr-2 h-5 w-5" />
                Ir para o Pagamento
              </Link>
              <div className="grid grid-cols-2 gap-2">
                <Button variant="ghost" className="text-sm" onClick={onClose}>
                  Continuar
                </Button>
                <Button variant="ghost" className="text-sm text-destructive hover:bg-destructive/10" onClick={clearCart}>
                  Limpar Carrinho
                </Button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
