import { Payment } from "@mercadopago/sdk-react"

import type { IAdditionalCardFormData, IPaymentFormData } from "@mercadopago/sdk-react/esm/bricks/payment/type";

export type PaymentProviderType = "mercadopago";

interface PaymentProviderSelectorProps {
  provider?: PaymentProviderType;
  amount: number;
}


export function PaymentProviderSelector({
  provider = "mercadopago",
  amount
}: PaymentProviderSelectorProps) {
  const onSubmit = async (_: IPaymentFormData, _a: IAdditionalCardFormData | null | undefined) => {
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
                ticket: "all",
                debitCard: "all",
                bankTransfer: "all",
                prepaidCard: "all",
              }
            }}
          />
        );
    }
  };

  return <div className="w-full">{renderProvider()}</div>;
}
