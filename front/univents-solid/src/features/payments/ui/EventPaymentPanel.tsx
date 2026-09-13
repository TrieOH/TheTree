import type { JSX } from "@solidjs/web";
import { For } from "solid-js";
import CreditCardIcon from "~icons/lucide/credit-card";
import { Button } from "@trieoh/ui-solid";
import type { PaymentProviderI } from "../api";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCreditCard = CreditCardIcon as unknown as IconComp;

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

export function EventPaymentPanel(props: {
  connected: boolean;
  disabled: boolean;
  onConnect: (provider: PaymentProviderI) => void;
  onDisconnect: () => void;
}): JSX.Element {
  return (
    <section class="order-3 space-y-3">
      <div class="flex items-center gap-3 px-1">
        <div class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <LucideCreditCard class="size-5" />
        </div>
        <div class="min-w-0">
          <h2 class="text-base font-semibold tracking-tight text-foreground">
            Pagamentos
          </h2>
          <p class="truncate text-xs text-muted-foreground">
            Conta que receberá as vendas deste evento.
          </p>
        </div>
      </div>

      <div class="flex flex-wrap gap-3">
        <For each={providers}>
          {(provider) => (
            <div class="flex w-full max-w-md flex-col justify-between gap-3 rounded-xl bg-card p-3 ring-1 ring-foreground/10 shadow-xs sm:flex-row sm:items-center">
              <div class="flex min-w-0 items-center gap-3">
                <div class="flex size-12 shrink-0 items-center justify-center rounded-lg bg-muted/50 p-2">
                  <img
                    src={provider.image}
                    alt={provider.name}
                    class="size-full object-contain"
                  />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-semibold text-foreground">
                    {provider.name}
                  </p>
                  <span
                    class={`mt-1 inline-block rounded-full px-2.5 py-0.5 text-[11px] font-semibold transition-colors ${
                      props.connected
                        ? "bg-primary text-primary-foreground shadow-xs"
                        : "bg-secondary text-secondary-foreground"
                    }`}
                  >
                    {props.connected ? "Conectado" : "Não conectado"}
                  </span>
                </div>
              </div>
              <Button
                class="h-9 w-full shrink-0 text-xs sm:w-auto"
                variant={props.connected ? "outline" : "default"}
                disabled={props.disabled}
                onClick={() =>
                  props.connected
                    ? props.onDisconnect()
                    : props.onConnect(provider.id)
                }
              >
                {props.connected ? "Desconectar conta" : "Conectar"}
              </Button>
            </div>
          )}
        </For>
      </div>
    </section>
  );
}
