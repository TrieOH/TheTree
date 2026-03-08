import z from "zod";

const EditionTypeSchema = z
  .enum(["year", "season", "number", "ordinal", "custom"], {error: "Invalid edition type"})
type EditionType = z.infer<typeof EditionTypeSchema>

export const editionCreateSchema = z.object({
  type: EditionTypeSchema,
  edition_name: z.string().min(3).max(256),
  tagline: z.string().max(512).optional().nullable(),
  description: z.string().max(8000).optional().nullable(),
  registration_opens_at: z.iso.datetime().optional().nullable(),
  registration_closes_at: z.iso.datetime().optional().nullable(),
  starts_at: z.iso.datetime(),
  ends_at: z.iso.datetime(),
  timezone: z.string(),
  location_name: z.string(),
  location_address: z.string(),
  logo_url: z.url().optional().nullable(),
  banner_url: z.url().optional().nullable(),
  contact_email: z.email().optional().nullable(),
  contact_phone: z.string().optional().nullable(),
  organizer_name: z.string().optional().nullable(),
})

export type EditionCreateI = z.infer<typeof editionCreateSchema>


export interface EditionI {
  id: string;
  event_id: string;
  goauth_scope_id: string;
  type: EditionType;
  edition_name: string;
  tagline: string | null;
  description: string | null;
  status: "draft" | "announced" | "open" | "ongoing" | "finished" |
    "completed" | "cancelled" | "postponed";
  monetary_type: "free" | "paid" | "mixed";
  registration_opens_at: string | null;
  registration_closes_at: string | null;
  starts_at: string;
  ends_at: string;
  timezone: string;
  location_name: string;
  location_address: string;
  logo_url: string | null;
  banner_url: string | null;
  contact_email: string | null;
  contact_phone: string | null;
  organizer_name: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}