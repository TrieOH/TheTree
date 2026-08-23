import { useQueries } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import type { Edition, Event } from "@trieoh/univents-api/schemas";
import { CalendarDays, Clock } from "lucide-react";
import { myParticipationsQueryOptions } from "../api";

export function UserActivitiesContent({
  editions,
  events,
}: {
  editions: Edition[];
  events: ReadonlyMap<string, Event>;
}) {
  const queries = useQueries({
    queries: editions.map((edition) =>
      myParticipationsQueryOptions(edition.id),
    ),
  });
  const activities = queries.flatMap((query, index) =>
    (query.data ?? []).map((participation) => ({
      participation,
      edition: editions[index],
    })),
  );

  if (queries.some((query) => query.isPending)) {
    return (
      <p className="text-sm text-muted-foreground">Carregando atividades…</p>
    );
  }

  if (activities.length === 0) {
    return (
      <div className="rounded-lg border border-dashed p-8 text-center">
        <CalendarDays className="mx-auto size-8 text-muted-foreground" />
        <p className="mt-3 font-medium">Nenhuma atividade inscrita</p>
        <p className="mt-1 text-sm text-muted-foreground">
          Suas inscrições na programação dos eventos aparecerão aqui.
        </p>
      </div>
    );
  }

  return (
    <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {activities.map(({ participation, edition }) => {
        const event = edition ? events.get(edition.event_id) : undefined;
        return (
          <article
            key={participation.id}
            className="rounded-lg border border-border bg-card p-4 shadow-sm"
          >
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0">
                <p className="text-xs font-medium text-primary">
                  {event?.full_name ?? edition?.name}
                </p>
                <h2 className="mt-1 truncate font-semibold">
                  {participation.program.name}
                </h2>
              </div>
              <span className="shrink-0 rounded-full bg-muted px-2 py-1 text-[11px] font-medium text-muted-foreground">
                {statusLabel(participation.status)}
              </span>
            </div>
            <p className="mt-3 flex items-center gap-1.5 text-sm text-muted-foreground">
              <Clock className="size-4" />
              {formatOccurrence(
                participation.occurrence.starts_at,
                participation.occurrence.ends_at,
              )}
            </p>
            {event ? (
              <Link
                to="/events/$slug/programs"
                params={{ slug: event.slug }}
                className="mt-4 inline-flex text-sm font-semibold text-primary hover:underline"
              >
                Ver programação
              </Link>
            ) : null}
          </article>
        );
      })}
    </div>
  );
}

function statusLabel(status: string) {
  if (status === "attended") return "Participou";
  if (status === "no_show") return "Não compareceu";
  return "Inscrito";
}

function formatOccurrence(startsAt: string, endsAt: string) {
  const start = new Date(startsAt);
  const end = new Date(endsAt);
  return `${start.toLocaleDateString("pt-BR")} · ${start.toLocaleTimeString("pt-BR", { hour: "2-digit", minute: "2-digit" })}–${end.toLocaleTimeString("pt-BR", { hour: "2-digit", minute: "2-digit" })}`;
}
