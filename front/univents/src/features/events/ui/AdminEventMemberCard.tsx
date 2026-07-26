import { ShieldCheck, Trash2, UserRound } from "lucide-react";
import { motion } from "motion/react";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import type { EventMemberI } from "../api/members";
import type { EventMemberRole } from "../model/member";

const roleConfig: Record<
  EventMemberRole,
  { label: string; pill: string; icon: typeof ShieldCheck }
> = {
  owner: {
    label: "Proprietário",
    pill: "border-violet-500/20 bg-violet-500/10 text-violet-700",
    icon: ShieldCheck,
  },
  admin: {
    label: "Administrador",
    pill: "border-blue-500/20 bg-blue-500/10 text-blue-700",
    icon: ShieldCheck,
  },
  staff: {
    label: "Equipe",
    pill: "border-emerald-500/20 bg-emerald-500/10 text-emerald-700",
    icon: UserRound,
  },
};

export interface AdminEventMemberCardProps {
  member: EventMemberI;
  index?: number;
  onRemove: (member: EventMemberI) => void;
}

export function AdminEventMemberCard({
  member,
  index = 0,
  onRemove,
}: AdminEventMemberCardProps) {
  const role = roleConfig[member.role];
  const RoleIcon = role.icon;
  const displayName = `Usuário ${member.user_id.slice(0, 8)}`;
  const initials = member.user_id.slice(0, 2).toUpperCase();

  return (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{
        delay: index * 0.05,
        duration: 0.35,
        ease: [0.25, 0.1, 0.25, 1],
      }}
      className={cn(
        "group relative flex w-full min-w-62.5 max-w-full flex-col overflow-hidden rounded-2xl bg-card text-left",
        "ring-1 ring-foreground/10 shadow-xs",
        "transform-gpu will-change-transform",
        "transition-all duration-300 ease-out",
        "hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm",
      )}
    >
      <div className="relative aspect-video overflow-hidden bg-muted">
        <div className="flex h-full w-full items-center justify-center bg-linear-to-br from-muted via-background to-muted/40">
          <div className="flex size-18 items-center justify-center rounded-full border border-border/70 bg-background/80 text-xl font-semibold text-muted-foreground/60 shadow-sm backdrop-blur-sm">
            {initials}
          </div>
        </div>

        <div className="absolute inset-0 bg-linear-to-t from-background/90 via-background/25 to-transparent" />

        <div className="absolute left-3 top-3">
          <span
            className={cn(
              "inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-medium backdrop-blur-sm",
              role.pill,
            )}
          >
            <RoleIcon className="size-3" />
            {role.label}
          </span>
        </div>

        <div className="absolute inset-x-0 bottom-0 p-3 sm:p-4">
          <h3 className="truncate text-base font-semibold leading-tight text-foreground transition-colors duration-300 group-hover:text-primary sm:text-lg">
            {displayName}
          </h3>
          <p className="mt-0.5 truncate text-xs leading-5 text-muted-foreground">
            Membro da equipe do evento
          </p>
        </div>
      </div>

      <div className="flex min-w-0 items-center justify-between gap-3 p-3 pt-2.5 sm:p-4 sm:pt-3">
        <span className="flex min-w-0 items-center gap-1 text-[11px] text-muted-foreground">
          <UserRound className="size-3 shrink-0" />
          <span className="truncate font-mono">{member.user_id}</span>
        </span>

        <Button
          type="button"
          variant="outline"
          size="sm"
          className="shrink-0 gap-1.5 text-destructive hover:text-destructive"
          onClick={() => onRemove(member)}
        >
          <Trash2 className="size-3.5" />
          Remover
        </Button>
      </div>
    </motion.article>
  );
}
