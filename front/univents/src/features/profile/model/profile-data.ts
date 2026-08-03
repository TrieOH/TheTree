import type { JsonValue, ProfileData } from "@trieoh/identityx-sdk-ts";

export interface ProfileSchemaNode {
  type?: string | string[];
  title?: string;
  description?: string;
  format?: string;
  default?: JsonValue;
  enum?: JsonValue[];
  properties?: Record<string, ProfileSchemaNode>;
  required?: string[];
  items?: ProfileSchemaNode;
}

export interface UniventsProfile {
  legalName?: string;
  preferredName?: string;
  role?: string;
  pfpUrl?: string | null;
  bannerUrl?: string | null;
  aboutMe?: string;
  tagline?: string;
  website?: string | null;
  pronouns?: string;
  timezone?: string;
  contactEmail?: string | null;
  organization?: string;
  languages?: string[];
  specializations?: string[];
  location?: {
    city?: string;
    region?: string;
    country?: string;
    countryCode?: string;
  };
  socials?: Record<string, string | null | undefined>;
  visibility?: {
    hideSocials?: boolean;
    hideLocation?: boolean;
    hideLegalName?: boolean;
    hideContactEmail?: boolean;
    hideOrganization?: boolean;
  };
  createdAt?: string;
  updatedAt?: string;
  profileCompleteness?: number;
}

export const SYSTEM_PROFILE_FIELDS = new Set([
  "createdAt",
  "updatedAt",
  "profileCompleteness",
]);

export function asUniventsProfile(profile: ProfileData): UniventsProfile {
  return profile as unknown as UniventsProfile;
}

export function profileDisplayName(profile: UniventsProfile): string {
  return profile.preferredName || profile.legalName || "Univents member";
}

export function profileCompleteness(profile: UniventsProfile): number {
  const hasSocial = Object.values(profile.socials ?? {}).some(Boolean);
  const checks = [
    Boolean(profile.preferredName || profile.legalName),
    Boolean(profile.pfpUrl),
    Boolean(profile.bannerUrl),
    Boolean(profile.aboutMe),
    Boolean(profile.role || profile.organization),
    Boolean(
      profile.location?.city ||
        profile.location?.region ||
        profile.location?.country,
    ),
    Boolean(profile.specializations?.length || profile.languages?.length),
    Boolean(profile.website || profile.contactEmail || hasSocial),
  ];

  return Math.round((checks.filter(Boolean).length / checks.length) * 100);
}

export function withProfileTimestamps(
  current: ProfileData,
  next: ProfileData,
  now = new Date().toISOString(),
): ProfileData {
  return {
    ...current,
    ...next,
    createdAt: current.createdAt ?? now,
    updatedAt: now,
  };
}

export function applyProfileSchemaDefaults(
  schema: ProfileSchemaNode,
  current: ProfileData,
): ProfileData {
  const result = structuredClone(current);
  applyObjectDefaults(schema, result);
  return result;
}

function applyObjectDefaults(
  schema: ProfileSchemaNode,
  target: Record<string, JsonValue>,
): void {
  for (const [name, field] of Object.entries(schema.properties ?? {})) {
    if (target[name] === undefined && field.default !== undefined) {
      target[name] = structuredClone(field.default);
    }
    if (field.properties) {
      const current = target[name];
      if (!current || Array.isArray(current) || typeof current !== "object") {
        target[name] = {};
      }
      applyObjectDefaults(field, target[name] as Record<string, JsonValue>);
    }
  }
}

export function socialHref(network: string, value: string): string {
  if (/^https?:\/\//i.test(value)) return value;
  const handle = value.replace(/^@/, "");
  const bases: Record<string, string> = {
    x: "https://x.com/",
    twitter: "https://x.com/",
    github: "https://github.com/",
    twitch: "https://twitch.tv/",
    bluesky: "https://bsky.app/profile/",
    youtube: "https://youtube.com/@",
    linkedin: "https://linkedin.com/in/",
    instagram: "https://instagram.com/",
  };
  return bases[network] ? `${bases[network]}${handle}` : value;
}
