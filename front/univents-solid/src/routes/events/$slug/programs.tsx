import { createFileRoute, Link } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import type { Edition } from "@trieoh/univents-api/schemas";
import type { JSX } from "@solidjs/web";
import { Loading, createMemo } from "solid-js";
import ArrowLeftIcon from "~icons/lucide/arrow-left";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import {
  myParticipationsQueryOptions,
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import { ProgramSection } from "@/features/programs/ui/ProgramSection";
import { myTicketQueryOptions } from "@/features/tickets/api";

const ArrowLeft = ArrowLeftIcon as unknown as () => JSX.Element;

export const Route = createFileRoute("/events/$slug/programs")({
  head: ({ params }) => ({ meta: [{ title: `Programação ${params.slug} - Univents` }] }),
  component: ProgramsPage,
});

function ProgramsPage() {
  const params = Route.useParams();
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuth();
  const data = createMemo(async () => {
    const event = await queryClient.fetchQuery(publicEventBySlugQueryOptions(params().slug));
    if (!event) return null;
    const edition = currentEdition(await queryClient.fetchQuery(allPublicEditionsQueryOptions(event.id)));
    if (!edition) return { event, edition: null };
    const authenticated = isAuthenticated();
    const [programs, occurrences, heldTicket, participations] = await Promise.all([
      queryClient.fetchQuery(programsQueryOptions(edition.id)),
      queryClient.fetchQuery(occurrencesQueryOptions(edition.id)),
      authenticated ? queryClient.fetchQuery(myTicketQueryOptions(edition.id)) : null,
      authenticated ? queryClient.fetchQuery(myParticipationsQueryOptions(edition.id)) : [],
    ]);
    return { event, edition, programs, occurrences, heldTicket, participations, authenticated };
  });

  return (
    <Loading fallback={<main class="min-h-screen animate-pulse bg-muted" />}>
      {(() => {
        const loaded = data();
        if (!loaded) return <main class="p-12 text-center">Evento não encontrado.</main>;
        if (!loaded.edition) return <main class="p-12 text-center">A programação não está disponível.</main>;
        return (
          <main class="min-h-screen bg-background px-4 pb-24 pt-4 md:pt-6">
            <div class="mx-auto w-full max-w-screen-2xl">
              <header class="mx-auto max-w-3xl">
                <Link
                  to="/events/$slug"
                  params={{ slug: loaded.event.slug }}
                  class="inline-flex items-center gap-2 rounded-lg px-2 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                  <ArrowLeft /> Evento
                </Link>
                <div class="mt-5 text-center">
                  <p class="text-xs font-semibold uppercase tracking-[0.2em] text-primary">Agenda do evento</p>
                  <h1 class="mt-3 text-3xl font-bold tracking-tight sm:text-4xl">Programação completa</h1>
                  <p class="mx-auto mt-3 max-w-xl text-muted-foreground">Confira todos os horários e atividades do evento.</p>
                </div>
              </header>
              <ProgramSection
                programs={loaded.programs}
                occurrences={loaded.occurrences}
                eventSlug={loaded.event.slug}
                editionId={loaded.edition.id}
                authenticated={loaded.authenticated}
                heldTicket={loaded.heldTicket}
                participations={loaded.participations}
                complete
              />
            </div>
          </main>
        );
      })()}
    </Loading>
  );
}

function currentEdition(editions: Edition[]) {
  const now = Date.now();
  const sorted = [...editions].sort((a, b) => a.starts_at.localeCompare(b.starts_at));
  return sorted.find((edition) => new Date(edition.starts_at).getTime() <= now && new Date(edition.ends_at).getTime() >= now)
    ?? sorted.find((edition) => new Date(edition.starts_at).getTime() > now)
    ?? sorted.at(-1);
}
