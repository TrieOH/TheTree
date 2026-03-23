import { useEffect, useRef, useState } from "react";
import { QrCode, Loader2, Lock, CreditCard } from "lucide-react";
import { loadMercadoPago } from "@mercadopago/sdk-js";
import type { SubmitPaymentPayloadI } from "../model";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import { Label } from "@/shared/ui/shadcn/label";
import { Input } from "@/shared/ui/shadcn/input";

declare global {
  interface Window {
    MercadoPago: new (
      publicKey: string,
      options?: Record<string, unknown>
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
        data: { id: string; name: string } | null
      ) => void;
    };
  }) => MercadoPagoCardForm;
}

interface PaymentPayload {
  card_token: string;
  payment_method_id: string;
  installments: number;
  issuer_id: string;
  payer: {
    email: string;
    identification: { type: string; number: string };
  };
}

type PaymentMethod = "credit_card" | "pix";

function formatBRL(cents: number) {
  return (cents / 100).toLocaleString("pt-BR", {
    style: "currency",
    currency: "BRL",
  });
}

const inputLike = [
  "flex items-center w-full rounded-md border border-input bg-background px-3",
  "text-sm text-foreground ring-offset-background transition-colors",
  "focus-within:outline-none focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-2",
  "disabled:cursor-not-allowed disabled:opacity-50",
].join(" ");

function IframeField({ id, label, className }: { id: string; label: string; className?: string }) {
  return (
    <div className={cn("space-y-1.5", className)}>
      <Label htmlFor={id} className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
        {label}
      </Label>
      <div
        id={id}
        className={cn(
          inputLike,
          "h-9",
          "[&>iframe]:w-full [&>iframe]:h-full [&>iframe]:border-none [&>iframe]:bg-transparent"
        )}
      />
    </div>
  );
}

function SelectField({ id, label, className }: { id: string; label: string; className?: string }) {
  return (
    <div className={cn("space-y-1.5", className)}>
      <Label htmlFor={id} className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
        {label}
      </Label>
      <select
        id={id}
        defaultValue=""
        className={cn(
          inputLike,
          "h-9 cursor-pointer appearance-none",
          "bg-[image:url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E\")] bg-no-repeat bg-position-[right_0.75rem_center] pr-8"
        )}
      >
        <option value="" disabled />
      </select>
    </div>
  );
}

function PixForm({ amount, onSubmit, loading }: {
  amount: number;
  onSubmit: (email: string, identificationType: string, identificationNumber: string) => void;
  loading: boolean;
}) {
  const [email, setEmail] = useState("");
  const [identificationType, setIdentificationType] = useState("CPF");
  const [identificationNumber, setIdentificationNumber] = useState("");

  const canSubmit = !loading && !!email && !!identificationNumber;

  return (
    <div className="space-y-5">
      <div className="flex items-start gap-3 pb-4 border-b border-border/50">
        <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary shrink-0 mt-0.5">
          <QrCode className="h-5 w-5 stroke-[1.5]" />
        </div>
        <div>
          <p className="text-sm font-medium leading-tight">Pix</p>
          <p className="text-xs text-muted-foreground mt-0.5 leading-relaxed">
            QR Code gerado na próxima etapa. Aprovação em segundos por qualquer banco ou carteira digital.
          </p>
        </div>
      </div>

      {/* E-mail */}
      <div className="space-y-1.5">
        <Label htmlFor="pix-email" className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          E-mail
        </Label>
        <Input
          id="pix-email"
          type="email"
          value={email}
          onChange={(e) => { setEmail(e.target.value) }}
          placeholder="email@exemplo.com"
        />
      </div>

      {/* Documento */}
      <div className="grid grid-cols-5 gap-3">
        <div className="space-y-1.5 col-span-2">
          <Label htmlFor="pix-identification-type" className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
            Tipo
          </Label>
          <select
            id="pix-identification-type"
            value={identificationType}
            onChange={(e) => { setIdentificationType(e.target.value); }}
            className={cn(
              inputLike,
              "h-9 cursor-pointer appearance-none",
              "bg-[image:url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E\")] bg-no-repeat bg-position-[right_0.75rem_center] pr-8"
            )}
          >
            <option value="CPF">CPF</option>
            <option value="CNPJ">CNPJ</option>
          </select>
        </div>
        <div className="space-y-1.5 col-span-3">
          <Label htmlFor="pix-identification-number" className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
            {identificationType === "CPF" ? "CPF" : "CNPJ"}
          </Label>
          <Input
            id="pix-identification-number"
            value={identificationNumber}
            onChange={(e) => { setIdentificationNumber(e.target.value); }}
            placeholder={identificationType === "CPF" ? "000.000.000-00" : "00.000.000/0000-00"}
            inputMode="numeric"
            className="text-sm font-mono"
          />
        </div>
      </div>

      <div className="flex items-center justify-between text-sm">
        <span className="text-muted-foreground">Total a pagar</span>
        <span className="font-bold text-foreground tabular-nums text-base">{formatBRL(amount)}</span>
      </div>

      <Button
        onClick={() => { onSubmit(email, identificationType, identificationNumber); }}
        disabled={!canSubmit}
      >
        {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <QrCode className="h-4 w-4" />}
        {loading ? "Gerando…" : "Gerar QR Code Pix"}
      </Button>
    </div>
  );
}

// ─── Cartão ───────────────────────────────────────────────────────────────────

function CreditCardForm({
  amount,
  onSubmit,
  loading,
  seller_public_key,
}: {
  amount: number;
  onSubmit: (data: PaymentPayload) => void;
  loading: boolean;
  seller_public_key: string;
}) {
  const cardFormRef = useRef<MercadoPagoCardForm | null>(null);
  const [fetching, setFetching] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [brandName, setBrandName] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    const init = async () => {
      await loadMercadoPago();
      if (cancelled) return;

      const mp = new window.MercadoPago(seller_public_key, { locale: "pt-BR" });

      const cardForm = mp.cardForm({
        amount: (amount / 100).toFixed(2),
        iframe: true,
        form: {
          id: "form-checkout",
          cardNumber: { id: "form-checkout__cardNumber", placeholder: "0000 0000 0000 0000" },
          expirationDate: { id: "form-checkout__expirationDate", placeholder: "MM/AA" },
          securityCode: { id: "form-checkout__securityCode", placeholder: "•••" },
          cardholderName: { id: "form-checkout__cardholderName", placeholder: "Como impresso no cartão" },
          identificationNumber: { id: "form-checkout__identificationNumber", placeholder: "000.000.000-00" },
          cardholderEmail: { id: "form-checkout__cardholderEmail", placeholder: "email@exemplo.com" },
          issuer: { id: "form-checkout__issuer", placeholder: "Banco emissor" },
          installments: { id: "form-checkout__installments", placeholder: "Parcelas" },
          identificationType: { id: "form-checkout__identificationType", placeholder: "Tipo de documento" },
        },
        callbacks: {
          onFormMounted: (error) => {
            if (error) { console.warn("MP CardForm mount error:", error); return; }
            setMounted(true);
          },
          onPaymentMethodReceived: (_err, data) => {
            setBrandName(data?.name ?? null);
          },
          onFetching: (resource) => {
            console.log("MP fetching:", resource);
            setFetching(true);
            return () => { setFetching(false); };
          },
          onSubmit: (event) => {
            event.preventDefault();
            const {
              token,
              issuerId,
              paymentMethodId,
              installments,
              identificationType,
              identificationNumber,
              cardholderEmail,
            } = cardForm.getCardFormData();

            onSubmit({
              card_token: token,
              payment_method_id: paymentMethodId,
              installments: Number(installments),
              issuer_id: issuerId,
              payer: {
                email: cardholderEmail,
                identification: { type: identificationType, number: identificationNumber },
              },
            });
          },
        },
      });

      cardFormRef.current = cardForm;
    };

    void init();
    return () => {
      cancelled = true;
      cardFormRef.current?.unmount();
    };
  }, [amount]);

  const isLoading = loading || fetching || !mounted;

  return (
    <form id="form-checkout" className="space-y-4">
      <div className="space-y-1.5">
        <div className="flex items-center justify-between">
          <Label htmlFor="form-checkout__cardNumber" className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
            Número do cartão
          </Label>
          {brandName && (
            <span className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground border border-border/60 rounded px-1.5 py-0.5 bg-muted/40 leading-none flex items-center h-5">
              {brandName}
            </span>
          )}
        </div>
        <div
          id="form-checkout__cardNumber"
          className={cn(
            inputLike,
            "h-9",
            "[&>iframe]:w-full [&>iframe]:h-full [&>iframe]:border-none [&>iframe]:bg-transparent"
          )}
        />
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="form-checkout__cardholderName" className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          Nome no cartão
        </Label>
        <Input
          id="form-checkout__cardholderName"
          name="cardholderName"
          placeholder="Como impresso no cartão"
          autoComplete="cc-name"
          className="uppercase tracking-widest text-sm"
          onChange={(e) => { e.target.value = e.target.value.toUpperCase(); }}
        />
      </div>

      <div className="grid grid-cols-2 gap-3">
        <IframeField id="form-checkout__expirationDate" label="Validade" />
        <IframeField id="form-checkout__securityCode" label="CVV" />
      </div>

      <SelectField id="form-checkout__installments" label="Parcelas" />

      <select id="form-checkout__issuer" className="hidden" defaultValue="">
        <option value="" disabled />
      </select>

      {/* Documento — capturado pelo cardForm e enviado no identification */}
      <div className="grid grid-cols-5 gap-3">
        <SelectField id="form-checkout__identificationType" label="Tipo" className="col-span-2" />
        <div className="space-y-1.5 col-span-3">
          <Label htmlFor="form-checkout__identificationNumber" className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
            CPF / CNPJ
          </Label>
          <Input
            id="form-checkout__identificationNumber"
            name="identificationNumber"
            placeholder="000.000.000-00"
            inputMode="numeric"
            className="text-sm font-mono"
          />
        </div>
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="form-checkout__cardholderEmail" className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          E-mail
        </Label>
        <Input
          id="form-checkout__cardholderEmail"
          name="cardholderEmail"
          type="email"
          placeholder="email@exemplo.com"
          autoComplete="email"
          className="text-sm"
        />
      </div>

      <div className="pt-1 space-y-2.5">
        <Button type="submit" className="w-full gap-2" disabled={isLoading}>
          {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Lock className="h-3.5 w-3.5" />}
          {loading ? "Processando…" : !mounted ? "Carregando…" : `Pagar ${formatBRL(amount)}`}
        </Button>
        <p className="text-center text-[10px] text-muted-foreground/50 leading-relaxed">
          Dados criptografados · Processado pelo Mercado Pago
        </p>
      </div>
    </form>
  );
}

// ─── Export ───────────────────────────────────────────────────────────────────

interface PropsI {
  amount: number;
  handleSubmit: (data: SubmitPaymentPayloadI) => void;
  seller_public_key: string;
}

export function MercadoPagoForm({ amount, handleSubmit, seller_public_key }: PropsI) {
  const isTooLowForCreditCard = amount < 100;
  const [method, setMethod] = useState<PaymentMethod>(isTooLowForCreditCard ? "pix" : "credit_card");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (isTooLowForCreditCard && method === "credit_card") {
      setMethod("pix");
    }
  }, [amount, isTooLowForCreditCard, method]);

  const handlePixSubmit = (
    email: string,
    identificationType: string,
    identificationNumber: string,
  ) => {
    setLoading(true);
    try {
      handleSubmit({
        payment_method_id: "pix",
        payer_email: email,
        payment_method_type: "bank_transfer",
        identification_type: identificationType,
        identification_number: identificationNumber,
      });
    } finally {
      setLoading(false);
    }
  };

  const handleCardSubmit = (data: PaymentPayload) => {
    setLoading(true);
    try {
      handleSubmit({
        ...data,
        payer_email: data.payer.email,
        payment_method_type: "credit_card",
        identification_type: data.payer.identification.type,
        identification_number: data.payer.identification.number,
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full space-y-4">
      <div className="flex border-b border-border">
        {(["credit_card", "pix"] as const).map((m) => {
          const isDisabled = m === "credit_card" && isTooLowForCreditCard;
          return (
            <button
              key={m}
              type="button"
              disabled={isDisabled}
              onClick={() => { setMethod(m); }}
              className={cn(
                "flex items-center gap-1.5 px-0 py-2 mr-5 text-xs font-semibold uppercase tracking-wide",
                "border-b-2 -mb-px transition-colors duration-150",
                method === m
                  ? "border-primary text-foreground"
                  : "border-transparent text-muted-foreground hover:text-foreground",
                isDisabled && "opacity-40 cursor-not-allowed"
              )}
            >
              {m === "credit_card" ? <CreditCard className="h-3.5 w-3.5" /> : <QrCode className="h-3.5 w-3.5" />}
              {m === "credit_card" ? "Cartão" : "Pix"}
            </button>
          );
        })}
      </div>

      {isTooLowForCreditCard && (
        <p className="text-[10px] text-amber-600 bg-amber-50 border border-amber-100 px-2 py-1.5 rounded">
          ⚠️ O valor mínimo para pagamento via cartão é de R$ 1,00. Use Pix para este valor.
        </p>
      )}

      <div className="pt-1">
        {method === "credit_card" ? (
          <CreditCardForm amount={amount} onSubmit={handleCardSubmit} loading={loading} seller_public_key={seller_public_key} />
        ) : (
          <PixForm amount={amount} onSubmit={handlePixSubmit} loading={loading} />
        )}
      </div>
    </div>
  );
}