import z from "zod"

export const ticketCreateSchema = z.object({
  name: z.string().min(3),
  description: z.string().optional().nullable(),
})

export type TicketCreateI = z.infer<typeof ticketCreateSchema>

export interface TicketI {
  id: string;
  scope_id: string;
  edition_id: string;
  name: string;
  description: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}