import type {
  CheckoutItem,
  CreateCheckoutRequest,
} from "@trieoh/univents-api/schemas";
import type { SubmitPaymentPayloadI } from "@/features/payments/model";
import type { CartItem } from "@/features/products/model/cart";

export function buildCheckoutRequest({
  items,
  totalCents,
  buyer,
  gift,
  payment,
}: {
  items: CartItem[];
  totalCents: number;
  buyer: { id: string; email: string };
  gift?: { name: string; email: string };
  payment?: SubmitPaymentPayloadI;
}): CreateCheckoutRequest {
  const checkoutItems = items.flatMap<CheckoutItem>((item) => {
    const itemType =
      item.type === "activity" ? "program_occurrence" : item.type;
    if (itemType !== "ticket") {
      return [
        { item_type: itemType, item_id: item.id, quantity: item.quantity },
      ];
    }

    return Array.from({ length: item.quantity }, () => ({
      item_type: "ticket" as const,
      item_id: item.id,
      quantity: 1,
      attendee: {
        user_id: gift ? undefined : buyer.id,
        email: gift?.email ?? buyer.email,
        name: gift?.name ?? buyer.email,
      },
    }));
  });

  if (totalCents === 0) return { items: checkoutItems };
  const pix = !payment || payment.payment_method_id === "pix";
  return {
    payment_method: pix ? "pix" : "credit_card",
    card_token: payment?.card_token,
    payment_method_id: pix ? undefined : payment.payment_method_id,
    issuer_id: payment?.issuer_id,
    installments: payment?.installments,
    payer: {
      email: payment?.payer_email ?? buyer.email,
      identification_type: payment?.identification_type ?? "",
      identification_number: payment?.identification_number ?? "",
    },
    items: checkoutItems,
  };
}
