import { createFileRoute } from "@tanstack/react-router";
import { Wifi, AlertTriangle } from "lucide-react";
import { useState, useEffect } from "react";
import { useCart } from "@/features/products/hooks/use-cart";
import { PaymentProviderSelector } from "@/features/payments/ui/PaymentProviderSelector";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  CardFooter,
} from "@/shared/ui/shadcn/card";

export const Route = createFileRoute(
  "/events/$eventId/editions/$editionId/checkout",
)({
  component: CheckoutPage,
});

function OrderSummary({ editionId }: { editionId: string }) {
  const { items, totalCents } = useCart(editionId);
  const [timeLeft, setTimeLeft] = useState(15 * 60); // 15 minutes in seconds

  useEffect(() => {
    if (timeLeft === 0) return;

    const timer = setInterval(() => {
      setTimeLeft((prevTime) => prevTime - 1);
    }, 1000);

    return () => { clearInterval(timer); };
  }, [timeLeft]);

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes.toString().padStart(2, "0")}:${remainingSeconds
      .toString()
      .padStart(2, "0")}`;
  };

  const priceFormatted = (cents: number) =>
    new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(cents / 100);

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Resumo do Pedido</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {items.map((item) => (
            <div key={item.id} className="flex justify-between items-start">
              <div>
                <p className="font-bold">
                  {item.name}{" "}
                  <span className="font-normal text-muted-foreground">
                    x{item.quantity}
                  </span>
                </p>
                <p className="text-sm text-muted-foreground">
                  {priceFormatted(item.price_cents)}
                </p>
              </div>
              <p className="font-bold text-lg">
                {priceFormatted(item.price_cents * item.quantity)}
              </p>
            </div>
          ))}
        </CardContent>
        <CardFooter className="flex justify-between items-baseline bg-muted/50 py-4 px-6">
          <span className="text-lg font-bold">Total</span>
          <span className="text-2xl font-black text-primary">
            {priceFormatted(totalCents)}
          </span>
        </CardFooter>
      </Card>

      <Card className="border-amber-500/50 border-2 bg-amber-500/10 text-amber-900 dark:text-amber-200">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-lg">Sua reserva expira em</CardTitle>
          <AlertTriangle className="w-5 h-5 text-amber-600" />
        </CardHeader>
        <CardContent>
          <div className="text-5xl font-mono font-bold text-center tracking-tight">
            {formatTime(timeLeft)}
          </div>
          <CardDescription className="text-center mt-2 text-amber-800 dark:text-amber-300">
            Seus itens estão reservados. Finalize a compra antes que o tempo
            acabe.
          </CardDescription>
        </CardContent>
      </Card>
    </div>
  );
}

function CheckoutPage() {
  const { editionId } = Route.useParams();
  const { totalCents } = useCart(editionId);

  return (
    <div className="container max-w-4xl mx-auto py-12 px-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
        {/* Left Column: Order Summary & Timer */}
        <OrderSummary editionId={editionId} />

        {/* Right Column: Payment */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Pagamento</CardTitle>
              <CardDescription>
                Selecione seu método de pagamento preferido.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <PaymentProviderSelector amount={totalCents} />
            </CardContent>
          </Card>
          <div className="flex items-center justify-end gap-2 text-xs text-muted-foreground">
            <Wifi size={14} className="text-green-500" />
            <span>Conexão segura em tempo real estabelecida.</span>
          </div>
        </div>
      </div>
    </div>
  );
}
