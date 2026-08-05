import type { Event, EventStatus } from "@trieoh/univents-api/schemas";
import z from "zod";

export const eventCreateSchema = z.object({
  full_name: z
    .string({ error: "Nome é obrigatório" })
    .min(2, "O nome deve ter pelo menos 2 caracteres."),
  acronym: z
    .string()
    .optional()
    .nullable()
    .transform((val) => (val === "" ? null : val)),
  slug: z
    .string({ error: "Slug é obrigatório" })
    .min(2, "O slug deve ter pelo menos 2 caracteres."),
  description: z
    .string()
    .optional()
    .nullable()
    .transform((val) => (val === "" ? null : val)),
  contact_email: z
    .email("Informe um e-mail válido")
    .optional()
    .nullable()
    .or(z.literal(""))
    .transform((val) => (val === "" ? null : val)),
  logo_url: z
    .string()
    .optional()
    .nullable()
    .or(z.literal(""))
    .transform((val) => (val === "" ? null : val)),
  banner_url: z
    .string()
    .optional()
    .nullable()
    .or(z.literal(""))
    .transform((val) => (val === "" ? null : val)),
});

export type EventCreateInputI = z.input<typeof eventCreateSchema>;
export type EventCreateOutputI = z.output<typeof eventCreateSchema>;

export type EventStatusI = EventStatus;
export type EventI = Event;
