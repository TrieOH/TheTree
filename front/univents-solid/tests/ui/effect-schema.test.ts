import { describe, expect, it } from "vitest";
import {
  Schema,
  createSchema,
  minLength,
  maxLength,
  pattern,
  uuid,
  min,
  max,
  positive,
  nonNegative,
} from "../../src/shared/lib/effect-schema";

describe("effect-schema bridge", () => {
  const UserSchema = createSchema(
    Schema.Struct({
      name: Schema.String.pipe(minLength(2, { message: "Nome deve ter pelo menos 2 caracteres" })),
      age: Schema.Number.pipe(min(18, { message: "Idade mínima é 18" })),
      email: Schema.optional(Schema.String),
    }),
  );

  it("parses valid input with .parse()", () => {
    const input = { name: "Maria", age: 28 };
    const result = UserSchema.parse(input);
    expect(result).toEqual({ name: "Maria", age: 28 });
  });

  it("throws on invalid input with .parse()", () => {
    expect(() => UserSchema.parse({ name: "Maria", age: "invalid" })).toThrow();
  });

  it("returns { success: true, data } on valid input with .safeParse()", () => {
    const res = UserSchema.safeParse({ name: "Carlos", age: 34, email: "carlos@trieoh.com" });
    expect(res.success).toBe(true);
    if (res.success) {
      expect(res.data).toEqual({
        name: "Carlos",
        age: 34,
        email: "carlos@trieoh.com",
      });
    }
  });

  it("returns { success: false, error } on invalid input with .safeParse()", () => {
    const res = UserSchema.safeParse({ name: 123, age: "invalid" });
    expect(res.success).toBe(false);
    if (!res.success) {
      expect(res.error.message).toBeDefined();
      expect(res.error.issues.length).toBeGreaterThan(0);
    }
  });

  it("custom error messages are returned in issues and error.message", () => {
    const res = UserSchema.safeParse({ name: "A", age: 15 });
    expect(res.success).toBe(false);
    if (!res.success) {
      const messages = res.error.issues.map((i) => i.message);
      expect(messages).toContain("Nome deve ter pelo menos 2 caracteres");
    }
  });

  it("validates combinators with custom messages", () => {
    const ProductSchema = createSchema(
      Schema.Struct({
        code: Schema.String.pipe(
          maxLength(5, { message: "Código muito longo" }),
          pattern(/^[A-Z]+$/, { message: "Apenas letras maiúsculas" }),
        ),
        id: Schema.String.pipe(uuid({ message: "ID deve ser UUID" })),
        price: Schema.Number.pipe(positive({ message: "Preço deve ser maior que zero" })),
        stock: Schema.Number.pipe(nonNegative({ message: "Estoque não pode ser negativo" })),
        rating: Schema.Number.pipe(max(5, { message: "Nota máxima é 5" })),
      }),
    );

    const res = ProductSchema.safeParse({
      code: "TOOLONG123",
      id: "invalid-uuid",
      price: 0,
      stock: -5,
      rating: 10,
    });

    expect(res.success).toBe(false);
    if (!res.success) {
      const messages = res.error.issues.map((i) => i.message);
      expect(messages).toContain("Código muito longo");
      expect(messages).toContain("ID deve ser UUID");
      expect(messages).toContain("Preço deve ser maior que zero");
      expect(messages).toContain("Estoque não pode ser negativo");
      expect(messages).toContain("Nota máxima é 5");
    }
  });

  it("implements StandardSchemaV1 via ~standard", () => {
    expect(UserSchema["~standard"]).toBeDefined();
    expect(UserSchema["~standard"].version).toBe(1);

    const validResult = UserSchema["~standard"].validate({ name: "Ana", age: 22 });
    expect("value" in validResult).toBe(true);

    const invalidResult = UserSchema["~standard"].validate({ name: "Ana", age: "wrong" });
    expect("issues" in invalidResult).toBe(true);
  });
});
