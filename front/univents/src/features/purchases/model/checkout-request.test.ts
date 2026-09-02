import { describe, expect, it } from "vitest";
import type { CartItem } from "@/features/products/model/cart";
import { buildCheckoutRequest } from "./checkout-request";

const item = (overrides: Partial<CartItem>): CartItem => ({
  id: "item-1",
  type: "product",
  name: "Item",
  price_cents: 1_000,
  quantity: 1,
  inventory_remaining: 10,
  has_inventory: true,
  ...overrides,
});

describe("buildCheckoutRequest", () => {
  it("maps products, activities and one attendee per ticket", () => {
    const request = buildCheckoutRequest({
      items: [
        item({ id: "product", quantity: 2 }),
        item({ id: "activity", type: "activity" }),
        item({ id: "ticket", type: "ticket", quantity: 2 }),
      ],
      totalCents: 4_000,
      buyer: { id: "actor-1", email: "buyer@example.com" },
    });

    expect(request.items).toEqual([
      { item_type: "product", item_id: "product", quantity: 2 },
      {
        item_type: "program_occurrence",
        item_id: "activity",
        quantity: 1,
      },
      {
        item_type: "ticket",
        item_id: "ticket",
        quantity: 1,
        attendee: {
          user_id: "actor-1",
          email: "buyer@example.com",
          name: "buyer@example.com",
        },
      },
      {
        item_type: "ticket",
        item_id: "ticket",
        quantity: 1,
        attendee: {
          user_id: "actor-1",
          email: "buyer@example.com",
          name: "buyer@example.com",
        },
      },
    ]);
  });

  it("uses gift attendee data and omits payment for a free checkout", () => {
    expect(
      buildCheckoutRequest({
        items: [item({ type: "ticket" })],
        totalCents: 0,
        buyer: { id: "actor-1", email: "buyer@example.com" },
        gift: { name: "Guest", email: "guest@example.com" },
      }),
    ).toEqual({
      items: [
        {
          item_type: "ticket",
          item_id: "item-1",
          quantity: 1,
          attendee: {
            user_id: undefined,
            email: "guest@example.com",
            name: "Guest",
          },
        },
      ],
    });
  });

  it("builds credit card payment data", () => {
    const request = buildCheckoutRequest({
      items: [item({})],
      totalCents: 1_000,
      buyer: { id: "actor-1", email: "buyer@example.com" },
      payment: {
        card_token: "token",
        payment_method_id: "visa",
        payment_method_type: "credit_card",
        issuer_id: "issuer",
        installments: 2,
        payer_email: "payer@example.com",
        identification_type: "CPF",
        identification_number: "123",
      },
    });

    expect(request).toMatchObject({
      payment_method: "credit_card",
      card_token: "token",
      payment_method_id: "visa",
      payer: {
        email: "payer@example.com",
        identification_type: "CPF",
        identification_number: "123",
      },
    });
  });
});
