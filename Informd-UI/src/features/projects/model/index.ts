import z from "zod";

export const projectCreateSchema = z.object({
  name: z.string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
});

export type ProjectCreateI = z.infer<typeof projectCreateSchema>;

export interface ProjectI {
  id: string;
  name: string;
  owner_id: string;
  scope_id: string;
  created_at: string;
  updated_at: string;
}