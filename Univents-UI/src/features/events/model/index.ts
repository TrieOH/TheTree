import z from "zod";

export const eventCreateSchema = z.object({
  organization_id: z.uuid().optional().nullable(),
  name: z.string().min(2),
  acronym: z.string().optional().nullable(),
  slug: z.string().min(2),
  tagline: z.string().optional().nullable(),
  description: z.string().optional().nullable(),
  is_series: z.boolean(),
  logo_url: z.string().optional().nullable(),
  banner_url: z.string().optional().nullable(),
  contact_email: z.email(),
})

export type EventCreateI = z.infer<typeof eventCreateSchema>


interface SocialLinks {
  twitter?: string
  instagram?: string
  linkedin?: string
  website?: string
}

export interface EventI {
  id: string;
  owner_id: string | null;
  organization_id: string | null;
  goauth_scope_id: string;
  name: string;
  acronym: string | null;
  slug: string;
  tagline: string | null;
  description: string | null;
  is_series: boolean;
  editions_count: number;
  logo_url: string | null;
  banner_url: string | null;
  has_gallery: boolean;
  gallery_urls: string[];
  contact_email: string | null;
  social_links: SocialLinks | null;
  status: "draft" | "active" | "archived" | "discontinued";
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}