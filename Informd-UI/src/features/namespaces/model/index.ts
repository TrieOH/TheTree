import z from "zod";

export const namespaceCreateSchema = z.object({
  name: z.string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
});

export type NamespaceCreateI = z.infer<typeof namespaceCreateSchema>;

export interface NamespaceI {
  id: string;
  name: string;
  owner_id: string;
  scope_id: string;
  created_at: string;
  updated_at: string;
}