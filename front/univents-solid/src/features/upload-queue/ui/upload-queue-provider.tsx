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
import "@/features/events/api/upload-association";
import "@/features/editions/api/upload-association";
import "@/features/products/api/upload-association";

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
  });

  createEffect(
    () => accountId(),
    (id) => {
      if (id) {
        if (notificationAccountId !== id) {
          notificationAccountId = id;
          previousStatuses.clear();
        }
        uploadQueueStore.setActiveAccount(id);
        void uploadQueueProcessor.start();
      } else {
        notificationAccountId = undefined;
        previousStatuses.clear();
        uploadQueueStore.setActiveAccount(undefined);
        uploadQueueProcessor.stop();
      }
    },
  );

  createEffect(
    () =>
      queue.tasks.map((task) => ({
        id: task.id,
        status: task.status,
        label: task.label,
        task,
      })),
    (items) => {
      for (const item of items) {
        const prev = previousStatuses.get(item.id);
        if (prev && prev !== item.status && notificationStatuses.has(item.status)) {
          if (item.status === "completed") {
            toast.success(`Upload concluído: ${item.label}`);
          } else if (item.status === "rejected") {
            toast.error(`A imagem foi rejeitada: ${item.label}`, {
              action: {
                label: "Corrigir",
                onClick: () => handleCorrectImage(item.task),
              },
            });
          } else if (item.status === "failed") {
            toast.error(`Falha no upload: ${item.label}`);
          }
        }
        previousStatuses.set(item.id, item.status);
      }
    },
  );

  return <>{props.children}</>;
}
