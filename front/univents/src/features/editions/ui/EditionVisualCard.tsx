import { Trash2, Upload } from "lucide-react";
import type { DragEvent, RefObject } from "react";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { usePatchEditionMutation } from "@/features/editions/api/mutations";
import type { EditionI } from "@/features/editions/model";
import { useUploadQueue } from "@/features/upload-queue";
import { Button } from "@/shared/ui/shadcn/button";

type ImageField = "logo_url" | "banner_url";

export function EditionVisualCard({
  edition,
  eventId,
}: {
  edition: EditionI;
  eventId: string;
}) {
  const [dragging, setDragging] = useState<ImageField>();
  const [hovered, setHovered] = useState<ImageField>();
  const bannerInput = useRef<HTMLInputElement>(null);
  const logoInput = useRef<HTMLInputElement>(null);
  const { enqueue, tasks } = useUploadQueue();
  const patchMutation = usePatchEditionMutation();

  const isUploading = (field: ImageField) =>
    tasks.some(
      (task) =>
        task.owner.type === "edition" &&
        task.owner.id === edition.id &&
        task.association?.handlerKey === "edition-image" &&
        task.association.input?.field === field &&
        !["completed", "failed", "rejected"].includes(task.status),
    );

  const upload = async (field: ImageField, file?: File) => {
    if (!file) return;
    const label = field === "logo_url" ? "logo" : "banner";
    try {
      await enqueue({
        file,
        owner: { type: "edition", id: edition.id, label: edition.name },
        mediaType: label,
        label: `${edition.name} — ${label}`,
        storagePath: `events/${eventId}/editions/${edition.id}/${label}`,
        correctionPath: `/admin/events/${eventId}/editions/${edition.id}`,
        association: {
          handlerKey: "edition-image",
          input: { field, eventId },
        },
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

  const patch = (field: ImageField, url: string | null) =>
    patchMutation.mutate({
      eventId,
      editionId: edition.id,
      data: {
        name: edition.name,
        slug: edition.slug,
        starts_at: edition.starts_at,
        ends_at: edition.ends_at,
        tagline: edition.tagline,
        description: edition.description,
        registration_opens_at: edition.registration_opens_at,
        location_name: edition.location_name,
        location_description: edition.location_description,
        contact_email: edition.contact_email,
        logo_url: field === "logo_url" ? url : edition.logo_url,
        banner_url: field === "banner_url" ? url : edition.banner_url,
      },
    });

  const dropProps = (field: ImageField) => ({
    onDragEnter: (event: DragEvent) => {
      event.preventDefault();
      setDragging(field);
    },
    onDragOver: (event: DragEvent) => event.preventDefault(),
    onDragLeave: () => setDragging(undefined),
    onDrop: (event: DragEvent) => {
      event.preventDefault();
      setDragging(undefined);
      void upload(field, event.dataTransfer.files[0]);
    },
  });

  useEffect(() => {
    if (
      tasks.some(
        (task) =>
          task.owner.id === edition.id &&
          task.association?.handlerKey === "edition-image" &&
          !["completed", "failed", "rejected"].includes(task.status),
      )
    )
      toast.loading("Enviando imagem da edição…", {
        id: `edition-image-${edition.id}`,
        duration: Infinity,
      });
    else toast.dismiss(`edition-image-${edition.id}`);
  }, [edition.id, tasks]);

  const input = (
    field: ImageField,
    ref: RefObject<HTMLInputElement | null>,
  ) => (
    <input
      id={`edition-${edition.id}-${field.replace("_url", "")}-upload`}
      ref={ref}
      type="file"
      accept="image/*"
      className="hidden"
      disabled={isUploading(field) || patchMutation.isPending}
      onChange={(event) => {
        void upload(field, event.target.files?.[0]);
        event.currentTarget.value = "";
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
        {edition.banner_url && !isUploading("banner_url") ? (
          <img
            src={edition.banner_url}
            alt="Banner da edição"
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
            onClick={(event) => {
              event.stopPropagation();
              bannerInput.current?.click();
            }}
            aria-label="Adicionar ou trocar banner"
          >
            <Upload className="size-4" />
          </Button>
          {edition.banner_url && !isUploading("banner_url") ? (
            <Button
              type="button"
              size="icon"
              variant="destructive"
              className="absolute bottom-3 right-3 bg-destructive text-destructive-foreground shadow-xl hover:bg-destructive/90"
              disabled={patchMutation.isPending}
              onClick={(event) => {
                event.stopPropagation();
                patch("banner_url", null);
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
        {edition.logo_url && !isUploading("logo_url") ? (
          <img
            src={edition.logo_url}
            alt="Logo da edição"
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
            onClick={(event) => {
              event.stopPropagation();
              logoInput.current?.click();
            }}
            aria-label="Adicionar ou trocar logo"
          >
            <Upload className="size-3.5" />
          </Button>
          {edition.logo_url && !isUploading("logo_url") ? (
            <Button
              type="button"
              size="icon"
              variant="destructive"
              className="absolute bottom-0 right-0 size-8 rounded-full border-2 border-background bg-destructive text-destructive-foreground shadow-xl ring-1 ring-border hover:bg-destructive/90!"
              disabled={patchMutation.isPending}
              onClick={(event) => {
                event.stopPropagation();
                patch("logo_url", null);
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
    </div>
  );
}
