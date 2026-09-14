import type { JSX } from "@solidjs/web";

import UserPlusIcon from "~icons/lucide/user-plus";

import { Reveal } from "@/shared/ui/Reveal";

const UserPlus = UserPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminAddMemberCardProps {
  index?: number;
  animate?: boolean;
  onAdd: () => void;
}

export function AdminAddMemberCard(props: AdminAddMemberCardProps): JSX.Element {
  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate} class="h-full">
      <button
        type="button"
        onClick={() => props.onAdd()}
        class="group relative flex h-full w-full min-w-0 flex-col justify-between rounded-xl bg-card p-3 text-left ring-1 ring-border transition-colors hover:ring-foreground/25 focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        <div class="flex items-center gap-2.5">
          <div class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground transition-colors group-hover:text-foreground">
            <UserPlus class="size-5" />
          </div>

          <div class="min-w-0 flex-1 space-y-0.5">
            <h3 class="truncate text-sm font-medium leading-tight text-foreground">
              Adicionar membro
            </h3>

            <span class="block truncate text-xs text-muted-foreground/80">
              Convidar para a equipe
            </span>
          </div>
        </div>

        <div class="mt-2.5 flex items-center justify-between border-t border-border/50 pt-2 text-[11px] text-muted-foreground/70">
          <span class="font-mono text-[11px]">Novo membro</span>
          <span>Acesso imediato</span>
        </div>
      </button>
    </Reveal>
  );
}
