import { Trash2, Upload } from "lucide-react";
import type { DragEvent, RefObject } from "react";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";
import { useUploadQueue } from "../../upload-queue";
import { usePatchEventMutation } from "../api/mutations";
import type { EventI } from "../model";

type ImageField = "logo_url" | "banner_url";

export function EventVisualCard({ event }: { event: EventI }) {
  const [dragging, setDragging] = useState<ImageField>();
  const [hovered, setHovered] = useState<ImageField>();
  const [removeField, setRemoveField] = useState<ImageField>();
  const bannerInput = useRef<HTMLInputElement>(null);
  const logoInput = useRef<HTMLInputElement>(null);
  const { enqueue, tasks } = useUploadQueue();
  const patchMutation = usePatchEventMutation();

  const upload = async (field: ImageField, file?: File) => {
    if (!file) return;
    const label = field === "logo_url" ? "logo" : "banner";
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

  const remove = (field: ImageField) =>
    patchMutation.mutateAsync({
      eventId: event.id,
      data: {
        full_name: event.full_name,
        slug: event.slug,
        acronym: event.acronym,
        description: event.description,
        contact_email: event.contact_email,
        logo_url: field === "logo_url" ? null : event.logo_url,
        banner_url: field === "banner_url" ? null : event.banner_url,
      },
    });

  const isUploading = (field: ImageField) =>
    tasks.some(
      (task) =>
        task.owner.type === "event" &&
        task.owner.id === event.id &&
        task.association?.handlerKey === "event-image" &&
        task.association.input?.field === field &&
        !["completed", "failed", "rejected"].includes(task.status),
    );

  const dropProps = (field: ImageField) => ({
    onDragEnter: (e: DragEvent) => {
      e.preventDefault();
      setDragging(field);
    },
    onDragOver: (e: DragEvent) => e.preventDefault(),
    onDragLeave: () => setDragging(undefined),
    onDrop: (e: DragEvent) => {
      e.preventDefault();
      setDragging(undefined);
      void upload(field, e.dataTransfer.files[0]);
    },
  });

  useEffect(() => {
    if (
      tasks.some(
        (task) =>
          task.owner.id === event.id &&
          task.association?.handlerKey === "event-image" &&
          !["completed", "failed", "rejected"].includes(task.status),
      )
    )
      toast.loading("Enviando imagem do evento…", {
        id: `event-image-${event.id}`,
        duration: Infinity,
      });
    else toast.dismiss(`event-image-${event.id}`);
  }, [event.id, tasks]);

  const input = (
    field: ImageField,
    ref: RefObject<HTMLInputElement | null>,
  ) => (
    <input
      id={`event-${field.replace("_url", "")}-upload`}
      ref={ref}
      type="file"
      accept="image/*"
      className="hidden"
      disabled={isUploading(field) || patchMutation.isPending}
      onChange={(e) => {
        void upload(field, e.target.files?.[0]);
        e.currentTarget.value = "";
      }}
    />
  );

  return (
    <div className="relative">
      {input("banner_url", bannerInput)}
      {input("logo_url", logoInput)}
      <div
        {...dropProps("banner_url")}
        className={`group relative flex h-56 cursor-pointer items-center justify-center rounded-md border border-dashed bg-muted/20 transition-colors ${dragging === "banner_url" ? "border-primary bg-primary/10" : hovered === "banner_url" ? "border-primary" : "border-border/60"}`}
        onClick={() => {
          if (!isUploading("banner_url") && !patchMutation.isPending)
            bannerInput.current?.click();
        }}
        onMouseEnter={() => setHovered("banner_url")}
        onMouseLeave={() => setHovered(undefined)}
      >
        {event.banner_url && !isUploading("banner_url") ? (
          <img
            src={event.banner_url}
            alt="Banner do evento"
            className="size-full rounded-md object-cover"
          />
        ) : null}
        <div
          className={`pointer-events-none absolute inset-0 rounded-md bg-primary/10 transition-opacity ${hovered === "banner_url" ? "opacity-100" : "opacity-0"}`}
        />
        <div className="absolute inset-0">
          <Button
            type="button"
            size="icon"
            variant="secondary"
            className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 active:-translate-y-1/2!"
            disabled={isUploading("banner_url") || patchMutation.isPending}
            onClick={(e) => {
              e.stopPropagation();
              bannerInput.current?.click();
            }}
            aria-label="Adicionar ou trocar banner"
          >
            <Upload className="size-4" />
          </Button>
          {event.banner_url && !isUploading("banner_url") ? (
            <Button
              type="button"
              size="icon"
              variant="destructive"
              className="absolute bottom-3 right-3 bg-destructive text-destructive-foreground shadow-xl hover:bg-destructive/90"
              disabled={patchMutation.isPending}
              onClick={(e) => {
                e.stopPropagation();
                setRemoveField("banner_url");
              }}
              onMouseEnter={() => setHovered(undefined)}
              onMouseLeave={() => setHovered("banner_url")}
              aria-label="Remover banner"
            >
              <Trash2 className="size-4" />
            </Button>
          ) : null}
        </div>
      </div>
      <div
        {...dropProps("logo_url")}
        className={`group absolute -bottom-8 left-5 z-10 flex size-24 cursor-pointer items-center justify-center rounded-full border-4 border-card bg-muted shadow-xl transition-all md:size-28 ${dragging === "logo_url" ? "ring-2 ring-primary/70" : hovered === "logo_url" ? "ring-4 ring-primary/50" : ""}`}
        onClick={() => {
          if (!isUploading("logo_url") && !patchMutation.isPending)
            logoInput.current?.click();
        }}
        onMouseEnter={() => setHovered("logo_url")}
        onMouseLeave={() => setHovered(undefined)}
      >
        {event.logo_url && !isUploading("logo_url") ? (
          <img
            src={event.logo_url}
            alt="Logo do evento"
            className="size-full rounded-full object-cover"
          />
        ) : null}
        <div
          className={`pointer-events-none absolute inset-0 rounded-full bg-primary/15 transition-opacity ${hovered === "logo_url" ? "opacity-100" : "opacity-0"}`}
        />
        <div className="absolute inset-0">
          <Button
            type="button"
            size="icon"
            variant="secondary"
            className="absolute left-1/2 top-1/2 size-8 -translate-x-1/2 -translate-y-1/2 rounded-full active:-translate-y-1/2!"
            disabled={isUploading("logo_url") || patchMutation.isPending}
            onClick={(e) => {
              e.stopPropagation();
              logoInput.current?.click();
            }}
            aria-label="Adicionar ou trocar logo"
          >
            <Upload className="size-3.5" />
          </Button>
          {event.logo_url && !isUploading("logo_url") ? (
            <Button
              type="button"
              size="icon"
              variant="destructive"
              className="absolute bottom-0 right-0 size-8 rounded-full border-2 border-background bg-destructive text-destructive-foreground shadow-xl ring-1 ring-border hover:bg-destructive/90!"
              disabled={patchMutation.isPending}
              onClick={(e) => {
                e.stopPropagation();
                setRemoveField("logo_url");
              }}
              onMouseEnter={() => setHovered(undefined)}
              onMouseLeave={() => setHovered("logo_url")}
              aria-label="Remover logo"
            >
              <Trash2 className="size-3.5" />
            </Button>
          ) : null}
        </div>
      </div>
      <AlertModal
        open={Boolean(removeField)}
        onOpenChange={(open) => !open && setRemoveField(undefined)}
        title={`Remover ${removeField === "logo_url" ? "logo" : "banner"}?`}
        description="Essa imagem será removida do evento. Essa ação não pode ser desfeita."
        confirmLabel="Remover imagem"
        variant="destructive"
        loading={patchMutation.isPending}
        onConfirm={async () => {
          if (!removeField) return;
          await remove(removeField);
          setRemoveField(undefined);
        }}
      />
    </div>
  );
}
