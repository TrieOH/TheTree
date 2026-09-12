import { createFileRoute } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core-solid";
import { Loading, createMemo, untrack } from "solid-js";
import type { JSX } from "@solidjs/web";
import CalendarIcon from "~icons/lucide/calendar";
import MapPinIcon from "~icons/lucide/map-pin";
import ShareIcon from "~icons/lucide/share-2";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import { ContactSection } from "@/features/events/ui/ContactSection";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import type { EditionI } from "@/features/editions/model";
import { EditionSummaryCard } from "@/features/editions/ui/EditionSummaryCard";
import { EventCatalog } from "@/features/events/ui/EventCatalog";
import { EventCart } from "@/features/products/ui/EventCart";
import { handleShare } from "@/shared/lib/share";

const Calendar = CalendarIcon as unknown as () => JSX.Element;
const MapPin = MapPinIcon as unknown as () => JSX.Element;
const Share = ShareIcon as unknown as () => JSX.Element;

export const Route = createFileRoute("/events/$slug/")({
  head: ({ params }) => ({
    meta: [{ title: `${params.slug} - Univents` }],
  }),
  component: EventPage,
});

function EventPage() {
  const params = Route.useParams();
  const queryClient = useQueryClient();
  const event = createMemo(() =>
    queryClient.fetchQuery(publicEventBySlugQueryOptions(params().slug)),
  );
  return (
    <Loading fallback={<div class="min-h-screen animate-pulse bg-muted" />}>
      <EventContent event={event()} />
    </Loading>
  );
}

function EventContent(props: { event: EventI | null }) {
  const event = untrack(() => props.event);
  if (!event) return <div class="p-12 text-center">Evento não encontrado.</div>;
  return (
    <main class="min-h-screen bg-background pb-24">
      <div class="relative">
        <div class="relative h-40 w-full border-b-4 border-b-accent min-[300px]:h-48 sm:h-52 md:h-64">
          {event.banner_url ? (
            <img
              src={event.banner_url}
              alt={event.full_name}
              class="h-full w-full object-cover"
            />
          ) : (
            <div class="h-full w-full bg-linear-to-br from-muted via-primary/25 to-secondary/25" />
          )}
        </div>
        <button
          type="button"
          aria-label="Compartilhar evento"
          class="absolute right-4 top-4 z-10 rounded-full border border-border/50 bg-background/80 p-2 shadow-sm backdrop-blur-sm transition-all duration-200 hover:scale-105 hover:bg-background sm:right-6 sm:top-6 sm:p-2.5"
          onClick={() => void handleShare(event.full_name)}
        >
          <Share />
        </button>
        <div class="absolute inset-x-0 top-full z-10 flex -translate-y-1/2 justify-center">
          <div class="flex size-37.5 items-center justify-center rounded-full border-4 border-accent bg-primary shadow-lg sm:size-40">
            <span class="text-xl font-bold text-primary-foreground sm:text-2xl md:text-3xl">
              {event.acronym ?? event.full_name.slice(0, 2).toUpperCase()}
            </span>
          </div>
        </div>
      </div>
      <div class="mx-auto max-w-6xl bg-background px-4 pt-24 sm:px-6 sm:pt-28 md:px-8">
        <h1 class="text-center text-xl font-bold leading-tight text-foreground min-[300px]:text-2xl sm:text-3xl md:text-4xl">
          {event.full_name}
        </h1>
        <EditionDetails event={event} />
      </div>
    </main>
  );
}

function EditionDetails(props: { event: EventI }) {
  const queryClient = useQueryClient();
  const editions = createMemo(() =>
    queryClient.fetchQuery(allPublicEditionsQueryOptions(props.event.id)),
  );
  return (
    <Loading
      fallback={
        <div class="mt-6 grid gap-4 md:grid-cols-2">
          <div class="h-28 animate-pulse rounded-xl bg-muted" />
          <div class="h-28 animate-pulse rounded-xl bg-muted" />
        </div>
      }
    >
      <EditionBody event={props.event} editions={editions()} />
    </Loading>
  );
}

function EditionBody(props: { event: EventI; editions: EditionI[] }) {
  const now = Date.now();
  const editions = untrack(() => props.editions);
  const sorted = [...editions].sort((a, b) =>
    a.starts_at.localeCompare(b.starts_at),
  );
  const active =
    sorted.find(
      (edition) =>
        new Date(edition.starts_at).getTime() <= now &&
        new Date(edition.ends_at).getTime() >= now,
    ) ??
    sorted.find((edition) => new Date(edition.starts_at).getTime() > now) ??
    sorted.at(-1);
  const otherEditions = sorted.filter((edition) => edition.id !== active?.id);

  return (
    <>
      {active && (
        <div class="mt-3 flex flex-wrap items-center justify-center gap-x-4 gap-y-2 sm:mt-4">
          <div class="flex items-center gap-1.5 text-muted-foreground">
            <span class="size-4">
              <Calendar />
            </span>
            <span class="text-xs font-medium min-[300px]:text-sm">
              {formatDateRange(active.starts_at, active.ends_at)}
            </span>
          </div>
          {active.location_name && (
            <div class="flex items-center gap-1.5 text-muted-foreground">
              <span class="size-4">
                <MapPin />
              </span>
              <span class="text-xs font-medium min-[300px]:text-sm">
                {active.location_name}
              </span>
            </div>
          )}
        </div>
      )}
      {props.event.description && (
        <p class="mx-auto mt-5 border-l-2 border-primary/50 pl-4 text-left text-xs leading-relaxed text-muted-foreground min-[300px]:text-sm sm:mt-6 sm:text-base md:mx-4!">
          {props.event.description}
        </p>
      )}
      {active && (
        <EventCatalog editionId={active.id} eventSlug={props.event.slug} />
      )}
      {otherEditions.length > 0 && (
        <section class="mt-12 w-full">
          <div class="mb-5 sm:mb-6">
            <h2 class="text-2xl font-normal tracking-tight text-primary sm:text-3xl">
              Outras Edições
            </h2>
            <p class="mt-1 text-sm text-muted-foreground sm:text-base">
              Conheça os próximos eventos ou relembre os destaques das edições
              anteriores.
            </p>
          </div>
          <div class="flex flex-wrap gap-4">
            {otherEditions.slice(0, 5).map((edition) => (
              <EditionSummaryCard edition={edition} />
            ))}
          </div>
          <a
            class="mt-5 inline-flex items-center gap-1.5 text-sm font-semibold text-primary transition-all duration-200 hover:gap-2.5"
            href={`/events/${props.event.slug}/editions`}
          >
            Ver todas as Edições <span aria-hidden="true">→</span>
          </a>
        </section>
      )}
      {active?.location_name && (
        <section class="relative z-0 mt-8 w-full sm:mt-10">
          <div class="mb-4">
            <h2 class="text-xl font-semibold text-foreground sm:text-2xl">
              Local do evento
            </h2>
            <p class="mt-1 text-sm text-muted-foreground">
              {active.location_name}
              {active.location_description
                ? ` — ${active.location_description}`
                : ""}
            </p>
          </div>
          <iframe
            title={`Mapa de ${active.location_name}`}
            src={`https://www.google.com/maps?q=${encodeURIComponent(
              [active.location_name, active.location_description]
                .filter(Boolean)
                .join(", "),
            )}&output=embed`}
            loading="lazy"
            class="h-80 w-full overflow-hidden rounded-xl border border-border shadow-sm"
          />
        </section>
      )}
      <ContactSection event={props.event} />
      {active && (
        <EventCart
          editionId={active.id}
          checkoutHref={`/events/${props.event.slug}/checkout`}
        />
      )}
    </>
  );
}

function formatDateRange(startsAt: string, endsAt: string) {
  const format = (value: string) =>
    new Date(value).toLocaleDateString("pt-BR", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    });
  return `${format(startsAt)} — ${format(endsAt)}`;
}
