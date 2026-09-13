import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import CalendarIcon from "~icons/lucide/calendar";
import type { EventI } from "../model";
import { EventVisualCard } from "./EventVisualCard";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCalendar = CalendarIcon as unknown as IconComp;

const statusConfig = {
  draft: {
    label: "Rascunho",
    className: "bg-amber-500/10 text-amber-700",
  },
  active: {
    label: "Ativo",
    className: "bg-emerald-500/10 text-emerald-700",
  },
  archived: {
    label: "Arquivado",
    className: "bg-slate-500/10 text-slate-700",
  },
  discontinued: {
    label: "Descontinuado",
    className: "bg-rose-500/10 text-rose-700",
  },
} as const;

export function EventOverviewHeader(props: { event: EventI | null }): JSX.Element {
  const event = () => props.event;

  const formattedDate = () => {
    const e = event();
    if (!e?.created_at) return "";
    return new Date(e.created_at)
      .toLocaleDateString("pt-BR", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      })
      .replace(".", "");
  };

  const status = () => {
    const e = event();
    if (!e) return statusConfig.draft;
    return statusConfig[e.status] ?? statusConfig.draft;
  };

  return (
    <Show when={event()}>
      {(ev) => (
        <>
          <EventVisualCard event={ev()} />
          <div class="space-y-1 px-1 text-center">
            <h1 class="text-xl font-medium tracking-tight text-foreground/90">
              {ev().full_name}
            </h1>
            <p class="flex items-center justify-center gap-1.5 text-xs text-muted-foreground">
              <LucideCalendar class="size-3.5" />
              Criado em {formattedDate()}
            </p>
            <span
              class={`inline-flex w-fit items-center justify-center rounded-4xl border-0 px-2 py-0.5 text-xs font-normal ${status().className}`}
            >
              {status().label}
            </span>
          </div>
        </>
      )}
    </Show>
  );
}
