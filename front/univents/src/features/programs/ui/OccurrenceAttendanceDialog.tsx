import { useQuery } from "@tanstack/react-query";
import type { ProgramParticipant } from "@trieoh/univents-api/schemas";
import { Camera, Check, Search, UserCheck } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/shadcn/dialog";
import { Input } from "@/shared/ui/shadcn/input";
import { useSidebar } from "@/widgets/sidebar/hooks/use-sidebar";
import { occurrenceParticipantsQueryOptions } from "../api";
import { useMarkParticipationAttendedMutation } from "../api/mutations";
import { actorIdFromQr } from "../lib/actor-id-from-qr";

type BarcodeDetectorLike = {
  detect: (source: HTMLVideoElement) => Promise<Array<{ rawValue: string }>>;
};

export function OccurrenceAttendanceDialog({
  occurrenceId,
  open,
  onOpenChange,
}: {
  occurrenceId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const { collapsed } = useSidebar();
  const { data: participants = [], isPending } = useQuery(
    occurrenceParticipantsQueryOptions(occurrenceId),
  );
  const mutation = useMarkParticipationAttendedMutation(occurrenceId);
  const [search, setSearch] = useState("");
  const [scanning, setScanning] = useState(false);

  const markByQr = (value: string) => {
    const actorId = actorIdFromQr(value);
    const participant = participants.find(
      (item) => item.attendee_user_id.toLowerCase() === actorId.toLowerCase(),
    );
    if (!participant) {
      toast.error(
        "Este QR Code não pertence a um participante desta atividade",
      );
      return;
    }
    if (participant.status === "attended") {
      toast.info("Presença já confirmada");
      return;
    }
    setScanning(false);
    mutation.mutate(participant.id);
  };

  const filtered = participants.filter((participant) =>
    [
      participant.attendee_name,
      participant.attendee_email,
      participant.attendee_user_id,
    ].some((value) => value.toLowerCase().includes(search.toLowerCase())),
  );

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className={cn(
          "box-border min-w-0 w-[calc(100vw-1rem)] max-w-[calc(100vw-1rem)] max-h-[90dvh] overflow-x-hidden overflow-y-auto p-4 sm:max-w-2xl sm:p-6",
          collapsed
            ? "lg:left-[calc(50%+2.25rem)] lg:max-w-[min(42rem,calc(100vw-6rem))]"
            : "lg:left-[calc(50%+9rem)] lg:max-w-[min(42rem,calc(100vw-20rem))]",
        )}
      >
        <DialogHeader className="min-w-0">
          <DialogTitle>Presença na atividade</DialogTitle>
          <DialogDescription>
            Leia o QR do crachá ou marque o participante pela lista.
          </DialogDescription>
        </DialogHeader>

        {scanning ? (
          <QrScanner onDetected={markByQr} onStop={() => setScanning(false)} />
        ) : (
          <Button
            type="button"
            className="gap-2"
            onClick={() => setScanning(true)}
          >
            <Camera className="size-4" /> Ler QR Code
          </Button>
        )}

        <div className="relative min-w-0">
          <Search className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Buscar por nome, e-mail ou ID…"
            className="pl-9!"
          />
        </div>

        <div className="min-w-0 space-y-2">
          {isPending ? (
            <p className="py-6 text-center text-sm text-muted-foreground">
              Carregando participantes…
            </p>
          ) : filtered.length === 0 ? (
            <p className="py-6 text-center text-sm text-muted-foreground">
              Nenhum participante encontrado.
            </p>
          ) : (
            filtered.map((participant) => (
              <ParticipantRow
                key={participant.id}
                participant={participant}
                pending={
                  mutation.isPending && mutation.variables === participant.id
                }
                onAttend={() => mutation.mutate(participant.id)}
              />
            ))
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}

function ParticipantRow({
  participant,
  pending,
  onAttend,
}: {
  participant: ProgramParticipant;
  pending: boolean;
  onAttend: () => void;
}) {
  const attended = participant.status === "attended";
  return (
    <div className="flex min-w-0 items-center gap-2 rounded-lg border border-border/60 p-2.5 sm:gap-3 sm:p-3">
      <div className="flex size-9 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
        <UserCheck className="size-4" />
      </div>
      <div className="min-w-0 flex-1">
        <p className="truncate text-sm font-medium">
          {participant.attendee_name}
        </p>
        <p className="truncate text-xs text-muted-foreground">
          {participant.attendee_email}
        </p>
      </div>
      <Button
        type="button"
        size="sm"
        variant={attended ? "outline" : "default"}
        className="min-w-0 shrink-0 px-2 text-xs sm:px-3 sm:text-sm"
        disabled={attended || pending}
        onClick={onAttend}
      >
        {attended ? (
          <>
            <Check className="mr-1.5 size-4" /> Presente
          </>
        ) : pending ? (
          "Salvando…"
        ) : (
          "Dar presença"
        )}
      </Button>
    </div>
  );
}

function QrScanner({
  onDetected,
  onStop,
}: {
  onDetected: (value: string) => void;
  onStop: () => void;
}) {
  const videoRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    let stopped = false;
    let stream: MediaStream | undefined;
    let frame = 0;
    let stopFallback: (() => void) | undefined;
    let detected = false;
    const Detector = (
      window as unknown as {
        BarcodeDetector?: new (options: {
          formats: string[];
        }) => BarcodeDetectorLike;
      }
    ).BarcodeDetector;

    const finishDetection = (value: string) => {
      if (stopped || detected) return;
      detected = true;
      stopped = true;
      cancelAnimationFrame(frame);
      stopFallback?.();
      stream?.getTracks().forEach((track) => {
        track.stop();
      });
      onDetected(value);
    };

    const startNativeScanner = async () => {
      if (!Detector || !videoRef.current) return false;

      const detector = new Detector({ formats: ["qr_code"] });
      const scan = async () => {
        if (stopped || !videoRef.current) return;
        const codes = await detector.detect(videoRef.current).catch(() => []);
        if (codes[0]?.rawValue) finishDetection(codes[0].rawValue);
        else frame = requestAnimationFrame(scan);
      };

      stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: { ideal: "environment" } },
      });
      if (stopped || !videoRef.current) return false;
      videoRef.current.srcObject = stream;
      await videoRef.current.play();
      void scan();
      return true;
    };

    const startFallbackScanner = async () => {
      const { BrowserQRCodeReader } = await import("@zxing/browser");
      if (stopped || !videoRef.current) return;

      const reader = new BrowserQRCodeReader();
      const controls = await reader.decodeFromConstraints(
        { video: { facingMode: { ideal: "environment" } } },
        videoRef.current,
        (result) => {
          if (!stopped && result?.getText()) finishDetection(result.getText());
        },
      );
      stopFallback = () => controls.stop();
    };

    void (async () => {
      try {
        const nativeStarted = await startNativeScanner();
        if (!nativeStarted) await startFallbackScanner();
      } catch {
        toast.error("Não foi possível acessar a câmera");
        onStop();
      }
    })();

    return () => {
      stopped = true;
      cancelAnimationFrame(frame);
      stopFallback?.();
      stream?.getTracks().forEach((track) => {
        track.stop();
      });
    };
  }, [onDetected, onStop]);

  return (
    <div className="space-y-2">
      <video
        ref={videoRef}
        muted
        playsInline
        className="mx-auto aspect-square w-full max-w-sm rounded-lg bg-black object-cover"
      />
      <Button
        type="button"
        variant="outline"
        className="w-full"
        onClick={onStop}
      >
        Cancelar leitura
      </Button>
    </div>
  );
}
