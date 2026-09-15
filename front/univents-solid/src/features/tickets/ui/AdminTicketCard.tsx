import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

import InfinityIcon from "~icons/lucide/infinity";
import KeyRoundIcon from "~icons/lucide/key-round";
import PencilIcon from "~icons/lucide/pencil";
import TicketIcon from "~icons/lucide/ticket";
import UsersIcon from "~icons/lucide/users";

import { Button, cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";
import type { TicketI } from "../model";

const InfinityLucide = InfinityIcon as unknown as (props: { class?: string }) => JSX.Element;
const KeyRound = KeyRoundIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (props: { class?: string }) => JSX.Element;
const Ticket = TicketIcon as unknown as (props: { class?: string }) => JSX.Element;
const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminTicketCardProps {
  ticket: TicketI;
  index?: number;
  animate?: boolean;
  onEdit: (ticket: TicketI) => void;
}

function formatPrice(cents: number): string {
  if (cents === 0) return "Gratuito";
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(cents / 100);
}

export function AdminTicketCard(props: AdminTicketCardProps): JSX.Element {
  const isFree = () => props.ticket.price_cents === 0;

  const handleEdit = () => props.onEdit(props.ticket);

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate} class="h-full">
      <div
        role="button"
        tabindex={0}
        onClick={handleEdit}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            handleEdit();
          }
        }}
        class={cn(
          "group relative flex h-full w-full min-w-0 flex-col justify-between overflow-hidden rounded-xl border border-border border-t-[3px] border-t-primary bg-card bg-linear-to-b from-primary/[0.04] via-card to-card p-4 text-left shadow-xs transition-all duration-200",
          "hover:-translate-y-0.5 hover:border-foreground/20 hover:border-t-primary hover:shadow-md cursor-pointer",
          "focus:outline-none focus-visible:ring-2 focus-visible:ring-ring",
        )}
      >
        {/* Header: Title + Price + Edit button */}
        <div class="space-y-2">
          <div class="flex items-start justify-between gap-3">
            <div class="flex items-center gap-2.5 min-w-0 flex-1">
              <div class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary transition-transform group-hover:scale-105">
                <Ticket class="size-4.5" />
              </div>
              <div class="min-w-0 flex-1">
                <h3
                  class="truncate text-sm font-semibold tracking-tight text-foreground transition-colors group-hover:text-primary"
                  title={props.ticket.name}
                >
                  {props.ticket.name}
                </h3>
                <p class={cn("text-xs font-semibold mt-0.5", isFree() ? "text-emerald-500 dark:text-emerald-400 font-medium" : "text-foreground")}>
                  {formatPrice(props.ticket.price_cents)}
                </p>
              </div>
            </div>

            <Button
              type="button"
              variant="ghost"
              size="icon"
              aria-label={`Editar ${props.ticket.name}`}
              title={`Editar ${props.ticket.name}`}
              onClick={(e) => {
                e.stopPropagation();
                handleEdit();
              }}
              class="size-7 shrink-0 text-muted-foreground opacity-60 transition-opacity hover:bg-muted hover:text-foreground group-hover:opacity-100"
            >
              <Pencil class="size-3.5" />
            </Button>
          </div>

          <p
            class="line-clamp-2 text-xs leading-relaxed text-muted-foreground min-h-[2rem]"
            title={props.ticket.description ?? ""}
          >
            {props.ticket.description || "Nenhuma descrição fornecida para este ingresso."}
          </p>
        </div>

        {/* Clean, sober specs footer */}
        <div class="mt-3 flex items-center justify-between border-t border-border/60 pt-2.5 text-[11px] text-muted-foreground">
          <div class="flex items-center gap-1.5 font-medium text-foreground/80">
            <KeyRound class="size-3 text-muted-foreground shrink-0" />
            <span>Acesso nível {props.ticket.access_level}</span>
          </div>

          <div class="flex items-center gap-1.5">
            <Show
              when={props.ticket.max_quantity != null}
              fallback={
                <span class="inline-flex items-center gap-1 text-muted-foreground" title="Vagas ilimitadas">
                  <InfinityLucide class="size-3" />
                  <span>Ilimitado</span>
                </span>
              }
            >
              <span class="inline-flex items-center gap-1 text-muted-foreground">
                <Users class="size-3 shrink-0" />
                <span>{props.ticket.max_quantity} vagas</span>
              </span>
            </Show>
          </div>
        </div>
      </div>
    </Reveal>
  );
}
