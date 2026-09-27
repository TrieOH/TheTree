import { describe, expect, it } from "vitest";
import { formatMoney, formatPrice } from "../../src/shared/lib/money";

describe("money formatting", () => {
  // Intl puts a non-breaking space after "R$", so match with \s instead of ==.
  it("formats cents as BRL", () => {
    expect(formatMoney(12345)).toMatch(/^R\$\s123,45$/);
  });

  it("labels a zero price as free", () => {
    expect(formatPrice(0)).toBe("Gratuito");
    expect(formatPrice(500)).toMatch(/^R\$\s5,00$/);
  });
});
