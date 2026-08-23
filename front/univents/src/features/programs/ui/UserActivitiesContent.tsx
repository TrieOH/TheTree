import { useQuery } from "@tanstack/react-query";
import type { Edition } from "@trieoh/univents-api/schemas";
import { CalendarDays, Clock, X } from "lucide-react";
import { useState } from "react";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";
import { myParticipationsQueryOptions } from "../api";
import { useOccurrenceRegistrationMutation } from "../api/mutations";

export function UserActivitiesContent({ edition }: { edition: Edition }) {
  const query = useQuery(myParticipationsQueryOptions(edition.id));
  const activities = query.data ?? [];
  const mutation = useOccurrenceRegistrationMutation(edition.id);
  const [occurrenceToCancel, setOccurrenceToCancel] = useState<string>();

  if (query.isPending) {
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
    <div className="grid gap-2 sm:grid-cols-2 xl:grid-cols-3">
      {activities.map((participation) => {
        return (
          <article
            key={participation.id}
            className="flex min-w-0 items-center gap-3 rounded-lg border border-border/70 bg-card p-3 shadow-xs"
          >
            <div className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <CalendarDays className="size-5" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="flex min-w-0 items-center gap-2">
                <h2 className="truncate text-sm font-semibold">
                  {participation.program.name}
                </h2>
                <span className="shrink-0 rounded-full bg-muted px-2 py-0.5 text-[10px] font-medium text-muted-foreground">
                  {statusLabel(participation.status)}
                </span>
              </div>
              <p className="mt-1 flex items-center gap-1.5 truncate text-xs text-muted-foreground">
                <Clock className="size-3.5 shrink-0" />
                {formatOccurrence(
                  participation.occurrence.starts_at,
                  participation.occurrence.ends_at,
                )}
              </p>
            </div>
            {participation.status === "registered" ? (
              <Button
                type="button"
                size="icon-sm"
                variant="ghost"
                className="shrink-0 text-muted-foreground hover:text-destructive"
                aria-label={`Cancelar inscrição em ${participation.program.name}`}
                onClick={() =>
                  setOccurrenceToCancel(participation.occurrence_id)
                }
              >
                <X className="size-4" />
              </Button>
            ) : null}
          </article>
        );
      })}
      <AlertModal
        open={Boolean(occurrenceToCancel)}
        onOpenChange={(open) => !open && setOccurrenceToCancel(undefined)}
        title="Cancelar inscrição?"
        description="Sua vaga será liberada para outra pessoa. Você poderá se inscrever novamente enquanto houver disponibilidade."
        confirmLabel="Cancelar inscrição"
        variant="destructive"
        loading={mutation.isPending}
        onConfirm={async () => {
          if (!occurrenceToCancel) return;
          await mutation.mutateAsync({
            occurrenceId: occurrenceToCancel,
            registered: true,
          });
          setOccurrenceToCancel(undefined);
        }}
      />
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
