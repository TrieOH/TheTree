import type { TicketType } from "@trieoh/univents-api/schemas";
import {
  Schema,
  createSchema,
  minLength,
  nonNegative,
  positive,
  refine,
  type Infer,
} from "@/shared/lib/effect-schema";

export const ticketStep1Schema = createSchema(
  Schema.Struct({
    name: Schema.Trim.pipe(
      minLength(2, { message: "O nome deve ter pelo menos 2 caracteres." }),
    ),
  }),
);

export const ticketStep2Schema = createSchema(
  Schema.Struct({
    price_cents: Schema.Unknown.pipe(
      refine(
        (v) => typeof v === "number" && !isNaN(v) && v >= 0,
        { message: "Informe um preço válido." },
      ),
    ),
    access_level: Schema.Unknown.pipe(
      refine(
        (v) => String(v ?? "").trim() !== "",
        { message: "O nível de acesso é obrigatório." },
      ),
      refine(
        (v) => {
          const n = Number(v);
          return !isNaN(n) && Number.isInteger(n) && n >= 0;
        },
        { message: "O nível de acesso deve ser um número inteiro >= 0." },
      ),
    ),
    max_quantity: Schema.optional(
      Schema.Unknown.pipe(
        refine(
          (v) => {
            if (v === null || v === undefined || String(v).trim() === "") return true;
            const n = Number(v);
            return !isNaN(n) && Number.isInteger(n) && n > 0;
          },
          { message: "Quantidade deve ser um número inteiro maior que 0." },
        ),
      ),
    ),
  }),
);

export const ticketCreateSchema = createSchema(
  Schema.Struct({
    name: Schema.Trim.pipe(
      minLength(2, { message: "O nome deve ter pelo menos 2 caracteres." }),
    ),
    description: Schema.optional(Schema.NullOr(Schema.String)),
    access_level: Schema.Number.pipe(
      nonNegative({ message: "O nível de acesso deve ser um número inteiro >= 0." }),
    ),
    price_cents: Schema.Number.pipe(
      nonNegative({ message: "Informe um preço válido." }),
    ),
    max_quantity: Schema.optional(
      Schema.NullOr(
        Schema.Number.pipe(
          positive({ message: "Quantidade deve ser um número inteiro maior que 0." }),
        ),
      ),
    ),
  }),
);

export type TicketCreateInputI = Infer<typeof ticketCreateSchema>;
export type TicketCreateOutputI = Infer<typeof ticketCreateSchema>;

export type TicketI = TicketType;
