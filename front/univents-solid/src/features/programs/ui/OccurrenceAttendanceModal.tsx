import type { JSX } from "@solidjs/web";
import { Portal } from "@solidjs/web";
import { useQuery } from "@trieoh/front-core-solid";
import type { ProgramParticipant } from "@trieoh/univents-api/schemas";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
} from "solid-js";

import ArrowLeftIcon from "~icons/lucide/arrow-left";
import CameraIcon from "~icons/lucide/camera";
import CameraOffIcon from "~icons/lucide/camera-off";
import CheckIcon from "~icons/lucide/check";
import ClockIcon from "~icons/lucide/clock";
import QrCodeIcon from "~icons/lucide/qr-code";
import SearchIcon from "~icons/lucide/search";
import UserCheckIcon from "~icons/lucide/user-check";
import UsersIcon from "~icons/lucide/users";
import XIcon from "~icons/lucide/x";

import { Button, Input, cn } from "@trieoh/ui-solid";
import { occurrenceParticipantsQueryOptions } from "../api";
import {
  useCheckpointCheckInMutation,
  useMarkParticipationAttendedMutation,
} from "../api/mutations";
import { actorIdFromQr } from "../lib/actor-id-from-qr";
import { formatOccurrenceSchedule } from "../lib/format-schedule";
import type { OccurrenceI } from "../model";
import { toast } from "@/shared/ui/toast";

const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const Camera = CameraIcon as unknown as (props: { class?: string }) => JSX.Element;
const CameraOff = CameraOffIcon as unknown as (props: { class?: string }) => JSX.Element;
const Check = CheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const QrCode = QrCodeIcon as unknown as (props: { class?: string }) => JSX.Element;
const Search = SearchIcon as unknown as (props: { class?: string }) => JSX.Element;
const UserCheck = UserCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;
const X = XIcon as unknown as (props: { class?: string }) => JSX.Element;

type BarcodeDetectorLike = {
  detect: (
    source: HTMLVideoElement,
  ) => Promise<Array<{ rawValue: string }>> | Array<{ rawValue: string }>;
};

type FilterStatus = "all" | "attended" | "pending";

export interface OccurrenceAttendanceModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  occurrence: OccurrenceI | null;
  occurrenceId?: string;
  programKind?: "activity" | "checkpoint";
  programName?: string;
}

export function OccurrenceAttendanceModal(
  props: OccurrenceAttendanceModalProps,
): JSX.Element {
  const [search, setSearch] = createSignal("");
  const [statusFilter, setStatusFilter] = createSignal<FilterStatus>("all");
  const [scanning, setScanning] = createSignal(false);
  const [savingId, setSavingId] = createSignal<string | null>(null);

  const occId = () => props.occurrence?.id ?? props.occurrenceId ?? "";
  const kind = () => props.programKind ?? "activity";

  const scheduleText = createMemo(() => {
    return formatOccurrenceSchedule(
      props.occurrence?.starts_at,
      props.occurrence?.ends_at,
    );
  });

  const participantsQuery = useQuery(() => ({
    ...occurrenceParticipantsQueryOptions(occId()),
    enabled: Boolean(props.open && occId()),
  }));

  const markMutation = useMarkParticipationAttendedMutation(occId);
  const checkInMutation = useCheckpointCheckInMutation(occId);

  const participants = createMemo(
    () => (participantsQuery().data ?? []) as ProgramParticipant[],
  );

  const isPending = () => participantsQuery().isPending;

  // Lock body scroll while open
  createEffect(
    () => props.open,
    (isOpen) => {
      if (!isOpen) return;
      const prev = document.body.style.overflow;
      document.body.style.overflow = "hidden";
      return () => {
        document.body.style.overflow = prev;
      };
    },
  );

  // Escape key handler: closes scanner if active, otherwise closes modal
  createEffect(
    () => props.open,
    (isOpen) => {
      if (!isOpen) return;
      const onKeyDown = (event: KeyboardEvent) => {
        if (event.key === "Escape") {
          if (scanning()) {
            setScanning(false);
          } else {
            props.onOpenChange(false);
          }
        }
      };
      window.addEventListener("keydown", onKeyDown);
      return () => {
        window.removeEventListener("keydown", onKeyDown);
      };
    },
  );

  const handleAttend = async (participantId: string) => {
    setSavingId(participantId);
    try {
      await markMutation.mutateAsync(participantId);
      toast.success("Presença confirmada!");
    } catch {
      toast.error("Não foi possível marcar presença");
    } finally {
      setSavingId(null);
    }
  };

  const markByQr = async (value: string) => {
    const actorId = actorIdFromQr(value);

    if (kind() === "checkpoint") {
      setScanning(false);
      try {
        await checkInMutation.mutateAsync(actorId);
        toast.success("Presença confirmada!");
      } catch {
        toast.error("Não foi possível fazer check-in");
      }
      return;
    }

    const participant = participants().find(
      (item) =>
        item.attendee_user_id?.toLowerCase() === actorId.toLowerCase() ||
        item.id.toLowerCase() === actorId.toLowerCase() ||
        item.attendee_email?.toLowerCase() === value.trim().toLowerCase(),
    );

    if (!participant) {
      toast.error("Este QR Code não pertence a um participante desta atividade");
      return;
    }

    if (participant.status === "attended") {
      toast.info("Presença já confirmada");
      return;
    }

    setScanning(false);
    await handleAttend(participant.id);
  };

  const counts = createMemo(() => {
    const list = participants();
    let attended = 0;
    for (let i = 0; i < list.length; i++) {
      if (list[i]!.status === "attended") attended++;
    }
    return {
      total: list.length,
      attended,
      pending: list.length - attended,
      percentage: list.length > 0 ? Math.round((attended / list.length) * 100) : 0,
    };
  });

  const filtered = createMemo(() => {
    const term = search().trim().toLowerCase();
    const st = statusFilter();
    const list = participants();
    const result: ProgramParticipant[] = [];

    for (let i = 0; i < list.length; i++) {
      const p = list[i]!;
      if (st === "attended" && p.status !== "attended") continue;
      if (st === "pending" && p.status === "attended") continue;

      if (!term) {
        result.push(p);
        continue;
      }
      const name = p.attendee_name?.toLowerCase() ?? "";
      const email = p.attendee_email?.toLowerCase() ?? "";
      const userId = p.attendee_user_id?.toLowerCase() ?? "";
      if (name.includes(term) || email.includes(term) || userId.includes(term)) {
        result.push(p);
      }
    }
    return result;
  });

  return (
    <Show when={props.open && occId()}>
      <Portal>
        <main class="fixed inset-0 z-100 h-dvh w-screen bg-background text-foreground overflow-hidden flex flex-col">
          {/* Main Layout: Desktop 2-column, Mobile 1-column */}
          <div class="flex-1 flex flex-col lg:flex-row min-h-0 overflow-hidden">
            {/* Sidebar (Desktop left / Mobile top) */}
            <aside class="w-full lg:w-80 xl:w-96 shrink-0 border-b lg:border-b-0 lg:border-r border-border flex flex-col bg-muted/20">
              {/* Header section with back button and program details */}
              <div class="p-4 sm:p-5 border-b border-border flex flex-col gap-3">
                <div class="flex items-center justify-between gap-2">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    class="-ml-2 gap-1.5 text-muted-foreground hover:text-foreground cursor-pointer"
                    onClick={() => props.onOpenChange(false)}
                  >
                    <ArrowLeft class="size-4" />
                    <span class="text-sm font-medium">Voltar</span>
                  </Button>

                  {/* Mobile-only camera toggle button */}
                  <div class="lg:hidden">
                    <Button
                      type="button"
                      size="sm"
                      variant={scanning() ? "destructive" : "default"}
                      class="gap-1.5 h-8 px-3 text-xs cursor-pointer rounded-lg"
                      onClick={() => setScanning(!scanning())}
                    >
                      <Show when={scanning()} fallback={<Camera class="size-3.5" />}>
                        <CameraOff class="size-3.5" />
                      </Show>
                      <span>{scanning() ? "Fechar câmera" : "Ler QR Code"}</span>
                    </Button>
                  </div>
                </div>

                <div class="min-w-0">
                  <p class="text-xs font-semibold uppercase tracking-wider text-primary">
                    {kind() === "checkpoint" ? "Check-in" : "Lista de Presença"}
                  </p>
                  <h1 class="text-base sm:text-lg font-bold text-foreground truncate mt-0.5">
                    {props.programName || "Programação"}
                  </h1>
                  <Show when={scheduleText()}>
                    <div class="flex items-center gap-1.5 text-xs text-muted-foreground mt-1">
                      <Clock class="size-3.5 text-primary shrink-0" />
                      <span class="font-medium">{scheduleText()}</span>
                    </div>
                  </Show>
                </div>

                {/* Progress bar and counter */}
                <div class="pt-1">
                  <div class="flex items-center justify-between text-xs font-semibold mb-1.5">
                    <span class="text-muted-foreground">Presentes</span>
                    <span class="text-foreground">
                      {counts().attended} / {counts().total}{" "}
                      <span class="text-primary font-bold">({counts().percentage}%)</span>
                    </span>
                  </div>
                  <div class="h-2 w-full rounded-full bg-muted overflow-hidden">
                    <div
                      class="h-full bg-primary transition-all duration-300 ease-out rounded-full"
                      style={{ width: `${counts().percentage}%` }}
                    />
                  </div>
                </div>
              </div>

              {/* QR Scanner for Desktop (always accessible) & Mobile (shown when toggle is on) */}
              <div
                class={cn(
                  "p-4 sm:p-5 flex flex-col items-center justify-center gap-3 overflow-y-auto",
                  !scanning() && "hidden lg:flex",
                )}
              >
                <Show
                  when={scanning()}
                  fallback={
                    <div class="w-full flex flex-col items-center justify-center p-6 border-2 border-dashed border-border rounded-2xl text-center bg-card">
                      <div class="flex size-12 items-center justify-center rounded-xl bg-primary/10 text-primary mb-3">
                        <QrCode class="size-6" />
                      </div>
                      <p class="text-sm font-semibold text-foreground">Leitor de QR Code</p>
                      <p class="text-xs text-muted-foreground mt-1 max-w-50">
                        Aponte o crachá do participante para confirmar presença instantaneamente.
                      </p>
                      <Button
                        type="button"
                        class="mt-4 gap-2 cursor-pointer w-full"
                        onClick={() => setScanning(true)}
                      >
                        <Camera class="size-4" />
                        <span>Abrir Câmera</span>
                      </Button>
                    </div>
                  }
                >
                  <div class="w-full flex flex-col items-center">
                    <QrScanner
                      onDetected={markByQr}
                      onStop={() => setScanning(false)}
                    />
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      class="w-full mt-3 cursor-pointer"
                      onClick={() => setScanning(false)}
                    >
                      Fechar câmera
                    </Button>
                  </div>
                </Show>
              </div>
            </aside>

            {/* Right Workstation: Search, Filters & Participant List */}
            <section class="flex-1 flex flex-col min-h-0 bg-background overflow-hidden">
              {/* Search & Filter Toolbar */}
              <div class="p-3 sm:p-5 border-b border-border flex flex-col sm:flex-row items-stretch sm:items-center gap-3 bg-background shrink-0">
                <div class="relative flex-1">
                  <Search class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    value={search()}
                    onInput={(e) => setSearch(e.currentTarget.value)}
                    placeholder="Buscar por nome, e-mail ou código..."
                    class="pl-9! pr-8! h-10 text-sm w-full"
                  />
                  <Show when={search()}>
                    <button
                      type="button"
                      onClick={() => setSearch("")}
                      class="absolute right-2.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground cursor-pointer"
                    >
                      <X class="size-3.5" />
                    </button>
                  </Show>
                </div>

                {/* Filter Tabs */}
                <div class="inline-flex rounded-lg border border-border p-1 bg-muted/50 text-xs font-medium shrink-0">
                  <button
                    type="button"
                    onClick={() => setStatusFilter("all")}
                    class={cn(
                      "px-3 py-1.5 rounded-md transition-colors cursor-pointer",
                      statusFilter() === "all"
                        ? "bg-card text-foreground shadow-xs font-semibold"
                        : "text-muted-foreground hover:text-foreground",
                    )}
                  >
                    Todos ({counts().total})
                  </button>
                  <button
                    type="button"
                    onClick={() => setStatusFilter("attended")}
                    class={cn(
                      "px-3 py-1.5 rounded-md transition-colors cursor-pointer",
                      statusFilter() === "attended"
                        ? "bg-card text-foreground shadow-xs font-semibold"
                        : "text-muted-foreground hover:text-foreground",
                    )}
                  >
                    Presentes ({counts().attended})
                  </button>
                  <button
                    type="button"
                    onClick={() => setStatusFilter("pending")}
                    class={cn(
                      "px-3 py-1.5 rounded-md transition-colors cursor-pointer",
                      statusFilter() === "pending"
                        ? "bg-card text-foreground shadow-xs font-semibold"
                        : "text-muted-foreground hover:text-foreground",
                    )}
                  >
                    Pendentes ({counts().pending})
                  </button>
                </div>
              </div>

              {/* Scrollable Participants List */}
              <div class="flex-1 overflow-y-auto p-3 sm:p-5">
                <Show
                  when={!isPending()}
                  fallback={
                    <div class="py-16 text-center text-sm text-muted-foreground flex flex-col items-center justify-center">
                      <div class="size-8 animate-spin rounded-full border-2 border-primary border-t-transparent mb-3" />
                      Carregando participantes...
                    </div>
                  }
                >
                  <Show
                    when={filtered().length > 0}
                    fallback={
                      <div class="py-16 text-center text-muted-foreground flex flex-col items-center justify-center">
                        <Users class="size-10 opacity-40 mb-3" />
                        <p class="text-sm font-semibold text-foreground">Nenhum participante encontrado</p>
                        <p class="text-xs text-muted-foreground mt-1">
                          {search()
                            ? "Nenhum resultado corresponde à busca."
                            : "Não há participantes para este filtro."}
                        </p>
                      </div>
                    }
                  >
                    <div class="divide-y divide-border/60 border border-border rounded-xl bg-card overflow-hidden">
                      <For each={filtered()}>
                        {(participant) => {
                          const attended = () => participant.status === "attended";
                          const pending = () => savingId() === participant.id;

                          return (
                            <div class="flex items-center justify-between gap-3 p-3 sm:p-3.5 transition-colors hover:bg-muted/30">
                              <div class="flex items-center gap-3 min-w-0 flex-1">
                                <div
                                  class={cn(
                                    "flex size-9 shrink-0 items-center justify-center rounded-full text-xs font-bold",
                                    attended()
                                      ? "bg-primary/10 text-primary"
                                      : "bg-muted text-muted-foreground",
                                  )}
                                >
                                  <Show
                                    when={attended()}
                                    fallback={<UserCheck class="size-4" />}
                                  >
                                    <Check class="size-4" />
                                  </Show>
                                </div>

                                <div class="min-w-0 flex-1">
                                  <p class="truncate text-sm font-semibold text-foreground">
                                    {participant.attendee_name || "Participante"}
                                  </p>
                                  <p class="truncate text-xs text-muted-foreground">
                                    {participant.attendee_email || participant.attendee_user_id}
                                  </p>
                                </div>
                              </div>

                              <div class="shrink-0">
                                <Show
                                  when={!attended()}
                                  fallback={
                                    <span class="inline-flex items-center gap-1 text-xs font-semibold text-primary bg-primary/10 px-2.5 py-1 rounded-md border border-primary/20">
                                      <Check class="size-3.5" /> Presente
                                    </span>
                                  }
                                >
                                  <Button
                                    type="button"
                                    size="sm"
                                    class="h-8 px-3 text-xs font-medium cursor-pointer"
                                    disabled={pending()}
                                    onClick={() => handleAttend(participant.id)}
                                  >
                                    <Show when={pending()} fallback="Dar presença">
                                      Salvando...
                                    </Show>
                                  </Button>
                                </Show>
                              </div>
                            </div>
                          );
                        }}
                      </For>
                    </div>
                  </Show>
                </Show>
              </div>
            </section>
          </div>
        </main>
      </Portal>
    </Show>
  );
}

function QrScanner(props: {
  onDetected: (value: string) => void;
  onStop: () => void;
}): JSX.Element {
  let videoRef: HTMLVideoElement | null = null;

  createEffect(
    () => true,
    () => {
      const handleDetected = props.onDetected;
      const handleStop = props.onStop;
      let stopped = false;
      let stream: MediaStream | undefined;
      let frame = 0;
      let stopFallback: (() => void) | undefined;
      let detected = false;

      const finishDetection = (value: string) => {
        if (stopped || detected) return;
        detected = true;
        stopped = true;
        if (frame) window.cancelAnimationFrame(frame);
        stopFallback?.();
        stream?.getTracks().forEach((track) => track.stop());
        handleDetected(value);
      };

      const startNativeScanner = async () => {
        const Detector = (
          window as unknown as {
            BarcodeDetector?: new (options: {
              formats: string[];
            }) => BarcodeDetectorLike;
          }
        ).BarcodeDetector;

        if (!Detector || !videoRef) return false;

        try {
          const detector = new Detector({ formats: ["qr_code"] });
          stream = await navigator.mediaDevices.getUserMedia({
            video: { facingMode: { ideal: "environment" } },
          });

          if (stopped || !videoRef) {
            stream.getTracks().forEach((track) => track.stop());
            return false;
          }

          videoRef.srcObject = stream;
          await videoRef.play().catch(() => { });

          const scan = async () => {
            if (stopped || !videoRef) return;
            try {
              if (videoRef.readyState >= 2) {
                const codes = await detector.detect(videoRef);
                if (codes.length > 0 && codes[0]?.rawValue) {
                  finishDetection(codes[0].rawValue);
                  return;
                }
              }
            } catch {
              // Frame detection error, continue next tick
            }
            frame = window.requestAnimationFrame(scan);
          };

          frame = window.requestAnimationFrame(scan);
          return true;
        } catch {
          return false;
        }
      };

      const startFallbackScanner = async () => {
        if (stopped || !videoRef) return;
        const { BrowserQRCodeReader } = await import("@zxing/browser");
        if (stopped || !videoRef) return;

        const reader = new BrowserQRCodeReader();
        const controls = await reader.decodeFromConstraints(
          { video: { facingMode: { ideal: "environment" } } },
          videoRef,
          (result) => {
            if (!stopped && result?.getText()) {
              finishDetection(result.getText());
            }
          },
        );
        stopFallback = () => controls.stop();
      };

      void (async () => {
        try {
          const nativeStarted = await startNativeScanner();
          if (!nativeStarted && !stopped) {
            await startFallbackScanner();
          }
        } catch {
          toast.error("Não foi possível acessar a câmera.");
          handleStop();
        }
      })();

      return () => {
        stopped = true;
        if (frame) window.cancelAnimationFrame(frame);
        stopFallback?.();
        stream?.getTracks().forEach((track) => track.stop());
      };
    },
  );

  return (
    <div class="flex flex-col items-center w-full">
      <div class="relative aspect-square w-full max-w-60 sm:max-w-70 rounded-2xl overflow-hidden bg-black flex items-center justify-center border border-border shadow-lg">
        <video
          ref={(el) => {
            videoRef = el;
          }}
          class="h-full w-full object-cover"
          muted
          playsinline
        />
        {/* Modern Viewfinder Corners */}
        <div class="pointer-events-none absolute inset-6 flex items-center justify-center">
          <div class="relative h-full w-full">
            <span class="absolute top-0 left-0 size-4 border-t-2 border-l-2 border-primary rounded-tl-sm" />
            <span class="absolute top-0 right-0 size-4 border-t-2 border-r-2 border-primary rounded-tr-sm" />
            <span class="absolute bottom-0 left-0 size-4 border-b-2 border-l-2 border-primary rounded-bl-sm" />
            <span class="absolute bottom-0 right-0 size-4 border-b-2 border-r-2 border-primary rounded-br-sm" />
            <div class="absolute inset-x-0 top-1/2 -translate-y-1/2 h-0.5 bg-primary shadow-[0_0_8px_rgba(var(--primary),0.8)] animate-pulse" />
          </div>
        </div>
      </div>
      <p class="text-xs text-muted-foreground text-center mt-2.5 max-w-60">
        Aponte a câmera para o QR Code do crachá.
      </p>
    </div>
  );
}
