import { createFileRoute } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core/solid";
import { Loading, createMemo, untrack } from "solid-js";
import type { JSX } from "@solidjs/web";
import ShareIcon from "~icons/lucide/share-2";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import { ContactSection } from "@/features/events/ui/ContactSection";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import { EditionSummaryCard } from "@/features/editions/ui/EditionSummaryCard";
import { handleShare } from "@/shared/lib/share";

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
        <h1 class="text-xl font-bold leading-tight text-foreground text-center">
          {event.full_name}
        </h1>
        {event.description && (
          <p class="mx-auto mt-5 border-l-2 border-primary/50 pl-4 text-left text-xs leading-relaxed text-muted-foreground sm:mt-6 md:mx-4!">
            {event.description}
          </p>
        )}
        <EditionList eventId={event.id} eventSlug={event.slug} />
        <ContactSection event={event} />
      </div>
    </main>
  );
}

function EditionList(props: { eventId: string; eventSlug: string }) {
  const queryClient = useQueryClient();
  const editions = createMemo(() =>
    queryClient.fetchQuery(allPublicEditionsQueryOptions(props.eventId)),
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
          {editions().length > 0 ? (
            editions().map((edition) => (
              <EditionSummaryCard edition={edition} />
            ))
          ) : (
            <div class="w-full rounded-xl border border-dashed border-border p-8 text-center text-sm text-muted-foreground">
              Nenhuma edição publicada no momento.
            </div>
          )}
        </div>
        <a
          class="mt-5 inline-flex items-center gap-1.5 text-sm font-semibold text-primary transition-all duration-200 hover:gap-2.5"
          href={`/events/${props.eventSlug}/editions`}
        >
          Ver todas as Edições <span aria-hidden="true">→</span>
        </a>
      </section>
    </Loading>
  );
}
