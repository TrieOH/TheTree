import { cn } from "#/shared/lib/utils";
import { User2, Shield, Eye, Pencil, Crown, UserMinus } from "lucide-react";
import type { NamespaceMemberI } from "../model";
import { timeAgo } from "#/shared/lib/helpers/date-utils";
import { Button } from "#/shared/ui/shadcn/button";
import {
  NamespaceMemberRoleEditor,
  NamespaceMemberRoleAdmin,
  NamespaceMemberRoleOwner,
  NamespaceMemberRoleViewer
} from "@trieoh/informd-models";

interface PropsI {
  data: NamespaceMemberI;
  onRemove?: (member: NamespaceMemberI | null) => void;
}

const roleConfig = {
  [NamespaceMemberRoleOwner]: {
    label: "Owner",
    icon: Crown,
    className: "text-amber-500 bg-amber-500/10",
  },
  [NamespaceMemberRoleAdmin]: {
    label: "Admin",
    icon: Shield,
    className: "text-blue-500 bg-blue-500/10",
  },
  [NamespaceMemberRoleEditor]: {
    label: "Editor",
    icon: Pencil,
    className: "text-emerald-500 bg-emerald-500/10",
  },
  [NamespaceMemberRoleViewer]: {
    label: "Viewer",
    icon: Eye,
    className: "text-muted-foreground bg-muted",
  },
} as const;

export function MemberCard({ data, onRemove }: PropsI) {
  const role = roleConfig[data.role] ?? roleConfig[NamespaceMemberRoleViewer];
  const RoleIcon = role.icon;
  const isOwner = data.role === NamespaceMemberRoleOwner;

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
            {data.user_id}
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
          Added {timeAgo(data.added_at)}
        </p>
      </div>

      {/* Remove button */}
      {!isOwner && (
        <Button
          variant="ghost"
          size="icon"
          className={cn(
            "shrink-0 text-muted-foreground",
            "hover:text-destructive hover:bg-destructive/10",
            "duration-150 cursor-pointer outline-0"
          )}
          onClick={() => onRemove?.(data)}
          title="Remove member"
        >
          <UserMinus className="size-4" />
        </Button>
      )}
    </div>
  );
}