import z from "zod"

const accessModeSchema = z
  .enum(
    ["open", "ticket", "staff_only"],
    { error: "Invalid access mode" }
  ).default("open")
type AccessMode = z.infer<typeof accessModeSchema>


const checkpointTypeSchema = z
  .enum(
    ["entry", "zone", "amenity", "session", "exit"],
    { error: "Invalid checkpoint type" }
  ).default("entry")
type CheckpointType = z.infer<typeof checkpointTypeSchema>

export const checkpointCreateSchema = z.object({
  name: z.string(),
  access_mode: accessModeSchema,
  type: checkpointTypeSchema,
  starts_at: z.iso.datetime().optional().nullable(),
  ends_at: z.iso.datetime().optional().nullable(),
})

export type CheckpointCreateI = z.infer<typeof checkpointCreateSchema>

export interface CheckpointI {
  id: string;
  scope_id: string;
  edition_id: string;
  name: string;
  type: CheckpointType;
  access_mode: AccessMode;
  starts_at: string | null;
  ends_at: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}