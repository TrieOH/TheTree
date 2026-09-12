import { loadMercadoPago } from "@mercadopago/sdk-js";
import { createSignal, createEffect, onCleanup, Show, type Component } from "solid-js";
import Loader2Raw from "~icons/lucide/loader-2";
import LockRaw from "~icons/lucide/lock";
import QrCodeRaw from "~icons/lucide/qr-code";
import AlertCircleRaw from "~icons/lucide/alert-circle";
import {
  formatCNPJ,
  formatCPF,
  validateCNPJ,
  validateCPF,
} from "@/shared/lib/masks";
import type { CheckoutPayment } from "@/features/purchases/model/checkout";

const Loader2 = Loader2Raw as Component<{ class?: string }>;
const Lock = LockRaw as Component<{ class?: string }>;
const QrCode = QrCodeRaw as Component<{ class?: string }>;
const AlertCircle = AlertCircleRaw as Component<{ class?: string }>;

export type PaymentMethodI = 'credit_card' | 'pix';

declare global {
  interface Window {
    MercadoPago: new (
      publicKey: string,
      options?: Record<string, unknown>,
    ) => MercadoPagoInstance;
  }
}

interface MercadoPagoCardFormData {
  token: string;
  issuerId: string;
  paymentMethodId: string;
  amount: string;
  installments: string;
  identificationNumber: string;
  identificationType: string;
  cardholderEmail: string;
}

interface MercadoPagoCardForm {
  getCardFormData: () => MercadoPagoCardFormData;
  unmount: () => void;
}

interface MercadoPagoInstance {
  cardForm: (config: {
    amount: string;
    iframe: boolean;
    form: {
      id: string;
      cardNumber: { id: string; placeholder: string };
      expirationDate: { id: string; placeholder: string };
      securityCode: { id: string; placeholder: string };
      cardholderName: { id: string; placeholder: string };
      issuer: { id: string; placeholder: string };
      installments: { id: string; placeholder: string };
      identificationType: { id: string; placeholder: string };
      identificationNumber: { id: string; placeholder: string };
      cardholderEmail: { id: string; placeholder: string };
    };
    callbacks: {
      onFormMounted?: (error: unknown) => void;
      onSubmit?: (event: Event) => void;
      onFetching?: (resource: string) => (() => void) | undefined;
      onPaymentMethodReceived?: (
        error: unknown,
        data: { id: string; name: string } | null,
      ) => void;
    };
  }) => MercadoPagoCardForm;
}

function formatBRL(cents: number) {
  return (cents / 100).toLocaleString('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  });
}

const inputLike = [
  'flex items-center w-full rounded-md border border-input bg-background px-3',
  'text-sm text-foreground ring-offset-background transition-colors',
  'focus-within:outline-none focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-2',
  'disabled:cursor-not-allowed disabled:opacity-50',
].join(' ');

const selectLike =
  inputLike +
  ' h-10 cursor-pointer appearance-none' +
  " bg-[image:url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E\")] bg-no-repeat bg-[position:right_0.75rem_center] pr-8" +
  ' [color-scheme:light] dark:[color-scheme:dark] [&>option]:bg-popover [&>option]:text-popover-foreground';

function IframeField(props: { id: string; label: string; class?: string }) {
  return (
    <div class={'space-y-1.5' + (props.class ? ' ' + props.class : '')}>
      <label
        for={props.id}
        class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
      >
        {props.label}
      </label>
      <div
        id={props.id}
        class={inputLike + ' h-9 [&>iframe]:w-full [&>iframe]:h-full [&>iframe]:border-none [&>iframe]:bg-transparent'}
      />
    </div>
  );
}

function SelectField(props: { id: string; label: string; class?: string }) {
  return (
    <div class={'space-y-1.5' + (props.class ? ' ' + props.class : '')}>
      <label
        for={props.id}
        class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
      >
        {props.label}
      </label>
      <select id={props.id} class={selectLike}>
        <option value="" disabled />
      </select>
    </div>
  );
}

function PixForm(props: {
  amount: number;
  onSubmit: (email: string, identificationType: string, identificationNumber: string) => void;
  loading: boolean;
}) {
  const [email, setEmail] = createSignal('');
  const [identificationType, setIdentificationType] = createSignal('CPF');
  const [identificationNumber, setIdentificationNumber] = createSignal('');

  const isEmailValid = () => email().includes('@') && email().includes('.');
  const isDocValid = () =>
    identificationType() === 'CPF'
      ? validateCPF(identificationNumber())
      : validateCNPJ(identificationNumber());

  const canSubmit = () => !props.loading && isEmailValid() && isDocValid();

  const handleDocChange = (val: string) => {
    const formatted =
      identificationType() === 'CPF' ? formatCPF(val) : formatCNPJ(val);
    setIdentificationNumber(formatted);
  };

  return (
    <div class="space-y-5">
      {/* E-mail */}
      <div class="space-y-1.5">
        <label
          for="pix-email"
          class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
        >
          E-mail
        </label>
        <input
          id="pix-email"
          type="email"
          value={email()}
          onInput={(e) => setEmail(e.currentTarget.value)}
          placeholder="email@exemplo.com"
          class={
            inputLike +
            ' h-10' +
            (!isEmailValid() && email().length > 0
              ? ' border-destructive focus-within:ring-destructive'
              : '')
          }
        />
        <Show when={!isEmailValid() && email().length > 0}>
          <p class="text-[10px] text-destructive flex items-center gap-1">
            <AlertCircle class="w-3 h-3" /> E-mail inválido
          </p>
        </Show>
      </div>

      {/* Documento */}
      <div class="grid grid-cols-5 gap-3">
        <div class="space-y-1.5 col-span-2">
          <label
            for="pix-identification-type"
            class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
          >
            Tipo
          </label>
          <select
            id="pix-identification-type"
            value={identificationType()}
            onChange={(e) => {
              setIdentificationType(e.currentTarget.value);
              setIdentificationNumber('');
            }}
            class={selectLike}
          >
            <option value="CPF">CPF</option>
            <option value="CNPJ">CNPJ</option>
          </select>
        </div>
        <div class="space-y-1.5 col-span-3">
          <label
            for="pix-identification-number"
            class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
          >
            {identificationType() === 'CPF' ? 'CPF' : 'CNPJ'}
          </label>
          <input
            id="pix-identification-number"
            value={identificationNumber()}
            onInput={(e) => handleDocChange(e.currentTarget.value)}
            placeholder={identificationType() === 'CPF' ? '000.000.000-00' : '00.000.000/0000-00'}
            inputmode="numeric"
            class={
              inputLike +
              ' h-10 text-sm font-mono' +
              (!isDocValid() && identificationNumber().length > 0
                ? ' border-destructive focus-within:ring-destructive'
                : '')
            }
          />
          <Show when={!isDocValid() && identificationNumber().length > 0}>
            <p class="text-[10px] text-destructive flex items-center gap-1">
              <AlertCircle class="w-3 h-3" /> {identificationType()} inválido
            </p>
          </Show>
        </div>
      </div>

      <div class="flex items-center justify-between text-sm">
        <span class="text-muted-foreground">Total a pagar</span>
        <span class="font-bold text-foreground tabular-nums text-base">
          {formatBRL(props.amount)}
        </span>
      </div>

      <button
        type="button"
        onClick={() => props.onSubmit(email(), identificationType(), identificationNumber())}
        disabled={!canSubmit()}
        class="w-full h-11 flex items-center justify-center gap-2 bg-primary text-primary-foreground font-semibold rounded-md disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        <Show
          when={props.loading}
          fallback={<QrCode class="h-4 w-4 mr-2" />}
        >
          <Loader2 class="h-4 w-4 animate-spin" />
        </Show>
        {props.loading ? 'Gerando…' : 'Gerar QR Code Pix'}
      </button>
    </div>
  );
}

interface CardPayload {
  card_token: string;
  payment_method_id: string;
  installments: number;
  issuer_id: string;
  payer: { email: string; identification: { type: string; number: string } };
}

function CreditCardForm(props: {
  amount: number;
  onSubmit: (data: CardPayload) => void;
  loading: boolean;
  sellerPublicKey: string;
}) {
  let cardFormRef: MercadoPagoCardForm | null = null;

  const [fetching, setFetching] = createSignal(false);
  const [mounted, setMounted] = createSignal(false);
  const [brandName, setBrandName] = createSignal<string | null>(null);

  const [identificationType, setIdentificationType] = createSignal('CPF');
  const [identificationNumber, setIdentificationNumber] = createSignal('');
  const [cardholderEmail, setCardholderEmail] = createSignal('');

  const isEmailValid = () =>
    cardholderEmail().includes('@') && cardholderEmail().includes('.');
  const isDocValid = () =>
    identificationType() === 'CPF'
      ? validateCPF(identificationNumber())
      : validateCNPJ(identificationNumber());

  const isLoading = () => props.loading || fetching() || !mounted();

  const handleDocChange = (val: string) => {
    const formatted =
      identificationType() === 'CPF' ? formatCPF(val) : formatCNPJ(val);
    setIdentificationNumber(formatted);
  };

  // Run once on mount — compute returns undefined (no tracking), effect runs the init
  createEffect(
    () => undefined,
    () => {
      let cancelled = false;

      const init = async () => {
        try {
          if (!props.sellerPublicKey) {
            console.warn('MercadoPago: public key missing');
            return;
          }

          await loadMercadoPago();
          await new Promise<void>((resolve) => setTimeout(resolve, 300));

          const mp = new window.MercadoPago(props.sellerPublicKey, { locale: 'pt-BR' });

          const formEl = document.getElementById('form-checkout');
          if (!formEl) {
            console.warn('MercadoPago: form-checkout element not found');
            return;
          }

          const cardFormInstance = mp.cardForm({
            amount: (props.amount / 100).toFixed(2),
            iframe: true,
            form: {
              id: 'form-checkout',
              cardNumber: { id: 'form-checkout__cardNumber', placeholder: '0000 0000 0000 0000' },
              expirationDate: { id: 'form-checkout__expirationDate', placeholder: 'MM/AA' },
              securityCode: { id: 'form-checkout__securityCode', placeholder: '•••' },
              cardholderName: { id: 'form-checkout__cardholderName', placeholder: 'Como impresso no cartão' },
              identificationNumber: { id: 'form-checkout__identificationNumber', placeholder: 'Número do documento' },
              cardholderEmail: { id: 'form-checkout__cardholderEmail', placeholder: 'email@exemplo.com' },
              issuer: { id: 'form-checkout__issuer', placeholder: 'Banco emissor' },
              installments: { id: 'form-checkout__installments', placeholder: 'Parcelas' },
              identificationType: { id: 'form-checkout__identificationType', placeholder: 'Tipo' },
            },
            callbacks: {
              onFormMounted: (error) => {
                if (cancelled) return;
                if (error) {
                  console.warn('MercadoPago CardForm mount error:', error);
                  return;
                }
                setMounted(true);
              },
              onPaymentMethodReceived: (_err, data) => {
                if (!cancelled) setBrandName(data?.name ?? null);
              },
              onFetching: (_) => {
                setFetching(true);
                return () => {
                  setFetching(false);
                };
              },
              onSubmit: (event) => {
                event.preventDefault();
                if (!cardFormRef) {
                  console.warn('MercadoPago: CardForm instance not ready');
                  return;
                }

                if (!isEmailValid() || !isDocValid()) return;

                const {
                  token,
                  issuerId,
                  paymentMethodId,
                  installments,
                  identificationType: type,
                  identificationNumber: number,
                  cardholderEmail: email,
                } = cardFormRef.getCardFormData();

                props.onSubmit({
                  card_token: token,
                  payment_method_id: paymentMethodId,
                  installments: Number(installments) || 1,
                  issuer_id: issuerId,
                  payer: {
                    email,
                    identification: {
                      type,
                      number: number.replace(/\D/g, ''),
                    },
                  },
                });
              },
            },
          });

          cardFormRef = cardFormInstance;
        } catch (err) {
          console.error('MercadoPago init error:', err);
        }
      };

      void init();

      onCleanup(() => {
        cancelled = true;
        if (cardFormRef) {
          try {
            cardFormRef.unmount();
          } catch (err) {
            if (err instanceof Error)
              console.error('Unexpected error while unmount: ', err);
          }
          cardFormRef = null;
        }
        setMounted(false);
        setBrandName(null);
        setFetching(false);
      });
    },
  );

  return (
    <form id="form-checkout" class="space-y-4">
      <div class="space-y-1.5">
        <div class="flex items-center justify-between">
          <label
            for="form-checkout__cardNumber"
            class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
          >
            Número do cartão
          </label>
          <Show when={brandName()}>
            <span class="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground border border-border/60 rounded px-1.5 py-0.5 bg-muted/40 leading-none flex items-center h-5">
              {brandName()}
            </span>
          </Show>
        </div>
        <div
          id="form-checkout__cardNumber"
          class={inputLike + ' h-9 [&>iframe]:w-full [&>iframe]:h-full [&>iframe]:border-none [&>iframe]:bg-transparent'}
        />
      </div>

      <div class="space-y-1.5">
        <label
          for="form-checkout__cardholderName"
          class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
        >
          Nome no cartão
        </label>
        <input
          id="form-checkout__cardholderName"
          name="cardholderName"
          placeholder="Como impresso no cartão"
          autocomplete="cc-name"
          class={inputLike + ' h-10 uppercase tracking-widest text-sm'}
          onInput={(e) => {
            e.currentTarget.value = e.currentTarget.value.toUpperCase();
          }}
        />
      </div>

      <div class="grid grid-cols-2 gap-3">
        <IframeField id="form-checkout__expirationDate" label="Validade" />
        <IframeField id="form-checkout__securityCode" label="CVV" />
      </div>

      <SelectField id="form-checkout__installments" label="Parcelas" />

      {/* Hidden issuer select — not display:none to avoid SDK errors */}
      <select
        id="form-checkout__issuer"
        class="absolute opacity-0 pointer-events-none -z-10"
      >
        <option value="" disabled />
      </select>

      <div class="grid grid-cols-5 gap-3">
        <div class="space-y-1.5 col-span-2">
          <label
            for="form-checkout__identificationType"
            class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
          >
            Tipo
          </label>
          <select
            id="form-checkout__identificationType"
            name="identificationType"
            value={identificationType()}
            onChange={(e) => {
              setIdentificationType(e.currentTarget.value);
              setIdentificationNumber('');
            }}
            class={selectLike}
          >
            <option value="CPF">CPF</option>
            <option value="CNPJ">CNPJ</option>
          </select>
        </div>
        <div class="space-y-1.5 col-span-3">
          <label
            for="doc-display"
            class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
          >
            {identificationType()}
          </label>
          {/* Visible input with mask */}
          <input
            id="doc-display"
            value={identificationNumber()}
            onInput={(e) => handleDocChange(e.currentTarget.value)}
            placeholder={
              identificationType() === 'CPF' ? '000.000.000-00' : '00.000.000/0000-00'
            }
            inputmode="numeric"
            class={
              inputLike +
              ' h-10 text-sm font-mono' +
              (!isDocValid() && identificationNumber().length > 0
                ? ' border-destructive focus-within:ring-destructive'
                : '')
            }
          />
          {/* Hidden input for MP SDK with unmasked value */}
          <input
            type="hidden"
            id="form-checkout__identificationNumber"
            name="identificationNumber"
            value={identificationNumber().replace(/\D/g, '')}
          />
          <Show when={!isDocValid() && identificationNumber().length > 0}>
            <p class="text-[10px] text-destructive flex items-center gap-1">
              <AlertCircle class="w-3 h-3" /> {identificationType()} inválido
            </p>
          </Show>
        </div>
      </div>

      <div class="space-y-1.5">
        <label
          for="form-checkout__cardholderEmail"
          class="text-xs font-medium text-muted-foreground uppercase tracking-wide"
        >
          E-mail
        </label>
        <input
          id="form-checkout__cardholderEmail"
          name="cardholderEmail"
          type="email"
          value={cardholderEmail()}
          onInput={(e) => setCardholderEmail(e.currentTarget.value)}
          placeholder="email@exemplo.com"
          autocomplete="email"
          class={
            inputLike +
            ' h-10 text-sm' +
            (!isEmailValid() && cardholderEmail().length > 0
              ? ' border-destructive focus-within:ring-destructive'
              : '')
          }
        />
        <Show when={!isEmailValid() && cardholderEmail().length > 0}>
          <p class="text-[10px] text-destructive flex items-center gap-1">
            <AlertCircle class="w-3 h-3" /> E-mail inválido
          </p>
        </Show>
      </div>

      <div class="pt-1 space-y-2.5">
        <button
          type="submit"
          class="w-full h-11 flex items-center justify-center gap-2 bg-primary text-primary-foreground font-semibold rounded-md disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          disabled={isLoading() || !isEmailValid() || !isDocValid()}
        >
          <Show when={isLoading()} fallback={<Lock class="h-3.5 w-3.5" />}>
            <Loader2 class="h-4 w-4 animate-spin" />
          </Show>
          {props.loading
            ? 'Processando…'
            : !mounted()
              ? 'Carregando…'
              : `Pagar ${formatBRL(props.amount)}`}
        </button>
        <p class="text-center text-[10px] text-muted-foreground/50 leading-relaxed">
          Dados criptografados · Processado pelo Mercado Pago
        </p>
      </div>
    </form>
  );
}

interface PropsI {
  amount: number;
  method: PaymentMethodI;
  handleSubmit: (data: CheckoutPayment) => void;
  sellerPublicKey: string;
}

export function MercadoPagoForm(props: PropsI) {
  const [loading, setLoading] = createSignal(false);

  const handlePixSubmit = (
    email: string,
    identificationType: string,
    identificationNumber: string,
  ) => {
    setLoading(true);
    try {
      props.handleSubmit({
        method: 'pix',
        email,
        identificationType,
        identificationNumber: identificationNumber.replace(/\D/g, ''),
      });
    } finally {
      setLoading(false);
    }
  };

  const handleCardSubmit = (data: CardPayload) => {
    setLoading(true);
    try {
      props.handleSubmit({
        method: 'credit_card',
        cardToken: data.card_token,
        paymentMethodId: data.payment_method_id,
        issuerId: data.issuer_id,
        installments: data.installments,
        email: data.payer.email,
        identificationType: data.payer.identification.type,
        identificationNumber: data.payer.identification.number,
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="w-full space-y-4">
      <Show
        when={props.method === 'credit_card'}
        fallback={
          <PixForm
            amount={props.amount}
            onSubmit={handlePixSubmit}
            loading={loading()}
          />
        }
      >
        <CreditCardForm
          amount={props.amount}
          onSubmit={handleCardSubmit}
          loading={loading()}
          sellerPublicKey={props.sellerPublicKey}
        />
      </Show>
    </div>
  );
}
