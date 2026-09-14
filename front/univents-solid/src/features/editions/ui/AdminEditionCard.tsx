import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

import CalendarDaysIcon from "~icons/lucide/calendar-days";
import MapPinIcon from "~icons/lucide/map-pin";
import PencilIcon from "~icons/lucide/pencil";

import { Button, cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";
import type { EditionI, EditionStatus } from "../model";

const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;
const MapPin = MapPinIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (props: { class?: string }) => JSX.Element;

const STATUS_CONFIG: Record<
  EditionStatus,
  { label: string; dot: string; pill: string }
> = {
  draft: {
    label: "Rascunho",
    dot: "bg-amber-500",
    pill: "border-amber-500/25 bg-amber-500/10 text-amber-700 dark:text-amber-300",
  },
  future: {
    label: "Futura",
    dot: "bg-sky-500",
    pill: "border-sky-500/25 bg-sky-500/10 text-sky-700 dark:text-sky-300",
  },
  active: {
    label: "Ativa",
    dot: "bg-emerald-500",
    pill: "border-emerald-500/25 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300",
  },
  past: {
    label: "Encerrada",
    dot: "bg-slate-500",
    pill: "border-slate-500/25 bg-slate-500/10 text-slate-700 dark:text-slate-300",
  },
};

export interface AdminEditionCardProps {
  edition: EditionI;
  eventId: string;
  index?: number;
  animate?: boolean;
  onEdit: (edition: EditionI) => void;
}

function formatDateRange(startsAt: string, endsAt: string): string {
  try {
    const start = new Date(startsAt);
    const end = new Date(endsAt);
    const startStr = new Intl.DateTimeFormat("pt-BR", {
      day: "2-digit",
      month: "short",
    }).format(start).replace(".", "");
    const endStr = new Intl.DateTimeFormat("pt-BR", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    }).format(end).replace(".", "");
    return `${startStr} – ${endStr}`;
  } catch {
    return "";
  }
}

export function AdminEditionCard(props: AdminEditionCardProps): JSX.Element {
  const status = () => STATUS_CONFIG[props.edition.status] ?? STATUS_CONFIG.draft;
  const dateRange = () => formatDateRange(props.edition.starts_at, props.edition.ends_at);

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate} class="h-full">
      <div class="group relative flex h-full w-full min-w-0 flex-col justify-between rounded-xl bg-card p-3 ring-1 ring-border transition-colors hover:ring-foreground/25">
        {/* Main identity row: Banner/Logo + Name & Location + Status & Actions */}
        <div class="flex items-center gap-2.5">
          {/* Visual with Status Dot */}
          <div class="relative size-10 shrink-0 select-none">
            <div class="flex size-10 items-center justify-center overflow-hidden rounded-lg bg-muted text-xs font-semibold text-muted-foreground/70">
              <Show
                when={props.edition.banner_url ?? props.edition.logo_url}
                fallback={
                  <span class="text-sm font-semibold text-muted-foreground/50">
                    {props.edition.name.slice(0, 2).toUpperCase()}
                  </span>
                }
              >
                {(src) => (
                  <img
                    src={src()}
                    alt={props.edition.name}
                    class="size-full object-cover"
                  />
                )}
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

          {/* Edition Name + Subtitle (location or status) */}
          <div class="min-w-0 flex-1 space-y-0.5">
            <h3 class="truncate text-sm font-medium leading-tight text-foreground">
              {props.edition.name}
            </h3>

            <div class="flex min-w-0 items-center gap-1 text-xs text-muted-foreground/80">
              <Show
                when={props.edition.location_name}
                fallback={<span class="truncate">{status().label}</span>}
              >
                {(loc) => (
                  <span class="inline-flex min-w-0 max-w-full items-center gap-1 truncate">
                    <MapPin class="size-3 shrink-0" />
                    <span class="truncate">{loc()}</span>
                  </span>
                )}
              </Show>
            </div>
          </div>

          {/* Edit button */}
          <Button
            type="button"
            variant="ghost"
            size="icon"
            aria-label={`Editar ${props.edition.name}`}
            title={`Editar ${props.edition.name}`}
            onClick={(e) => {
              e.stopPropagation();
              props.onEdit(props.edition);
            }}
            class="relative z-10 size-7 shrink-0 text-muted-foreground opacity-0 transition-opacity hover:bg-muted group-hover:opacity-100 focus-visible:opacity-100 [@media(hover:none)]:opacity-100"
          >
            <Pencil class="size-3.5" />
          </Button>
        </div>

        {/* Link covering the card */}
        <a
          href={`/admin/events/${props.eventId}/editions/${props.edition.id}`}
          aria-label={props.edition.name}
          class="absolute inset-0 z-0 rounded-xl focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />

        {/* Footer info: Date range + Slug */}
        <div class="mt-2.5 flex items-center justify-between border-t border-border/50 pt-2 text-[11px] text-muted-foreground">
          <span class="inline-flex items-center gap-1">
            <CalendarDays class="size-3 shrink-0 opacity-70" />
            <span class="truncate">{dateRange()}</span>
          </span>

          <span class="truncate font-mono text-[11px]">
            /{props.edition.slug}
          </span>
        </div>
      </div>
    </Reveal>
  );
}
