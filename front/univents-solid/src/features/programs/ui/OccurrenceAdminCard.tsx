import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

import CalendarDaysIcon from "~icons/lucide/calendar-days";
import ClockIcon from "~icons/lucide/clock";
import GiftIcon from "~icons/lucide/gift";
import InfinityIcon from "~icons/lucide/infinity";
import PencilIcon from "~icons/lucide/pencil";
import TrashIcon from "~icons/lucide/trash";
import UserCheckIcon from "~icons/lucide/user-check";
import UsersIcon from "~icons/lucide/users";

import { Button, cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";
import type { OccurrenceI } from "../model";

const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;
const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Gift = GiftIcon as unknown as (props: { class?: string }) => JSX.Element;
const InfinityLucide = InfinityIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash = TrashIcon as unknown as (props: { class?: string }) => JSX.Element;
const UserCheck = UserCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface OccurrenceAdminCardProps {
  occurrence: OccurrenceI;
  programKind?: "activity" | "checkpoint";
  index?: number;
  animate?: boolean;
  onEdit: (occurrence: OccurrenceI) => void;
  onDelete: (occurrence: OccurrenceI) => void;
  onAttendance: (occurrence: OccurrenceI) => void;
  onDraw?: (occurrence: OccurrenceI) => void;
}

export function OccurrenceAdminCard(props: OccurrenceAdminCardProps): JSX.Element {
  const startDate = () => new Date(props.occurrence.starts_at);
  const endDate = () => new Date(props.occurrence.ends_at);

  const formattedDate = () => {
    try {
      return new Intl.DateTimeFormat("pt-BR", {
        weekday: "short",
        day: "2-digit",
        month: "short",
      }).format(startDate()).replace(".", "");
    } catch {
      return startDate().toLocaleDateString("pt-BR");
    }
  };

  const formattedTime = () => {
    try {
      const s = startDate().toLocaleTimeString("pt-BR", {
        hour: "2-digit",
        minute: "2-digit",
      });
      const e = endDate().toLocaleTimeString("pt-BR", {
        hour: "2-digit",
        minute: "2-digit",
      });
      return `${s} – ${e}`;
    } catch {
      return "";
    }
  };

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate} class="h-full">
      <div
        role="region"
        aria-label={`Ocorrência ${formattedDate()}`}
        class={cn(
          "group relative flex h-full w-full min-w-0 flex-col justify-between overflow-hidden rounded-xl border border-border border-t-[3px] border-t-primary bg-card bg-linear-to-b from-primary/4 via-card to-card p-4 text-left shadow-xs transition-all duration-200",
          "hover:-translate-y-0.5 hover:border-foreground/20 hover:border-t-primary hover:shadow-md",
          "focus-within:ring-2 focus-within:ring-ring",
        )}
      >
        <div class="space-y-3">
          {/* Header with icon, date, time and top actions */}
          <div class="flex items-start justify-between gap-3">
            <div class="flex items-center gap-2.5 min-w-0 flex-1">
              <div class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary transition-transform group-hover:scale-105">
                <CalendarDays class="size-4.5" />
              </div>
              <div class="min-w-0 flex-1">
                <h3
                  class="truncate text-sm font-semibold tracking-tight text-foreground capitalize transition-colors group-hover:text-primary"
                  title={formattedDate()}
                >
                  {formattedDate()}
                </h3>
                <div class="mt-0.5 flex items-center gap-1.5 text-xs text-muted-foreground">
                  <Clock class="size-3 shrink-0" />
                  <span>{formattedTime()}</span>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-0.5 shrink-0">
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label="Editar ocorrência"
                title="Editar horário"
                onClick={() => props.onEdit(props.occurrence)}
                class="size-7 text-muted-foreground opacity-60 transition-opacity hover:bg-muted hover:text-foreground group-hover:opacity-100 cursor-pointer"
              >
                <Pencil class="size-3.5" />
              </Button>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label="Excluir ocorrência"
                title="Excluir horário"
                onClick={() => props.onDelete(props.occurrence)}
                class="size-7 text-muted-foreground opacity-60 transition-opacity hover:bg-destructive/15 hover:text-destructive group-hover:opacity-100 cursor-pointer"
              >
                <Trash class="size-3.5" />
              </Button>
            </div>
          </div>
        </div>

        {/* Specs & actions footer */}
        <div class="mt-4 flex items-center justify-between border-t border-border/60 pt-3 text-xs">
          <div class="flex items-center gap-1.5 text-[11px] text-muted-foreground">
            <Show
              when={props.occurrence.max_capacity != null && props.occurrence.max_capacity > 0}
              fallback={
                <span class="inline-flex items-center gap-1 text-muted-foreground" title="Vagas ilimitadas">
                  <InfinityLucide class="size-3" />
                  <span>Ilimitado</span>
                </span>
              }
            >
              <span class="inline-flex items-center gap-1 font-medium text-foreground/80">
                <Users class="size-3 text-muted-foreground shrink-0" />
                <span>{props.occurrence.max_capacity} vagas</span>
              </span>
            </Show>
          </div>

          <div class="flex items-center gap-1.5">
            <Show when={props.onDraw}>
              <Button
                type="button"
                variant="outline"
                size="sm"
                class="h-7 gap-1 px-2 text-xs cursor-pointer"
                onClick={() => props.onDraw?.(props.occurrence)}
                title="Sortear participantes"
              >
                <Gift class="size-3" />
                <span>Sortear</span>
              </Button>
            </Show>

            <Button
              type="button"
              variant="default"
              size="sm"
              class="h-7 gap-1 px-2.5 text-xs font-medium cursor-pointer"
              onClick={() => props.onAttendance(props.occurrence)}
              title="Lista e controle de presença"
            >
              <UserCheck class="size-3" />
              <span>Presença</span>
            </Button>
          </div>
        </div>
      </div>
    </Reveal>
  );
}
