import { EventMemberRole } from "@trieoh/univents-api/schemas";
import type { EventMember } from "@trieoh/univents-api/schemas";

export const eventMemberRoles = [
  EventMemberRole.owner,
  EventMemberRole.admin,
  EventMemberRole.staff,
] as const;

export type { EventMemberRole };

export interface EventMemberWithEmailI extends EventMember {
  email?: string;
  pfp_url?: string | null;
}

export interface MemberCardViewModel {
  readonly primaryLabel: string;
  readonly shortId: string;
  readonly initials: string;
  readonly formattedDate: string;
  readonly roleLabel: string;
  readonly roleBadgeClass: string;
  readonly roleDotClass: string;
  readonly pfpUrl: string | null;
}

const ROLE_DISPLAY: Record<
  EventMemberRole,
  { label: string; badge: string; dot: string }
> = {
  owner: {
    label: "Proprietário",
    badge: "border-violet-500/25 bg-violet-500/10 text-violet-700 dark:text-violet-300",
    dot: "bg-violet-500",
  },
  admin: {
    label: "Administrador",
    badge: "border-blue-500/25 bg-blue-500/10 text-blue-700 dark:text-blue-300",
    dot: "bg-blue-500",
  },
  staff: {
    label: "Equipe",
    badge: "border-emerald-500/25 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300",
    dot: "bg-emerald-500",
  },
};

/**
 * Pure helper function to format member display rules cleanly outside the JSX tree.
 */
export function formatMemberViewModel(member: EventMemberWithEmailI): MemberCardViewModel {
  const email = member.email?.trim();
  const shortId = member.user_id.slice(0, 8);
  const primaryLabel = email || `Usuário ${shortId}`;

  let initials = member.user_id.slice(0, 2).toUpperCase();
  if (email) {
    const parts = email.split("@")[0].replace(/[^a-zA-Z0-9]/g, "");
    initials = (parts.slice(0, 2) || "MB").toUpperCase();
  }

  let formattedDate = "";
  try {
    formattedDate = new Intl.DateTimeFormat("pt-BR", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    }).format(new Date(member.created_at));
  } catch {
    formattedDate = "";
  }

  const role = ROLE_DISPLAY[member.role] ?? ROLE_DISPLAY.staff;

  return {
    primaryLabel,
    shortId,
    initials,
    formattedDate,
    roleLabel: role.label,
    roleBadgeClass: role.badge,
    roleDotClass: role.dot,
    pfpUrl: member.pfp_url ?? null,
  };
}
