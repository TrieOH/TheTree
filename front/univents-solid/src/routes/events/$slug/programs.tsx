import { createFileRoute, Link } from "@tanstack/solid-router";
import { useQueryClient, useQuery } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import ArrowLeftIcon from "~icons/lucide/arrow-left";
import {
  programPageQueryOptions,
} from "@/features/programs/api";
import { ProgramSection } from "@/features/programs/ui/ProgramSection";
import { appQueryClient } from "@/shared/lib/query-client";
import { warmQuery } from "@/shared/lib/route-data";
import { PageNotice } from "@/shared/ui/PageNotice";
import { QueryBoundary } from "@/shared/ui/QueryBoundary";

const ArrowLeft = ArrowLeftIcon as unknown as () => JSX.Element;

export const Route = createFileRoute("/events/$slug/programs")({
  loader: ({ params, context }) =>
    warmQuery(
      appQueryClient,
      programPageQueryOptions(
        params.slug,
        context.session?.isAuthenticated === true,
        appQueryClient,
      ),
    ),
  head: ({ params }) => ({ meta: [{ title: `Programação ${params.slug} - Univents` }] }),
  component: ProgramsPage,
});

function ProgramsPage() {
  const params = Route.useParams();
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuth();
  const dataQuery = useQuery(() => programPageQueryOptions(
    params().slug,
    isAuthenticated(),
    queryClient,
  ));
  return (
    <QueryBoundary
      query={dataQuery}
      what="a programação"
      skeleton={<ProgramsSkeleton />}
      notFound={<PageNotice title="Evento não encontrado" />}
    >
      {(data) => (
        <Show
          when={data()}
          keyed
          fallback={<PageNotice title="Evento não encontrado" />}
        >
          {(loaded) => (
            <Show
              when={loaded.edition}
              keyed
              fallback={<PageNotice title="A programação não está disponível." />}
            >
              {(edition) => (
                <div class="min-h-screen bg-background px-4 pb-24 pt-4 md:pt-6">
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
                      editionId={edition.id}
                      authenticated={loaded.authenticated}
                      heldTicket={loaded.heldTicket}
                      participations={loaded.participations}
                      complete
                    />
                  </div>
                </div>
              )}
            </Show>
          )}
        </Show>
      )}
    </QueryBoundary>
  );
}

/** The page header at its loaded size, with the schedule's rows as placeholders. */
function ProgramsSkeleton() {
  return (
    <div class="min-h-screen px-4 pb-24 pt-4 md:pt-6">
      <div class="mx-auto w-full max-w-screen-2xl">
        <div class="mx-auto max-w-3xl space-y-3 pt-12 text-center">
          <div class="mx-auto h-4 w-32 animate-pulse rounded bg-muted" />
          <div class="mx-auto h-10 w-3/4 animate-pulse rounded bg-muted" />
          <div class="mx-auto h-5 w-2/3 animate-pulse rounded bg-muted" />
        </div>
        <div class="mt-10 space-y-3">
          <For each={[1, 2, 3, 4]}>
            {() => <div class="h-20 animate-pulse rounded-xl bg-muted" />}
          </For>
        </div>
      </div>
    </div>
  );
}
