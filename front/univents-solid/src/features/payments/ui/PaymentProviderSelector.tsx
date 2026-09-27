import { createSignal, onCleanup } from "solid-js";
import type { CheckoutPayment } from "@/features/purchases/model/checkout";
import { PaymentFormSheet } from "./checkout/PaymentFormSheet";
import { PaymentMethodSelector } from "./checkout/PaymentMethodSelector";
import type { PaymentMethodI } from "./checkout/PaymentMethodSelector";
import { MercadoPagoForm } from "./MercadoPago";

interface Props {
  provider?: 'mercadopago';
  amount: number;
  handleSubmit: (data: CheckoutPayment) => void;
  sellerPublicKey: string;
}

export function PaymentProviderSelector(props: Props) {
  const isTooLowForCreditCard = () => props.amount < 100;
  const [method, setMethod] = createSignal<PaymentMethodI | null>(null);
  const [sheetOpen, setSheetOpen] = createSignal(false);
  const [formReady, setFormReady] = createSignal(false);
  let closeTimeout: ReturnType<typeof setTimeout> | null = null;

  onCleanup(() => {
    if (closeTimeout) clearTimeout(closeTimeout);
  });

  const handleSelectMethod = (m: PaymentMethodI) => {
    if (m === 'credit_card' && isTooLowForCreditCard()) return;

    if (closeTimeout) {
      clearTimeout(closeTimeout);
      closeTimeout = null;
    }

    setMethod(m);
    setFormReady(false);
    setSheetOpen(true);
  };

  const handleClose = () => {
    setSheetOpen(false);
    setFormReady(false);

    if (closeTimeout) clearTimeout(closeTimeout);

    closeTimeout = setTimeout(() => {
      setMethod(null);
      closeTimeout = null;
    }, 400);
  };

  return (
    <>
      <PaymentMethodSelector
        amountCents={props.amount}
        selectedMethod={method()}
        onSelectMethod={handleSelectMethod}
        isTooLowForCreditCard={isTooLowForCreditCard()}
      />

      <PaymentFormSheet
        open={sheetOpen()}
        method={method()}
        onClose={handleClose}
        onReady={() => setFormReady(true)}
      >
        {method() && formReady() && (
          <MercadoPagoForm
            amount={props.amount}
            method={method()!}
            handleSubmit={(data) => {
              props.handleSubmit(data);
              handleClose();
            }}
            sellerPublicKey={props.sellerPublicKey}
          />
        )}
      </PaymentFormSheet>
    </>
  );
}
