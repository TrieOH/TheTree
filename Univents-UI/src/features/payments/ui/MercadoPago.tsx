import { useEffect, useRef, useState } from "react";
import { QrCode, Loader2, Lock, AlertCircle } from "lucide-react";
import { loadMercadoPago } from "@mercadopago/sdk-js";
import type { PaymentMethodI, SubmitPaymentPayloadI } from "../model";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import { Label } from "@/shared/ui/shadcn/label";
import { Input } from "@/shared/ui/shadcn/input";
import { formatCPF, formatCNPJ, validateCPF, validateCNPJ } from "@/shared/lib/masks";

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

  const isEmailValid = email.includes("@") && email.includes(".");
  const isDocValid = identificationType === "CPF"
    ? validateCPF(identificationNumber)
    : validateCNPJ(identificationNumber);

  const canSubmit = !loading && isEmailValid && isDocValid;

  const handleDocChange = (val: string) => {
    const formatted = identificationType === "CPF" ? formatCPF(val) : formatCNPJ(val);
    setIdentificationNumber(formatted);
  };

  return (
    <div className="space-y-5">
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
          className={cn(!isEmailValid && email.length > 0 && "border-destructive focus-visible:ring-destructive")}
        />
        {!isEmailValid && email.length > 0 && (
          <p className="text-[10px] text-destructive flex items-center gap-1">
            <AlertCircle className="w-3 h-3" /> E-mail inválido
          </p>
        )}
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
            onChange={(e) => {
              setIdentificationType(e.target.value);
              setIdentificationNumber("");
            }}
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
            onChange={(e) => { handleDocChange(e.target.value); }}
            placeholder={identificationType === "CPF" ? "000.000.000-00" : "00.000.000/0000-00"}
            inputMode="numeric"
            className={cn(
              "text-sm font-mono",
              !isDocValid && identificationNumber.length > 0 && "border-destructive focus-visible:ring-destructive"
            )}
          />
          {!isDocValid && identificationNumber.length > 0 && (
            <p className="text-[10px] text-destructive flex items-center gap-1">
              <AlertCircle className="w-3 h-3" /> {identificationType} inválido
            </p>
          )}
        </div>
      </div>

      <div className="flex items-center justify-between text-sm">
        <span className="text-muted-foreground">Total a pagar</span>
        <span className="font-bold text-foreground tabular-nums text-base">{formatBRL(amount)}</span>
      </div>

      <Button
        onClick={() => { onSubmit(email, identificationType, identificationNumber); }}
        disabled={!canSubmit}
        className="w-full h-11"
      >
        {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <QrCode className="h-4 w-4 mr-2" />}
        {loading ? "Gerando…" : "Gerar QR Code Pix"}
      </Button>
    </div>
  );
}

function CreditCardForm({
  amount,
  onSubmit,
  loading,
  sellerPublicKey,
}: {
  amount: number;
  onSubmit: (data: PaymentPayload) => void;
  loading: boolean;
  sellerPublicKey: string;
}) {
  const cardFormRef = useRef<MercadoPagoCardForm | null>(null);
  const onSubmitRef = useRef(onSubmit);
  onSubmitRef.current = onSubmit;

  const [fetching, setFetching] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [brandName, setBrandName] = useState<string | null>(null);

  // Form states for validation and masking
  const [identificationType, setIdentificationType] = useState("CPF");
  const [identificationNumber, setIdentificationNumber] = useState("");
  const [cardholderEmail, setCardholderEmail] = useState("");

  const isEmailValid = cardholderEmail.includes("@") && cardholderEmail.includes(".");
  const isDocValid = identificationType === "CPF"
    ? validateCPF(identificationNumber)
    : validateCNPJ(identificationNumber);

  const validationRef = useRef({ isEmailValid, isDocValid });
  useEffect(() => {
    validationRef.current = { isEmailValid, isDocValid };
  }, [isEmailValid, isDocValid]);

  useEffect(() => {
    let cancelled = false;

    const init = async () => {
      try {
        if (!sellerPublicKey) {
          console.warn("MercadoPago: public key missing");
          return;
        }

        await loadMercadoPago();

        await new Promise<void>((resolve) => setTimeout(resolve, 300));

        const mp = new window.MercadoPago(sellerPublicKey, { locale: "pt-BR" });

        const formEl = document.getElementById("form-checkout");
        if (!formEl) {
          console.warn("MercadoPago: form-checkout element not found");
          return;
        }

        const cardFormInstance = mp.cardForm({
          amount: (amount / 100).toFixed(2),
          iframe: true,
          form: {
            id: "form-checkout",
            cardNumber: { id: "form-checkout__cardNumber", placeholder: "0000 0000 0000 0000" },
            expirationDate: { id: "form-checkout__expirationDate", placeholder: "MM/AA" },
            securityCode: { id: "form-checkout__securityCode", placeholder: "•••" },
            cardholderName: { id: "form-checkout__cardholderName", placeholder: "Como impresso no cartão" },
            identificationNumber: { id: "form-checkout__identificationNumber", placeholder: "Número do documento" },
            cardholderEmail: { id: "form-checkout__cardholderEmail", placeholder: "email@exemplo.com" },
            issuer: { id: "form-checkout__issuer", placeholder: "Banco emissor" },
            installments: { id: "form-checkout__installments", placeholder: "Parcelas" },
            identificationType: { id: "form-checkout__identificationType", placeholder: "Tipo" },
          },
          callbacks: {
            onFormMounted: (error) => {
              if (cancelled) return;
              if (error) {
                console.warn("MercadoPago CardForm mount error:", error);
                return;
              }
              setMounted(true);
            },
            onPaymentMethodReceived: (_err, data) => {
              if (!cancelled) setBrandName(data?.name ?? null);
            },
            onFetching: (_) => {
              setFetching(true);
              return () => { setFetching(false); };
            },
            onSubmit: (event) => {
              event.preventDefault();
              if (!cardFormRef.current) {
                console.warn("MercadoPago: CardForm instance not ready");
                return;
              }

              // Validation before submit using ref to get current value
              if (!validationRef.current.isEmailValid || !validationRef.current.isDocValid) return;

              const {
                token, issuerId, paymentMethodId, installments,
                identificationType: type, identificationNumber: number, cardholderEmail: email,
              } = cardFormRef.current.getCardFormData();

              onSubmitRef.current({
                card_token: token,
                payment_method_id: paymentMethodId,
                installments: Number(installments),
                issuer_id: issuerId,
                payer: {
                  email,
                  identification: {
                    type,
                    number: number.replace(/\D/g, "")
                  },
                },
              });
            },
          },
        });

        cardFormRef.current = cardFormInstance;
      } catch (err) {
        console.error("MercadoPago init error:", err);
      }
    };

    void init();

    return () => {
      cancelled = true;
      if (cardFormRef.current) {
        try {
          cardFormRef.current.unmount();
        } catch (err) {
          if (err instanceof Error)
            console.error("Unexpected error while unmount: ", err)
        }
        cardFormRef.current = null;
      }
      setMounted(false);
      setBrandName(null);
      setFetching(false);
    };
  }, [amount, sellerPublicKey]);

  const isLoading = loading || fetching || !mounted;

  const handleDocChange = (val: string) => {
    const formatted = identificationType === "CPF" ? formatCPF(val) : formatCNPJ(val);
    setIdentificationNumber(formatted);
  };

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

      {/* Hidden issuer select but not with display:none to avoid node not found errors */}
      <select id="form-checkout__issuer" className="absolute opacity-0 pointer-events-none -z-10" defaultValue="">
        <option value="" disabled />
      </select>

      <div className="grid grid-cols-5 gap-3">
        <div className="space-y-1.5 col-span-2">
          <Label htmlFor="form-checkout__identificationType" className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
            Tipo
          </Label>
          <select
            id="form-checkout__identificationType"
            name="identificationType"
            value={identificationType}
            onChange={(e) => {
              setIdentificationType(e.target.value);
              setIdentificationNumber("");
            }}
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
          <Label htmlFor="doc-display" className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
            {identificationType}
          </Label>
          {/* Visible input with mask for the user */}
          <Input
            id="doc-display"
            value={identificationNumber}
            onChange={(e) => { handleDocChange(e.target.value); }}
            placeholder={identificationType === "CPF" ? "000.000.000-00" : "00.000.000/0000-00"}
            inputMode="numeric"
            className={cn(
              "text-sm font-mono",
              !isDocValid && identificationNumber.length > 0 && "border-destructive focus-visible:ring-destructive"
            )}
          />
          {/* Hidden input for Mercado Pago SDK with UNMASKED value */}
          <input
            type="hidden"
            id="form-checkout__identificationNumber"
            name="identificationNumber"
            value={identificationNumber.replace(/\D/g, "")}
          />
          {!isDocValid && identificationNumber.length > 0 && (
            <p className="text-[10px] text-destructive flex items-center gap-1">
              <AlertCircle className="w-3 h-3" /> {identificationType} inválido
            </p>
          )}
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
          value={cardholderEmail}
          onChange={(e) => { setCardholderEmail(e.target.value); }}
          placeholder="email@exemplo.com"
          autoComplete="email"
          className={cn(
            "text-sm",
            !isEmailValid && cardholderEmail.length > 0 && "border-destructive focus-visible:ring-destructive"
          )}
        />
        {!isEmailValid && cardholderEmail.length > 0 && (
          <p className="text-[10px] text-destructive flex items-center gap-1">
            <AlertCircle className="w-3 h-3" /> E-mail inválido
          </p>
        )}
      </div>

      <div className="pt-1 space-y-2.5">
        <Button type="submit" className="w-full h-11 gap-2" disabled={isLoading || !isEmailValid || !isDocValid}>
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

interface PropsI {
  amount: number;
  method: PaymentMethodI;
  handleSubmit: (data: SubmitPaymentPayloadI) => void;
  sellerPublicKey: string;
}

export function MercadoPagoForm({ amount, handleSubmit, method, sellerPublicKey }: PropsI) {
  const [loading, setLoading] = useState(false);

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
        identification_number: identificationNumber.replace(/\D/g, ""),
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
        identification_number: data.payer.identification.number.replace(/\D/g, ""),
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full space-y-4">
      {method === "credit_card" ? (
        <CreditCardForm
          amount={amount}
          onSubmit={handleCardSubmit}
          loading={loading}
          sellerPublicKey={sellerPublicKey}
        />
      ) : (
        <PixForm amount={amount} onSubmit={handlePixSubmit} loading={loading} />
      )}
    </div>
  );
}
