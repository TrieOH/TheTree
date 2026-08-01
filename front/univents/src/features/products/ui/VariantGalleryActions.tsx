import { ImagePlus, Trash2 } from "lucide-react";
import { useEffect, useRef } from "react";
import { toast } from "sonner";
import { useUploadQueue } from "@/features/upload-queue";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/shared/ui/shadcn/dialog";
import { patchVariantFn } from "../api";
import type { VariantI } from "../model";

export function VariantGalleryActions({ variant }: { variant: VariantI }) {
  const inputRef = useRef<HTMLInputElement>(null);
  const { enqueue, tasks } = useUploadQueue();
  const images = variant.gallery_urls ?? [];
  const active = tasks.some(
    (task) =>
      task.owner.type === "variant" &&
      task.owner.id === variant.id &&
      task.association?.handlerKey === "variant-gallery" &&
      !["completed", "failed", "rejected"].includes(task.status),
  );
  useEffect(() => {
    const id = `variant-gallery-upload-${variant.id}`;
    if (active)
      toast.loading("Enviando imagens da variante…", {
        id,
        duration: Infinity,
      });
    else toast.dismiss(id);
    return () => {
      toast.dismiss(id);
    };
  }, [active, variant.id]);
  const patch = (gallery_urls: string[]) =>
    void patchVariantFn(variant.id, {
      vendor_code: variant.vendor_code,
      name: variant.name,
      description: variant.description,
      price: variant.price,
      stock: variant.stock,
      gallery_urls,
    });
  const addFiles = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = [...(event.target.files ?? [])];
    event.target.value = "";
    for (const file of files)
      void enqueue({
        file,
        owner: { type: "variant", id: variant.id, label: variant.name },
        mediaType: "gallery",
        storagePath: `products/${variant.product_id}/variants/${variant.id}/gallery`,
        association: {
          handlerKey: "variant-gallery",
          input: { productId: variant.product_id },
        },
      });
  };
  return (
    <div onClick={(event) => event.stopPropagation()}>
      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        multiple
        className="hidden"
        onChange={addFiles}
      />
      <Dialog>
        <DialogTrigger
          render={
            <Button type="button" size="sm" variant="outline" disabled={active}>
              <ImagePlus className="size-4" />
              Galeria ({images.length})
            </Button>
          }
        />
        <DialogContent className="w-[calc(100vw-1rem)] max-w-2xl overflow-hidden p-4 sm:p-6">
          <DialogHeader>
            <DialogTitle>Galeria da variante</DialogTitle>
            <DialogDescription>
              {variant.name} · {images.length}{" "}
              {images.length === 1 ? "imagem" : "imagens"}
            </DialogDescription>
          </DialogHeader>
          <div className="grid min-w-0 grid-cols-2 gap-2 sm:grid-cols-4 sm:gap-3">
            {images.map((url, index) => (
              <div
                key={url}
                className="group relative aspect-square overflow-hidden rounded-xl border bg-muted"
              >
                <img
                  src={url}
                  alt={`Imagem ${index + 1}`}
                  className="h-full w-full object-cover"
                />
                <button
                  type="button"
                  title={`Remover imagem ${index + 1}`}
                  onClick={() => patch(images.filter((item) => item !== url))}
                  className="absolute right-2 top-2 rounded-full bg-black/75 p-2 text-white opacity-100 transition-opacity sm:opacity-0 sm:group-hover:opacity-100"
                >
                  <Trash2 className="size-4" />
                </button>
              </div>
            ))}
            <button
              type="button"
              disabled={active}
              onClick={() => inputRef.current?.click()}
              className="flex aspect-square flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-border text-muted-foreground hover:bg-muted"
            >
              <ImagePlus className="size-6" />
              <span className="text-xs font-medium">Adicionar imagens</span>
            </button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
