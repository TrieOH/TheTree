import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

import MailIcon from "~icons/lucide/mail";
import PencilIcon from "~icons/lucide/pencil";

import { cn } from "@trieoh/ui-solid";

import { Reveal } from "@/shared/ui/Reveal";

import type { EventI } from "../model";

const Mail = MailIcon as unknown as () => JSX.Element;
const Pencil = PencilIcon as unknown as () => JSX.Element;

/** Same labels and accent colours as the React card. */
const STATUS_CONFIG: Record<EventI["status"], { label: string; dot: string }> = {
  draft: { label: "Rascunho", dot: "bg-amber-500" },
  active: { label: "Ativo", dot: "bg-emerald-500" },
  discontinued: { label: "Descontinuado", dot: "bg-rose-500" },
};

export interface AdminEventCardProps {
  event: EventI;
  index?: number;
  onEdit: (event: EventI) => void;
  /** Opens the event dashboard. Omitted while that route does not exist yet. */
  onOpen?: (event: EventI) => void;
  /** Forwarded to `Reveal`; false when the row appeared because of a resize. */
  animate?: boolean;
}

export function AdminEventCard(props: AdminEventCardProps): JSX.Element {
  const status = () => STATUS_CONFIG[props.event.status];
  const visual = () => props.event.banner_url ?? props.event.logo_url ?? "";
  const open = () => props.onOpen?.(props.event);

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate}>
      <article
        class={cn(
          "group relative flex w-full min-w-60 max-w-full flex-col overflow-hidden rounded-2xl bg-card text-left",
          "ring-1 ring-foreground/10 shadow-xs",
          "transform-gpu will-change-transform",
          "transition-all duration-300 ease-out",
          "hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm",
          props.onOpen && "cursor-pointer focus:outline-none",
        )}
        role={props.onOpen ? "button" : undefined}
        tabindex={props.onOpen ? 0 : undefined}
        onClick={() => open()}
        onKeyDown={(event) => {
          if (!props.onOpen) return;
          if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            open();
          }
        }}
      >
        <div class="relative aspect-video overflow-hidden bg-muted">
          <Show
            when={visual()}
            fallback={
              <div class="flex h-full w-full items-center justify-center bg-linear-to-br from-muted via-background to-muted/40">
                <div class="flex size-20 items-center justify-center rounded-full border border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
                  <span class="text-2xl font-semibold text-muted-foreground/40">
                    {props.event.acronym ?? props.event.full_name.charAt(0)}
                  </span>
                </div>
              </div>
            }
          >
            {(src) => (
              <img
                src={src()}
                alt="Representação Visual do Evento"
                loading={(props.index ?? 0) < 4 ? "eager" : "lazy"}
                class="h-full w-full object-cover transition-transform duration-700 ease-out group-hover:scale-105"
              />
            )}
          </Show>

          <div class="absolute left-4 top-4 flex flex-wrap items-center gap-2">
            <span class="inline-flex items-center gap-1 rounded-full border border-white/30 bg-black/75 px-2.5 py-1 text-[11px] font-semibold text-white shadow-lg backdrop-blur-sm drop-shadow-[0_1px_2px_rgba(0,0,0,0.9)]">
              <span class={cn("size-1.5 rounded-full", status().dot)} />
              {status().label}
            </span>
          </div>

          <div class="absolute right-4 top-4 z-20">
            <button
              type="button"
              aria-label={`Editar ${props.event.full_name}`}
              onClick={(event) => {
                event.stopPropagation();
                props.onEdit(props.event);
              }}
              class="inline-flex items-center gap-1.5 rounded-lg bg-background/85 px-3 py-1.5 text-xs font-medium text-foreground shadow-sm backdrop-blur-sm transition-colors hover:bg-background"
            >
              <Pencil />
              Editar
            </button>
          </div>

          <div class="absolute inset-x-0 bottom-0 flex items-end justify-between gap-3 p-4 sm:p-5">
            <div class="min-w-0 space-y-1">
              <h3 class="line-clamp-2 text-balance text-lg font-semibold leading-snug text-white drop-shadow-[0_1px_3px_rgba(0,0,0,0.9)] transition-colors duration-300 group-hover:text-white sm:text-xl">
                {props.event.full_name}
              </h3>
              <Show when={props.event.description}>
                <p class="line-clamp-1 max-w-[min(65%,32rem)] truncate text-xs text-white/90 drop-shadow-[0_1px_2px_rgba(0,0,0,0.9)]">
                  {props.event.description}
                </p>
              </Show>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between gap-3 p-4 pt-3 sm:p-5 sm:pt-4">
          <div class="min-w-0 flex-1 space-y-1">
            <div class="flex min-w-0 items-center text-xs text-muted-foreground">
              <span
                class="inline-flex min-w-0 max-w-full items-center gap-1.5"
                title={props.event.contact_email ?? "Sem contato cadastrado"}
              >
                <Mail />
                <span class="truncate">
                  {props.event.contact_email ?? "Sem contato"}
                </span>
              </span>
            </div>
            <code class="block truncate font-mono text-[11px] text-muted-foreground/80">
              {props.event.slug}
            </code>
          </div>
        </div>
      </article>
    </Reveal>
  );
}

