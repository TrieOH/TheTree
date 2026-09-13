import type { JSX } from "@solidjs/web";
import { Match, Show, Switch, createEffect, createMemo, createSignal, onCleanup } from "solid-js";
import { Button } from "@trieoh/ui-solid";

import AlertTriangleIcon from "~icons/lucide/alert-triangle";
import CheckIcon from "~icons/lucide/check";
import CheckCircle2Icon from "~icons/lucide/check-circle-2";
import Clock3Icon from "~icons/lucide/clock-3";
import FileImageIcon from "~icons/lucide/file-image";
import Link2Icon from "~icons/lucide/link-2";
import Loader2Icon from "~icons/lucide/loader-2";
import RefreshCwIcon from "~icons/lucide/refresh-cw";
import ReplaceIcon from "~icons/lucide/replace";
import Trash2Icon from "~icons/lucide/trash-2";
import UploadCloudIcon from "~icons/lucide/upload-cloud";
import WifiOffIcon from "~icons/lucide/wifi-off";

import { uploadQueueConfig } from "../lib/config";
import type { UploadTask, UploadTaskStatus } from "../model/types";

const AlertTriangle = AlertTriangleIcon as unknown as (props: { class?: string }) => JSX.Element;
const Check = CheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const CheckCircle2 = CheckCircle2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Clock3 = Clock3Icon as unknown as (props: { class?: string }) => JSX.Element;
const FileImage = FileImageIcon as unknown as (props: { class?: string }) => JSX.Element;
const Link2 = Link2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;
const RefreshCw = RefreshCwIcon as unknown as (props: { class?: string }) => JSX.Element;
const Replace = ReplaceIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;
const UploadCloud = UploadCloudIcon as unknown as (props: { class?: string }) => JSX.Element;
const WifiOff = WifiOffIcon as unknown as (props: { class?: string }) => JSX.Element;

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
  queued: "bg-muted/10",
  uploading: "bg-primary/[0.02]",
  associating: "bg-primary/[0.02]",
  waiting_retry: "bg-amber-500/[0.02]",
  completed: "bg-emerald-500/[0.02]",
  failed: "bg-destructive/[0.02]",
  rejected: "bg-destructive/[0.02]",
  paused: "bg-amber-500/[0.02]",
};

function TaskPreview(props: { task: UploadTask }) {
  const [objectUrl, setObjectUrl] = createSignal<string | undefined>(undefined);
  let currentObjectUrl: string | undefined;

  onCleanup(() => {
    if (currentObjectUrl) URL.revokeObjectURL(currentObjectUrl);
  });

  createEffect(
    () => ({ url: props.task.uploadedUrl, file: props.task.file }),
    ({ url, file }) => {
      if (currentObjectUrl) {
        URL.revokeObjectURL(currentObjectUrl);
        currentObjectUrl = undefined;
      }
      if (url || file.size === 0) {
        setObjectUrl(undefined);
        return;
      }
      const created = URL.createObjectURL(file);
      currentObjectUrl = created;
      setObjectUrl(created);
    },
  );

  const source = () => props.task.uploadedUrl ?? objectUrl();

  return (
    <Show
      when={source()}
      fallback={
        <div class="flex h-full min-h-32 items-center justify-center bg-muted text-muted-foreground">
          <FileImage class="size-8" />
        </div>
      }
    >
      {(src) => (
        <img
          src={src()}
          alt={props.task.label}
          class="h-full min-h-32 w-full object-cover"
        />
      )}
    </Show>
  );
}

function StatusIcon(props: { status: UploadTaskStatus }) {
  return (
    <Switch fallback={<AlertTriangle class="size-4 text-destructive" />}>
      <Match when={props.status === "completed"}>
        <CheckCircle2 class="size-4 text-emerald-600" />
      </Match>
      <Match when={props.status === "uploading" || props.status === "associating"}>
        <Loader2 class="size-4 animate-spin text-primary" />
      </Match>
      <Match when={props.status === "waiting_retry"}>
        <WifiOff class="size-4 text-amber-600" />
      </Match>
      <Match when={props.status === "queued"}>
        <Clock3 class="size-4" />
      </Match>
    </Switch>
  );
}

function StageStep(props: {
  icon: (p: { class?: string }) => JSX.Element;
  label: string;
  state: "pending" | "active" | "done";
}) {
  return (
    <div
      class={`flex flex-1 items-center gap-2 rounded-lg border px-3 py-1.5 text-xs transition-colors ${props.state === "done"
          ? "border-emerald-500/20 bg-emerald-500/5 text-emerald-700"
          : props.state === "active"
            ? "border-primary/25 bg-primary/5 text-primary"
            : "border-border/70 text-muted-foreground"
        }`}
    >
      <span class="flex size-5 shrink-0 items-center justify-center rounded-full bg-background shadow-xs">
        <Show
          when={props.state === "done"}
          fallback={<props.icon class={`size-3 ${props.state === "active" ? "animate-pulse" : ""}`} />}
        >
          <Check class="size-3" />
        </Show>
      </span>
      <span class="font-medium">{props.label}</span>
    </div>
  );
}

export function UploadTaskCard(props: {
  task: UploadTask;
  highlighted: boolean;
  onRetry: () => void;
  onReplace: (file: File) => void;
  onRemove: () => void;
}) {
  let fileInputRef: HTMLInputElement | null = null;

  const isProcessing = createMemo(
    () => props.task.status === "uploading" || props.task.status === "associating",
  );
  const canRetry = createMemo(
    () => props.task.status === "failed" && Boolean(props.task.error?.retryable),
  );
  const needsReplacement = createMemo(
    () => props.task.status === "rejected" || Boolean(props.task.error?.requiresReplacement),
  );
  const uploadIsDone = createMemo(
    () => props.task.stage === "association" || props.task.status === "completed",
  );
  const uploadIsActive = createMemo(
    () =>
      props.task.stage === "upload" &&
      ["queued", "uploading", "waiting_retry"].includes(props.task.status),
  );
  const associationIsActive = createMemo(
    () => props.task.stage === "association" && props.task.status === "associating",
  );
  const nextAttempt = createMemo(() =>
    props.task.nextAttemptAt
      ? new Date(props.task.nextAttemptAt).toLocaleTimeString("pt-BR", {
        hour: "2-digit",
        minute: "2-digit",
      })
      : undefined,
  );

  return (
    <div
      id={`upload-${props.task.id}`}
      class={`rounded-xl border border-border bg-card text-card-foreground shadow-xs transition-colors ${statusClasses[props.task.status]
        } ${props.highlighted ? "ring-2 ring-primary ring-offset-2 ring-offset-background" : ""}`}
    >
      <div class="grid gap-0 p-0 lg:grid-cols-[12rem_minmax(0,1fr)]">
        <div class="relative overflow-hidden border-b border-border bg-muted lg:border-r lg:border-b-0 rounded-t-xl lg:rounded-tr-none lg:rounded-l-xl">
          <TaskPreview task={props.task} />
          <div class="absolute inset-x-0 bottom-0 bg-linear-to-t from-black/80 to-transparent px-3 pt-6 pb-2.5 text-white">
            <p class="truncate text-xs font-medium">{props.task.fileName}</p>
            <p class="text-[10px] text-white/70">
              {(props.task.size / 1024 / 1024).toFixed(2)} MB
            </p>
          </div>
        </div>

        <div class="flex min-w-0 flex-col p-3.5 sm:p-4">
          <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
            <div class="min-w-0">
              <div class="mb-1.5 flex flex-wrap gap-1.5">
                <span class="inline-flex items-center rounded-full border border-border px-2 py-0.5 text-[11px] font-medium text-muted-foreground">
                  {props.task.mediaType}
                </span>
                <span class="inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium text-muted-foreground">
                  {props.task.owner.type}
                </span>
              </div>
              <h2 class="truncate text-sm font-semibold sm:text-base">{props.task.label}</h2>
              <p class="mt-0.5 truncate text-xs text-muted-foreground">
                {props.task.owner.label ?? props.task.owner.id}
              </p>
            </div>

            <span
              class={`inline-flex items-center gap-1.5 rounded-full border border-border px-2.5 py-1 text-xs font-medium h-7 shrink-0 ${props.task.status === "completed"
                  ? "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
                  : props.task.status === "failed" || props.task.status === "rejected"
                    ? "bg-destructive/10 text-destructive"
                    : "bg-muted/40 text-foreground"
                }`}
            >
              <StatusIcon status={props.task.status} />
              {statusLabels[props.task.status]}
            </span>
          </div>

          <div class="my-3 flex items-center gap-2">
            <StageStep
              icon={UploadCloud}
              label="Upload e moderação"
              state={uploadIsDone() ? "done" : uploadIsActive() ? "active" : "pending"}
            />
            <div class="h-px w-3 shrink-0 bg-border" />
            <StageStep
              icon={Link2}
              label="Associação"
              state={
                props.task.status === "completed"
                  ? "done"
                  : associationIsActive()
                    ? "active"
                    : "pending"
              }
            />
          </div>

          <Show when={props.task.error}>
            {(err) => (
              <div class="mb-3 flex gap-2.5 rounded-lg border border-destructive/25 bg-destructive/5 p-2.5">
                <AlertTriangle class="mt-0.5 size-4 shrink-0 text-destructive" />
                <div class="min-w-0">
                  <p class="text-xs font-medium text-foreground">{err().message}</p>
                  <p class="mt-0.5 text-[10px] uppercase tracking-wide text-muted-foreground">
                    {err().code}
                  </p>
                </div>
              </div>
            )}
          </Show>

          <div class="mt-auto flex flex-col gap-2.5 border-t border-border pt-2.5 sm:flex-row sm:items-center sm:justify-between">
            <div class="text-xs text-muted-foreground">
              <Show
                when={props.task.retryCount > 0}
                fallback={<span>Nenhuma nova tentativa necessária</span>}
              >
                <span>
                  {props.task.retryCount >= uploadQueueConfig.maxRetries
                    ? `Limite de ${uploadQueueConfig.maxRetries} retries atingido`
                    : `${props.task.retryCount} de ${uploadQueueConfig.maxRetries} retries utilizados`}
                  {nextAttempt() ? ` · próxima tentativa às ${nextAttempt()}` : ""}
                </span>
              </Show>
            </div>

            <div class="flex flex-wrap justify-end gap-2">
              <Show when={canRetry()}>
                <Button type="button" variant="outline" onClick={props.onRetry}>
                  <RefreshCw class="size-3.5 mr-1.5" />
                  Tentar novamente
                </Button>
              </Show>

              <Show when={needsReplacement()}>
                <Button type="button" onClick={() => fileInputRef?.click()}>
                  <Replace class="size-3.5 mr-1.5" />
                  Trocar imagem
                </Button>
              </Show>

              <Button
                type="button"
                variant="ghost"
                onClick={props.onRemove}
                aria-label={isProcessing() ? "Cancelar upload" : "Remover tarefa"}
              >
                <Trash2 class="size-3.5 mr-1.5" />
                {isProcessing() ? "Cancelar" : "Remover"}
              </Button>
            </div>
          </div>
        </div>

        <input
          ref={(el) => (fileInputRef = el)}
          type="file"
          accept="image/*"
          class="hidden"
          onChange={(event) => {
            const file = event.currentTarget.files?.[0];
            if (file) props.onReplace(file);
            event.currentTarget.value = "";
          }}
        />
      </div>
    </div>
  );
}
