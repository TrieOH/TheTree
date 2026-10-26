import {
  Schema,
  createSchema,
  minLength,
  nonNegative,
  positive,
  refine,
  type Infer,
} from "@/shared/lib/effect-schema";

export const programStep1Schema = createSchema(
  Schema.Struct({
    name: Schema.Trim.pipe(
      minLength(2, { message: "Nome deve ter pelo menos 2 caracteres." }),
    ),
  }),
);

export const programStep2ActivitySchema = createSchema(
  Schema.Struct({
    price_cents: Schema.Unknown.pipe(
      refine(
        (v) => typeof v === "number" && !isNaN(v) && v >= 0,
        { message: "Informe um valor válido em centavos." },
      ),
    ),
    min_access_level: Schema.Unknown.pipe(
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
  }),
);

export const programStep2CheckpointSchema = createSchema(
  Schema.Struct({
    min_access_level: Schema.Unknown.pipe(
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
  }),
);

export const programStep2Schema = (kind: "activity" | "checkpoint") =>
  kind === "activity" ? programStep2ActivitySchema : programStep2CheckpointSchema;

export const programSchema = createSchema(
  Schema.Struct({
    kind: Schema.Literals(["activity", "checkpoint"]),
    name: Schema.Trim.pipe(
      minLength(2, { message: "Nome deve ter pelo menos 2 caracteres." }),
    ),
    description: Schema.optional(Schema.String),
    min_access_level: Schema.optional(
      Schema.Number.pipe(
        nonNegative({ message: "O nível de acesso deve ser um número inteiro >= 0." }),
      ),
    ),
    staff_only: Schema.optional(Schema.Boolean),
    price: Schema.optional(
      Schema.Number.pipe(
        nonNegative({ message: "Informe um valor válido em centavos." }),
      ),
    ),
    banner_url: Schema.optional(Schema.NullOr(Schema.String)),
  }),
);

export type ProgramCreateInput = Infer<typeof programSchema>;
export type ProgramCreateOutput = Infer<typeof programSchema>;

export type ProgramI = ProgramCreateOutput & {
  id: string;
  edition_id: string;
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
  banner_url: string | null;
  staff_only: boolean;
};

export interface OccurrenceI {
  id: string;
  program_id: string;
  edition_id: string;
  starts_at: string;
  ends_at: string;
  max_capacity: number | null;
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
}

export type ProgramOccurrenceI = OccurrenceI;

export const occurrenceSchema = createSchema(
  Schema.Struct({
    starts_at: Schema.String.pipe(minLength(1)),
    ends_at: Schema.String.pipe(minLength(1)),
    max_capacity: Schema.optional(Schema.NullOr(Schema.Number.pipe(positive()))),
  }),
);

export type OccurrenceCreateInput = Infer<typeof occurrenceSchema>;
export type OccurrenceCreateOutput = Infer<typeof occurrenceSchema>;
