import z from "zod";

export interface SocialLinks {
  twitter?: string | null
  instagram?: string | null
  linkedin?: string | null
  website?: string | null
}

export const eventCreateSchema = z.object({
  organization_id: z.uuid().optional().nullable(),
  name: z.string({ error: "Name is required" })
    .min(2, "Name must be at least 2 characters long."),
  acronym: z.string().optional().nullable().transform(val => val === "" ? null : val),
  slug: z.string({ error: "Slug is required" })
    .min(2, "Slug must be at least 2 characters long."),
  tagline: z.string().optional().nullable().transform(val => val === "" ? null : val),
  description: z.string().optional().nullable().transform(val => val === "" ? null : val),
  is_series: z.boolean().nullish().transform(val => val ?? false),
  logo_url: z.string().optional().nullable().transform(val => val === "" ? null : val),
  banner_url: z.string().optional().nullable().transform(val => val === "" ? null : val),
  gallery_urls: z.array(z.string()).nullish().transform(val => val ?? []),
  contact_email: z.email(),
  social_links: z.object({
    twitter: z.url("Twitter link must be a valid URL").optional().nullable().or(z.literal("")).transform(v => v === "" ? null : v),
    instagram: z.url("Instagram link must be a valid URL").optional().nullable().or(z.literal("")).transform(v => v === "" ? null : v),
    linkedin: z.url("LinkedIn link must be a valid URL").optional().nullable().or(z.literal("")).transform(v => v === "" ? null : v),
    website: z.url("Website link must be a valid URL").optional().nullable().or(z.literal("")).transform(v => v === "" ? null : v),
  }).partial().optional().nullable(),
})

export type EventCreateI = z.infer<typeof eventCreateSchema>

export type EventStatusI = 'draft' | 'active' | 'archived' | 'discontinued'

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
  gallery_urls: string[];
  contact_email: string | null;
  social_links: SocialLinks | null;
  status: EventStatusI;
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}