import type { JSX } from "@solidjs/web";
import { For, createSignal } from "solid-js";

import ImagePlusIcon from "~icons/lucide/image-plus";
import TrashIcon from "~icons/lucide/trash";

import { useUploadQueue } from "@/features/upload-queue";
import { resolveStorageUrl } from "@/shared/lib/storage-url";
import { Button, Dialog } from "@trieoh/ui-solid";
import { patchVariantFn } from "../api";
import type { VariantI } from "../model";

const ImagePlus = ImagePlusIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash = TrashIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface VariantGalleryDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  variant: VariantI;
  onUpdateGallery?: (urls: string[]) => void;
}

export function VariantGalleryDialog(props: VariantGalleryDialogProps): JSX.Element {
  let fileInputRef: HTMLInputElement | undefined = undefined;
  const queue = useUploadQueue();
  const [saving, setSaving] = createSignal(false);

  const images = () => props.variant.gallery_urls ?? [];

  const isUploading = () =>
    queue.tasks.some(
      (task) =>
        task.owner.type === "variant" &&
        task.owner.id === props.variant.id &&
        task.association?.handlerKey === "variant-gallery" &&
        !["completed", "failed", "rejected"].includes(task.status),
    );

  const handleFiles = (event: Event & { currentTarget: HTMLInputElement }) => {
    const files = Array.from(event.currentTarget.files ?? []);
    event.currentTarget.value = "";
    for (const file of files) {
      void queue.enqueue({
        file,
        owner: { type: "variant", id: props.variant.id, label: props.variant.name },
        mediaType: "gallery",
        storagePath: `products/${props.variant.product_id}/variants/${props.variant.id}/gallery`,
        association: {
          handlerKey: "variant-gallery",
          input: { productId: props.variant.product_id },
        },
      });
    }
  };

  const handleRemoveImage = async (urlToRemove: string) => {
    const nextUrls = images().filter((url) => url !== urlToRemove);
    setSaving(true);
    try {
      await patchVariantFn(props.variant.id, {
        vendor_code: props.variant.vendor_code,
        name: props.variant.name,
        description: props.variant.description,
        price: props.variant.price,
        stock: props.variant.stock,
        gallery_urls: nextUrls,
      });
      props.onUpdateGallery?.(nextUrls);
    } finally {
      setSaving(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title="Galeria da variação"
      description={`${props.variant.name} (${props.variant.vendor_code}) · ${images().length} ${images().length === 1 ? "foto" : "fotos"}`}
    >
      <input
        id={`variant-${props.variant.id}-gallery-upload`}
        ref={(el) => {
          fileInputRef = el;
        }}
        type="file"
        accept="image/*"
        multiple
        class="hidden"
        onChange={handleFiles}
      />

      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-2 sm:grid-cols-3 sm:gap-3">
          <For each={images()}>
            {(url, index) => (
              <div class="group relative aspect-square overflow-hidden rounded-xl border border-border bg-muted">
                <img
                  src={resolveStorageUrl(url)}
                  alt={`Foto ${index() + 1} de ${props.variant.name}`}
                  class="size-full object-cover"
                />
                <button
                  type="button"
                  aria-label={`Remover foto ${index() + 1}`}
                  disabled={saving()}
                  onClick={() => handleRemoveImage(url)}
                  class="absolute top-2 right-2 flex size-7 items-center justify-center rounded-full bg-black/60 text-white opacity-0 transition-opacity hover:bg-destructive group-hover:opacity-100 disabled:opacity-50 cursor-pointer"
                >
                  <Trash class="size-3.5" />
                </button>
              </div>
            )}
          </For>

          <button
            type="button"
            disabled={isUploading()}
            onClick={() => fileInputRef?.click()}
            class="flex aspect-square flex-col items-center justify-center gap-1.5 rounded-xl border border-dashed border-border bg-muted/30 text-muted-foreground transition-colors hover:border-primary hover:bg-muted/60 hover:text-foreground disabled:opacity-50 cursor-pointer"
          >
            <ImagePlus class="size-5" />
            <span class="text-xs font-medium">
              {isUploading() ? "Enviando..." : "Adicionar fotos"}
            </span>
          </button>
        </div>

        <div class="flex justify-end pt-2">
          <Button
            type="button"
            variant="outline"
            onClick={() => props.onOpenChange(false)}
          >
            Fechar
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
