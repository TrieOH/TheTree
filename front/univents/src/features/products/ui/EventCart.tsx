import { ShoppingCart } from "lucide-react";
import { useState } from "react";
import { Button } from "@/shared/ui/shadcn/button";
import { useCart } from "../hooks/use-cart";
import { Cart } from "./Cart";

interface EventCartProps {
  editionId: string;
  onCheckout: () => void;
  onExplore?: () => void;
}

export function EventCart({
  editionId,
  onCheckout,
  onExplore,
}: EventCartProps) {
  const [isOpen, setIsOpen] = useState(false);
  const { itemCount } = useCart(editionId);

  return (
    <>
      <Cart
        isOpen={isOpen}
        editionId={editionId}
        onClose={() => setIsOpen(false)}
        onCheckout={onCheckout}
        onExplore={onExplore}
      />
      <Button
        type="button"
        onClick={() => setIsOpen(true)}
        className="fixed bottom-24 right-4 z-40 h-13 rounded-full px-5 shadow-md shadow-primary/10 transition-transform hover:scale-105 md:right-8"
        aria-label="Abrir carrinho"
      >
        <ShoppingCart className="mr-2 h-5 w-5" />
        <span className="hidden sm:inline">Carrinho</span>
        {itemCount > 0 && (
          <span className="ml-2 rounded-full bg-background px-2 py-0.5 text-xs font-bold text-foreground">
            {itemCount}
          </span>
        )}
      </Button>
    </>
  );
}
