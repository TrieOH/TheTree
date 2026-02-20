import z from "zod";

export const scopeCRUDSchema = z.object({
  id: z.string(),
  project_id: z.string(),
  name: z.string().min(3, "Name must be at least 3 characters long"),
  external_id: z.string().optional()
    .refine(
      val => val === undefined || val === "" || val.length >= 3,
      "External ID must be at least 3 characters long or empty"
    )
});

export type ScopeCRUD = z.infer<typeof scopeCRUDSchema>;

export interface Scope {
  id: string;
  name: string;
  type: string;
  project_id: string;
  external_id?: string;
  created_at: string;
}