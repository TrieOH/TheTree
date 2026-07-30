import z from "zod";

export const programSchema = z.object({
  kind: z.enum(["activity", "checkpoint"]),
  name: z.string().trim().min(2, "Nome deve ter pelo menos 2 caracteres"),
  description: z.string().optional(),
  min_access_level: z.coerce.number().int().min(0).optional(),
  staff_only: z.boolean().default(false),
  price: z.coerce.number().int().min(0).optional(),
  banner_url: z.string().nullable().optional(),
});
export type ProgramCreateInput = z.input<typeof programSchema>;
export type ProgramCreateOutput = z.output<typeof programSchema>;
export interface ProgramI extends ProgramCreateOutput {
  id: string;
  edition_id: string;
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
  banner_url: string | null;
}
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
export const occurrenceSchema = z.object({
  starts_at: z.string().min(1, "Início é obrigatório"),
  ends_at: z.string().min(1, "Término é obrigatório"),
  max_capacity: z.coerce.number().int().positive().optional(),
});
export type OccurrenceCreateInput = z.input<typeof occurrenceSchema>;
export type OccurrenceCreateOutput = z.output<typeof occurrenceSchema>;
