import z from "zod";

const difficultyTypeSchema = z
  .enum(
    ["no_prerequisites", "beginner", "intermediate", "advanced", "expert"],
    { error: "Invalid difficulty type" }
  ).default("no_prerequisites")
type DifficultyType = z.infer<typeof difficultyTypeSchema>

export const activityCreateSchema = z.object({
  title: z.string().min(3),
  description: z.string().optional(),
  location: z.string(),
  starts_at: z.iso.datetime(),
  ends_at: z.iso.datetime(),
  presenter_name: z.string().optional().nullable(),
  token_cost: z.number().gte(0),
  has_capacity: z.boolean(),
  capacity: z.int().gte(0),
  difficulty: difficultyTypeSchema,
})

export type ActivityCreateI = z.infer<typeof activityCreateSchema>


export interface ActivityI {
  id: string;
  scope_id: string;
  edition_id: string;
  title: string;
  description: string | null;
  status: "draft" | "published" | "ongoing" | "completed" | "canceled";
  location: string;
  starts_at: string;
  ends_at: string;
  presenter_name: string | null;
  token_cost: number;
  has_capacity: boolean;
  capacity: number;
  remaining_capacity: number;
  difficulty: DifficultyType;
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}