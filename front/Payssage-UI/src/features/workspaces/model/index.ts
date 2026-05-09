import z from "zod";

export const workspaceCreateSchema = z.object({
  name: z.string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
});

export type WorkspaceCreateI = z.infer<typeof workspaceCreateSchema>;

export interface WorkspaceI {
  id: string;
  name: string;
  created_at: string;
  sandbox: boolean;
}