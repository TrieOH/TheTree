import type { JSX } from "@solidjs/web";

import PlusIcon from "~icons/lucide/plus";

import { Reveal } from "@/shared/ui/Reveal";

const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateEventCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateEventCard(
  props: AdminCreateEventCardProps,
): JSX.Element {
  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate}>
      <button
        type="button"
        onClick={props.onCreate}
        class="group relative flex w-full min-w-0 items-center gap-3 rounded-xl bg-card p-3 text-left ring-1 ring-border transition-colors hover:ring-foreground/25 focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        <div class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground transition-colors group-hover:text-foreground">
          <Plus class="size-5" />
        </div>

        <div class="min-w-0 flex-1 space-y-0.5">
          <h3 class="truncate text-sm font-medium leading-tight">
            Novo evento
          </h3>

          <span class="block truncate text-xs text-muted-foreground">
            Criar um novo evento
          </span>
        </div>
      </button>
    </Reveal>
  );
}