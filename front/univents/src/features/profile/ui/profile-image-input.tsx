import { ImagePlus, Upload } from "lucide-react";
import { useEffect, useId, useState } from "react";
import { cn } from "@/shared/lib/utils";

const ALLOWED_TYPES = ["image/png", "image/jpeg", "image/webp"];
const MAX_SIZE = 10 * 1024 * 1024;

interface ProfileImageInputProps {
  label: string;
  currentUrl?: string | null;
  file?: File;
  onSelect: (file: File) => void;
  variant: "banner" | "avatar";
  className?: string;
}

export function ProfileImageInput({
  label,
  currentUrl,
  file,
  onSelect,
  variant,
  className,
}: ProfileImageInputProps) {
  const id = useId();
  const [dragging, setDragging] = useState(false);
  const [error, setError] = useState<string>();
  const [previewUrl, setPreviewUrl] = useState<string>();

  useEffect(() => {
    if (!file) {
      setPreviewUrl(undefined);
      return;
    }
    const url = URL.createObjectURL(file);
    setPreviewUrl(url);
    return () => URL.revokeObjectURL(url);
  }, [file]);

  const select = (next?: File) => {
    if (!next) return;
    if (!ALLOWED_TYPES.includes(next.type)) {
      setError("Use uma imagem PNG, JPG ou WebP.");
      return;
    }
    if (next.size > MAX_SIZE) {
      setError("A imagem deve ter no máximo 10 MB.");
      return;
    }
    setError(undefined);
    onSelect(next);
  };

  const imageUrl = previewUrl ?? currentUrl ?? undefined;
  const isAvatar = variant === "avatar";

  return (
    <div className={cn("relative", className)}>
      <label
        htmlFor={id}
        className={cn(
          "group relative flex cursor-pointer overflow-hidden border border-dashed",
          "transition-colors focus-within:ring-2 focus-within:ring-ring",
          dragging ? "border-primary bg-primary/10" : "border-border",
          isAvatar
            ? "size-24 items-center justify-center rounded-full border-4 border-background shadow-xl md:size-32"
            : "h-full w-full items-center justify-center",
        )}
        onDragEnter={(event) => {
          event.preventDefault();
          setDragging(true);
        }}
        onDragOver={(event) => event.preventDefault()}
        onDragLeave={() => setDragging(false)}
        onDrop={(event) => {
          event.preventDefault();
          setDragging(false);
          select(event.dataTransfer.files[0]);
        }}
      >
        {imageUrl && (
          <img
            src={imageUrl}
            alt={`Prévia de ${label.toLowerCase()}`}
            className="absolute inset-0 size-full object-cover"
          />
        )}
        <span
          className={cn(
            "relative z-10 flex items-center gap-2 rounded-lg bg-background/90",
            "px-3 py-2 text-xs font-medium shadow-md backdrop-blur-sm",
            isAvatar &&
              "size-9 justify-center rounded-full p-0 opacity-90 md:opacity-0 md:group-hover:opacity-100",
          )}
        >
          {isAvatar ? (
            <ImagePlus className="size-4" />
          ) : (
            <Upload className="size-4" />
          )}
          {!isAvatar &&
            (file ? "Imagem selecionada" : `Alterar ${label.toLowerCase()}`)}
        </span>
        <input
          id={id}
          type="file"
          accept={ALLOWED_TYPES.join(",")}
          className="sr-only"
          onChange={(event) => select(event.target.files?.[0])}
        />
      </label>
      {error && (
        <p className="absolute left-2 top-full z-20 mt-1 rounded bg-destructive px-2 py-1 text-xs text-destructive-foreground shadow">
          {error}
        </p>
      )}
    </div>
  );
}
