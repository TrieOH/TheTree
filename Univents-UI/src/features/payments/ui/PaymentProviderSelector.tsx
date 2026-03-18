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
  const onSubmit = async (data: IPaymentFormData, addData: IAdditionalCardFormData | null | undefined) => {
    console.log("Data: ", data)
    console.log("Aditional Data: ", addData)
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
