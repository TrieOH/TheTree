import type { JSX } from "@solidjs/web";
import { useNavigate } from "@tanstack/solid-router";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { createEffect, onCleanup } from "solid-js";
import { toast } from "@/shared/ui/toast";
import { useUploadQueue } from "../hooks/use-upload-queue";
import { uploadQueueProcessor } from "../lib/processor";
import { uploadQueueStore } from "../lib/store";
import type { UploadTaskStatus } from "../model/types";

const notificationStatuses = new Set<UploadTaskStatus>([
  "completed",
  "failed",
  "rejected",
]);

export function UploadQueueProvider(props: { children?: JSX.Element }): JSX.Element {
  const navigate = useNavigate();
  const { auth, isAuthenticated } = useAuth();
  const accountId = () => (isAuthenticated() ? auth.profile()?.id : undefined);
  const queue = useUploadQueue();
  const previousStatuses = new Map<string, UploadTaskStatus>();
  let notificationAccountId: string | undefined;

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
        previousStatuses.clear();
        for (const task of tasks) {
          previousStatuses.set(task.id, task.status);
        }
        return;
      }

      if (previousStatuses.size === 0) {
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

        if (task.status === "rejected" || task.error?.requiresReplacement) {
          toast.error(
            task.error?.message ?? `É necessário substituir ${task.label}.`,
            {
              action: {
                label: "Corrigir imagem",
                onClick: () => {
                  if (task.correctionPath) {
                    void navigate({ to: task.correctionPath });
                  } else {
                    void navigate({
                      to: "/admin/uploads",
                      search: { task: task.id },
                    });
                  }
                },
              },
            },
          );
          continue;
        }

        toast.error(
          task.error?.message ?? `Não foi possível enviar ${task.label}.`,
          {
            action: task.error?.retryable
              ? {
                  label: "Tentar novamente",
                  onClick: () => void queue.retry(task.id),
                }
              : {
                  label: "Ver detalhes",
                  onClick: () =>
                    void navigate({
                      to: "/admin/uploads",
                      search: { task: task.id },
                    }),
                },
          },
        );
      }

      previousStatuses.clear();
      for (const task of tasks) {
        previousStatuses.set(task.id, task.status);
      }
    },
  );

  return <>{props.children}</>;
}
