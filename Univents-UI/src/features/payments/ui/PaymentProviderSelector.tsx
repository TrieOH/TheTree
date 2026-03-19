import { Payment } from "@mercadopago/sdk-react"
import type { PaymentPayload } from "@/features/products/hooks/use-checkout-socket";

import type { IPaymentFormData } from "@mercadopago/sdk-react/esm/bricks/payment/type";

export type PaymentProviderType = "mercadopago";

interface PaymentProviderSelectorProps {
  provider?: PaymentProviderType;
  amount: number;
  handleSubmit: (data: PaymentPayload) => void
}


export function PaymentProviderSelector({
  provider = "mercadopago",
  amount,
  handleSubmit
}: PaymentProviderSelectorProps) {


  const onSubmit = async (data: IPaymentFormData) => {
    handleSubmit({
      card_token: data.formData.token,
      payment_method_id: data.formData.payment_method_id,
      installments: data.formData.installments
    })
  };
  const renderProvider = () => {
    switch (provider) {
      default:
        return (
          <Payment
            initialization={{ amount: amount }}
            onSubmit={onSubmit}
            customization={{
              paymentMethods: {
                creditCard: "all",
                mercadoPago: "all",
                atm: "all",
                debitCard: "all",
                bankTransfer: "all",
                prepaidCard: "all",
              },
              visual: {
                hideFormTitle: true,
                style: {
                  customVariables: {
                    formBackgroundColor: "transparent",
                    baseColor: "var(--primary)",
                  }
                }
              }
            }}
          />
        );
    }
  };

  return <div className="w-full">{renderProvider()}</div>;
}
