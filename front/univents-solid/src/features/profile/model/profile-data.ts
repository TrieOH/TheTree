type ProfileData = Record<string, unknown>;

export type UniventsProfile = ProfileData & {
  legalName?: string;
  preferredName?: string;
  role?: string;
  pfpUrl?: string | null;
  bannerUrl?: string | null;
  aboutMe?: string;
  website?: string | null;
  pronouns?: string;
  contactEmail?: string | null;
  organization?: string;
  languages?: string[];
  socials?: Record<string, string | null | undefined>;
  visibility?: { hideLegalName?: boolean };
  createdAt?: string;
};

export const asUniventsProfile = (profile: ProfileData) =>
  profile as UniventsProfile;

export const profileDisplayName = (profile: UniventsProfile) =>
  profile.preferredName || profile.legalName || "Membro Univents";

export const profileCompleteness = (profile: UniventsProfile) => {
  const hasSocial = Object.values(profile.socials ?? {}).some(Boolean);
  const checks = [
    profile.preferredName || profile.legalName,
    profile.pfpUrl,
    profile.bannerUrl,
    profile.aboutMe,
    profile.role || profile.organization,
    profile.languages?.length,
    profile.website || profile.contactEmail || hasSocial,
  ];
  return Math.round((checks.filter(Boolean).length / checks.length) * 100);
};

export const socialHref = (network: string, value: string) => {
  if (/^https?:\/\//i.test(value)) return value;
  const bases: Record<string, string> = {
    x: "https://x.com/",
    github: "https://github.com/",
    twitch: "https://twitch.tv/",
    bluesky: "https://bsky.app/profile/",
    youtube: "https://youtube.com/@",
    linkedin: "https://linkedin.com/in/",
    instagram: "https://instagram.com/",
    discord: "https://discord.com/users/",
  };
  return bases[network] ? `${bases[network]}${value.replace(/^@/, "")}` : value;
};
