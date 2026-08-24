import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { activeEditionQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import { storeStockQueryOptions } from "@/features/products/api";
import { useInventoryStream } from "@/features/products/hooks/use-inventory-stream";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import { useProgramRegistration } from "@/features/programs/hooks/use-program-registration";
import { ProgramDayCard } from "@/features/programs/ui/ProgramDayCard";
import { groupByDay } from "@/features/programs/ui/ProgramSection";
import { Carousel } from "@/widgets/carousel/GenericCarousel";

export const Route = createFileRoute("/events/$slug/programs")({
  loader: async ({ context, params }) => {
    const event = await context.queryClient.ensureQueryData(
      publicEventBySlugQueryOptions(params.slug),
    );
    if (!event) throw redirect({ to: "/events" });
    return event;
  },
  component: ProgramsPage,
});

function ProgramsPage() {
  const event = Route.useLoaderData();
  const { data: edition } = useSuspenseQuery(
    activeEditionQueryOptions(event.id),
  );

  if (!edition) {
    return (
      <main className="mx-auto max-w-3xl px-4 py-16 text-center">
        A programação não está disponível no momento.
      </main>
    );
  }

  return (
    <Programs
      editionId={edition.id}
      eventId={event.id}
      eventSlug={event.slug}
    />
  );
}

function Programs({
  editionId,
  eventId,
  eventSlug,
}: {
  editionId: string;
  eventId: string;
  eventSlug: string;
}) {
  useInventoryStream(editionId);
  const { data: programs } = useSuspenseQuery(programsQueryOptions(editionId));
  const { data: occurrences } = useSuspenseQuery(
    occurrencesQueryOptions(editionId),
  );
  const { data: stock } = useSuspenseQuery(storeStockQueryOptions(editionId));
  const stockByOccurrence = new Map(
    stock
      .filter((item) => item.item_type === "program_occurrence")
      .map((item) => [item.id, item.stock]),
  );
  const days = groupByDay(programs, occurrences);
  const registration = useProgramRegistration(
    editionId,
    eventId,
    eventSlug,
    `/events/${eventSlug}/programs`,
  );

  return (
    <main className="min-h-screen bg-background px-4 pt-4 pb-24 md:pt-6">
      <div className="mx-auto w-full max-w-screen-2xl">
        <header className="mx-auto max-w-3xl">
          <Link
            to="/events/$slug"
            params={{ slug: eventSlug }}
            className="inline-flex items-center gap-2 rounded-lg px-2 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <ArrowLeft className="size-4" /> Evento
          </Link>
          <div className="mt-5 text-center">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-primary">
              Agenda do evento
            </p>
            <h1 className="mt-3 text-3xl font-bold tracking-tight sm:text-4xl">
              Programação completa
            </h1>
            <p className="mx-auto mt-3 max-w-xl text-muted-foreground">
              Confira todos os horários e atividades do evento.
            </p>
          </div>
        </header>

        {days.length > 0 ? (
          <section
            className="mx-auto mt-10 max-w-screen-2xl"
            aria-label="Ocorrências da programação"
          >
            <Carousel
              items={days}
              renderItem={(day) => (
                <ProgramDayCard
                  date={day.date}
                  items={day.items}
                  maxItems={day.items.length}
                  showFullDescription
                  registration={registration}
                  stockByOccurrence={stockByOccurrence}
                />
              )}
              itemMinWidth={320}
              itemMaxWidth={320}
              gap={24}
              loop
              scrollBy={1}
              arrowPosition="top"
            />
          </section>
        ) : (
          <p className="mt-8 text-muted-foreground">
            Nenhuma atividade programada.
          </p>
        )}
      </div>
    </main>
  );
}
