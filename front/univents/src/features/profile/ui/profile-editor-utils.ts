import type { JsonValue, ProfileData } from "@trieoh/identityx-sdk-ts";
import blueskyIcon from "@/shared/ui/social-icons/assets/bluesky.svg";
import discordIcon from "@/shared/ui/social-icons/assets/discord.svg";
import githubIcon from "@/shared/ui/social-icons/assets/github.svg";
import instagramIcon from "@/shared/ui/social-icons/assets/instagram.svg";
import linkedinIcon from "@/shared/ui/social-icons/assets/linkedin.svg";
import twitchIcon from "@/shared/ui/social-icons/assets/twitch.svg";
import xIcon from "@/shared/ui/social-icons/assets/x.svg";
import youtubeIcon from "@/shared/ui/social-icons/assets/youtube.svg";
import type { ProfileSchemaNode } from "../model/profile-data";

export function inputValue(
  field: ProfileSchemaNode,
  inputType: string,
  value: string,
): JsonValue {
  if (inputType === "number") return Number(value);
  if (!value && Array.isArray(field.type) && field.type.includes("null")) {
    return null;
  }
  return value;
}

export function readValue(
  values: ProfileData,
  path: string[],
): JsonValue | undefined {
  let current: JsonValue = values;
  for (const segment of path) {
    if (!current || Array.isArray(current) || typeof current !== "object")
      return undefined;
    current = current[segment] as JsonValue;
  }
  return current;
}

export function writeValue(
  values: ProfileData,
  path: string[],
  value: JsonValue,
): ProfileData {
  const next = structuredClone(values);
  let current: Record<string, JsonValue> = next;
  path.forEach((segment, index) => {
    if (index === path.length - 1) current[segment] = value;
    else {
      const child = current[segment];
      if (!child || Array.isArray(child) || typeof child !== "object")
        current[segment] = {};
      current = current[segment] as Record<string, JsonValue>;
    }
  });
  return next;
}

export const socialIconSources: Record<string, string> = {
  x: xIcon,
  github: githubIcon,
  twitch: twitchIcon,
  bluesky: blueskyIcon,
  discord: discordIcon,
  youtube: youtubeIcon,
  linkedin: linkedinIcon,
  instagram: instagramIcon,
};

export function humanize(value: string): string {
  return value
    .replace(/([a-z])([A-Z])/g, "$1 $2")
    .replace(/^./, (letter) => letter.toUpperCase());
}

const FIELD_LABELS: Record<string, string> = {
  preferredName: "Nome de exibição",
  legalName: "Nome completo",
  pronouns: "Pronomes",
  tagline: "Frase de apresentação",
  aboutMe: "Sobre mim",
  role: "Função",
  organization: "Organização",
  contactEmail: "E-mail de contato",
  website: "Site",
  socials: "Redes sociais",
  city: "Cidade",
  region: "Estado ou região",
  country: "País",
  countryCode: "Código do país",
  languages: "Idiomas",
  specializations: "Idiomas",
  timezone: "Fuso horário",
  hideSocials: "Ocultar redes sociais",
  hideLocation: "Ocultar localização",
  hideLegalName: "Ocultar nome completo",
  hideContactEmail: "Ocultar e-mail de contato",
  hideOrganization: "Ocultar organização",
};

const FIELD_PLACEHOLDERS: Record<string, string> = {
  preferredName: "Como você quer ser chamado",
  legalName: "Digite seu nome completo",
  pronouns: "Ex.: ela/dela",
  tagline: "Uma frase curta sobre você",
  aboutMe: "Conte um pouco sobre você, sua experiência e seus interesses…",
  role: "Ex.: Desenvolvedor de software",
  organization: "Ex.: Univents",
  contactEmail: "voce@exemplo.com",
  website: "https://seusite.com",
  city: "Ex.: São Paulo",
  region: "Ex.: São Paulo",
  country: "Ex.: Brasil",
  countryCode: "Ex.: BR",
  languages: "Digite um idioma e separe por vírgula",
  specializations: "Digite um idioma e separe por vírgula",
  x: "Ex.: @univents ou x.com/univents",
  twitter: "Ex.: @univents ou twitter.com/univents",
  instagram: "Ex.: @univents ou instagram.com/univents",
  linkedin: "Ex.: trieoh ou linkedin.com/in/trieoh",
  github: "Ex.: trieoh ou github.com/trieoh",
  twitch: "Ex.: univents ou twitch.tv/univents",
  bluesky: "Ex.: univents.bsky.social",
  youtube: "Ex.: @univents ou youtube.com/@univents",
};

export const TIMEZONES = [
  "America/Sao_Paulo",
  "America/Manaus",
  "America/Cuiaba",
  "America/Rio_Branco",
  "America/Noronha",
  "America/New_York",
  "America/Los_Angeles",
  "Europe/Lisbon",
  "Europe/London",
  "Europe/Paris",
  "UTC",
];

export function fieldLabel(name: string, field: ProfileSchemaNode): string {
  return FIELD_LABELS[name] ?? field.title ?? humanize(name);
}

export function fieldPlaceholder(name: string): string {
  return FIELD_PLACEHOLDERS[name] ?? "Ex.: seu perfil ou informação";
}

export function normalizeWebsite(value: string): string {
  const trimmed = value.trim();
  if (!trimmed) return "";
  return `https://${trimmed.replace(/^https?:\/\//i, "").replace(/^\/+/, "")}`;
}

export function normalizeSocialProfile(network: string, value: string): string {
  const trimmed = value.trim();
  if (!trimmed) return "";
  const withoutProtocol = trimmed.replace(/^https?:\/\//i, "");
  if (!withoutProtocol.includes("/")) return withoutProtocol.replace(/^@/, "");

  try {
    const url = new URL(`https://${withoutProtocol}`);
    const parts = url.pathname.split("/").filter(Boolean);
    const prefixIndex =
      network === "linkedin"
        ? parts.indexOf("in")
        : network === "bluesky"
          ? parts.indexOf("profile")
          : -1;
    const profile = prefixIndex >= 0 ? parts[prefixIndex + 1] : parts.at(0);
    return (profile ?? trimmed).replace(/^@/, "").replace(/\/$/, "");
  } catch {
    return trimmed.replace(/^@/, "").replace(/\/$/, "");
  }
}

export function normalizeProfileLinks(values: ProfileData): ProfileData {
  const next = structuredClone(values);
  if (typeof next.website === "string") {
    next.website = normalizeWebsite(next.website);
  }
  const socials = next.socials;
  if (socials && !Array.isArray(socials) && typeof socials === "object") {
    for (const [network, value] of Object.entries(socials)) {
      if (typeof value === "string") {
        socials[network] = normalizeSocialProfile(network, value);
      }
    }
  }
  return next;
}

export function projectProfileValues(
  schema: ProfileSchemaNode,
  values: ProfileData,
): ProfileData {
  const project = (node: ProfileSchemaNode, value: JsonValue): JsonValue => {
    if (
      value === null ||
      Array.isArray(value) ||
      !node.properties ||
      typeof value !== "object"
    )
      return value;
    const result: Record<string, JsonValue> = {};
    for (const [name, child] of Object.entries(node.properties)) {
      const childValue = value[name];
      if (childValue !== undefined) result[name] = project(child, childValue);
    }
    return result;
  };

  return project(schema, values) as ProfileData;
}

export function timezoneLabel(timezone: string): string {
  const labels: Record<string, string> = {
    "America/Sao_Paulo": "Brasília — São Paulo (UTC−03:00)",
    "America/Manaus": "Manaus (UTC−04:00)",
    "America/Cuiaba": "Cuiabá (UTC−04:00)",
    "America/Rio_Branco": "Rio Branco (UTC−05:00)",
    "America/Noronha": "Fernando de Noronha (UTC−02:00)",
    UTC: "UTC (tempo universal)",
  };
  return labels[timezone] ?? timezone.replaceAll("_", " ");
}

export function nameFromId(id: string): string {
  return id.split("-").at(-1) ?? id;
}
