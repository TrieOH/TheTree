import type { EditionI } from "@/features/editions/model";
import CalendarIcon from "~icons/lucide/calendar";
import MapPinIcon from "~icons/lucide/map-pin";
import ArrowRightIcon from "~icons/lucide/arrow-right";
import type { JSX } from "@solidjs/web";
const Calendar = CalendarIcon as unknown as () => JSX.Element;
const MapPin = MapPinIcon as unknown as () => JSX.Element;
const ArrowRight = ArrowRightIcon as unknown as () => JSX.Element;

export function EditionSummaryCard(props: { edition: EditionI }) {
  const closed = () => new Date(props.edition.ends_at).getTime() < Date.now();
  const date = () =>
    `${new Date(props.edition.starts_at).toLocaleDateString("pt-BR", { day: "2-digit", month: "short" }).replace(".", "")} – ${new Date(props.edition.ends_at).toLocaleDateString("pt-BR", { day: "2-digit", month: "short", year: "numeric" }).replace(".", "")}`;
  return (
    <article class="group flex min-w-64 max-w-96 flex-1 overflow-hidden rounded-xl border border-border/60 bg-card transition-all duration-200 hover:border-border hover:shadow-md">
      <div class="relative w-24 shrink-0 overflow-hidden bg-muted">
        {(props.edition.banner_url ?? props.edition.logo_url) ? (
          <img
            src={props.edition.banner_url ?? props.edition.logo_url ?? ""}
            alt={props.edition.name}
            class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
        ) : (
          <div class="flex h-full items-center justify-center bg-linear-to-br from-muted to-muted/50">
            <span class="text-xl font-bold text-muted-foreground/50">
              {props.edition.name.slice(0, 2).toUpperCase()}
            </span>
          </div>
        )}
      </div>
      <div class="flex min-w-0 flex-1 flex-col justify-between p-4">
        <div class="flex items-start justify-between gap-2">
          <h3 class="line-clamp-2 text-sm font-semibold leading-tight sm:text-base">
            {props.edition.name}
          </h3>
          <span
            class={`inline-block shrink-0 rounded px-1.5 py-0.5 text-[9px] font-bold tracking-wide sm:text-[10px] ${closed() ? "bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400" : "bg-blue-100 text-blue-700 dark:bg-blue-950/40 dark:text-blue-400"}`}
          >
            {closed() ? "CLOSED" : "UPCOMING"}
          </span>
        </div>
        <div class="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-muted-foreground">
          <span class="flex items-center gap-1">
            <Calendar />
            <span class="text-[11px] sm:text-xs">{date()}</span>
          </span>
          {props.edition.location_name && (
            <span class="flex min-w-0 items-center gap-1">
              <MapPin />
              <span class="max-w-35 truncate text-[11px] sm:text-xs">
                {props.edition.location_name}
              </span>
            </span>
          )}
        </div>
        <span class="mt-2 inline-flex items-center gap-1 text-xs font-semibold text-primary sm:mt-3 sm:text-sm">
          {closed() ? "View Archives" : "Learn More"}
          <ArrowRight />
        </span>
      </div>
    </article>
  );
}
