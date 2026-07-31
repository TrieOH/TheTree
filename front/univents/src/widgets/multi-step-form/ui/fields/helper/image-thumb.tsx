import { Loader2, X } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import type { ImageItem } from "@/widgets/multi-step-form/model/types";

export function ImageThumb({
  item,
  onRemove,
  className = "aspect-square h-auto w-full",
}: {
  item: ImageItem;
  onRemove: () => void;
  className?: string;
}) {
  const isProcessing = item.status === "processing";
  const hasError = item.status === "error";

  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-md border bg-muted",
        className,
      )}
    >
      <img src={item.url} alt="" className="h-full w-full object-cover" />

      {isProcessing ? (
        <div className="absolute inset-0 flex flex-col items-center justify-center gap-1 bg-background/85 text-[10px] font-medium">
          <Loader2 className="h-4 w-4 animate-spin" />
          Processando...
        </div>
      ) : null}

      {hasError ? (
        <div className="absolute inset-0 flex items-center justify-center bg-destructive/10 p-1 text-center text-[9px] font-medium text-destructive">
          {item.errorMessage ?? "Erro"}
        </div>
      ) : null}

      <Button
        type="button"
        onClick={onRemove}
        aria-label="Remover imagem"
        variant="secondary"
        size="icon-xs"
        className="absolute right-2 top-2 rounded-full shadow-sm"
      >
        <X className="h-3 w-3" />
      </Button>
    </div>
  );
}
