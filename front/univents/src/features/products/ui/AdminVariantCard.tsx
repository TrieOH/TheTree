import { MoreVertical, Pencil, Trash2 } from "lucide-react";
import { motion } from "motion/react";
import { Button } from "@/shared/ui/shadcn/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/shared/ui/shadcn/dropdown-menu";
import type { VariantI } from "../model";
import { VariantGalleryActions } from "./VariantGalleryActions";

export function AdminVariantCard({
  variant,
  index,
  onEdit,
  onDelete,
}: {
  variant: VariantI;
  index: number;
  onEdit: () => void;
  onDelete: () => void;
}) {
  const image = variant.gallery_urls?.[0];
  return (
    <motion.article
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.04, duration: 0.3 }}
      className="group flex min-w-0 flex-col overflow-hidden rounded-2xl bg-card text-left shadow-xs ring-1 ring-foreground/10 transition-all hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm"
    >
      <div className="relative aspect-video overflow-hidden bg-muted">
        {image ? (
          <img
            src={image}
            alt={variant.name}
            className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full items-center justify-center bg-linear-to-br from-muted via-background to-muted/40">
            <span className="text-3xl font-semibold text-muted-foreground/30">
              {variant.name.charAt(0).toUpperCase()}
            </span>
          </div>
        )}
        <div className="absolute inset-0 bg-linear-to-t from-black/70 via-transparent to-transparent" />
        <div className="absolute left-3 top-3 rounded-full border border-white/30 bg-black/75 px-2.5 py-1 text-[11px] font-semibold text-white shadow-sm">
          {variant.gallery_urls?.length ?? 0}{" "}
          {(variant.gallery_urls?.length ?? 0) === 1 ? "imagem" : "imagens"}
        </div>
        <div className="absolute right-3 top-3">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <button
                  type="button"
                  className="inline-flex size-8 items-center justify-center rounded-full bg-background/90 shadow-sm"
                  aria-label="Abrir ações da variante"
                >
                  <MoreVertical className="size-4" />
                </button>
              }
            />
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={onEdit}>
                <Pencil className="size-4" />
                Editar
              </DropdownMenuItem>
              <DropdownMenuItem variant="destructive" onClick={onDelete}>
                <Trash2 className="size-4" />
                Excluir
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
      <div className="flex min-w-0 flex-1 flex-col gap-3 p-4">
        <div className="min-w-0">
          <h3 className="truncate text-base font-semibold text-foreground">
            {variant.name}
          </h3>
          <p className="truncate text-xs text-muted-foreground">
            Código: {variant.vendor_code}
          </p>
        </div>
        {variant.description ? (
          <p className="line-clamp-2 min-h-10 text-sm text-muted-foreground">
            {variant.description}
          </p>
        ) : (
          <div className="min-h-10" />
        )}
        <div className="flex items-end justify-between gap-3">
          <div>
            <p className="text-sm font-bold text-foreground">
              R$ {(variant.price / 100).toFixed(2)}
            </p>
            <p className="text-xs text-muted-foreground">
              {variant.stock === null
                ? "Ilimitado"
                : `${variant.stock} em estoque`}
            </p>
          </div>
          <VariantGalleryActions variant={variant} />
        </div>
        <Button
          type="button"
          variant="outline"
          className="w-full"
          onClick={onEdit}
        >
          Editar variante
        </Button>
      </div>
    </motion.article>
  );
}
