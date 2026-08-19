import {
  AlertTriangle,
  Check,
  CheckCircle2,
  Clock3,
  FileImage,
  Link2,
  Loader2,
  RefreshCw,
  Replace,
  Trash2,
  UploadCloud,
  WifiOff,
} from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Button } from "@/shared/ui/shadcn/button";
import { Card, CardContent } from "@/shared/ui/shadcn/card";
import { uploadQueueConfig } from "../lib/config";
import type { UploadTask, UploadTaskStatus } from "../model/types";

const statusLabels: Record<UploadTaskStatus, string> = {
  queued: "Na fila",
  uploading: "Enviando imagem",
  associating: "Associando imagem",
  waiting_retry: "Retry agendado",
  completed: "Concluído",
  failed: "Envio interrompido",
  rejected: "Ação necessária",
  paused: "Processamento pausado",
};

const statusClasses: Record<UploadTaskStatus, string> = {
  queued: "border-border bg-muted/30",
  uploading: "border-primary/30 bg-primary/3",
  associating: "border-primary/30 bg-primary/3",
  waiting_retry: "border-amber-500/30 bg-amber-500/3",
  completed: "border-emerald-500/25 bg-emerald-500/3",
  failed: "border-destructive/30 bg-destructive/3",
  rejected: "border-destructive/30 bg-destructive/3",
  paused: "border-amber-500/30 bg-amber-500/3",
};

function TaskPreview({ task }: { task: UploadTask }) {
  const [source, setSource] = useState(task.uploadedUrl);

  useEffect(() => {
    if (task.uploadedUrl) {
      setSource(task.uploadedUrl);
      return;
    }
    if (task.file.size === 0) {
      setSource(undefined);
      return;
    }

    const objectUrl = URL.createObjectURL(task.file);
    setSource(objectUrl);
    return () => URL.revokeObjectURL(objectUrl);
  }, [task.file, task.uploadedUrl]);

  if (!source) {
    return (
      <div className="flex h-full min-h-40 items-center justify-center bg-muted text-muted-foreground">
        <FileImage className="size-9" />
      </div>
    );
  }

  return (
    <img
      src={source}
      alt={task.label}
      className="h-full min-h-40 w-full object-cover"
    />
  );
}

function StatusIcon({ status }: { status: UploadTaskStatus }) {
  if (status === "completed")
    return <CheckCircle2 className="size-4 text-emerald-600" />;
  if (status === "uploading" || status === "associating")
    return <Loader2 className="size-4 animate-spin text-primary" />;
  if (status === "waiting_retry")
    return <WifiOff className="size-4 text-amber-600" />;
  if (status === "queued") return <Clock3 className="size-4" />;
  return <AlertTriangle className="size-4 text-destructive" />;
}

function StageStep({
  icon: Icon,
  label,
  state,
}: {
  icon: typeof UploadCloud;
  label: string;
  state: "pending" | "active" | "done";
}) {
  return (
    <div
      className={cn(
        "flex flex-1 items-center gap-2 rounded-lg border px-3 py-2 text-xs transition-colors",
        state === "done" &&
          "border-emerald-500/20 bg-emerald-500/5 text-emerald-700",
        state === "active" && "border-primary/25 bg-primary/5 text-primary",
        state === "pending" && "border-border/70 text-muted-foreground",
      )}
    >
      <span className="flex size-6 shrink-0 items-center justify-center rounded-full bg-background shadow-sm">
        {state === "done" ? (
          <Check className="size-3.5" />
        ) : (
          <Icon
            className={cn("size-3.5", state === "active" && "animate-pulse")}
          />
        )}
      </span>
      <span className="font-medium">{label}</span>
    </div>
  );
}

export function UploadTaskCard({
  task,
  highlighted,
  onRetry,
  onReplace,
  onRemove,
}: {
  task: UploadTask;
  highlighted: boolean;
  onRetry: () => void;
  onReplace: (file: File) => void;
  onRemove: () => void;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const isProcessing =
    task.status === "uploading" || task.status === "associating";
  const canRetry = task.status === "failed" && task.error?.retryable;
  const needsReplacement =
    task.status === "rejected" || task.error?.requiresReplacement;
  const uploadIsDone =
    task.stage === "association" || task.status === "completed";
  const uploadIsActive =
    task.stage === "upload" &&
    ["queued", "uploading", "waiting_retry"].includes(task.status);
  const associationIsActive =
    task.stage === "association" && task.status === "associating";
  const nextAttempt = task.nextAttemptAt
    ? new Date(task.nextAttemptAt).toLocaleTimeString("pt-BR", {
        hour: "2-digit",
        minute: "2-digit",
      })
    : undefined;

  return (
    <Card
      id={`upload-${task.id}`}
      className={cn(
        "border py-0 transition-all",
        statusClasses[task.status],
        highlighted &&
          "ring-2 ring-primary ring-offset-2 ring-offset-background",
      )}
    >
      <CardContent className="grid gap-0 p-0 lg:grid-cols-[14rem_minmax(0,1fr)]">
        <div className="relative overflow-hidden border-b bg-muted lg:border-r lg:border-b-0">
          <TaskPreview task={task} />
          <div className="absolute inset-x-0 bottom-0 bg-linear-to-t from-black/75 to-transparent px-3 pt-8 pb-3 text-white">
            <p className="truncate text-xs font-medium">{task.fileName}</p>
            <p className="text-[10px] text-white/70">
              {(task.size / 1024 / 1024).toFixed(2)} MB
            </p>
          </div>
        </div>

        <div className="flex min-w-0 flex-col p-4 sm:p-5">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
            <div className="min-w-0">
              <div className="mb-2 flex flex-wrap gap-1.5">
                <Badge variant="outline">{task.mediaType}</Badge>
                <Badge variant="ghost">{task.owner.type}</Badge>
              </div>
              <h2 className="truncate text-base font-semibold">{task.label}</h2>
              <p className="mt-1 truncate text-xs text-muted-foreground">
                {task.owner.label ?? task.owner.id}
              </p>
            </div>

            <Badge
              variant={
                task.status === "completed"
                  ? "secondary"
                  : task.status === "failed" || task.status === "rejected"
                    ? "destructive"
                    : "outline"
              }
              className="h-7 gap-1.5 px-2.5"
            >
              <StatusIcon status={task.status} />
              {statusLabels[task.status]}
            </Badge>
          </div>

          <div className="my-4 flex items-center gap-2">
            <StageStep
              icon={UploadCloud}
              label="Upload e moderação"
              state={
                uploadIsDone ? "done" : uploadIsActive ? "active" : "pending"
              }
            />
            <div className="h-px w-3 shrink-0 bg-border" />
            <StageStep
              icon={Link2}
              label="Associação"
              state={
                task.status === "completed"
                  ? "done"
                  : associationIsActive
                    ? "active"
                    : "pending"
              }
            />
          </div>

          {task.error ? (
            <div className="mb-4 flex gap-3 rounded-lg border border-destructive/20 bg-background/70 p-3">
              <AlertTriangle className="mt-0.5 size-4 shrink-0 text-destructive" />
              <div className="min-w-0">
                <p className="text-xs font-medium text-foreground">
                  {task.error.message}
                </p>
                <p className="mt-1 text-[10px] uppercase tracking-wide text-muted-foreground">
                  {task.error.code}
                </p>
              </div>
            </div>
          ) : null}

          <div className="mt-auto flex flex-col gap-3 border-t border-border/60 pt-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="text-xs text-muted-foreground">
              {task.retryCount > 0 ? (
                <span>
                  {task.retryCount >= uploadQueueConfig.maxRetries
                    ? `Limite de ${uploadQueueConfig.maxRetries} retries atingido`
                    : `${task.retryCount} de ${uploadQueueConfig.maxRetries} retries utilizados`}
                  {nextAttempt ? ` · próxima tentativa às ${nextAttempt}` : ""}
                </span>
              ) : (
                <span>Nenhuma nova tentativa necessária</span>
              )}
            </div>

            <div className="flex flex-wrap justify-end gap-2">
              {canRetry ? (
                <Button type="button" variant="outline" onClick={onRetry}>
                  <RefreshCw />
                  Tentar novamente
                </Button>
              ) : null}

              {needsReplacement ? (
                <Button type="button" onClick={() => inputRef.current?.click()}>
                  <Replace />
                  Trocar imagem
                </Button>
              ) : null}

              <Button
                type="button"
                variant="ghost"
                onClick={onRemove}
                aria-label={isProcessing ? "Cancelar upload" : "Remover tarefa"}
              >
                <Trash2 />
                {isProcessing ? "Cancelar" : "Remover"}
              </Button>
            </div>
          </div>
        </div>

        <input
          ref={inputRef}
          type="file"
          accept="image/*"
          className="hidden"
          onChange={(event) => {
            const file = event.target.files?.[0];
            if (file) onReplace(file);
            event.target.value = "";
          }}
        />
      </CardContent>
    </Card>
  );
}
