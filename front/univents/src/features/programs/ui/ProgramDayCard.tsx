import type { ProgramParticipationStatus } from "@trieoh/univents-api/schemas";
import { Activity, Award, ChevronRight, LogIn } from "lucide-react";
import { useState } from "react";
import { cn } from "@/shared/lib/utils";
import { AlertModal } from "@/widgets/ui/alert-modal";
import type { OccurrenceI, ProgramI } from "../model";

function formatTimeRange(start: string, end: string): string {
  const fmt = (iso: string) =>
    new Date(iso).toLocaleTimeString("pt-BR", {
      hour: "2-digit",
      minute: "2-digit",
    });
  return `${fmt(start)} – ${fmt(end)}`;
}

function formatDayLabel(dateIso: string): string {
  const date = new Date(dateIso);
  const today = new Date();
  const tomorrow = new Date(today);
  tomorrow.setDate(tomorrow.getDate() + 1);

  if (date.toDateString() === today.toDateString()) return "Hoje";
  if (date.toDateString() === tomorrow.toDateString()) return "Amanhã";

  return date.toLocaleDateString("pt-BR", {
    weekday: "long",
    day: "numeric",
    month: "long",
  });
}

function getIconForProgram(program: ProgramI) {
  if (program.kind === "checkpoint")
    return { Icon: LogIn, bg: "bg-primary", text: "text-primary-foreground" };
  return { Icon: Activity, bg: "bg-muted", text: "text-muted-foreground" };
}

interface ProgramDayCardProps {
  date: string;
  items: { program: ProgramI; occurrence: OccurrenceI }[];
  maxItems?: number;
  certificateProgramIds?: ReadonlySet<string>;
  showFullDescription?: boolean;
  registration?: {
    isAuthenticated: boolean;
    hasTicket: boolean | null | undefined;
    ticketStatus?: "pending" | "confirmed";
    accessLevel?: number;
    isStaff: boolean;
    participationByOccurrence: ReadonlyMap<string, ProgramParticipationStatus>;
    pendingOccurrenceId?: string;
    onToggle: (occurrenceId: string, registered: boolean) => void;
  };
  stockByOccurrence?: ReadonlyMap<string, number | null>;
}

export function ProgramDayCard({
  date,
  items,
  maxItems = 3,
  certificateProgramIds,
  showFullDescription = false,
  registration,
  stockByOccurrence,
}: ProgramDayCardProps) {
  const [occurrenceToCancel, setOccurrenceToCancel] = useState<string>();
  const visible = items.slice(0, maxItems);
  const remaining = items.length - maxItems;

  return (
    <div className="flex w-80 max-w-full flex-col rounded-2xl p-6">
      {/* Day header */}
      <h3 className="w-full rounded-xl bg-primary/10 px-4 py-2 text-center text-sm font-bold tracking-wide text-primary capitalize mb-6">
        {formatDayLabel(date)}
      </h3>

      {/* Timeline */}
      <div className="flex flex-col">
        {visible.map(({ program, occurrence }, index) => {
          const isLast = index === visible.length - 1 && remaining <= 0;
          const { Icon, bg, text } = getIconForProgram(program);

          return (
            <div key={occurrence.id} className="flex gap-4">
              {/* Timeline column */}
              <div className="flex flex-col items-center shrink-0 relative">
                <div
                  className={cn(
                    "w-9 h-9 rounded-xl z-10 flex items-center justify-center shrink-0",
                    bg,
                  )}
                >
                  <Icon className={cn("w-4 h-4", text)} />
                </div>
                <div
                  className={cn(
                    "w-px rounded-full bg-border",
                    isLast
                      ? "absolute h-[calc(100%+0.5rem)]"
                      : "min-h-6 flex-1",
                  )}
                />
              </div>

              {/* Content */}
              <div className={cn("flex-1", isLast ? "pb-0" : "pb-6")}>
                {/* Time */}
                <span className="text-xs font-semibold text-primary">
                  {formatTimeRange(occurrence.starts_at, occurrence.ends_at)}
                </span>
                {stockByOccurrence?.has(occurrence.id) ? (
                  <span className="ml-2 text-[11px] text-muted-foreground">
                    {stockByOccurrence.get(occurrence.id) === null
                      ? "Vagas ilimitadas"
                      : stockByOccurrence.get(occurrence.id) === 0
                        ? "Esgotado"
                        : `${stockByOccurrence.get(occurrence.id)} vagas restantes`}
                  </span>
                ) : null}

                {/* Title */}
                <h4 className="mt-1 text-[15px] font-bold text-card-foreground leading-snug">
                  {program.name}
                </h4>
                {certificateProgramIds?.has(program.id) ? (
                  <span className="mt-1 inline-flex items-center gap-1 text-[11px] font-medium text-emerald-600">
                    <Award className="size-3.5" />
                    Certificado disponível
                  </span>
                ) : null}

                {/* Description */}
                {program.description && (
                  <p
                    className={cn(
                      "mt-1 text-xs leading-relaxed text-muted-foreground",
                      !showFullDescription && "line-clamp-2",
                    )}
                  >
                    {program.description}
                  </p>
                )}
                {program.kind === "activity" && registration ? (
                  <RegistrationButton
                    occurrenceId={occurrence.id}
                    program={program}
                    registration={registration}
                    onCancel={() => setOccurrenceToCancel(occurrence.id)}
                  />
                ) : null}
              </div>
            </div>
          );
        })}

        {/* "+ N" */}
        {remaining > 0 && (
          <div className="flex items-center gap-2 -mt-2 ml-13 text-xs font-medium text-muted-foreground">
            <span>
              +{remaining} ocorrência{remaining > 1 ? "s" : ""}
            </span>
            <ChevronRight className="w-3 h-3" />
          </div>
        )}
      </div>
      {registration ? (
        <AlertModal
          open={Boolean(occurrenceToCancel)}
          onOpenChange={(open) => !open && setOccurrenceToCancel(undefined)}
          title="Cancelar inscrição?"
          description="Sua vaga será liberada para outra pessoa. Você poderá se inscrever novamente enquanto houver disponibilidade."
          confirmLabel="Cancelar inscrição"
          variant="destructive"
          loading={registration.pendingOccurrenceId === occurrenceToCancel}
          onConfirm={async () => {
            if (!occurrenceToCancel) return;
            registration.onToggle(occurrenceToCancel, true);
            setOccurrenceToCancel(undefined);
          }}
        />
      ) : null}
    </div>
  );
}

function RegistrationButton({
  occurrenceId,
  program,
  registration,
  onCancel,
}: {
  occurrenceId: string;
  program: ProgramI;
  registration: NonNullable<ProgramDayCardProps["registration"]>;
  onCancel: () => void;
}) {
  const status = registration.participationByOccurrence.get(occurrenceId);
  const registered = status === "registered";
  const pending = registration.pendingOccurrenceId === occurrenceId;
  const insufficientLevel =
    program.min_access_level != null &&
    (registration.accessLevel ?? -1) < program.min_access_level;
  const blocked =
    status === "attended" ||
    status === "no_show" ||
    registration.ticketStatus === "pending" ||
    insufficientLevel ||
    (program.staff_only && !registration.isStaff);
  const label = !registration.isAuthenticated
    ? "Entrar para se inscrever"
    : registration.hasTicket === undefined
      ? "Verificando ingresso…"
      : registration.hasTicket === false
        ? "Ingresso necessário"
        : status === "attended"
          ? "Participou"
          : status === "no_show"
            ? "Não compareceu"
            : registration.ticketStatus === "pending"
              ? "Aguardando aprovação"
              : program.staff_only && !registration.isStaff
                ? "Somente equipe"
                : insufficientLevel
                  ? "Nível de ingresso insuficiente"
                  : registered
                    ? "Cancelar inscrição"
                    : "Inscrever-se";

  return (
    <button
      type="button"
      className={cn(
        "mt-3 inline-flex h-8 items-center rounded-md border px-3 text-xs font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-60",
        registered
          ? "border-border bg-background text-muted-foreground hover:bg-muted"
          : "border-primary bg-primary text-primary-foreground hover:bg-primary/90",
      )}
      disabled={pending || blocked || registration.hasTicket === undefined}
      onPointerDown={(event) => event.stopPropagation()}
      onClick={() =>
        registered ? onCancel() : registration.onToggle(occurrenceId, false)
      }
    >
      {pending ? "Salvando…" : label}
    </button>
  );
}
