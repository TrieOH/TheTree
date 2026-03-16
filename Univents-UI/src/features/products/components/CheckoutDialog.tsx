import { CheckCircle2, ShoppingCart } from "lucide-react";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/shared/ui/shadcn/dialog";

interface CheckoutDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CheckoutDialog({ isOpen, onClose }: CheckoutDialogProps) {
  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-lg p-0 gap-0 border-none overflow-hidden rounded-3xl">
        <DialogHeader className="pt-10 pb-6 text-center space-y-0">
          <div className="mx-auto bg-primary/10 w-20 h-20 rounded-full flex items-center justify-center mb-4">
            <CheckCircle2 className="h-10 w-10 text-primary" />
          </div>
          <DialogTitle className="text-2xl font-black">Finalizar Pedido</DialogTitle>
          <DialogDescription className="text-muted-foreground text-sm mt-1">
            O checkout será implementado nesta seção.
          </DialogDescription>
        </DialogHeader>

        <div className="px-10 pb-8 space-y-6 text-center">
          <div className="p-6 rounded-2xl bg-muted/50 border border-dashed border-primary/20">
            <p className="text-sm font-medium leading-relaxed text-muted-foreground">
              Informações de pagamento, formulários de contato e integração com gateways de pagamento
              serão adicionados aqui futuramente conforme as necessidades do projeto.
            </p>
          </div>

          {/* Placeholders animados */}
          <div className="space-y-3 opacity-50">
            <div className="h-10 w-full bg-muted/50 rounded-lg animate-pulse" />
            <div className="h-10 w-full bg-muted/50 rounded-lg animate-pulse" />
          </div>
        </div>

        <div className="p-6 pt-0 flex gap-3">
          <Button
            variant="outline"
            className="flex-1 rounded-xl h-14 text-base"
            onClick={onClose}
          >
            <ShoppingCart className="mr-2 h-4 w-4" />
            Voltar ao Carrinho
          </Button>
          <Button
            className="flex-1 rounded-xl h-14 text-base font-bold"
            onClick={() => {
              alert("Obrigado por testar o protótipo!");
              onClose();
            }}
          >
            Confirmar Mock
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}