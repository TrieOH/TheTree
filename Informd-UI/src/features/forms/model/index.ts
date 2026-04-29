import z from "zod";

export const formCreateSchema = z.object({
  name: z.string({ error: "Title is required" })
    .min(3, "Title must be at least 3 characters long"),
});

export type FormCreateI = z.infer<typeof formCreateSchema>;

export interface FormI {
  id: string;
  namespace_id: string | null;
  owner_id: string;
  name: string;
  status: "draft" | "open" | "closed" | "archived";
  created_at: string;
  updated_at: string;
  opened_at: string | null;
  closed_at: string | null;
  archived_at: string | null;
}