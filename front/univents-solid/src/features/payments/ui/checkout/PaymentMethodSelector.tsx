import { For, Show, type Component } from "solid-js";
import CreditCardRaw from "~icons/lucide/credit-card";
import QrCodeRaw from "~icons/lucide/qr-code";
import LockRaw from "~icons/lucide/lock";
import AlertTriangleRaw from "~icons/lucide/alert-triangle";

const CreditCard = CreditCardRaw as Component<{ class?: string }>;
const QrCode = QrCodeRaw as Component<{ class?: string }>;
const Lock = LockRaw as Component<{ class?: string }>;
const AlertTriangle = AlertTriangleRaw as Component<{ class?: string }>;

export type PaymentMethodI = 'credit_card' | 'pix';

interface MethodDef {
  id: PaymentMethodI;
  label: string;
  description: string;
  icon: Component<{ class?: string }>;
}

const methods: MethodDef[] = [
  {
    id: 'credit_card',
    label: 'Cartão de crédito',
    description: 'Até 12 parcelas',
    icon: CreditCard,
  },
  {
    id: 'pix',
    label: 'Pix',
    description: 'Aprovação imediata',
    icon: QrCode,
  },
];

function formatCurrency(cents: number) {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(cents / 100);
}

interface PaymentMethodSelectorProps {
  amountCents: number;
  selectedMethod: PaymentMethodI | null;
  onSelectMethod: (method: PaymentMethodI) => void;
  disabled?: boolean;
  isTooLowForCreditCard?: boolean;
}

export function PaymentMethodSelector(props: PaymentMethodSelectorProps) {
  return (
    <div class="w-full min-w-0 space-y-4">
      {/* Header */}
      <div class="flex items-center gap-2 pb-2 border-b border-border">
        <Lock class="w-4 h-4 text-primary" />
        <h2 class="text-xs font-bold uppercase tracking-wide text-muted-foreground">
          Pagamento seguro
        </h2>
      </div>

      {/* Methods */}
      <div class="space-y-2">
        <For each={methods}>
          {(method) => {
            const isSelected = () => props.selectedMethod === method.id;
            const isDisabled = () =>
              (props.disabled ?? false) ||
              (method.id === 'credit_card' && (props.isTooLowForCreditCard ?? false));

            const Icon = method.icon;

            return (
              <div class="space-y-1">
                <button
                  type="button"
                  data-state={isSelected() ? 'active' : 'inactive'}
                  onClick={() => {
                    if (!isDisabled()) props.onSelectMethod(method.id);
                  }}
                  disabled={isDisabled()}
                  class={
                    'w-full flex items-center gap-3 p-3 text-left transition-all rounded-none border' +
                    (isSelected() ? ' border-primary' : ' border-border') +
                    (isDisabled() ? ' opacity-50 cursor-not-allowed' : '')
                  }
                >
                  <Icon
                    class={
                      'w-5 h-5 shrink-0 transition-colors duration-200' +
                      (isSelected() ? ' text-primary' : ' text-muted-foreground')
                    }
                  />

                  <div class="flex-1 min-w-0">
                    <p
                      class={
                        'text-sm font-bold uppercase tracking-wide transition-colors duration-200' +
                        (isSelected() ? ' text-primary' : ' text-foreground')
                      }
                    >
                      {method.label}
                    </p>
                    <p class="text-xs text-muted-foreground font-medium normal-case tracking-normal">
                      {method.description}
                    </p>
                  </div>

                  <div
                    class={
                      'w-4 h-4 border shrink-0 relative overflow-hidden transition-colors duration-200' +
                      (isSelected()
                        ? ' border-primary bg-primary'
                        : ' border-muted-foreground/30')
                    }
                  >
                    <Show when={isSelected()}>
                      <div class="absolute inset-0 scale-50 bg-primary-foreground" />
                    </Show>
                  </div>
                </button>

                {/* Warning below credit card button */}
                <Show when={method.id === 'credit_card' && (props.isTooLowForCreditCard ?? false)}>
                  <div class="overflow-hidden">
                    <div class="flex items-start gap-2 px-3 py-2 bg-amber-50 border border-amber-200 dark:bg-amber-950/30 dark:border-amber-800">
                      <AlertTriangle class="w-3.5 h-3.5 text-amber-600 dark:text-amber-500 shrink-0 mt-0.5" />
                      <p class="text-[11px] text-amber-700 dark:text-amber-400 leading-relaxed">
                        Valor mínimo para cartão é{' '}
                        <span class="font-semibold">R$ 1,00</span>. Use Pix
                        para este valor.
                      </p>
                    </div>
                  </div>
                </Show>
              </div>
            );
          }}
        </For>
      </div>

      {/* Total */}
      <div class="pt-4 border-t-2 border-primary/10 space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-sm font-bold uppercase tracking-wide text-foreground">
            Total
          </span>
          <span class="text-xl font-bold text-primary tabular-nums">
            {formatCurrency(props.amountCents)}
          </span>
        </div>
      </div>
    </div>
  );
}
