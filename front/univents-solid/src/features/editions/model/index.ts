import type { Edition } from "@trieoh/univents-api/schemas";
import {
  Schema,
  createSchema,
  minLength,
  refine,
} from "@/shared/lib/effect-schema";

export type EditionStatus = "draft" | "future" | "active" | "past";

export interface EditionI extends Edition {
  status: EditionStatus;
}

export function inferEditionStatus(
  edition: Pick<Edition, "is_draft" | "starts_at" | "ends_at">,
  now = new Date(),
): EditionStatus {
  if (edition.is_draft) return "draft";

  const currentTime = now.getTime();
  if (currentTime < new Date(edition.starts_at).getTime()) return "future";
  if (currentTime > new Date(edition.ends_at).getTime()) return "past";
  return "active";
}

export function normalizeEdition(edition: Edition): EditionI {
  return {
    ...edition,
    status: inferEditionStatus(edition),
  };
}

export function editionRangesOverlap(
  left: Pick<EditionI, "starts_at" | "ends_at">,
  right: Pick<EditionI, "starts_at" | "ends_at">,
): boolean {
  return (
    new Date(left.starts_at).getTime() < new Date(right.ends_at).getTime() &&
    new Date(left.ends_at).getTime() > new Date(right.starts_at).getTime()
  );
}

export const editionStep1Schema = createSchema(
  Schema.Struct({
    name: Schema.Trim.pipe(
      minLength(2, { message: "Informe ao menos 2 caracteres." }),
    ),
    slug: Schema.Trim.pipe(
      minLength(2, { message: "Informe ao menos 2 caracteres." }),
    ),
  }),
);

const editionStep2DateSchema = createSchema(
  Schema.Struct({
    starts_at: Schema.Unknown.pipe(
      refine((v) => Boolean(v), { message: "Informe a data de início." }),
    ),
    ends_at: Schema.Unknown.pipe(
      refine((v) => Boolean(v), { message: "Informe a data de término." }),
    ),
  }),
);

const editionStep2RangeSchema = createSchema(
  Schema.Struct({
    starts_at: Schema.Unknown,
    ends_at: Schema.Unknown,
  }).pipe(
    refine(
      (data: { starts_at?: unknown; ends_at?: unknown }) => {
        if (!data.starts_at || !data.ends_at) return true;
        return (
          new Date(String(data.ends_at)).getTime() >
          new Date(String(data.starts_at)).getTime()
        );
      },
      { message: "O término deve ser posterior ao início." },
    ),
  ),
);

export const editionStep2Schema = {
  safeParse(val: { starts_at: unknown; ends_at: unknown }) {
    const base = editionStep2DateSchema.safeParse(val);
    if (!base.success) return base;
    const range = editionStep2RangeSchema.safeParse(val);
    if (!range.success) {
      return {
        success: false,
        error: {
          message: range.error.message,
          issues: [{ message: range.error.message, path: ["ends_at"] }],
        },
      };
    }
    return { success: true, data: val };
  },
};
