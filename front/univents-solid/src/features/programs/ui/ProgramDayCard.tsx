import type {
  MyTicket,
  Program,
  ProgramOccurrence,
  ProgramParticipationStatus,
} from "@trieoh/univents-api/schemas";
import type { JSX } from "@solidjs/web";
import ActivityIcon from "~icons/lucide/activity";
import LogInIcon from "~icons/lucide/log-in";

const Activity = ActivityIcon as unknown as () => JSX.Element;
const LogIn = LogInIcon as unknown as () => JSX.Element;

const time = (value: string) =>
  new Date(value).toLocaleTimeString("pt-BR", {
    hour: "2-digit",
    minute: "2-digit",
  });

export function ProgramDayCard(props: {
  date: string;
  items: { program: Program; occurrence: ProgramOccurrence }[];
  maxItems?: number;
  showFullDescription?: boolean;
  registration: {
    authenticated: boolean;
    heldTicket: MyTicket | null;
    statuses: () => ReadonlyMap<string, ProgramParticipationStatus>;
    pending: () => string | undefined;
    toggle: (occurrenceId: string, registered: boolean) => void;
  };
}) {
  return (
    <article class="flex w-80 max-w-full flex-col rounded-2xl p-6">
      <h3 class="mb-6 flex min-h-9 w-full items-center justify-center rounded-xl bg-primary/10 px-4 py-2 text-center text-sm font-bold capitalize leading-none tracking-wide text-primary">
        {new Date(props.date).toLocaleDateString("pt-BR", {
          weekday: "long",
          day: "numeric",
          month: "long",
        })}
      </h3>
      <div class="flex flex-col">
        {props.items.slice(0, props.maxItems ?? 3).map(({ program, occurrence }, index) => (
          <div class="flex gap-4">
            <div class="relative flex shrink-0 flex-col items-center">
              <span
                class={`z-10 flex size-9 shrink-0 items-center justify-center rounded-xl ${
                  program.kind === "checkpoint"
                    ? "bg-primary text-primary-foreground"
                    : "bg-muted text-muted-foreground"
                }`}
              >
                <span class="size-4">
                  {program.kind === "checkpoint" ? <LogIn /> : <Activity />}
                </span>
              </span>
              {index < Math.min(props.items.length, props.maxItems ?? 3) - 1 && (
                <span class="min-h-6 w-px flex-1 rounded-full bg-border" />
              )}
            </div>
            <div class="flex-1 pb-6">
              <span class="text-xs font-semibold text-primary">
                {time(occurrence.starts_at)} – {time(occurrence.ends_at)}
              </span>
              <h4 class="mt-1 text-[15px] font-bold leading-snug text-card-foreground">
                {program.name}
              </h4>
              {program.description && (
                <p class={`mt-1 text-xs leading-relaxed text-muted-foreground ${props.showFullDescription ? "" : "line-clamp-2"}`}>
                  {program.description}
                </p>
              )}
              {program.kind === "activity" && (
                <RegistrationButton
                  program={program}
                  occurrenceId={occurrence.id}
                  registration={props.registration}
                />
              )}
            </div>
          </div>
        ))}
      </div>
    </article>
  );
}

function RegistrationButton(props: {
  program: Program;
  occurrenceId: string;
  registration: ProgramDayCardProps["registration"];
}) {
  const status = () => props.registration.statuses().get(props.occurrenceId);
  const registered = () => status() === "registered";
  const pending = () => props.registration.pending() === props.occurrenceId;
  const insufficient = () =>
    props.program.min_access_level != null &&
    (props.registration.heldTicket?.ticket_type.access_level ?? -1) <
      props.program.min_access_level;
  const blocked = () =>
    props.registration.authenticated &&
    (status() === "attended" ||
      status() === "no_show" ||
      props.registration.heldTicket?.status === "pending" ||
      insufficient() ||
      props.program.staff_only);
  const label = () =>
    !props.registration.authenticated
      ? "Entrar para se inscrever"
      : !props.registration.heldTicket
        ? "Ingresso necessário"
        : status() === "attended"
          ? "Participou"
          : status() === "no_show"
            ? "Não compareceu"
            : props.registration.heldTicket.status === "pending"
              ? "Aguardando aprovação"
              : props.program.staff_only
                ? "Somente equipe"
                : insufficient()
                  ? "Nível de ingresso insuficiente"
                  : registered()
                    ? "Cancelar inscrição"
                    : "Inscrever-se";

  return (
    <button
      type="button"
      disabled={pending() || blocked()}
      class={`mt-3 inline-flex h-8 items-center rounded-md border px-3 text-xs font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-60 ${
        registered()
          ? "border-border bg-background text-muted-foreground hover:bg-muted"
          : "border-primary bg-primary text-primary-foreground hover:bg-primary/90"
      }`}
      onClick={() =>
        props.registration.toggle(props.occurrenceId, registered())
      }
    >
      {pending() ? "Salvando…" : label()}
    </button>
  );
}

type ProgramDayCardProps = Parameters<typeof ProgramDayCard>[0];
