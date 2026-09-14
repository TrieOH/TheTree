import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import CalendarIcon from "~icons/lucide/calendar";
import { cn } from "@trieoh/ui-solid";
import type { EditionI, EditionStatus } from "../model";
import { EditionVisualCard } from "./EditionVisualCard";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCalendar = CalendarIcon as unknown as IconComp;

const STATUS_CONFIG: Record<
  EditionStatus,
  { label: string; className: string }
> = {
  draft: {
    label: "Rascunho",
    className: "bg-amber-500/10 text-amber-700 dark:text-amber-300",
  },
  future: {
    label: "Futura",
    className: "bg-sky-500/10 text-sky-700 dark:text-sky-300",
  },
  active: {
    label: "Ativa",
    className: "bg-emerald-500/10 text-emerald-700 dark:text-emerald-300",
  },
  past: {
    label: "Encerrada",
    className: "bg-slate-500/10 text-slate-700 dark:text-slate-300",
  },
};

export function EditionOverviewHeader(props: {
  edition: EditionI | null;
  eventId: string;
}): JSX.Element {
  const edition = () => props.edition;

  const formattedDate = () => {
    const e = edition();
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
    const e = edition();
    if (!e) return STATUS_CONFIG.draft;
    return STATUS_CONFIG[e.status] ?? STATUS_CONFIG.draft;
  };

  return (
    <Show when={edition()}>
      {(ed) => (
        <>
          <EditionVisualCard edition={ed()} eventId={props.eventId} />
          <div class="space-y-1 px-1 text-center">
            <h1 class="text-xl font-medium tracking-tight text-foreground/90">
              {ed().name}
            </h1>
            <p class="flex items-center justify-center gap-1.5 text-xs text-muted-foreground">
              <LucideCalendar class="size-3.5" />
              Criada em {formattedDate()}
            </p>
            <div>
              <span
                class={cn(
                  "inline-flex w-fit items-center justify-center rounded-4xl border-0 px-2 py-0.5 text-xs font-normal",
                  status().className,
                )}
              >
                {status().label}
              </span>
            </div>
          </div>
        </>
      )}
    </Show>
  );
}
