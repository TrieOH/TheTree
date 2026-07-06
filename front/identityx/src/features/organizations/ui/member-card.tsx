import { cn } from "@/shared/lib/utils";
import { User2, Shield, Crown, UserMinus } from "lucide-react";
import type { OrganizationMemberI } from "../model";
import { timeAgo } from "@/shared/lib/date-utils";
import { OrganizationRoleAdmin, OrganizationRoleMember, OrganizationRoleOwner } from "@trieoh/identityx-models";


interface PropsI {
  data: OrganizationMemberI;
  onRemove?: (member: OrganizationMemberI | null) => void;
}

const roleConfig = {
  [OrganizationRoleOwner]: {
    label: "Owner",
    icon: Crown,
    className: "text-amber-500 bg-amber-500/10",
  },
  [OrganizationRoleAdmin]: {
    label: "Admin",
    icon: Shield,
    className: "text-blue-500 bg-blue-500/10",
  },
  [OrganizationRoleMember]: {
    label: "Member",
    icon: User2,
    className: "text-muted-foreground bg-muted",
  },
} as const;

export function MemberCard({ data, onRemove }: PropsI) {
  const role = roleConfig[data.role] ?? roleConfig[OrganizationRoleMember];
  const RoleIcon = role.icon;
  const isOwner = data.role === OrganizationRoleOwner;

  return (
    <div
      className={cn(
        "bg-card rounded-sm w-full cursor-default",
        "ring-1 ring-foreground/10 shadow-xs",
        "flex items-center gap-3 px-4 py-3",
        "hover:ring-foreground/20 duration-150"
      )}
    >
      {/* Avatar */}
      <div className="shrink-0 size-9 rounded-full bg-muted ring-1 ring-foreground/10 flex items-center justify-center">
        <User2 className="size-4 text-muted-foreground" />
      </div>

      {/* Info */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-sm font-semibold truncate">
            {data.actor_id}
          </span>
          <span
            className={cn(
              "inline-flex items-center gap-1 text-xs font-medium px-1.5 py-0.5 rounded-sm",
              role.className
            )}
          >
            <RoleIcon className="size-3" />
            {role.label}
          </span>
        </div>
        <p className="text-xs text-muted-foreground mt-0.5">
          Joined {timeAgo(data.joined_at)}
        </p>
      </div>

      {/* Remove button */}
      {!isOwner && (
        <button
          type="button"
          className={cn(
            "shrink-0 text-muted-foreground p-2 rounded-sm border border-transparent transition-all duration-150 cursor-pointer outline-none",
            "hover:text-destructive hover:bg-destructive/10"
          )}
          onClick={() => onRemove?.(data)}
          title="Remove member"
        >
          <UserMinus className="size-4" />
        </button>
      )}
    </div>
  );
}
