import type { Event, EventStatus } from "@trieoh/univents-api/schemas";
import {
  Schema,
  createSchema,
  minLength,
  refine,
} from "@/shared/lib/effect-schema";

export type EventI = Event;
export type EventStatusI = EventStatus;

export const eventStep1Schema = createSchema(
  Schema.Struct({
    full_name: Schema.Trim.pipe(
      minLength(2, { message: "Informe ao menos 2 caracteres." }),
    ),
    slug: Schema.Trim.pipe(
      minLength(2, { message: "Informe ao menos 2 caracteres." }),
    ),
  }),
);

export const eventStep2Schema = createSchema(
  Schema.Struct({
    contact_email: Schema.optional(
      Schema.Unknown.pipe(
        refine(
          (val) => {
            const str = String(val ?? "").trim();
            if (str.length === 0) return true;
            return str.includes("@");
          },
          { message: "E-mail inválido." },
        ),
      ),
    ),
  }),
);

export const eventMemberStep1Schema = createSchema(
  Schema.Struct({
    email: Schema.Unknown.pipe(
      refine(
        (val) => {
          const str = String(val ?? "").trim().toLowerCase();
          return Boolean(str) && str.includes("@");
        },
        { message: "Informe um e-mail válido." },
      ),
    ),
  }),
);

export const eventMemberStep2Schema = createSchema(
  Schema.Struct({
    role: Schema.Unknown.pipe(
      refine(
        (val) => Boolean(val),
        { message: "Selecione uma função para o membro." },
      ),
    ),
  }),
);
