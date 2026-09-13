import type { JSX } from "@solidjs/web";
import { For, Show, createMemo } from "solid-js";
import { createFileRoute } from "@tanstack/solid-router";
import { z } from "zod";
import { Button } from "@trieoh/ui-solid";

import AlertTriangleIcon from "~icons/lucide/alert-triangle";
import CheckCircle2Icon from "~icons/lucide/check-circle-2";
import ImagesIcon from "~icons/lucide/images";
import RefreshCwIcon from "~icons/lucide/refresh-cw";
import UploadCloudIcon from "~icons/lucide/upload-cloud";

import { useUploadQueue } from "@/features/upload-queue";
import { UploadTaskCard } from "@/features/upload-queue/ui/upload-task-card";

const AlertTriangle = AlertTriangleIcon as unknown as (props: { class?: string }) => JSX.Element;
const CheckCircle2 = CheckCircle2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Images = ImagesIcon as unknown as (props: { class?: string }) => JSX.Element;
const RefreshCw = RefreshCwIcon as unknown as (props: { class?: string }) => JSX.Element;
const UploadCloud = UploadCloudIcon as unknown as (props: { class?: string }) => JSX.Element;

const searchSchema = z.object({
  task: z.string().optional(),
});

export const Route = createFileRoute("/admin/uploads")({
  validateSearch: searchSchema,
  head: () => ({ meta: [{ title: "Central de Uploads - Admin - Univents" }] }),
  component: UploadsPage,
});

function SummaryCard(props: {
  icon: (p: { class?: string }) => JSX.Element;
  label: string;
  value: number;
}) {
  return (
    <div class="flex items-center gap-3 rounded-xl border border-border bg-card p-3.5 text-card-foreground shadow-xs">
      <div class="flex size-9 items-center justify-center rounded-lg bg-muted text-muted-foreground">
        <props.icon class="size-4" />
      </div>
      <div>
        <p class="text-xl font-semibold tracking-tight">{props.value}</p>
        <p class="text-xs text-muted-foreground">{props.label}</p>
      </div>
    </div>
  );
}

function UploadsPage(): JSX.Element {
  const search = Route.useSearch();
  const queue = useUploadQueue();

  const tasks = createMemo(() => queue.tasks);
  const initialized = createMemo(() => queue.initialized);

  const counts = createMemo(() => {
    const list = tasks();
    return {
      active: list.filter((task) =>
        ["queued", "uploading", "associating", "waiting_retry"].includes(
          task.status,
        ),
      ).length,
      problems: list.filter((task) =>
        ["failed", "rejected", "paused"].includes(task.status),
      ).length,
      completed: list.filter((task) => task.status === "completed").length,
    };
  });

  const retryableTasks = createMemo(() =>
    tasks().filter((task) => task.status === "failed" && task.error?.retryable),
  );

  return (
    <div class="w-full space-y-6">
      <header class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <span class="mb-2 inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-xs font-semibold">
            Central de uploads
          </span>
          <h1 class="text-2xl font-semibold tracking-tight">
            Uploads e processamento de mídia
          </h1>
          <p class="mt-1 max-w-2xl text-sm text-muted-foreground">
            Acompanhe envios em segundo plano, repita falhas recuperáveis e
            substitua imagens recusadas.
          </p>
        </div>

        <Show when={retryableTasks().length > 0}>
          <Button
            type="button"
            variant="outline"
            onClick={() =>
              void Promise.all(retryableTasks().map((task) => queue.retry(task.id)))
            }
          >
            <RefreshCw class="size-4 mr-2" />
            Tentar novamente todas
          </Button>
        </Show>
      </header>

      <section class="grid gap-3 sm:grid-cols-3">
        <SummaryCard
          icon={UploadCloud}
          label="Em andamento"
          value={counts().active}
        />
        <SummaryCard
          icon={AlertTriangle}
          label="Precisam de atenção"
          value={counts().problems}
        />
        <SummaryCard
          icon={CheckCircle2}
          label="Concluídos"
          value={counts().completed}
        />
      </section>

      <Show
        when={initialized()}
        fallback={
          <div class="rounded-xl border border-border bg-card p-10 text-center text-sm text-muted-foreground">
            Carregando fila persistente…
          </div>
        }
      >
        <Show
          when={tasks().length > 0}
          fallback={
            <div class="flex flex-col items-center rounded-xl border border-border bg-card py-12 px-4 text-center">
              <Images class="mb-4 size-10 text-muted-foreground" />
              <h2 class="font-medium text-foreground">Nenhum upload na fila</h2>
              <p class="mt-1 max-w-md text-sm text-muted-foreground">
                Quando uma integração adicionar imagens, o andamento e eventuais
                correções aparecerão aqui.
              </p>
            </div>
          }
        >
          <section class="grid gap-3.5">
            <For each={tasks()}>
              {(task) => (
                <UploadTaskCard
                  task={task}
                  highlighted={task.id === search().task}
                  onRetry={() => void queue.retry(task.id)}
                  onReplace={(file) => void queue.replaceFile(task.id, file)}
                  onRemove={() => void queue.remove(task.id)}
                />
              )}
            </For>
          </section>
        </Show>
      </Show>
    </div>
  );
}
