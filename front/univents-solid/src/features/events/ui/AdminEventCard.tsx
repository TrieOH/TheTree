import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { Show } from "solid-js";

import MailIcon from "~icons/lucide/mail";
import PencilIcon from "~icons/lucide/pencil";

import { Button, cn } from "@trieoh/ui-solid";

import { Reveal } from "@/shared/ui/Reveal";

import type { EventI } from "../model";

const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as () => JSX.Element;

const STATUS_CONFIG: Record<EventI["status"], { label: string; dot: string }> = {
  draft: { label: "Rascunho", dot: "bg-amber-500" },
  active: { label: "Ativo", dot: "bg-emerald-500" },
  discontinued: { label: "Descontinuado", dot: "bg-rose-500" },
};

export interface AdminEventCardProps {
  event: EventI;
  index?: number;
  onEdit: (event: EventI) => void;
  /** Forwarded to `Reveal`; false when the row appeared because of a resize. */
  animate?: boolean;
}

export function AdminEventCard(props: AdminEventCardProps): JSX.Element {
  const status = () => STATUS_CONFIG[props.event.status];

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate}>
      <div class="group relative flex w-full min-w-0 items-center gap-3 rounded-xl bg-card p-3 ring-1 ring-border transition-colors hover:ring-foreground/25">
        <div class="relative size-10 shrink-0">
          <div class="flex size-10 items-center justify-center overflow-hidden rounded-lg bg-muted">
            <Show
              when={props.event.logo_url ?? props.event.banner_url}
              fallback={
                <span class="text-sm font-semibold text-muted-foreground/50">
                  {props.event.acronym?.charAt(0) ?? props.event.full_name.charAt(0)}
                </span>
              }
            >
              {(src) => <img src={src()} alt="" class="size-full object-cover" />}
            </Show>
          </div>
          <span
            title={status().label}
            class={cn(
              "absolute -bottom-0.5 -right-0.5 size-2.5 rounded-full ring-2 ring-card",
              status().dot,
            )}
          />
        </div>

        <div class="min-w-0 flex-1 space-y-0.5">
          <h3 class="truncate text-sm font-medium leading-tight">
            {props.event.full_name}
          </h3>
          <div class="flex min-w-0 items-center gap-1 text-xs text-muted-foreground">
            <Show
              when={props.event.contact_email}
              fallback={<span class="truncate">/{props.event.slug}</span>}
            >
              {(email) => (
                <span class="inline-flex min-w-0 max-w-full items-center gap-1">
                  <Mail class="shrink-0" />
                  <span class="truncate">{email()}</span>
                </span>
              )}
            </Show>
          </div>
        </div>

        <Link
          to="/admin/events/$eventId"
          params={{ eventId: props.event.id }}
          aria-label={props.event.full_name}
          class="absolute inset-0 z-0 rounded-xl focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />

        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label={`Editar ${props.event.full_name}`}
          onClick={() => props.onEdit(props.event)}
          class="relative z-10 size-7 shrink-0 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100 focus-visible:opacity-100 [@media(hover:none)]:opacity-100"
        >
          <Pencil />
        </Button>
      </div>
    </Reveal>
  );
}