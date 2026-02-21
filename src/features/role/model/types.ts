import z from "zod";

export const roleCRUDSchema = z
  .object({
    id: z.string(),
    project_id: z.string(),
    name: z.string().min(3, "Name must be at least 3 characters long").optional(),
    description: z
      .string()
      .min(3, "Description must be at least 3 characters long")
      .optional(),
  })
  .superRefine((data, ctx) => {
    const isUpdate = data.id.trim().length > 0;

    // CREATE
    if (!isUpdate) {
      if (!data.name) {
        ctx.addIssue({
          code: "custom",
          path: ["name"],
          message: "Name is required",
        });
      }

      if (!data.description) {
        ctx.addIssue({
          code: "custom",
          path: ["description"],
          message: "Description is required",
        });
      }
    }

    // UPDATE
    if (isUpdate) {
      if (!data.description) {
        ctx.addIssue({
          code: "custom",
          path: ["description"],
          message: "Description is required",
        });
      }
    }
  }
);

export type RoleCRUD = z.infer<typeof roleCRUDSchema>;

export interface Role {
  id: string;
  name: string;
  description: string;
  external_id: string;
  project_id: string;
  scope_id: string;
  scope_name: string;
  created_at: string;
  updated_at: string;
}