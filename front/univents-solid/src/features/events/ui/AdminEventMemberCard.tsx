import type { JSX } from "@solidjs/web";
import { Show, createMemo, createSignal } from "solid-js";

import CheckIcon from "~icons/lucide/check";
import CopyIcon from "~icons/lucide/copy";
import TrashIcon from "~icons/lucide/trash";

import { Button, cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";
import { toast } from "@/shared/ui/toast";
import { formatMemberViewModel, type EventMemberWithEmailI } from "../model/member";

const Trash = TrashIcon as unknown as (props: { class?: string }) => JSX.Element;
const Copy = CopyIcon as unknown as (props: { class?: string }) => JSX.Element;
const Check = CheckIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminEventMemberCardProps {
  member: EventMemberWithEmailI;
  index?: number;
  animate?: boolean;
  onRemove: (member: EventMemberWithEmailI) => void;
}

export function AdminEventMemberCard(props: AdminEventMemberCardProps): JSX.Element {
  const [copied, setCopied] = createSignal(false);
  const vm = createMemo(() => formatMemberViewModel(props.member));

  const copyUserId = async (e: MouseEvent) => {
    e.stopPropagation();
    try {
      await navigator.clipboard.writeText(props.member.user_id);
      setCopied(true);
      toast.success("ID do usuário copiado!");
      setTimeout(() => setCopied(false), 2000);
    } catch {
      toast.error("Erro ao copiar ID.");
    }
  };

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate} class="h-full">
      <div class="group relative flex h-full w-full min-w-0 flex-col justify-between rounded-xl bg-card p-3 ring-1 ring-border transition-colors hover:ring-foreground/25">
        {/* Main identity row: Avatar + Identity Texts (Email + Role subtitle) + Trash */}
        <div class="flex items-center gap-2.5">
          {/* Avatar with Role Status Dot */}
          <div class="relative size-10 shrink-0 select-none">
            <div class="flex size-10 items-center justify-center overflow-hidden rounded-lg bg-muted text-xs font-semibold text-muted-foreground/70">
              <Show
                when={vm().pfpUrl}
                fallback={<span>{vm().initials}</span>}
              >
                {(src) => (
                  <img
                    src={src()}
                    alt={`Avatar de ${vm().primaryLabel}`}
                    class="size-full object-cover"
                  />
                )}
              </Show>
            </div>
            <span
              title={vm().roleLabel}
              class={cn(
                "absolute -bottom-0.5 -right-0.5 size-2.5 rounded-full ring-2 ring-card",
                vm().roleDotClass,
              )}
            />
          </div>

          {/* Member Name/Email + Role subtitle */}
          <div class="min-w-0 flex-1 space-y-0.5">
            <h3
              class="truncate text-sm font-medium leading-tight text-foreground"
              title={props.member.email ?? props.member.user_id}
            >
              {vm().primaryLabel}
            </h3>

            <span class="block truncate text-xs text-muted-foreground/80">
              {vm().roleLabel}
            </span>
          </div>

          {/* Remove Member Button */}
          <Button
            type="button"
            variant="ghost"
            size="icon"
            aria-label={`Remover ${vm().primaryLabel}`}
            title={`Remover ${vm().primaryLabel}`}
            onClick={() => props.onRemove(props.member)}
            class="relative z-10 size-7 shrink-0 text-muted-foreground opacity-0 transition-opacity hover:bg-destructive/15 hover:text-destructive group-hover:opacity-100 focus-visible:opacity-100 [@media(hover:none)]:opacity-100"
          >
            <Trash class="size-3.5" />
          </Button>
        </div>

        {/* Footer info: User ID copy + Addition date */}
        <div class="mt-2.5 flex items-center justify-between border-t border-border/50 pt-2 text-[11px] text-muted-foreground">
          <button
            type="button"
            onClick={copyUserId}
            class="group/id inline-flex items-center gap-1 font-mono text-muted-foreground transition-colors hover:text-foreground"
            title="Clique para copiar o ID completo"
          >
            <span>ID: {vm().shortId}</span>
            <Show when={copied()} fallback={<Copy class="size-2.5 opacity-40 group-hover/id:opacity-100" />}>
              <Check class="size-2.5 text-emerald-500" />
            </Show>
          </button>

          <span title={`Adicionado em ${props.member.created_at}`}>
            {vm().formattedDate}
          </span>
        </div>
      </div>
    </Reveal>
  );
}
