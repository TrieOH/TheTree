import type {
  MyParticipation,
  MyTicket,
  Program,
  ProgramOccurrence,
  ProgramParticipationStatus,
} from "@trieoh/univents-api/schemas";
import { For, createSignal, untrack } from "solid-js";
import {
  deregisterOccurrenceFn,
  registerOccurrenceFn,
} from "@/features/programs/api";
import { toast } from "@/shared/ui/toast";
import { ProgramDayCard } from "./ProgramDayCard";
import { Carousel } from "@/widgets/carousel/Carousel";
import { Link } from "@tanstack/solid-router";

export function ProgramSection(props: {
  programs: Program[];
  occurrences: ProgramOccurrence[];
  eventSlug: string;
  editionId: string;
  authenticated: boolean;
  heldTicket: MyTicket | null;
  participations: MyParticipation[];
  complete?: boolean;
}) {
  const {
    programItems,
    occurrenceItems,
    authenticated,
    heldTicket,
    participations,
    eventSlug,
    complete,
  } = untrack(() => ({
    programItems: props.programs,
    occurrenceItems: props.occurrences,
    authenticated: props.authenticated,
    heldTicket: props.heldTicket,
    participations: props.participations,
    eventSlug: props.eventSlug,
    complete: props.complete ?? false,
  }));
  const programs = new Map(
    programItems.map((program) => [program.id, program]),
  );
  const grouped = new Map<
    string,
    {
      date: string;
      items: { program: Program; occurrence: ProgramOccurrence }[];
    }
  >();
  for (const occurrence of occurrenceItems) {
    const program = programs.get(occurrence.program_id);
    if (!program) continue;
    const key = occurrence.starts_at.slice(0, 10);
    const group = grouped.get(key) ?? { date: occurrence.starts_at, items: [] };
    group.items.push({ program, occurrence });
    grouped.set(key, group);
  }
  const days = [...grouped.values()]
    .sort((a, b) => a.date.localeCompare(b.date))
    .map((day) => ({
      ...day,
      items: day.items.sort((a, b) =>
        a.occurrence.starts_at.localeCompare(b.occurrence.starts_at),
      ),
    }))
    .slice(0, complete ? undefined : 3);
  const [statuses, setStatuses] = createSignal<
    ReadonlyMap<string, ProgramParticipationStatus>
  >(
    new Map(
      participations.map((participation) => [
        participation.occurrence_id,
        participation.status,
      ]),
    ),
  );
  const [pending, setPending] = createSignal<string>();
  const toggle = async (occurrenceId: string, registered: boolean) => {
    if (!authenticated) {
      window.location.href = `/auth?redirect=${encodeURIComponent(
        `/events/${eventSlug}`,
      )}`;
      return;
    }
    if (!heldTicket) {
      window.location.href = `/events/${eventSlug}/store?tab=tickets`;
      return;
    }
    setPending(occurrenceId);
    try {
      const participation = registered
        ? await deregisterOccurrenceFn(occurrenceId)
        : await registerOccurrenceFn(occurrenceId);
      setStatuses((current) => {
        const next = new Map(current);
        next.set(occurrenceId, participation.status);
        return next;
      });
      toast.success(registered ? "Inscrição cancelada" : "Inscrição realizada");
    } catch {
      toast.error(
        registered
          ? "Não foi possível cancelar a inscrição"
          : "Não foi possível realizar a inscrição",
      );
    } finally {
      setPending(undefined);
    }
  };
  const registration = {
    authenticated,
    heldTicket,
    statuses,
    pending,
    toggle: (occurrenceId: string, registered: boolean) =>
      void toggle(occurrenceId, registered),
  };

  const renderDay = (day: (typeof days)[number]) => (
    <ProgramDayCard
      date={day.date}
      items={day.items}
      maxItems={complete ? undefined : 3}
      showFullDescription={complete}
      registration={registration}
    />
  );

  return (
    <section class="w-full py-10">
      {!complete && <div class="mb-8 text-center">
        <h2 class="text-3xl font-semibold tracking-tight text-foreground">
          Programação
        </h2>
        <p class="mx-auto mt-2 max-w-xl text-sm leading-relaxed text-muted-foreground">
          Confira as atividades e checkpoints do evento.
        </p>
      </div>}
      {days.length === 0 ? (
        <p class="mt-8 text-muted-foreground">Nenhuma atividade programada.</p>
      ) : complete ? (
        <Carousel
          items={days}
          itemMinWidth={320}
          itemMaxWidth={320}
          gap={24}
          scrollBy={"page"}
          arrowPosition="top"
          renderItem={renderDay}
          className="mx-auto max-w-7xl"
        />
      ) : <div class="mx-auto flex flex-nowrap justify-center gap-5">
        <For each={days}>
          {(day, index) => (
            <div
              class={
                index() === 1
                  ? "hidden items-start gap-5 md:flex"
                  : index() === 2
                    ? "hidden items-start gap-5 lg:flex"
                    : "flex items-start gap-5"
              }
            >
              {renderDay(day)}
              {index() < days.length - 1 && (
                <div
                  class={
                    index() === 0
                      ? "hidden w-px self-stretch rounded-full bg-border md:block"
                      : "hidden w-px self-stretch rounded-full bg-border lg:block"
                  }
                />
              )}
            </div>
          )}
        </For>
      </div>}
      {!complete && <div class="mt-6 text-center">
        <Link
          class="text-sm font-semibold text-primary"
          to="/events/$slug/programs"
          params={{ slug: props.eventSlug }}
        >
          Programação Completa →
        </Link>
      </div>}
    </section>
  );
}
