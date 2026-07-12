import { useRef, useState } from "react";
import type { FieldValues } from "react-hook-form";
import { Camera, UploadCloud } from "lucide-react";
import type { FieldConfig, FieldFormApi } from "../../model/types";
import { useImageUploadField } from "../../hooks/use-image-upload-field";
import { ImageThumb } from "./helper/image-thumb";
import { Button } from "@/shared/ui/shadcn/button";
import { Label } from "@/shared/ui/shadcn/label";
import { cn } from "@/shared/lib/utils";

export interface GalleryFieldRendererProps<TFieldValues extends FieldValues> {
  field: FieldConfig<TFieldValues>;
  form: FieldFormApi<TFieldValues>;
}

export function GalleryFieldRenderer<TFieldValues extends FieldValues>({
  field,
  form,
}: GalleryFieldRendererProps<TFieldValues>) {
  if (field.kind !== "gallery") return null;

  const inputRef = useRef<HTMLInputElement>(null);
  const dragDepthRef = useRef(0);
  const [isDragging, setIsDragging] = useState(false);

  const [initialUrls] = useState<string[]>(() => {
    const current = form.getValues(field.name);
    return Array.isArray(current) ? current.filter((url: unknown): url is string => typeof url === "string") : [];
  });

  const { items, addFiles, removeItem, canAddMore } = useImageUploadField({
    fieldKey: String(field.name),
    initialUrls,
    maxItems: field.maxItems ?? Number.POSITIVE_INFINITY,
    accept: field.accept,
    maxSizeMB: field.maxSizeMB,
    onValueChange: (urls) => {
      form.setValue(field.name, urls as never, { shouldDirty: true });
    },
    onTrackingChange: field.onTrackingChange,
  });

  return (
    <div className="space-y-1.5">
      <Label className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
        {field.label}
        {field.optional ? (
          <span className="ml-1 font-normal normal-case text-muted-foreground/70">(opcional)</span>
        ) : null}
      </Label>

      <div
        onDragEnter={(event) => {
          event.preventDefault();
          dragDepthRef.current += 1;
          setIsDragging(true);
        }}
        onDragOver={(event) => {
          event.preventDefault();
          setIsDragging(true);
        }}
        onDragLeave={(event) => {
          event.preventDefault();
          dragDepthRef.current = Math.max(dragDepthRef.current - 1, 0);
          if (dragDepthRef.current === 0) setIsDragging(false);
        }}
        onDrop={(event) => {
          event.preventDefault();
          dragDepthRef.current = 0;
          setIsDragging(false);
          if (event.dataTransfer.files.length > 0) addFiles(event.dataTransfer.files);
        }}
        className={cn(
          "rounded-xl border border-dashed p-4 transition-colors",
          isDragging ? "border-primary bg-primary/5" : "border-border bg-muted/10",
        )}
      >
        <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
          <div className="space-y-1">
            <p className="text-sm font-medium text-foreground">
              Arraste e solte imagens aqui
            </p>
            <p className="text-xs text-muted-foreground">
              ou use o botão para selecionar várias imagens
            </p>
          </div>

          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => inputRef.current?.click()}
            className="w-full gap-2 sm:w-auto"
          >
            <Camera className="h-4 w-4" />
            Adicionar imagens
          </Button>
        </div>

        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
          {items.map((item) => (
            <ImageThumb key={item.id} item={item} onRemove={() => removeItem(item.id)} className="aspect-square h-auto w-full" />
          ))}

          {canAddMore ? (
            <Button
              type="button"
              variant="outline"
              size="none"
              onClick={() => inputRef.current?.click()}
              className="flex aspect-square w-full flex-col items-center justify-center gap-1 rounded-md border-dashed px-0 text-[10px] font-semibold uppercase tracking-wide text-muted-foreground transition-colors hover:border-primary hover:bg-primary/5 hover:text-foreground"
            >
              <UploadCloud className="h-4 w-4" />
              Adicionar
            </Button>
          ) : null}
        </div>

        <input
          ref={inputRef}
          type="file"
          accept={field.accept}
          multiple
          className="hidden"
          onChange={(event) => {
            if (event.target.files && event.target.files.length > 0) {
              addFiles(event.target.files);
            }
            event.target.value = "";
          }}
        />
      </div>

      {field.hint ? <p className="text-xs text-muted-foreground">{field.hint}</p> : null}
    </div>
  );
}
