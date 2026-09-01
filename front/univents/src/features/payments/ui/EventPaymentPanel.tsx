import { CreditCard } from "lucide-react";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Button } from "@/shared/ui/shadcn/button";
import type { PaymentProviderI } from "../model";

const providers: Array<{
  id: PaymentProviderI;
  name: string;
  image: string;
}> = [
  {
    id: "mercadopago",
    name: "Mercado Pago",
    image: "/mercado-pago.svg",
  },
];

export function EventPaymentPanel({
  connected,
  disabled,
  onConnect,
  onDisconnect,
}: {
  connected: boolean;
  disabled: boolean;
  onConnect: (provider: PaymentProviderI) => void;
  onDisconnect: () => void;
}) {
  return (
    <section className="order-3 space-y-3">
      <div className="flex items-center gap-3 px-1">
        <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <CreditCard className="size-5" />
        </div>
        <div className="min-w-0">
          <h2 className="text-base font-semibold tracking-tight">Pagamentos</h2>
          <p className="truncate text-xs text-muted-foreground">
            Conta que receberá as vendas deste evento.
          </p>
        </div>
      </div>
      <div className="flex flex-wrap gap-3">
        {providers.map((provider) => (
          <div
            key={provider.id}
            className="flex w-full max-w-md items-center gap-3 rounded-xl bg-card px-3 py-3 ring-1 ring-foreground/10 shadow-xs"
          >
            <div className="flex size-12 shrink-0 items-center justify-center rounded-lg bg-muted/50 p-2">
              <img
                src={provider.image}
                alt={provider.name}
                className="size-full object-contain"
              />
            </div>
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-semibold">{provider.name}</p>
              <Badge
                className="mt-1"
                variant={connected ? "default" : "secondary"}
              >
                {connected ? "Conectado" : "Não conectado"}
              </Badge>
            </div>
            <Button
              className="h-9 shrink-0 text-xs"
              variant={connected ? "outline" : "default"}
              disabled={disabled}
              onClick={() =>
                connected ? onDisconnect() : onConnect(provider.id)
              }
            >
              {connected ? "Desconectar conta" : "Conectar"}
            </Button>
          </div>
        ))}
      </div>
    </section>
  );
}
