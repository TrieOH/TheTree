import { Image as ImageIcon, ImagePlus, Trash2, Upload } from "lucide-react";
import { useEffect, useRef } from "react";
import { toast } from "sonner";
import { Button } from "@/shared/ui/shadcn/button";
import { useUploadQueue } from "../../upload-queue";
import { usePatchEventMutation } from "../api/mutations";
import type { EventI } from "../model";

type ImageField = "logo_url" | "banner_url";

export function EventImageActions({
  event,
  field,
  compact = false,
}: {
  event: EventI;
  field: ImageField;
  compact?: boolean;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const { enqueue, tasks } = useUploadQueue();
  const patchMutation = usePatchEventMutation();
  const currentUrl = event[field];
  const label = field === "logo_url" ? "logo" : "banner";
  const activeTask = tasks.find(
    (task) =>
      task.owner.type === "event" &&
      task.owner.id === event.id &&
      task.association?.handlerKey === "event-image" &&
      task.association.input?.field === field &&
      !["completed", "failed", "rejected"].includes(task.status),
  );

  useEffect(() => {
    const toastId = `event-image-upload-${event.id}-${field}`;
    if (activeTask) {
      toast.loading(`Enviando ${label} do evento…`, {
        id: toastId,
        duration: Infinity,
        description: "Você poderá continuar quando o upload terminar.",
      });
    } else {
      toast.dismiss(toastId);
    }
    return () => {
      toast.dismiss(toastId);
    };
  }, [activeTask, event.id, field, label]);

  const patchImage = async (url: string | null) => {
    const response = await patchMutation.mutateAsync({
      eventId: event.id,
      data: {
        full_name: event.full_name,
        slug: event.slug,
        acronym: event.acronym,
        description: event.description,
        contact_email: event.contact_email,
        logo_url: field === "logo_url" ? url : event.logo_url,
        banner_url: field === "banner_url" ? url : event.banner_url,
      },
    });
    if (!response.success) throw new Error(response.message);
  };

  const handleFile = async (file?: File) => {
    if (!file) return;
    try {
      await enqueue({
        file,
        owner: { type: "event", id: event.id, label: event.full_name },
        mediaType: label,
        label: `${event.full_name} — ${label}`,
        storagePath: `events/${event.id}/${label}`,
        correctionPath: `/admin/events/${event.id}`,
        association: { handlerKey: "event-image", input: { field } },
      });
      toast.success("Imagem adicionada à fila de upload.");
    } catch (error) {
      toast.error(
        error instanceof Error
          ? error.message
          : "Não foi possível iniciar o upload.",
      );
    }
  };

  const imageLabel = field === "logo_url" ? "Logo" : "Banner";
  const ImageTypeIcon = field === "logo_url" ? ImageIcon : ImagePlus;

  return (
    <div
      className={compact ? "flex items-center gap-1" : "flex flex-wrap gap-2"}
    >
      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={(e) => {
          void handleFile(e.target.files?.[0]);
          e.currentTarget.value = "";
        }}
      />
      <Button
        type="button"
        size="sm"
        variant="outline"
        title={`${currentUrl ? "Trocar" : "Adicionar"} ${imageLabel}`}
        disabled={Boolean(activeTask) || patchMutation.isPending}
        onClick={() => inputRef.current?.click()}
      >
        {compact ? (
          <ImageTypeIcon className="size-4" />
        ) : currentUrl ? (
          <ImagePlus className="size-4" />
        ) : (
          <Upload className="size-4" />
        )}
        {compact
          ? null
          : `${currentUrl ? "Trocar" : "Adicionar"} ${imageLabel}`}
      </Button>
      {currentUrl ? (
        <Button
          type="button"
          size="sm"
          variant="ghost"
          disabled={Boolean(activeTask) || patchMutation.isPending}
          onClick={() => void patchImage(null)}
          title={`Remover ${imageLabel}`}
        >
          <Trash2 className="size-4" />
          {compact ? null : "Remover"}
        </Button>
      ) : null}
    </div>
  );
}
