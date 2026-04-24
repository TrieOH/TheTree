import z from "zod";

export const formCreateSchema = z.object({
  title: z.string({ error: "Title is required" })
    .min(3, "Title must be at least 3 characters long"),
});

export type FormCreateI = z.infer<typeof formCreateSchema>;

export interface FormI {
  id: string;
  project_id: string;
  owner_id: string;
  scope_id: string;
  title: string;
  status: "draft" | "open" | "closed" | "archived";
  current_version_id: string | null;
  created_at: string;
  updated_at: string;
  opened_at: string | null;
  closed_at: string | null;
  archived_at: string | null;
}