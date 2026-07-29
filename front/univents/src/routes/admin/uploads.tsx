import { createFileRoute } from "@tanstack/react-router";
import {
  AlertTriangle,
  CheckCircle2,
  Images,
  RefreshCw,
  UploadCloud,
} from "lucide-react";
import { useMemo } from "react";
import { z } from "zod";
import { useUploadQueue } from "@/features/upload-queue";
import { UploadTaskCard } from "@/features/upload-queue/ui/upload-task-card";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Button } from "@/shared/ui/shadcn/button";
import { Card, CardContent } from "@/shared/ui/shadcn/card";

const searchSchema = z.object({
  task: z.string().optional(),
});

export const Route = createFileRoute("/admin/uploads")({
  validateSearch: searchSchema,
  component: UploadsPage,
});

function UploadsPage() {
  const { task: highlightedTaskId } = Route.useSearch();
  const { tasks, initialized, retry, replaceFile, remove } = useUploadQueue();

  const counts = useMemo(
    () => ({
      active: tasks.filter((task) =>
        ["queued", "uploading", "associating", "waiting_retry"].includes(
          task.status,
        ),
      ).length,
      problems: tasks.filter((task) =>
        ["failed", "rejected", "paused"].includes(task.status),
      ).length,
      completed: tasks.filter((task) => task.status === "completed").length,
    }),
    [tasks],
  );

  const retryableTasks = tasks.filter(
    (task) => task.status === "failed" && task.error?.retryable,
  );

  return (
    <div className="mx-auto w-full max-w-6xl space-y-6 p-4 pb-28">
      <header className="flex flex-col gap-4">
        <div>
          <Badge variant="outline" className="mb-2">
            Central de uploads
          </Badge>
          <h1 className="text-2xl font-semibold tracking-tight">
            Uploads e processamento de mídia
          </h1>
          <p className="mt-1 max-w-2xl text-sm text-muted-foreground">
            Acompanhe envios em segundo plano, repita falhas recuperáveis e
            substitua imagens recusadas.
          </p>
        </div>

        {retryableTasks.length > 0 ? (
          <Button
            type="button"
            variant="outline"
            onClick={() =>
              void Promise.all(retryableTasks.map((task) => retry(task.id)))
            }
          >
            <RefreshCw />
            Tentar novamente todas
          </Button>
        ) : null}
      </header>

      <section className="grid gap-3 sm:grid-cols-3">
        <SummaryCard
          icon={UploadCloud}
          label="Em andamento"
          value={counts.active}
        />
        <SummaryCard
          icon={AlertTriangle}
          label="Precisam de atenção"
          value={counts.problems}
        />
        <SummaryCard
          icon={CheckCircle2}
          label="Concluídos"
          value={counts.completed}
        />
      </section>

      {!initialized ? (
        <Card>
          <CardContent className="py-10 text-center text-sm text-muted-foreground">
            Carregando fila persistente…
          </CardContent>
        </Card>
      ) : tasks.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center py-14 text-center">
            <Images className="mb-4 size-10 text-muted-foreground" />
            <h2 className="font-medium">Nenhum upload na fila</h2>
            <p className="mt-1 max-w-md text-sm text-muted-foreground">
              Quando uma integração adicionar imagens, o andamento e eventuais
              correções aparecerão aqui.
            </p>
          </CardContent>
        </Card>
      ) : (
        <section className="grid gap-4">
          {tasks.map((task) => (
            <UploadTaskCard
              key={task.id}
              task={task}
              highlighted={task.id === highlightedTaskId}
              onRetry={() => void retry(task.id)}
              onReplace={(file) => void replaceFile(task.id, file)}
              onRemove={() => void remove(task.id)}
            />
          ))}
        </section>
      )}
    </div>
  );
}

function SummaryCard({
  icon: Icon,
  label,
  value,
}: {
  icon: typeof UploadCloud;
  label: string;
  value: number;
}) {
  return (
    <Card size="sm">
      <CardContent className="flex items-center gap-3">
        <div className="flex size-9 items-center justify-center rounded-lg bg-muted">
          <Icon className="size-4" />
        </div>
        <div>
          <p className="text-xl font-semibold">{value}</p>
          <p className="text-xs text-muted-foreground">{label}</p>
        </div>
      </CardContent>
    </Card>
  );
}
