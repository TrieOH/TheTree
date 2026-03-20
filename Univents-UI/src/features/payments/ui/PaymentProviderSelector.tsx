import { MercadoPagoForm } from "./MercadoPago";
import type { PaymentProviderI, SubmitPaymentPayloadI } from "../model";

interface PaymentProviderSelectorProps {
  provider?: PaymentProviderI;
  amount: number;
  handleSubmit: (data: SubmitPaymentPayloadI) => void;
}

export function PaymentProviderSelector({
  provider = "mercadopago",
  amount,
  handleSubmit,
}: PaymentProviderSelectorProps) {

  switch (provider) {
    default: <MercadoPagoForm amount={amount} handleSubmit={handleSubmit} />
  }
  return <MercadoPagoForm amount={amount} handleSubmit={handleSubmit} />
}