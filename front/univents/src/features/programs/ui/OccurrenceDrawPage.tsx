import { useQuery } from "@tanstack/react-query";
import type { ProgramParticipant } from "@trieoh/univents-api/schemas";
import {
  ArrowLeft,
  CheckCircle2,
  Gift,
  ListChecks,
  LoaderCircle,
  Users,
  X,
} from "lucide-react";
import { AnimatePresence, motion } from "motion/react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useActorDisplayNames } from "@/features/profile/api/actor-display-names";
import { Button } from "@/shared/ui/shadcn/button";
import { occurrenceParticipantsQueryOptions } from "../api";
import { drawSequence, drawTimeline, randomItem } from "../lib/draw-sequence";

type Audience = "attended" | "registered";

export function OccurrenceDrawPage({
  occurrenceId,
  programName,
  onBack,
}: {
  occurrenceId: string;
  programName: string;
  onBack: () => void;
}) {
  const { data: participants = [], isPending } = useQuery(
    occurrenceParticipantsQueryOptions(occurrenceId),
  );
  const actorIds = useMemo(
    () =>
      participants.flatMap((participant) =>
        participant.attendee_user_id ? [participant.attendee_user_id] : [],
      ),
    [participants],
  );
  const { data: displayNames = {}, isFetching: isLoadingNames } =
    useActorDisplayNames(actorIds);
  const [audience, setAudience] = useState<Audience>("attended");
  const [winner, setWinner] = useState<ProgramParticipant>();
  const [displayed, setDisplayed] = useState<ProgramParticipant>();
  const [drawing, setDrawing] = useState(false);
  const [immersive, setImmersive] = useState(false);
  const [countdown, setCountdown] = useState<number>();
  const [rollDurationMs, setRollDurationMs] = useState(0);
  const timers = useRef<number[]>([]);

  const eligible = useMemo(
    () =>
      participants.filter((participant) =>
        audience === "attended"
          ? participant.status === "attended"
          : participant.status !== "cancelled",
      ),
    [audience, participants],
  );
  const clearTimers = () => {
    timers.current.forEach(window.clearTimeout);
    timers.current = [];
  };
  const participantName = (participant: ProgramParticipant) => {
    const actorId = participant.attendee_user_id;
    const profileName = actorId ? displayNames[actorId] : undefined;
    return profileName && profileName !== actorId
      ? profileName
      : participant.attendee_email;
  };

  useEffect(() => clearTimers, []);

  useEffect(() => {
    if (!immersive) return;
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape" && !drawing) setImmersive(false);
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [drawing, immersive]);

  const changeAudience = (next: Audience) => {
    clearTimers();
    setAudience(next);
    setWinner(undefined);
    setDisplayed(undefined);
    setDrawing(false);
    setCountdown(undefined);
  };

  const draw = () => {
    const pool =
      eligible.length > 1 && winner
        ? eligible.filter((participant) => participant.id !== winner.id)
        : eligible;
    const selected = randomItem(pool);
    if (!selected) return;

    clearTimers();
    setWinner(undefined);
    setDisplayed(undefined);
    setDrawing(true);
    setImmersive(true);
    setCountdown(3);
    const sequence = drawSequence(eligible);
    const { delays, durationMs } = drawTimeline(sequence.length);
    setRollDurationMs(durationMs);

    [2, 1].forEach((value, index) => {
      timers.current.push(
        window.setTimeout(() => setCountdown(value), 550 * (index + 1)),
      );
    });
    timers.current.push(window.setTimeout(() => setCountdown(undefined), 1650));

    sequence.forEach((participant, index) => {
      timers.current.push(
        window.setTimeout(() => {
          setDisplayed(participant);
        }, 1650 + delays[index]),
      );
    });
    timers.current.push(
      window.setTimeout(
        () => {
          setDisplayed(selected);
          setWinner(selected);
          setDrawing(false);
        },
        1650 + durationMs + 180,
      ),
    );
  };

  return (
    <main className="h-dvh overflow-y-auto bg-white text-[#17201b] dark:bg-[#141618] dark:text-[#f3f5f2]">
      <div className="grid min-h-full min-w-0 lg:grid-cols-[21rem_minmax(0,1fr)]">
        <aside className="space-y-6 border-b border-slate-200 p-5 dark:border-white/10 lg:border-r lg:border-b-0 lg:p-6">
          <div className="flex min-w-0 items-center gap-2">
            <Button
              type="button"
              variant="ghost"
              size="icon-sm"
              className="-ml-2 shrink-0 text-muted-foreground"
              onClick={onBack}
            >
              <ArrowLeft className="size-4" />
              <span className="sr-only">Voltar para ocorrências</span>
            </Button>
            <div className="min-w-0">
              <p className="text-xs text-muted-foreground">Sorteio</p>
              <p className="truncate text-sm font-medium">{programName}</p>
            </div>
          </div>

          <fieldset>
            <legend className="mb-4">
              <span className="block text-xs font-semibold text-blue-700 uppercase dark:text-[#d7ff43]">
                Público do sorteio
              </span>
              <span className="mt-1 block text-lg font-semibold">
                Quem pode concorrer?
              </span>
            </legend>
            <div className="space-y-2">
              <AudienceButton
                selected={audience === "attended"}
                icon={<CheckCircle2 className="size-5" />}
                title="Só presentes"
                description="Presença confirmada"
                onClick={() => changeAudience("attended")}
              />
              <AudienceButton
                selected={audience === "registered"}
                icon={<ListChecks className="size-5" />}
                title="Todos inscritos"
                description="Com ou sem presença"
                onClick={() => changeAudience("registered")}
              />
            </div>
          </fieldset>

          <div className="border-y border-border/60 py-4">
            <div className="flex items-center gap-3">
              <div className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-blue-100 text-blue-700 dark:bg-[#d7ff43]/10 dark:text-[#d7ff43]">
                <Users className="size-5" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-xs text-muted-foreground">
                  Total no sorteio
                </p>
                <p className="mt-0.5 text-2xl font-semibold tabular-nums">
                  {isPending ? "–" : eligible.length}
                </p>
              </div>
              <span className="text-right text-xs text-muted-foreground">
                {audience === "attended" ? "presentes" : "inscritos"}
              </span>
            </div>
          </div>
        </aside>

        <motion.section
          layout
          transition={{
            layout: { duration: 0.58, ease: [0.22, 1, 0.36, 1] },
          }}
          className={`flex min-w-0 flex-col items-center justify-center overflow-hidden bg-slate-50 px-5 py-10 text-center text-[#111827] transition-[border-radius] duration-500 dark:bg-[#101214] dark:text-[#f4f5f6] sm:px-10 ${
            immersive
              ? "fixed inset-0 z-100 min-h-dvh border-0"
              : "relative min-h-120 sm:min-h-144"
          }`}
        >
          <div className="absolute inset-x-0 top-0 h-1 bg-blue-600 dark:bg-[#d7ff43]" />
          {immersive ? (
            <div className="absolute inset-x-5 top-5 flex items-center justify-between gap-4 sm:inset-x-8 sm:top-7">
              <div className="min-w-0 text-left">
                <p className="truncate text-sm font-medium">{programName}</p>
                <p className="mt-0.5 text-xs opacity-55">
                  {eligible.length}{" "}
                  {eligible.length === 1 ? "participante" : "participantes"}
                </p>
              </div>
              {!drawing ? (
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="shrink-0 text-slate-700 hover:bg-slate-100 hover:text-slate-950 dark:text-[#f4f5f6] dark:hover:bg-white/10 dark:hover:text-white"
                  onClick={() => setImmersive(false)}
                >
                  <X className="size-5" />
                  <span className="sr-only">Sair da tela cheia</span>
                </Button>
              ) : null}
            </div>
          ) : null}
          {isPending || isLoadingNames ? (
            <LoaderCircle className="size-10 animate-spin opacity-70" />
          ) : eligible.length === 0 ? (
            <>
              <Users className="mb-4 size-12 opacity-60" />
              <p className="text-xl font-semibold">
                Ninguém elegível para o sorteio
              </p>
              <p className="mt-2 max-w-md text-sm opacity-65">
                {audience === "attended"
                  ? "Confirme presenças ou escolha todos os inscritos."
                  : "Ainda não há inscrições nesta atividade."}
              </p>
            </>
          ) : countdown ? (
            <AnimatePresence mode="wait">
              <motion.div
                key={countdown}
                initial={{ opacity: 0, scale: 0.82, y: 24 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 1.12, y: -20 }}
                transition={{ duration: 0.28, ease: [0.22, 1, 0.36, 1] }}
                className="text-center"
              >
                <p className="mb-5 text-xs font-medium text-slate-500 uppercase dark:text-white/50">
                  Preparando sorteio
                </p>
                <p className="text-8xl font-semibold tabular-nums text-blue-600 dark:text-[#d7ff43] sm:text-9xl">
                  {countdown}
                </p>
              </motion.div>
            </AnimatePresence>
          ) : winner ? (
            <>
              <Confetti />
              <motion.div
                initial={{ opacity: 0, y: 32 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.45, ease: [0.22, 1, 0.36, 1] }}
                className="relative z-10 max-w-full"
              >
                <motion.div
                  className="mx-auto mb-8 h-1 bg-blue-600 dark:bg-[#d7ff43]"
                  initial={{ width: 0 }}
                  animate={{ width: 72 }}
                  transition={{ delay: 0.2, duration: 0.35 }}
                />
                <p className="text-sm font-semibold text-blue-700 uppercase dark:text-[#d7ff43] sm:text-base">
                  Vencedor
                </p>
                <p className="mt-4 max-w-[90vw] wrap-break-word text-4xl font-bold sm:text-7xl">
                  {participantName(winner)}
                </p>
                <p className="mt-3 break-all text-sm opacity-65 sm:text-base">
                  {winner.attendee_email}
                </p>
              </motion.div>
            </>
          ) : displayed ? (
            <div className="max-w-full">
              <div className="relative flex h-32 w-[min(92vw,70rem)] items-center justify-center overflow-hidden sm:h-44">
                <AnimatePresence initial={false}>
                  <motion.p
                    key={displayed.id}
                    initial={{ opacity: 0, y: 45, filter: "blur(12px)" }}
                    animate={{ opacity: 1, y: 0, filter: "blur(0px)" }}
                    exit={{ opacity: 0, y: -45, filter: "blur(12px)" }}
                    transition={{ duration: 0.075 }}
                    className="absolute inset-0 flex items-center justify-center overflow-hidden wrap-break-word px-3 text-3xl leading-tight font-semibold sm:text-6xl"
                  >
                    <span className="line-clamp-2">
                      {participantName(displayed)}
                    </span>
                  </motion.p>
                </AnimatePresence>
              </div>
              {drawing ? (
                <div className="mx-auto mt-10 h-px w-48 overflow-hidden bg-slate-200 dark:bg-white/15">
                  <motion.div
                    className="h-full bg-blue-600 dark:bg-[#d7ff43]"
                    initial={{ width: "0%" }}
                    animate={{ width: "100%" }}
                    transition={{
                      duration: rollDurationMs / 1000,
                      ease: "easeInOut",
                    }}
                  />
                </div>
              ) : null}
            </div>
          ) : (
            <>
              <Gift className="mb-5 size-12 text-blue-600 dark:text-[#d7ff43]" />
              <p className="text-3xl font-bold sm:text-5xl">
                Pronto para sortear
              </p>
              <p className="mt-3 text-sm opacity-65 sm:text-base">
                A seleção será feita entre {eligible.length}{" "}
                {eligible.length === 1 ? "pessoa" : "pessoas"}.
              </p>
            </>
          )}

          {!drawing ? (
            <div className="relative z-10 mt-10 flex flex-wrap justify-center gap-2">
              <Button
                type="button"
                size="lg"
                className="min-w-52 gap-2 bg-blue-600 text-white hover:bg-blue-700 dark:bg-[#d7ff43] dark:text-[#111407] dark:hover:bg-[#e3ff77]"
                disabled={isPending || isLoadingNames || eligible.length === 0}
                onClick={draw}
              >
                <Gift className="size-4" />
                {winner ? "Sortear novamente" : "Começar sorteio"}
              </Button>
              {winner && immersive ? (
                <Button
                  type="button"
                  size="lg"
                  variant="outline"
                  className="border-slate-300 bg-transparent text-slate-700 hover:bg-slate-100 hover:text-slate-950 dark:border-white/25 dark:text-[#f4f5f6] dark:hover:bg-white/10 dark:hover:text-white"
                  onClick={() => setImmersive(false)}
                >
                  Encerrar
                </Button>
              ) : null}
            </div>
          ) : null}
        </motion.section>
      </div>
    </main>
  );
}

const confettiPieces = Array.from({ length: 28 }, (_, index) => ({
  left: `${(index * 37) % 100}%`,
  delay: (index % 7) * 0.07,
  duration: 1.7 + (index % 5) * 0.18,
  rotate: 240 + (index % 4) * 100,
}));

function Confetti() {
  const colors = [
    "bg-blue-600 dark:bg-[#d7ff43]",
    "bg-slate-900 dark:bg-white",
    "bg-sky-300 dark:bg-[#7dd3fc]",
  ];

  return (
    <div
      className="pointer-events-none absolute inset-0 overflow-hidden"
      aria-hidden="true"
    >
      {confettiPieces.map((piece, index) => (
        <motion.span
          key={`${piece.left}-${piece.delay}`}
          className={`absolute top-0 h-3 w-1.5 ${colors[index % colors.length]}`}
          style={{ left: piece.left }}
          initial={{ y: "-5vh", rotate: 0, opacity: 1 }}
          animate={{ y: "105vh", rotate: piece.rotate, opacity: [1, 1, 0] }}
          transition={{
            duration: piece.duration,
            delay: piece.delay,
            repeat: Infinity,
            repeatDelay: 0.35,
            ease: "easeIn",
          }}
        />
      ))}
    </div>
  );
}

function AudienceButton({
  selected,
  icon,
  title,
  description,
  onClick,
}: {
  selected: boolean;
  icon: React.ReactNode;
  title: string;
  description: string;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      role="radio"
      aria-checked={selected}
      className={`grid w-full grid-cols-[2.5rem_minmax(0,1fr)_auto] items-center gap-3 rounded-lg border p-3 text-left transition-colors ${
        selected
          ? "border-blue-600 bg-blue-50 ring-1 ring-blue-600 dark:border-[#d7ff43] dark:bg-[#d7ff43]/8 dark:ring-[#d7ff43]"
          : "border-border/60 bg-background/60 hover:bg-accent/50"
      }`}
      onClick={onClick}
    >
      <span
        className={`flex size-10 items-center justify-center rounded-md ${
          selected
            ? "bg-blue-600 text-white dark:bg-[#d7ff43] dark:text-[#111407]"
            : "bg-muted text-muted-foreground"
        }`}
      >
        {icon}
      </span>
      <span className="min-w-0">
        <span className="block text-sm font-semibold">{title}</span>
        <span className="mt-0.5 block truncate text-xs text-muted-foreground">
          {description}
        </span>
      </span>
      <span className="flex min-w-4 items-center justify-end">
        <span
          className={`flex size-4 items-center justify-center rounded-full border ${
            selected
              ? "border-blue-600 dark:border-[#d7ff43]"
              : "border-muted-foreground/40"
          }`}
        >
          {selected ? (
            <span className="size-2 rounded-full bg-blue-600 dark:bg-[#d7ff43]" />
          ) : null}
        </span>
      </span>
    </button>
  );
}
