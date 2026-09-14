import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { For, Show } from "solid-js";
import ChevronRightIcon from "~icons/lucide/chevron-right";
import Layers3Icon from "~icons/lucide/layers-3";
import type { EditionI } from "@/features/editions/model";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideChevronRight = ChevronRightIcon as unknown as IconComp;
const LucideLayers3 = Layers3Icon as unknown as IconComp;

const statusConfig = {
  draft: {
    label: "Rascunho",
    className: "border-amber-500/20 bg-amber-500/10 text-amber-700 dark:text-amber-400",
    dot: "bg-amber-500",
  },
  future: {
    label: "Futura",
    className: "border-sky-500/20 bg-sky-500/10 text-sky-700 dark:text-sky-400",
    dot: "bg-sky-500",
  },
  active: {
    label: "Ativa",
    className: "border-emerald-500/20 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400",
    dot: "bg-emerald-500",
  },
  past: {
    label: "Encerrada",
    className: "border-slate-500/20 bg-slate-500/10 text-slate-600 dark:text-slate-400",
    dot: "bg-slate-500",
  },
} as const;

export function EventEditionsList(props: {
  eventId: string;
  editions: EditionI[];
}): JSX.Element {
  return (
    <section class="order-7 space-y-3">
      <div class="flex min-w-0 items-center justify-between gap-3 px-1">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <div class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
            <LucideLayers3 class="size-5" />
          </div>
          <div class="min-w-0 flex-1">
            <h2 class="text-base font-semibold tracking-tight text-foreground">
              Edições
            </h2>
            <p class="truncate text-xs text-muted-foreground">
              Acesse rapidamente cada edição do evento.
            </p>
          </div>
        </div>
        <span class="flex size-8 shrink-0 items-center justify-center rounded-full border border-border bg-muted/50 p-0 text-xs font-medium text-muted-foreground">
          {props.editions.length}
        </span>
      </div>

      <div class="overflow-hidden rounded-xl border border-border bg-card shadow-xs">
        <Show
          when={props.editions.length > 0}
          fallback={
            <p class="p-5 text-sm text-muted-foreground">
              Nenhuma edição cadastrada.
            </p>
          }
        >
          <For each={props.editions.slice(0, 6)}>
            {(edition, index) => {
              const status = () =>
                statusConfig[edition.status] ?? statusConfig.draft;

              const formattedDate = () => {
                const date = new Date(edition.starts_at);
                return date.toLocaleDateString("pt-BR", {
                  day: "2-digit",
                  month: "short",
                  year: "numeric",
                });
              };

              return (
                <Link
                  to="/admin/events/$eventId/editions/$editionId"
                  params={{ eventId: props.eventId, editionId: edition.id }}
                  class={`group flex items-center justify-between gap-4 px-5 py-4 transition-colors hover:bg-muted/40 ${index() > 0 ? "border-t border-border" : ""
                    }`}
                >
                  <div class="flex min-w-0 items-center gap-3">
                    <span class="flex size-8 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">
                      {index() + 1}
                    </span>
                    <div class="min-w-0">
                      <p class="truncate text-sm font-semibold text-foreground">
                        {edition.name}
                      </p>
                      <p class="mt-0.5 truncate text-xs text-muted-foreground">
                        {formattedDate()}
                        {edition.location_name
                          ? ` · ${edition.location_name}`
                          : ""}
                      </p>
                    </div>
                  </div>
                  <div class="flex shrink-0 items-center gap-3">
                    <span
                      class={`hidden items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-[11px] font-medium sm:inline-flex ${status().className
                        }`}
                    >
                      <span class={`size-1.5 rounded-full ${status().dot}`} />
                      {status().label}
                    </span>
                    <LucideChevronRight class="size-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
                  </div>
                </Link>
              );
            }}
          </For>
        </Show>
      </div>
    </section>
  );
}
