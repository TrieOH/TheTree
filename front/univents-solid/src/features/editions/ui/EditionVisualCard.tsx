import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal } from "solid-js";
import AlertTriangleIcon from "~icons/lucide/alert-triangle";
import Loader2Icon from "~icons/lucide/loader-2";
import Trash2Icon from "~icons/lucide/trash-2";
import UploadIcon from "~icons/lucide/upload";
import { Button, cn } from "@trieoh/ui-solid";
import { AlertModal } from "@/widgets/ui/AlertModal";
import { useUploadQueue } from "@/features/upload-queue";
import { usePatchEditionMutation } from "../api/mutations";
import type { EditionI } from "../model";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideAlertTriangle = AlertTriangleIcon as unknown as IconComp;
const LucideLoader2 = Loader2Icon as unknown as IconComp;
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

  const getTask = (field: ImageField) =>
    uploadQueue.tasks.find(
      (task) =>
        task.owner.type === "edition" &&
        task.owner.id === props.edition.id &&
        task.association?.handlerKey === "edition-image" &&
        task.association.input?.field === field &&
        !["completed", "failed", "rejected"].includes(task.status),
    );

  const isUploading = (field: ImageField) => Boolean(getTask(field));

  const currentBannerUrl = () =>
    getTask("banner_url")?.uploadedUrl ?? props.edition.banner_url;

  const currentLogoUrl = () =>
    getTask("logo_url")?.uploadedUrl ?? props.edition.logo_url;

  const upload = async (field: ImageField, file?: File) => {
    if (!file || isUploading(field) || isPending()) return;
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

  const handleDrop = (field: ImageField, e: DragEvent) => {
    e.preventDefault();
    setDragging(undefined);
    if (isUploading(field) || isPending()) return;
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
          if (!isUploading("banner_url")) setDragging("banner_url");
        }}
        onDragOver={(e) => e.preventDefault()}
        onDragLeave={() => setDragging(undefined)}
        onDrop={(e) => handleDrop("banner_url", e)}
        class={cn(
          "group relative flex h-56 items-center justify-center overflow-hidden rounded-xl border border-dashed transition-colors",
          isUploading("banner_url")
            ? "cursor-not-allowed border-border/80 bg-muted/40"
            : "cursor-pointer bg-muted/20",
          !isUploading("banner_url") && dragging() === "banner_url"
            ? "border-primary bg-primary/10"
            : !isUploading("banner_url") && hovered() === "banner_url"
              ? "border-primary"
              : "border-border/60",
        )}
        onClick={() => {
          if (!isUploading("banner_url") && !isPending()) {
            bannerInput?.click();
          }
        }}
        onMouseEnter={() => {
          if (!isUploading("banner_url")) setHovered("banner_url");
        }}
        onMouseLeave={() => setHovered(undefined)}
      >
        <Show when={currentBannerUrl()}>
          {(url) => (
            <img
              src={url()}
              alt="Banner da edição"
              class="size-full rounded-xl object-cover"
            />
          )}
        </Show>

        <Show when={!isUploading("banner_url")}>
          <div
            class={cn(
              "pointer-events-none absolute inset-0 rounded-xl bg-primary/10 transition-opacity",
              hovered() === "banner_url" ? "opacity-100" : "opacity-0",
            )}
          />
        </Show>

        <Show
          when={isUploading("banner_url")}
          fallback={
            <div class="absolute inset-0">
              <Button
                type="button"
                size="icon"
                variant="outline"
                class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 border-border bg-background/90 shadow-sm"
                disabled={isPending()}
                onClick={(e) => {
                  e.stopPropagation();
                  bannerInput?.click();
                }}
                aria-label="Adicionar ou trocar banner"
              >
                <LucideUpload class="size-4" />
              </Button>
              <Show when={props.edition.banner_url}>
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
          }
        >
          <div class="absolute inset-0 z-10 flex flex-col items-center justify-center gap-2 bg-background/80 p-4 text-center backdrop-blur-xs">
            <Show
              when={getTask("banner_url")?.status === "paused"}
              fallback={
                <>
                  <LucideLoader2 class="size-6 animate-spin text-primary" />
                  <div class="space-y-0.5">
                    <p class="text-xs font-semibold text-foreground">
                      Enviando banner...
                    </p>
                    <p class="text-[11px] text-muted-foreground">
                      Aguarde a conclusão do upload para realizar novas alterações.
                    </p>
                  </div>
                </>
              }
            >
              <LucideAlertTriangle class="size-6 text-amber-500" />
              <div class="space-y-0.5">
                <p class="text-xs font-semibold text-amber-600 dark:text-amber-400">
                  Upload pausado
                </p>
                <p class="text-[11px] text-muted-foreground">
                  Acesse a fila de uploads para retomar ou remover.
                </p>
              </div>
            </Show>
          </div>
        </Show>
      </div>

      {/* Logo Dropzone */}
      <div
        onDragEnter={(e) => {
          e.preventDefault();
          if (!isUploading("logo_url")) setDragging("logo_url");
        }}
        onDragOver={(e) => e.preventDefault()}
        onDragLeave={() => setDragging(undefined)}
        onDrop={(e) => handleDrop("logo_url", e)}
        class={cn(
          "group absolute -bottom-8 left-5 z-30 flex size-24 items-center justify-center rounded-full border-4 border-card bg-muted shadow-xl transition-all md:size-28",
          isUploading("logo_url")
            ? "cursor-not-allowed opacity-90"
            : "cursor-pointer",
          !isUploading("logo_url") && dragging() === "logo_url"
            ? "ring-2 ring-primary/70"
            : !isUploading("logo_url") && hovered() === "logo_url"
              ? "ring-4 ring-primary/50"
              : "",
        )}
        onClick={() => {
          if (!isUploading("logo_url") && !isPending()) {
            logoInput?.click();
          }
        }}
        onMouseEnter={() => {
          if (!isUploading("logo_url")) setHovered("logo_url");
        }}
        onMouseLeave={() => setHovered(undefined)}
      >
        <div class="absolute inset-0 overflow-hidden rounded-full">
          <Show when={currentLogoUrl()}>
            {(url) => (
              <img
                src={url()}
                alt="Logo da edição"
                class="size-full object-cover"
              />
            )}
          </Show>

          <Show when={!isUploading("logo_url")}>
            <div
              class={cn(
                "pointer-events-none absolute inset-0 bg-primary/15 transition-opacity",
                hovered() === "logo_url" ? "opacity-100" : "opacity-0",
              )}
            />
          </Show>

          <Show when={isUploading("logo_url")}>
            <div class="absolute inset-0 z-10 flex flex-col items-center justify-center bg-background/85 p-1 text-center backdrop-blur-xs">
              <Show
                when={getTask("logo_url")?.status === "paused"}
                fallback={
                  <>
                    <LucideLoader2 class="size-5 animate-spin text-primary" />
                    <span class="mt-1 text-[9px] font-semibold text-foreground">Enviando...</span>
                  </>
                }
              >
                <LucideAlertTriangle class="size-5 text-amber-500" />
                <span class="mt-1 text-[9px] font-semibold text-amber-600">Pausado</span>
              </Show>
            </div>
          </Show>
        </div>

        <Show when={!isUploading("logo_url")}>
          <div class="pointer-events-none absolute inset-0">
            <Button
              type="button"
              size="icon"
              variant="outline"
              class="pointer-events-auto absolute left-1/2 top-1/2 size-8 -translate-x-1/2 -translate-y-1/2 rounded-full border-border bg-background/90 shadow-sm"
              disabled={isPending()}
              onClick={(e) => {
                e.stopPropagation();
                logoInput?.click();
              }}
              aria-label="Adicionar ou trocar logo"
            >
              <LucideUpload class="size-3.5" />
            </Button>
            <Show when={props.edition.logo_url}>
              <div
                class="pointer-events-auto absolute -bottom-1 -right-1 z-10"
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
        </Show>
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
