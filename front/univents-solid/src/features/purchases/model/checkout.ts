import type { CheckoutItem, CreateCheckoutRequest } from "@trieoh/univents-api/schemas";
import type { CartItem } from "@/features/products/model/cart";

export type CheckoutPayment = {
  method: "pix" | "credit_card";
  cardToken?: string;
  paymentMethodId?: string;
  issuerId?: string;
  installments?: number;
  email: string;
  identificationType: string;
  identificationNumber: string;
};

export function checkoutRequest(items: CartItem[], actor: { id: string; email: string }, payment?: CheckoutPayment, gift?: { name: string; email: string }): CreateCheckoutRequest {
  return {
    ...(items.some((item) => item.price_cents > 0)
      ? {
          payment_method: payment?.method ?? "pix",
          card_token: payment?.cardToken,
          payment_method_id: payment?.paymentMethodId,
          issuer_id: payment?.issuerId,
          installments: payment?.installments,
          payer: {
            email: payment?.email ?? actor.email,
            identification_type: payment?.identificationType ?? "",
            identification_number: payment?.identificationNumber ?? "",
          },
        }
      : {}),
    items: items.flatMap<CheckoutItem>((item) =>
      item.type === "ticket"
        ? [{ item_type: "ticket" as const, item_id: item.id, quantity: 1, attendee: { user_id: gift ? undefined : actor.id, email: gift?.email ?? actor.email, name: gift?.name ?? actor.email } }]
        : [{ item_type: "product" as const, item_id: item.id, quantity: item.quantity }],
    ),
  };
}
