import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal } from "solid-js";
import Trash2Icon from "~icons/lucide/trash-2";
import UploadIcon from "~icons/lucide/upload";
import { Button, cn } from "@trieoh/ui-solid";
import { AlertModal } from "@/widgets/ui/AlertModal";
import { useUploadQueue } from "@/features/upload-queue";
import { usePatchEventMutation } from "../api/mutations";
import type { EventI } from "../model";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideTrash2 = Trash2Icon as unknown as IconComp;
const LucideUpload = UploadIcon as unknown as IconComp;

type ImageField = "logo_url" | "banner_url";

export function EventVisualCard(props: { event: EventI }): JSX.Element {
  const [dragging, setDragging] = createSignal<ImageField | undefined>(
    undefined,
  );
  const [hovered, setHovered] = createSignal<ImageField | undefined>(undefined);
  const [removeField, setRemoveField] = createSignal<ImageField | undefined>(
    undefined,
  );

  let bannerInput: HTMLInputElement | undefined;
  let logoInput: HTMLInputElement | undefined;

  const uploadQueue = useUploadQueue();
  const patchMutation = usePatchEventMutation();

  const isPending = () => patchMutation.result().status === "pending";

  const upload = async (field: ImageField, file?: File) => {
    if (!file) return;
    const label = field === "logo_url" ? "logo" : "banner";
    try {
      await uploadQueue.enqueue({
        file,
        owner: {
          type: "event",
          id: props.event.id,
          label: props.event.full_name,
        },
        mediaType: label,
        label: `${props.event.full_name} — ${label}`,
        storagePath: `events/${props.event.id}/${label}`,
        correctionPath: `/admin/events/${props.event.id}`,
        association: {
          handlerKey: "event-image",
          input: { field },
        },
      });
    } catch {
      // Handled silently
    }
  };

  const remove = (field: ImageField) =>
    patchMutation.mutateAsync({
      eventId: props.event.id,
      data: {
        full_name: props.event.full_name,
        slug: props.event.slug,
        acronym: props.event.acronym,
        description: props.event.description,
        contact_email: props.event.contact_email,
        logo_url: field === "logo_url" ? null : props.event.logo_url,
        banner_url: field === "banner_url" ? null : props.event.banner_url,
      },
    });

  const isUploading = (field: ImageField) =>
    uploadQueue.tasks.some(
      (task) =>
        task.owner.type === "event" &&
        task.owner.id === props.event.id &&
        task.association?.handlerKey === "event-image" &&
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
    () => props.event.id,
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
        id="event-banner-upload"
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
        id="event-logo-upload"
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

      {/* Banner */}
      <div
        onDragEnter={(e) => {
          e.preventDefault();
          setDragging("banner_url");
        }}
        onDragOver={(e) => e.preventDefault()}
        onDragLeave={() => setDragging(undefined)}
        onDrop={(e) => handleDrop("banner_url", e)}
        onClick={() => {
          if (!isUploading("banner_url") && !isPending()) {
            bannerInput?.click();
          }
        }}
        onMouseEnter={() => setHovered("banner_url")}
        onMouseLeave={() => setHovered(undefined)}
        class={cn(
          "group relative flex h-56 cursor-pointer items-center justify-center rounded-md border border-dashed bg-muted/20 transition-colors",
          dragging() === "banner_url"
            ? "border-primary bg-primary/10"
            : hovered() === "banner_url"
              ? "border-primary"
              : "border-border/60",
        )}
      >
        <Show when={props.event.banner_url && !isUploading("banner_url")}>
          <img
            src={props.event.banner_url ?? undefined}
            alt="Banner do evento"
            class="size-full rounded-md object-cover"
          />
        </Show>
        <div
          class={cn(
            "pointer-events-none absolute inset-0 rounded-md bg-primary/10 transition-opacity",
            hovered() === "banner_url" ? "opacity-100" : "opacity-0",
          )}
        />

        <div class="absolute inset-0">
          <Button
            type="button"
            size="icon"
            variant="outline"
            class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 border-border bg-background/90 shadow-sm active:-translate-y-1/2!"
            disabled={isUploading("banner_url") || isPending()}
            onClick={(e) => {
              e.stopPropagation();
              bannerInput?.click();
            }}
            aria-label="Adicionar ou trocar banner"
          >
            <LucideUpload class="size-4" />
          </Button>

          <Show when={props.event.banner_url && !isUploading("banner_url")}>
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

      {/* Logo */}
      <div
        onDragEnter={(e) => {
          e.preventDefault();
          setDragging("logo_url");
        }}
        onDragOver={(e) => e.preventDefault()}
        onDragLeave={() => setDragging(undefined)}
        onDrop={(e) => handleDrop("logo_url", e)}
        onClick={() => {
          if (!isUploading("logo_url") && !isPending()) {
            logoInput?.click();
          }
        }}
        onMouseEnter={() => setHovered("logo_url")}
        onMouseLeave={() => setHovered(undefined)}
        class={cn(
          "group absolute -bottom-8 left-5 z-10 flex size-24 cursor-pointer items-center justify-center rounded-full border-4 border-card bg-muted shadow-xl transition-all md:size-28",
          dragging() === "logo_url"
            ? "ring-2 ring-primary/70"
            : hovered() === "logo_url"
              ? "ring-4 ring-primary/50"
              : "",
        )}
      >
        <Show when={props.event.logo_url && !isUploading("logo_url")}>
          <img
            src={props.event.logo_url ?? undefined}
            alt="Logo do evento"
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
            class="absolute left-1/2 top-1/2 size-8 -translate-x-1/2 -translate-y-1/2 rounded-full border-border bg-background/90 shadow-sm active:-translate-y-1/2!"
            disabled={isUploading("logo_url") || isPending()}
            onClick={(e) => {
              e.stopPropagation();
              logoInput?.click();
            }}
            aria-label="Adicionar ou trocar logo"
          >
            <LucideUpload class="size-3.5" />
          </Button>

          <Show when={props.event.logo_url && !isUploading("logo_url")}>
            <div
              class="absolute bottom-0 right-0"
              onMouseEnter={() => setHovered(undefined)}
              onMouseLeave={() => setHovered("logo_url")}
            >
              <Button
                type="button"
                size="icon"
                variant="destructive"
                class="size-8 rounded-full border-2 border-background shadow-xl hover:bg-destructive/90"
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
        description="Esta imagem será desassociada do evento."
        confirmLabel="Remover"
        variant="destructive"
        loading={isPending()}
        onConfirm={handleConfirmRemove}
      />
    </div>
  );
}
