import type { JSX } from "@solidjs/web";
import { Portal } from "@solidjs/web";
import { useQuery } from "@trieoh/front-core-solid";
import type { ProgramParticipant } from "@trieoh/univents-api/schemas";
import {
  For,
  Show,
  createMemo,
  createSignal,
  onSettled,
} from "solid-js";

import ArrowLeftIcon from "~icons/lucide/arrow-left";
import CheckCircle2Icon from "~icons/lucide/check-circle-2";
import ClockIcon from "~icons/lucide/clock";
import GiftIcon from "~icons/lucide/gift";
import HistoryIcon from "~icons/lucide/history";
import ListChecksIcon from "~icons/lucide/list-checks";
import LoaderCircleIcon from "~icons/lucide/loader-circle";
import RotateCcwIcon from "~icons/lucide/rotate-ccw";
import UsersIcon from "~icons/lucide/users";
import XIcon from "~icons/lucide/x";

import { Button, cn } from "@trieoh/ui-solid";
import { occurrenceParticipantsQueryOptions } from "../api";
import { drawSequence, drawTimeline, randomItem } from "../lib/draw-sequence";
import { formatOccurrenceSchedule } from "../lib/format-schedule";
import type { OccurrenceI } from "../model";

const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const CheckCircle2 = CheckCircle2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Gift = GiftIcon as unknown as (props: { class?: string }) => JSX.Element;
const History = HistoryIcon as unknown as (props: { class?: string }) => JSX.Element;
const ListChecks = ListChecksIcon as unknown as (props: { class?: string }) => JSX.Element;
const LoaderCircle = LoaderCircleIcon as unknown as (props: { class?: string }) => JSX.Element;
const RotateCcw = RotateCcwIcon as unknown as (props: { class?: string }) => JSX.Element;
const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;
const X = XIcon as unknown as (props: { class?: string }) => JSX.Element;

type Audience = "attended" | "registered";

export interface OccurrenceDrawPageProps {
  occurrenceId: string;
  occurrence?: OccurrenceI | null;
  programName: string;
  programKind?: "activity" | "checkpoint";
  onBack: () => void;
}

const confettiPieces = Array.from({ length: 36 }, (_, index) => ({
  left: `${(index * 37) % 100}%`,
  delay: `${(index % 7) * 0.08}s`,
  duration: `${1.7 + (index % 5) * 0.18}s`,
  rotate: `${240 + (index % 4) * 100}deg`,
  color: [
    "bg-blue-600 dark:bg-[#d7ff43]",
    "bg-slate-900 dark:bg-white",
    "bg-sky-300 dark:bg-[#7dd3fc]",
  ][index % 3]!,
}));

export function OccurrenceDrawPage(props: OccurrenceDrawPageProps): JSX.Element {
  let sectionRef: HTMLElement | undefined;

  const participantsQuery = useQuery(() => ({
    ...occurrenceParticipantsQueryOptions(props.occurrenceId),
    enabled: Boolean(props.occurrenceId),
  }));

  const participants = createMemo(
    () => (participantsQuery().data ?? []) as ProgramParticipant[],
  );

  const isPending = () => participantsQuery().isPending;

  const [audience, setAudience] = createSignal<Audience>("attended");
  const [winner, setWinner] = createSignal<ProgramParticipant | undefined>(undefined);
  const [displayed, setDisplayed] = createSignal<ProgramParticipant | undefined>(undefined);
  const [drawing, setDrawing] = createSignal(false);
  const [immersive, setImmersive] = createSignal(false);
  const [countdown, setCountdown] = createSignal<number | undefined>(undefined);
  const [_singleParticipantDraw, setSingleParticipantDraw] = createSignal(false);
  const [allowRepeatWinners, setAllowRepeatWinners] = createSignal(false);
  const [winnerHistory, setWinnerHistory] = createSignal<ProgramParticipant[]>([]);

  let timers: number[] = [];

  const clearTimers = () => {
    timers.forEach((t) => window.clearTimeout(t));
    timers = [];
  };

  const participantName = (participant: ProgramParticipant) => {
    return participant.attendee_name || participant.attendee_email;
  };

  const scheduleText = createMemo(() => {
    return formatOccurrenceSchedule(
      props.occurrence?.starts_at,
      props.occurrence?.ends_at,
    );
  });

  // FLIP transition to animate section smoothly between side column and full screen
  function animateLayoutTransition(el: HTMLElement, firstRect: DOMRect) {
    const lastRect = el.getBoundingClientRect();
    const deltaX = firstRect.left - lastRect.left;
    const deltaY = firstRect.top - lastRect.top;
    const scaleX = firstRect.width / (lastRect.width || 1);
    const scaleY = firstRect.height / (lastRect.height || 1);

    if (
      Math.abs(deltaX) < 1 &&
      Math.abs(deltaY) < 1 &&
      Math.abs(scaleX - 1) < 0.01 &&
      Math.abs(scaleY - 1) < 0.01
    ) {
      return;
    }

    el.animate(
      [
        {
          transformOrigin: "top left",
          transform: `translate(${deltaX}px, ${deltaY}px) scale(${scaleX}, ${scaleY})`,
        },
        {
          transformOrigin: "top left",
          transform: "none",
        },
      ],
      {
        duration: 580,
        easing: "cubic-bezier(0.22, 1, 0.36, 1)",
      },
    );
  }

  const toggleImmersive = (next: boolean) => {
    if (immersive() === next) return;
    if (!sectionRef) {
      setImmersive(next);
      return;
    }
    const firstRect = sectionRef.getBoundingClientRect();
    setImmersive(next);
    requestAnimationFrame(() => {
      if (sectionRef) {
        animateLayoutTransition(sectionRef, firstRect);
      }
    });
  };

  // Lock body scroll and register escape handler
  onSettled(() => {
    const prev = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        if (drawing()) return; // Don't cancel mid-draw animation
        if (immersive()) {
          toggleImmersive(false);
        } else {
          props.onBack();
        }
      }
    };

    window.addEventListener("keydown", onKeyDown);

    return () => {
      document.body.style.overflow = prev;
      window.removeEventListener("keydown", onKeyDown);
      clearTimers();
    };
  });

  const eligible = createMemo(() => {
    const list = participants();
    const aud = audience();
    const res: ProgramParticipant[] = [];
    for (let i = 0; i < list.length; i++) {
      const p = list[i]!;
      if (aud === "attended" ? p.status === "attended" : p.status !== "cancelled") {
        res.push(p);
      }
    }
    return res;
  });

  const available = createMemo(() => {
    if (allowRepeatWinners()) return eligible();
    const history = winnerHistory();
    const winnerIds = new Set(history.map((p) => p.id));
    return eligible().filter((p) => !winnerIds.has(p.id));
  });

  const changeAudience = (next: Audience) => {
    clearTimers();
    setAudience(next);
    setWinner(undefined);
    setDisplayed(undefined);
    setDrawing(false);
    setCountdown(undefined);
    setSingleParticipantDraw(false);
  };

  const draw = () => {
    const list = available();
    const selected = randomItem(list);
    if (!selected) return;

    clearTimers();
    setWinner(undefined);
    setDisplayed(undefined);
    toggleImmersive(true);

    if (list.length === 1) {
      setSingleParticipantDraw(true);
      setDisplayed(selected);
      setWinner(selected);
      setWinnerHistory((prev) => [...prev, selected]);
      setDrawing(false);
      setCountdown(undefined);
      return;
    }

    setSingleParticipantDraw(false);
    setDrawing(true);
    setCountdown(3);
    const sequence = drawSequence(list);
    const { delays, durationMs } = drawTimeline(sequence.length, 6170);

    [2, 1].forEach((value, index) => {
      timers.push(
        window.setTimeout(() => setCountdown(value), 550 * (index + 1)),
      );
    });
    timers.push(window.setTimeout(() => setCountdown(undefined), 1650));

    sequence.forEach((participant, index) => {
      timers.push(
        window.setTimeout(() => {
          setDisplayed(participant);
        }, 1650 + delays[index]!),
      );
    });

    timers.push(
      window.setTimeout(
        () => {
          setDisplayed(selected);
          setWinner(selected);
          setWinnerHistory((prev) => [...prev, selected]);
          setDrawing(false);
        },
        1650 + durationMs + 180,
      ),
    );
  };

  return (
    <Portal>
      <main class="fixed inset-0 z-100 h-dvh w-screen overflow-y-auto bg-white text-[#17201b] dark:bg-[#141618] dark:text-[#f3f5f2] lg:overflow-hidden select-none">
        <style>{`
          @keyframes confetti-fall {
            0% {
              transform: translateY(-5vh) rotate(0deg);
              opacity: 1;
            }
            75% {
              opacity: 1;
            }
            100% {
              transform: translateY(105vh) rotate(var(--rot, 360deg));
              opacity: 0;
            }
          }
          @keyframes slot-flip {
            0% {
              opacity: 0.25;
              transform: translateY(42px);
              filter: blur(10px);
            }
            100% {
              opacity: 1;
              transform: translateY(0);
              filter: blur(0px);
            }
          }
          @keyframes countdown-pop {
            0% {
              opacity: 0;
              transform: scale(0.82) translateY(24px);
            }
            100% {
              opacity: 1;
              transform: scale(1) translateY(0);
            }
          }
          @keyframes winner-reveal {
            0% {
              opacity: 0;
              transform: translateY(32px);
            }
            100% {
              opacity: 1;
              transform: translateY(0);
            }
          }
          @keyframes winner-bar {
            0% {
              width: 0px;
            }
            100% {
              width: 72px;
            }
          }
          .animate-slot-flip {
            animation: slot-flip 0.075s cubic-bezier(0.22, 1, 0.36, 1) forwards;
          }
          .animate-countdown-pop {
            animation: countdown-pop 0.28s cubic-bezier(0.22, 1, 0.36, 1) forwards;
          }
          .animate-winner-reveal {
            animation: winner-reveal 0.45s cubic-bezier(0.22, 1, 0.36, 1) forwards;
          }
          .animate-winner-bar {
            animation: winner-bar 0.35s 0.2s cubic-bezier(0.22, 1, 0.36, 1) forwards;
          }
        `}</style>

        <div class="grid min-h-full min-w-0 lg:h-full lg:grid-cols-[21rem_minmax(0,1fr)]">
          {/* Left Column: Config Sidebar */}
          <aside class="space-y-6 border-b border-slate-200 p-5 dark:border-white/10 lg:h-full lg:overflow-y-auto lg:border-r lg:border-b-0 lg:p-6">
            {/* Top Back / Header */}
            <div class="flex min-w-0 items-center gap-2">
              <Button
                type="button"
                variant="ghost"
                size="sm"
                class="-ml-2 shrink-0 text-muted-foreground hover:text-foreground cursor-pointer"
                onClick={props.onBack}
              >
                <ArrowLeft class="size-4" />
                <span class="sr-only">Voltar para ocorrências</span>
              </Button>
              <div class="min-w-0 flex-1">
                <p class="text-xs text-muted-foreground font-medium">Sorteio</p>
                <p class="truncate text-sm font-semibold text-foreground">
                  {props.programName}
                </p>
                <Show when={scheduleText()}>
                  <div class="flex items-center gap-1 text-[11px] text-muted-foreground mt-0.5">
                    <Clock class="size-3 text-primary shrink-0" />
                    <span class="truncate font-medium">{scheduleText()}</span>
                  </div>
                </Show>
              </div>
            </div>

            {/* Audience Fieldset */}
            <fieldset>
              <legend class="mb-4">
                <span class="block text-xs font-semibold text-blue-700 uppercase dark:text-[#d7ff43]">
                  Público do sorteio
                </span>
                <span class="mt-1 block text-lg font-semibold">
                  Quem pode concorrer?
                </span>
              </legend>
              <div class="space-y-2">
                <AudienceButton
                  selected={audience() === "attended"}
                  icon={<CheckCircle2 class="size-5" />}
                  title="Só presentes"
                  description="Presença confirmada"
                  onClick={() => changeAudience("attended")}
                />
                <AudienceButton
                  selected={audience() === "registered"}
                  icon={<ListChecks class="size-5" />}
                  title="Todos inscritos"
                  description="Com ou sem presença"
                  onClick={() => changeAudience("registered")}
                />
              </div>
            </fieldset>

            {/* Total Box */}
            <div class="border-y border-border/60 py-4">
              <div class="flex items-center gap-3">
                <div class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-blue-100 text-blue-700 dark:bg-[#d7ff43]/10 dark:text-[#d7ff43]">
                  <Users class="size-5" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-xs text-muted-foreground">Total no sorteio</p>
                  <p class="mt-0.5 text-2xl font-semibold tabular-nums">
                    {isPending() ? "–" : eligible().length}
                  </p>
                </div>
                <span class="text-right text-xs text-muted-foreground">
                  {audience() === "attended" ? "presentes" : "inscritos"}
                </span>
              </div>
            </div>

            {/* Repetir Vencedores & Histórico */}
            <section class="space-y-3">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <p class="text-sm font-semibold">Repetir vencedores</p>
                  <p class="text-xs text-muted-foreground">
                    Permitir quem já ganhou
                  </p>
                </div>
                <button
                  type="button"
                  role="switch"
                  aria-checked={allowRepeatWinners() ? "true" : "false"}
                  disabled={drawing()}
                  onClick={() => setAllowRepeatWinners(!allowRepeatWinners())}
                  class={cn(
                    "relative inline-flex h-6 w-11 shrink-0 items-center rounded-full border-2 border-transparent transition-colors cursor-pointer disabled:opacity-50",
                    allowRepeatWinners()
                      ? "bg-blue-600 dark:bg-[#d7ff43]"
                      : "bg-slate-300 dark:bg-white/20",
                  )}
                >
                  <span
                    class={cn(
                      "pointer-events-none inline-block size-5 rounded-full bg-white shadow-sm transition-transform dark:bg-[#111407]",
                      allowRepeatWinners() ? "translate-x-5" : "translate-x-0",
                    )}
                  />
                  <span class="sr-only">Permitir vencedores repetidos</span>
                </button>
              </div>

              <div class="flex items-center justify-between gap-3 border-t border-border/60 pt-3">
                <div class="flex items-center gap-2">
                  <History class="size-4 text-muted-foreground" />
                  <p class="text-sm font-semibold">Vencedores</p>
                  <Show when={winnerHistory().length > 0}>
                    <span class="text-xs tabular-nums text-muted-foreground">
                      {winnerHistory().length}
                    </span>
                  </Show>
                </div>
                <Show when={winnerHistory().length > 0}>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    class="h-7 gap-1.5 px-2 text-xs text-muted-foreground cursor-pointer"
                    disabled={drawing()}
                    onClick={() => {
                      setWinnerHistory([]);
                      setWinner(undefined);
                      setDisplayed(undefined);
                      setSingleParticipantDraw(false);
                    }}
                  >
                    <RotateCcw class="size-3.5" />
                    Limpar
                  </Button>
                </Show>
              </div>

              <Show
                when={winnerHistory().length > 0}
                fallback={
                  <p class="text-xs text-muted-foreground">
                    Os resultados aparecerão aqui.
                  </p>
                }
              >
                <ol class="space-y-1 max-h-48 overflow-y-auto pr-1">
                  <For each={[...winnerHistory()].reverse()}>
                    {(participant, index) => (
                      <li class="flex min-w-0 items-center gap-2 py-1.5 text-sm">
                        <span class="w-5 shrink-0 text-right text-xs tabular-nums text-muted-foreground">
                          {winnerHistory().length - index()}.
                        </span>
                        <span class="truncate font-medium">
                          {participantName(participant)}
                        </span>
                      </li>
                    )}
                  </For>
                </ol>
              </Show>
            </section>
          </aside>

          {/* Right Column / Immersive Screen with FLIP animation */}
          <section
            ref={(el) => {
              sectionRef = el;
            }}
            class={cn(
              "flex min-w-0 flex-col items-center justify-center overflow-hidden bg-slate-50 px-5 py-10 text-center text-[#111827] dark:bg-[#101214] dark:text-[#f4f5f6] sm:px-10 relative transition-[border-radius] duration-500",
              immersive()
                ? "fixed inset-0 z-100 min-h-dvh border-0"
                : "relative min-h-120 sm:min-h-144",
            )}
          >
            {/* Top blue/lime accent bar */}
            <div class="absolute inset-x-0 top-0 h-1 bg-blue-600 dark:bg-[#d7ff43]" />

            {/* Immersive Top Bar */}
            <Show when={immersive()}>
              <div class="absolute inset-x-5 top-5 flex items-center justify-between gap-4 sm:inset-x-8 sm:top-7 z-20">
                <div class="min-w-0 text-left">
                  <p class="truncate text-sm font-medium">{props.programName}</p>
                  <Show when={scheduleText()}>
                    <p class="text-xs text-primary font-medium">{scheduleText()}</p>
                  </Show>
                  <p class="mt-0.5 text-xs opacity-55">
                    {eligible().length}{" "}
                    {eligible().length === 1 ? "participante" : "participantes"}
                  </p>
                </div>
                <Show when={!drawing()}>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    class="shrink-0 text-slate-700 hover:bg-slate-100 hover:text-slate-950 dark:text-[#f4f5f6] dark:hover:bg-white/10 dark:hover:text-white cursor-pointer"
                    onClick={() => toggleImmersive(false)}
                    title="Sair da tela cheia (Esc)"
                  >
                    <X class="size-5" />
                    <span class="sr-only">Sair da tela cheia</span>
                  </Button>
                </Show>
              </div>
            </Show>

            {/* Interactive Draw Content */}
            <Show
              when={!isPending()}
              fallback={<LoaderCircle class="size-10 animate-spin opacity-70 text-blue-600 dark:text-[#d7ff43]" />}
            >
              <Show
                when={eligible().length > 0}
                fallback={
                  <div class="flex flex-col items-center">
                    <Users class="mb-4 size-12 opacity-60 text-muted-foreground" />
                    <p class="text-xl font-semibold">
                      Ninguém elegível para o sorteio
                    </p>
                    <p class="mt-2 max-w-md text-sm opacity-65">
                      {audience() === "attended"
                        ? "Confirme presenças ou escolha todos os inscritos."
                        : "Ainda não há inscrições nesta atividade."}
                    </p>
                  </div>
                }
              >
                <Show
                  when={countdown() !== undefined}
                  fallback={
                    <Show
                      when={winner()}
                      fallback={
                        <Show
                          when={displayed()}
                          fallback={
                            <Show
                              when={available().length > 0}
                              fallback={
                                <div class="flex flex-col items-center">
                                  <History class="mb-4 size-12 opacity-60" />
                                  <p class="text-xl font-semibold">
                                    Todos já foram sorteados
                                  </p>
                                  <p class="mt-2 max-w-md text-sm opacity-65">
                                    Limpe o histórico ou permita repetir vencedores.
                                  </p>
                                </div>
                              }
                            >
                              <div class="flex flex-col items-center">
                                <Gift class="mb-5 size-12 text-blue-600 dark:text-[#d7ff43]" />
                                <p class="text-3xl font-bold sm:text-5xl">
                                  Pronto para sortear
                                </p>
                                <p class="mt-3 text-sm opacity-65 sm:text-base">
                                  A seleção será feita entre {eligible().length}{" "}
                                  {eligible().length === 1 ? "pessoa" : "pessoas"}.
                                </p>
                              </div>
                            </Show>
                          }
                        >
                          {/* Slot Machine Display */}
                          <div class="max-w-full">
                            <div class="relative flex h-32 w-[min(92vw,70rem)] items-center justify-center overflow-hidden sm:h-44">
                              <Show when={displayed()} keyed>
                                {(item) => (
                                  <p class="absolute inset-0 flex items-center justify-center overflow-hidden wrap-break-word px-3 text-3xl leading-tight font-semibold sm:text-6xl text-foreground animate-slot-flip">
                                    <span class="line-clamp-2">
                                      {participantName(item)}
                                    </span>
                                  </p>
                                )}
                              </Show>
                            </div>
                          </div>
                        </Show>
                      }
                    >
                      {/* Confetti across the entire section viewport */}
                      <Confetti />

                      {/* Winner Display */}
                      <div class="relative z-20 max-w-full flex flex-col items-center animate-winner-reveal">
                        <div class="mx-auto mb-8 h-1 bg-blue-600 dark:bg-[#d7ff43] animate-winner-bar" />
                        <p class="text-sm font-semibold text-blue-700 uppercase dark:text-[#d7ff43] sm:text-base">
                          Vencedor
                        </p>
                        <p class="mt-4 max-w-[90vw] wrap-break-word text-4xl font-bold sm:text-7xl">
                          {participantName(winner()!)}
                        </p>
                        <p class="mt-3 break-all text-sm opacity-65 sm:text-base">
                          {winner()!.attendee_email}
                        </p>
                      </div>
                    </Show>
                  }
                >
                  {/* Countdown Display */}
                  <Show when={countdown()} keyed>
                    {(num) => (
                      <div class="text-center animate-countdown-pop">
                        <p class="mb-5 text-xs font-medium text-slate-500 uppercase dark:text-white/50">
                          Preparando sorteio
                        </p>
                        <p class="text-8xl font-semibold tabular-nums text-blue-600 dark:text-[#d7ff43] sm:text-9xl">
                          {num}
                        </p>
                      </div>
                    )}
                  </Show>
                </Show>
              </Show>
            </Show>

            {/* Bottom Actions */}
            <Show when={!drawing()}>
              <div class="relative z-10 mt-10 flex flex-wrap justify-center gap-3">
                <Button
                  type="button"
                  size="lg"
                  class="min-w-52 gap-2 bg-blue-600 text-white hover:bg-blue-700 dark:bg-[#d7ff43] dark:text-[#111407] dark:hover:bg-[#e3ff77] cursor-pointer shadow-md"
                  disabled={
                    isPending() ||
                    eligible().length === 0 ||
                    available().length === 0
                  }
                  onClick={draw}
                >
                  <Gift class="size-4" />
                  <span>{winner() ? "Sortear novamente" : "Começar sorteio"}</span>
                </Button>

                <Show when={winner() && immersive()}>
                  <Button
                    type="button"
                    size="lg"
                    variant="outline"
                    class="border-slate-300 bg-transparent text-slate-700 hover:bg-slate-100 hover:text-slate-950 dark:border-white/25 dark:text-[#f4f5f6] dark:hover:bg-white/10 dark:hover:text-white cursor-pointer"
                    onClick={() => toggleImmersive(false)}
                  >
                    Encerrar
                  </Button>
                </Show>
              </div>
            </Show>
          </section>
        </div>
      </main>
    </Portal>
  );
}

function AudienceButton(props: {
  selected: boolean;
  icon: JSX.Element;
  title: string;
  description: string;
  onClick: () => void;
}): JSX.Element {
  return (
    <button
      type="button"
      onClick={() => props.onClick()}
      class={cn(
        "flex w-full items-start gap-3 rounded-xl border p-3.5 text-left transition-all cursor-pointer",
        props.selected
          ? "border-blue-600 bg-blue-50 text-blue-950 dark:border-[#d7ff43] dark:bg-[#d7ff43]/10 dark:text-[#f3f5f2]"
          : "border-slate-200 bg-white text-slate-600 hover:bg-slate-50 dark:border-white/10 dark:bg-transparent dark:text-white/70 dark:hover:bg-white/5",
      )}
    >
      <div
        class={cn(
          "mt-0.5 shrink-0",
          props.selected
            ? "text-blue-600 dark:text-[#d7ff43]"
            : "text-muted-foreground",
        )}
      >
        {props.icon}
      </div>
      <div class="min-w-0 flex-1">
        <p class="text-sm font-semibold text-foreground">{props.title}</p>
        <p class="text-xs text-muted-foreground">{props.description}</p>
      </div>
    </button>
  );
}

function Confetti(): JSX.Element {
  return (
    <div
      class="pointer-events-none absolute inset-0 overflow-hidden z-10"
      aria-hidden="true"
    >
      <For each={confettiPieces}>
        {(piece) => (
          <span
            class={cn("absolute top-0 h-3.5 w-2 rounded-[1px]", piece.color)}
            style={{
              left: piece.left,
              "--rot": piece.rotate,
              "animation-name": "confetti-fall",
              "animation-duration": piece.duration,
              "animation-delay": piece.delay,
              "animation-iteration-count": "infinite",
              "animation-timing-function": "cubic-bezier(0.4, 0, 0.9, 1)",
            }}
          />
        )}
      </For>
    </div>
  );
}
