import type { JSX } from "@solidjs/web";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { useNavigate } from "@tanstack/solid-router";
import { createEffect, onCleanup } from "solid-js";
import { toast } from "@/shared/ui/toast";
import { uploadQueueProcessor } from "../lib/processor";
import { uploadQueueStore } from "../lib/store";
import { promptImageReplacement } from "../lib/correction";
import { useUploadQueue } from "../hooks/use-upload-queue";
import type { UploadTask, UploadTaskStatus } from "../model/types";

const notificationStatuses = new Set<UploadTaskStatus>([
  "completed",
  "rejected",
  "failed",
]);

export function UploadQueueProvider(props: { children?: JSX.Element }): JSX.Element {
  const { auth, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const accountId = () => (isAuthenticated() ? auth.profile()?.id : undefined);
  const queue = useUploadQueue();
  const previousStatuses = new Map<string, UploadTaskStatus>();
  let notificationAccountId: string | undefined;

  const handleCorrectImage = (task: UploadTask) => {
    if (
      task.correctionPath &&
      typeof window !== "undefined" &&
      window.location.pathname !== task.correctionPath
    ) {
      void navigate({ to: task.correctionPath });
    }

    promptImageReplacement(task, (file) => {
      void queue.replaceFile(task.id, file);
    });
  };

  onCleanup(() => {
    uploadQueueProcessor.stop();
    uploadQueueStore.setActiveAccount(undefined);
  });

  createEffect(
    () => accountId(),
    (id) => {
      uploadQueueStore.setActiveAccount(id);
      previousStatuses.clear();

      if (!id) {
        uploadQueueProcessor.stop();
        return;
      }

      void uploadQueueProcessor.start();
    },
  );

  createEffect(
    () => ({ id: accountId(), initialized: queue.initialized, tasks: queue.tasks }),
    ({ id, initialized, tasks }) => {
      if (!initialized || !id) return;

      if (notificationAccountId !== id) {
        notificationAccountId = id;
        for (const task of tasks) {
          previousStatuses.set(task.id, task.status);
        }
        return;
      }

      for (const task of tasks) {
        const prevStatus = previousStatuses.get(task.id);
        if (prevStatus === task.status || !notificationStatuses.has(task.status)) {
          continue;
        }

        if (task.status === "completed") {
          toast.success(`${task.label} foi enviada com sucesso.`);
          continue;
        }

        if (task.status === "rejected") {
          toast.error(
            task.error?.message ?? `${task.label} foi rejeitada.`,
            {
              action: task.correctionPath
                ? {
                  label: "Corrigir imagem",
                  onClick: () => handleCorrectImage(task),
                }
                : undefined,
            },
          );
          continue;
        }

        if (task.status === "failed") {
          toast.error(
            task.error?.message ?? `Falha no envio de ${task.label}.`,
            {
              action: task.correctionPath
                ? {
                  label: "Corrigir imagem",
                  onClick: () => handleCorrectImage(task),
                }
                : undefined,
            },
          );
        }
      }

      previousStatuses.clear();
      for (const task of tasks) {
        previousStatuses.set(task.id, task.status);
      }
    },
  );

  return <>{props.children}</>;
}
