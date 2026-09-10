import type { JSX } from "@solidjs/web";
import { createSignal, Show } from "solid-js";
import ImagePlusIcon from "~icons/lucide/image-plus";
import UploadIcon from "~icons/lucide/upload";

const ImagePlus = ImagePlusIcon as unknown as () => JSX.Element;
const Upload = UploadIcon as unknown as () => JSX.Element;

export function ProfileImageInput(props: {
  label: string;
  currentUrl?: string | null;
  variant: "banner" | "avatar";
  onSelect: (file: File) => void;
}) {
  const [preview, setPreview] = createSignal<string>();
  const [error, setError] = createSignal<string>();
  const [dragging, setDragging] = createSignal(false);
  let dragDepth = 0;
  const select = (file?: File) => {
    if (!file) return;
    if (!["image/png", "image/jpeg", "image/webp"].includes(file.type))
      return setError("Use uma imagem PNG, JPG ou WebP.");
    if (file.size > 10 * 1024 * 1024)
      return setError("A imagem deve ter no máximo 10 MB.");
    setError();
    setPreview(URL.createObjectURL(file));
    props.onSelect(file);
  };
  return (
    <div
      class={`relative ${props.variant === "banner" ? "absolute inset-0 h-full" : ""}`}
    >
      <label
        class={`group relative flex cursor-pointer overflow-hidden border border-dashed bg-background transition-colors hover:border-primary ${dragging() ? "border-primary bg-primary/10" : "border-border"} ${props.variant === "avatar" ? "size-24 items-center justify-center rounded-full border-4 border-background shadow-xl md:size-32" : "absolute inset-0 size-full items-center justify-center"}`}
        onDragEnter={(event) => {
          event.preventDefault();
          dragDepth += 1;
          setDragging(true);
        }}
        onDragOver={(event) => event.preventDefault()}
        onDragLeave={(event) => {
          event.preventDefault();
          dragDepth -= 1;
          if (dragDepth <= 0) {
            dragDepth = 0;
            setDragging(false);
          }
        }}
        onDrop={(event) => {
          event.preventDefault();
          dragDepth = 0;
          setDragging(false);
          select(event.dataTransfer?.files[0]);
        }}
      >
        <Show when={preview() || props.currentUrl}>
          <img
            src={preview() ?? props.currentUrl ?? ""}
            alt={`Prévia de ${props.label.toLowerCase()}`}
            class="absolute inset-0 size-full object-cover"
          />
        </Show>
        <span
          class={`relative z-10 flex items-center gap-2 rounded-lg bg-background/90 px-3 py-2 text-xs font-medium shadow-md backdrop-blur-sm ${props.variant === "avatar" ? "size-9 justify-center rounded-full p-0 opacity-90" : ""}`}
        >
          {props.variant === "avatar" ? <ImagePlus /> : <Upload />}
          {props.variant === "banner" && `Alterar ${props.label.toLowerCase()}`}
        </span>
        <input
          type="file"
          accept="image/png,image/jpeg,image/webp"
          class="sr-only"
          onChange={(event) => select(event.currentTarget.files?.[0])}
        />
      </label>
      <Show when={error()}>
        <p class="absolute left-2 top-full z-20 mt-1 rounded bg-destructive px-2 py-1 text-xs text-destructive-foreground">
          {error()}
        </p>
      </Show>
    </div>
  );
}
