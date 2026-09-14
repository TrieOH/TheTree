import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal } from "solid-js";
import Trash2Icon from "~icons/lucide/trash-2";
import UploadIcon from "~icons/lucide/upload";
import { Button, cn } from "@trieoh/ui-solid";
import { AlertModal } from "@/widgets/ui/AlertModal";
import { useUploadQueue } from "@/features/upload-queue";
import { usePatchEditionMutation } from "../api/mutations";
import type { EditionI } from "../model";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideTrash2 = Trash2Icon as unknown as IconComp;
const LucideUpload = UploadIcon as unknown as IconComp;

type ImageField = "logo_url" | "banner_url";

export function EditionVisualCard(props: {
  edition: EditionI;
  eventId: string;
}): JSX.Element {
  const [dragging, setDragging] = createSignal<ImageField | undefined>(undefined);
  const [hovered, setHovered] = createSignal<ImageField | undefined>(undefined);
  const [removeField, setRemoveField] = createSignal<ImageField | undefined>(undefined);

  let bannerInput: HTMLInputElement | undefined;
  let logoInput: HTMLInputElement | undefined;

  const uploadQueue = useUploadQueue();
  const patchMutation = usePatchEditionMutation();

  const isPending = () => patchMutation.result().status === "pending";

  const upload = async (field: ImageField, file?: File) => {
    if (!file) return;
    const label = field === "logo_url" ? "logo" : "banner";
    try {
      await uploadQueue.enqueue({
        file,
        owner: {
          type: "edition",
          id: props.edition.id,
          label: props.edition.name,
        },
        mediaType: label,
        label: `${props.edition.name} — ${label}`,
        storagePath: `events/${props.eventId}/editions/${props.edition.id}/${label}`,
        correctionPath: `/admin/events/${props.eventId}/editions/${props.edition.id}`,
        association: {
          handlerKey: "edition-image",
          input: { field, eventId: props.eventId },
        },
      });
    } catch {
      // Handled silently
    }
  };

  const remove = (field: ImageField) =>
    patchMutation.mutateAsync({
      eventId: props.eventId,
      editionId: props.edition.id,
      data: {
        name: props.edition.name,
        slug: props.edition.slug,
        starts_at: props.edition.starts_at,
        ends_at: props.edition.ends_at,
        tagline: props.edition.tagline,
        description: props.edition.description,
        registration_opens_at: props.edition.registration_opens_at,
        location_name: props.edition.location_name,
        location_description: props.edition.location_description,
        contact_email: props.edition.contact_email,
        logo_url: field === "logo_url" ? null : props.edition.logo_url,
        banner_url: field === "banner_url" ? null : props.edition.banner_url,
      },
    });

  const isUploading = (field: ImageField) =>
    uploadQueue.tasks.some(
      (task) =>
        task.owner.type === "edition" &&
        task.owner.id === props.edition.id &&
        task.association?.handlerKey === "edition-image" &&
        task.association.input?.field === field &&
        !["completed", "failed", "rejected"].includes(task.status),
    );

  const handleDrop = (field: ImageField, e: DragEvent) => {
    e.preventDefault();
    setDragging(undefined);
    if (e.dataTransfer?.files?.[0]) {
      void upload(field, e.dataTransfer.files[0]);
    }
  };

  createEffect(
    () => props.edition.id,
    (_id) => {
      return () => {
        setDragging(undefined);
        setHovered(undefined);
        setRemoveField(undefined);
      };
    },
  );

  const handleConfirmRemove = async () => {
    const field = removeField();
    if (!field) return;
    await remove(field);
    setRemoveField(undefined);
  };

  return (
    <div class="relative">
      <input
        id={`edition-${props.edition.id}-banner-upload`}
        ref={(el) => {
          bannerInput = el;
        }}
        type="file"
        accept="image/*"
        class="hidden"
        disabled={isUploading("banner_url")}
        onChange={(e) => {
          void upload("banner_url", e.currentTarget.files?.[0]);
          e.currentTarget.value = "";
        }}
      />
      <input
        id={`edition-${props.edition.id}-logo-upload`}
        ref={(el) => {
          logoInput = el;
        }}
        type="file"
        accept="image/*"
        class="hidden"
        disabled={isUploading("logo_url")}
        onChange={(e) => {
          void upload("logo_url", e.currentTarget.files?.[0]);
          e.currentTarget.value = "";
        }}
      />

      {/* Banner Dropzone */}
      <div
        onDragEnter={(e) => {
          e.preventDefault();
          setDragging("banner_url");
        }}
        onDragOver={(e) => e.preventDefault()}
        onDragLeave={() => setDragging(undefined)}
        onDrop={(e) => handleDrop("banner_url", e)}
        class={cn(
          "group relative flex h-56 cursor-pointer items-center justify-center rounded-xl border border-dashed bg-muted/20 transition-colors",
          dragging() === "banner_url"
            ? "border-primary bg-primary/10"
            : hovered() === "banner_url"
              ? "border-primary"
              : "border-border/60",
        )}
        onClick={() => {
          if (!isUploading("banner_url") && !isPending()) {
            bannerInput?.click();
          }
        }}
        onMouseEnter={() => setHovered("banner_url")}
        onMouseLeave={() => setHovered(undefined)}
      >
        <Show when={props.edition.banner_url && !isUploading("banner_url")}>
          <img
            src={props.edition.banner_url!}
            alt="Banner da edição"
            class="size-full rounded-xl object-cover"
          />
        </Show>
        <div
          class={cn(
            "pointer-events-none absolute inset-0 rounded-xl bg-primary/10 transition-opacity",
            hovered() === "banner_url" ? "opacity-100" : "opacity-0",
          )}
        />
        <div class="absolute inset-0">
          <Button
            type="button"
            size="icon"
            variant="outline"
            class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 border-border bg-background/90 shadow-sm"
            disabled={isUploading("banner_url") || isPending()}
            onClick={(e) => {
              e.stopPropagation();
              bannerInput?.click();
            }}
            aria-label="Adicionar ou trocar banner"
          >
            <LucideUpload class="size-4" />
          </Button>
          <Show when={props.edition.banner_url && !isUploading("banner_url")}>
            <div
              class="absolute bottom-3 right-3"
              onMouseEnter={() => setHovered(undefined)}
              onMouseLeave={() => setHovered("banner_url")}
            >
              <Button
                type="button"
                size="icon"
                variant="destructive"
                class="bg-destructive text-destructive-foreground shadow-xl hover:bg-destructive/90"
                disabled={isPending()}
                onClick={(e) => {
                  e.stopPropagation();
                  setRemoveField("banner_url");
                }}
                aria-label="Remover banner"
              >
                <LucideTrash2 class="size-4" />
              </Button>
            </div>
          </Show>
        </div>
      </div>

      {/* Logo Dropzone */}
      <div
        onDragEnter={(e) => {
          e.preventDefault();
          setDragging("logo_url");
        }}
        onDragOver={(e) => e.preventDefault()}
        onDragLeave={() => setDragging(undefined)}
        onDrop={(e) => handleDrop("logo_url", e)}
        class={cn(
          "group absolute -bottom-8 left-5 z-10 flex size-24 cursor-pointer items-center justify-center rounded-full border-4 border-card bg-muted shadow-xl transition-all md:size-28",
          dragging() === "logo_url"
            ? "ring-2 ring-primary/70"
            : hovered() === "logo_url"
              ? "ring-4 ring-primary/50"
              : "",
        )}
        onClick={() => {
          if (!isUploading("logo_url") && !isPending()) {
            logoInput?.click();
          }
        }}
        onMouseEnter={() => setHovered("logo_url")}
        onMouseLeave={() => setHovered(undefined)}
      >
        <Show when={props.edition.logo_url && !isUploading("logo_url")}>
          <img
            src={props.edition.logo_url!}
            alt="Logo da edição"
            class="size-full rounded-full object-cover"
          />
        </Show>
        <div
          class={cn(
            "pointer-events-none absolute inset-0 rounded-full bg-primary/15 transition-opacity",
            hovered() === "logo_url" ? "opacity-100" : "opacity-0",
          )}
        />
        <div class="absolute inset-0">
          <Button
            type="button"
            size="icon"
            variant="outline"
            class="absolute left-1/2 top-1/2 size-8 -translate-x-1/2 -translate-y-1/2 rounded-full border-border bg-background/90 shadow-sm"
            disabled={isUploading("logo_url") || isPending()}
            onClick={(e) => {
              e.stopPropagation();
              logoInput?.click();
            }}
            aria-label="Adicionar ou trocar logo"
          >
            <LucideUpload class="size-3.5" />
          </Button>
          <Show when={props.edition.logo_url && !isUploading("logo_url")}>
            <div
              class="absolute bottom-0 right-0"
              onMouseEnter={() => setHovered(undefined)}
              onMouseLeave={() => setHovered("logo_url")}
            >
              <Button
                type="button"
                size="icon"
                variant="destructive"
                class="size-8 rounded-full border-2 border-background shadow-xl"
                disabled={isPending()}
                onClick={(e) => {
                  e.stopPropagation();
                  setRemoveField("logo_url");
                }}
                aria-label="Remover logo"
              >
                <LucideTrash2 class="size-3.5" />
              </Button>
            </div>
          </Show>
        </div>
      </div>

      <AlertModal
        open={Boolean(removeField())}
        onOpenChange={(open) => {
          if (!open) setRemoveField(undefined);
        }}
        title={`Remover ${removeField() === "banner_url" ? "banner" : "logo"}?`}
        description="Esta imagem será desassociada da edição."
        confirmLabel="Remover"
        variant="destructive"
        loading={isPending()}
        onConfirm={handleConfirmRemove}
      />
    </div>
  );
}
