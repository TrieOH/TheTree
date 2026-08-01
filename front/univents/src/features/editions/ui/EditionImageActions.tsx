import { ImagePlus, Trash2 } from "lucide-react";
import { useEffect, useRef } from "react";
import { toast } from "sonner";
import { useUploadQueue } from "@/features/upload-queue";
import { Button } from "@/shared/ui/shadcn/button";
import { usePatchEditionMutation } from "../api/mutations";
import type { EditionI } from "../model";

export function EditionImageActions({
  edition,
  eventId,
  field,
  compact = false,
}: {
  edition: EditionI;
  eventId: string;
  field: "logo_url" | "banner_url";
  compact?: boolean;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const { enqueue, tasks } = useUploadQueue();
  const mutation = usePatchEditionMutation();
  const active = tasks.some(
    (task) =>
      task.owner.type === "edition" &&
      task.owner.id === edition.id &&
      task.association?.handlerKey === "edition-image" &&
      task.association.input?.field === field &&
      !["completed", "failed", "rejected"].includes(task.status),
  );
  const label = field === "logo_url" ? "Logo" : "Banner";

  useEffect(() => {
    const id = `edition-image-upload-${edition.id}-${field}`;
    if (active)
      toast.loading(`Enviando ${label} da edição…`, { id, duration: Infinity });
    else toast.dismiss(id);
    return () => {
      toast.dismiss(id);
    };
  }, [active, edition.id, field, label]);

  const patch = (url: string | null) =>
    mutation.mutate({
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
  return (
    <div className="flex items-center gap-1">
      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={(e) => {
          const file = e.target.files?.[0];
          e.target.value = "";
          if (file)
            void enqueue({
              file,
              owner: { type: "edition", id: edition.id, label: edition.name },
              mediaType: label.toLowerCase(),
              storagePath: `events/${eventId}/editions/${edition.id}/${label.toLowerCase()}`,
              association: {
                handlerKey: "edition-image",
                input: { field, eventId },
              },
            });
        }}
      />
      <Button
        type="button"
        size="sm"
        variant="outline"
        title={`${edition[field] ? "Trocar" : "Adicionar"} ${label}`}
        disabled={active || mutation.isPending}
        onClick={() => inputRef.current?.click()}
      >
        <ImagePlus className="size-4" />
        {compact ? null : label}
      </Button>
      {edition[field] ? (
        <Button
          type="button"
          size="sm"
          variant="ghost"
          title={`Remover ${label}`}
          disabled={active || mutation.isPending}
          onClick={() => patch(null)}
        >
          <Trash2 className="size-4" />
          {compact ? null : "Remover"}
        </Button>
      ) : null}
    </div>
  );
}
