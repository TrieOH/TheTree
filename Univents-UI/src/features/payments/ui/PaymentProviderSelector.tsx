import { MercadoPagoForm } from "./MercadoPago";
import type { PaymentProviderI, SubmitPaymentPayloadI } from "../model";

interface PaymentProviderSelectorProps {
  provider?: PaymentProviderI;
  amount: number;
  handleSubmit: (data: SubmitPaymentPayloadI) => void;
  seller_credential_id: string;
  seller_public_key: string;
}

export function PaymentProviderSelector({
  provider = "mercadopago",
  amount,
  handleSubmit,
  seller_credential_id,
  seller_public_key,
}: PaymentProviderSelectorProps) {

  switch (provider) {
    default:
      return (
        <MercadoPagoForm
          amount={amount}
          handleSubmit={handleSubmit}
          seller_credential_id={seller_credential_id}
          seller_public_key={seller_public_key}
        />
      )
  }
}