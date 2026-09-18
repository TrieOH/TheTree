import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

import CalendarDaysIcon from "~icons/lucide/calendar-days";
import ClockIcon from "~icons/lucide/clock";
import FlagIcon from "~icons/lucide/flag";
import PencilIcon from "~icons/lucide/pencil";
import TrashIcon from "~icons/lucide/trash";

import { Button, cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";
import type { OccurrenceI, ProgramI } from "../model";

const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;
const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Flag = FlagIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash = TrashIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminProgramCardProps {
  program: ProgramI;
  occurrences?: OccurrenceI[];
  index?: number;
  animate?: boolean;
  onEdit: (program: ProgramI) => void;
  onDelete: (program: ProgramI) => void;
  onOpenCalendar?: (program: ProgramI) => void;
  onManageOccurrences?: (program: ProgramI) => void;
  occurrencesHref?: string;
}

export function AdminProgramCard(props: AdminProgramCardProps): JSX.Element {
  const isDeleted = () => props.program.deleted_at !== null;
  const occurrencesCount = () => props.occurrences?.length ?? 0;
  const isCheckpoint = () => props.program.kind === "checkpoint";

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate} class="h-full">
      <div
        role="region"
        aria-label={`Programa ${props.program.name}`}
        class={cn(
          "group relative flex h-full w-full min-w-0 flex-col justify-between overflow-hidden rounded-xl border border-border border-t-[3px] border-t-primary bg-card bg-linear-to-b from-primary/4 via-card to-card p-4 text-left shadow-xs transition-all duration-200",
          "hover:-translate-y-0.5 hover:border-foreground/20 hover:border-t-primary hover:shadow-md",
          "focus-within:ring-2 focus-within:ring-ring",
          isDeleted() && "opacity-60 grayscale",
        )}
      >
        <div class="space-y-2">
          {/* Header with thumbnail, name, metadata and actions */}
          <div class="flex items-start justify-between gap-3">
            <div class="flex items-center gap-2.5 min-w-0 flex-1">
              <div class="relative flex size-11 shrink-0 items-center justify-center overflow-hidden rounded-lg border border-border/70 bg-muted/60 transition-transform group-hover:scale-105">
                <Show
                  when={props.program.banner_url}
                  fallback={
                    <div class="flex size-full items-center justify-center bg-primary/10 text-primary">
                      <Show
                        when={isCheckpoint()}
                        fallback={<CalendarDays class="size-5" />}
                      >
                        <Flag class="size-5" />
                      </Show>
                    </div>
                  }
                >
                  {(url) => (
                    <img
                      src={url()}
                      alt={props.program.name}
                      class="size-full object-cover"
                    />
                  )}
                </Show>
              </div>

              <div class="min-w-0 flex-1">
                <h3
                  class="truncate text-sm font-semibold tracking-tight text-foreground transition-colors group-hover:text-primary"
                  title={props.program.name}
                >
                  {props.program.name}
                </h3>
                <div class="mt-0.5 flex min-w-0 items-center gap-1.5 text-xs text-muted-foreground">
                  <span class="shrink-0">{isCheckpoint() ? "Checkpoint" : "Atividade"}</span>
                  <Show when={props.program.min_access_level != null && props.program.min_access_level > 0}>
                    <span class="text-muted-foreground/60">•</span>
                    <span class="truncate">Nível {props.program.min_access_level}</span>
                  </Show>
                  <Show when={props.program.staff_only}>
                    <span class="text-muted-foreground/60">•</span>
                    <span class="font-medium text-destructive shrink-0">Staff</span>
                  </Show>
                </div>
              </div>
            </div>

            {/* Actions */}
            <div class="flex items-center gap-0.5 shrink-0">
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label={`Editar ${props.program.name}`}
                title="Editar programa"
                onClick={() => props.onEdit(props.program)}
                class="size-7 text-muted-foreground opacity-60 transition-opacity hover:bg-muted hover:text-foreground group-hover:opacity-100 cursor-pointer"
              >
                <Pencil class="size-3.5" />
              </Button>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label={`Excluir ${props.program.name}`}
                title="Excluir programa"
                onClick={() => props.onDelete(props.program)}
                class="size-7 text-muted-foreground opacity-60 transition-opacity hover:bg-destructive/15 hover:text-destructive group-hover:opacity-100 cursor-pointer"
              >
                <Trash class="size-3.5" />
              </Button>
            </div>
          </div>

          <p
            class="line-clamp-2 text-xs leading-relaxed text-muted-foreground min-h-8"
            title={props.program.description ?? ""}
          >
            {props.program.description || "Nenhuma descrição fornecida para este programa."}
          </p>
        </div>

        {/* Specs footer */}
        <div class="mt-3 flex items-center justify-between border-t border-border/60 pt-2.5 text-[11px] text-muted-foreground">
          <div class="flex items-center gap-1.5 font-medium text-foreground/80">
            <Clock class="size-3 text-muted-foreground shrink-0" />
            <span>
              {occurrencesCount()} {occurrencesCount() === 1 ? "ocorrência" : "ocorrências"}
            </span>
          </div>

          <Show
            when={props.occurrencesHref}
            fallback={
              <button
                type="button"
                onClick={() => props.onManageOccurrences?.(props.program)}
                class="inline-flex items-center gap-1 text-[11px] font-medium text-primary hover:underline cursor-pointer"
                title="Gerenciar ocorrências"
              >
                <CalendarDays class="size-3 shrink-0" />
                <span>Ocorrências</span>
              </button>
            }
          >
            {(href) => (
              <a
                href={href()}
                onClick={(e) => {
                  if (props.onManageOccurrences) {
                    e.preventDefault();
                    props.onManageOccurrences(props.program);
                  }
                }}
                class="inline-flex items-center gap-1 text-[11px] font-medium text-primary hover:underline cursor-pointer"
                title="Gerenciar ocorrências"
              >
                <CalendarDays class="size-3 shrink-0" />
                <span>Ocorrências</span>
              </a>
            )}
          </Show>
        </div>
      </div>
    </Reveal>
  );
}
