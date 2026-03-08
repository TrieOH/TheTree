import z from "zod";

export const eventCreateSchema = z.object({
  organization_id: z.uuid().optional(),
  name: z.string().min(2),
  acronym: z.string().optional(),
  slug: z.string().min(2),
  tagline: z.string().optional(),
  description: z.string().optional(),
  is_series: z.boolean(),
  logo_url: z.string().optional(),
  banner_url: z.string().optional(),
  contact_email: z.email(),
})

export type EventCreateI = z.infer<typeof eventCreateSchema>


interface SocialLinks {
  twitter?: string
  instagram?: string
  linkedin?: string
  website?: string
}
// type SocialLinks = Record<string, string>

export interface EventI {
  id: string;
  owner_id?: string;
  organization_id?: string;
  goauth_scope_id: string;
  name: string;
  acronym?: string;
  slug: string;
  tagline?: string;
  description?: string;
  is_series: boolean;
  editions_count: number;
  logo_url?: string;
  banner_url?: string;
  has_gallery: boolean;
  gallery_urls: string[];
  contact_email?: string;
  social_links?: SocialLinks;
  status: "draft" | "active" | "archived" | "discontinued";
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}