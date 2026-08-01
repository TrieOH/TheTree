import {
  Layers,
  MoreVertical,
  Package,
  Pencil,
  ShieldCheck,
  Tag,
  Trash2,
} from "lucide-react";
import { motion } from "motion/react";
import type React from "react";
import type { ProductI } from "@/features/products/model";
import { cn } from "@/shared/lib/utils";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from "@/shared/ui/shadcn/context-menu";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/ui/shadcn/dropdown-menu";

interface AdminProductCardProps {
  product: ProductI;
  index: number;
  onEdit: (product: ProductI) => void;
  onDelete: (product: ProductI) => void;
  onManageVariants: (product: ProductI) => void;
}

function MenuItems({
  isContext = false,
  onEdit,
  onDelete,
  onManageVariants,
}: {
  isContext?: boolean;
  onEdit: () => void;
  onDelete: () => void;
  onManageVariants: () => void;
}) {
  const Item = isContext ? ContextMenuItem : DropdownMenuItem;
  const Separator = isContext ? ContextMenuSeparator : DropdownMenuSeparator;
  const stop =
    (action: () => void) => (e: React.MouseEvent | React.KeyboardEvent) => {
      e.preventDefault();
      e.stopPropagation();
      action();
    };

  return (
    <>
      <Item onClick={stop(onEdit)}>
        <Pencil className="size-4" />
        <span>Editar</span>
      </Item>
      <Item onClick={stop(onManageVariants)}>
        <Layers className="size-4" />
        <span>Variações</span>
      </Item>
      <Separator />
      <Item onClick={stop(onDelete)}>
        <Trash2 className="size-4" />
        <span>Excluir</span>
      </Item>
    </>
  );
}

function getRandomGradient(index: number): string {
  const gradients = [
    "from-violet-500/20 via-fuchsia-500/10 to-background",
    "from-emerald-500/20 via-teal-500/10 to-background",
    "from-amber-500/20 via-orange-500/10 to-background",
    "from-blue-500/20 via-cyan-500/10 to-background",
    "from-rose-500/20 via-pink-500/10 to-background",
    "from-indigo-500/20 via-purple-500/10 to-background",
  ];
  return gradients[index % gradients.length];
}

export function AdminProductCard({
  product,
  index,
  onEdit,
  onDelete,
  onManageVariants,
}: AdminProductCardProps) {
  const handleAction = (type: "edit" | "delete" | "variants") => {
    if (type === "edit") onEdit(product);
    if (type === "delete") onDelete(product);
    if (type === "variants") onManageVariants(product);
  };

  const isDeleted = product.deleted_at !== null;

  const article = (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{
        delay: index * 0.05,
        duration: 0.35,
        ease: [0.25, 0.1, 0.25, 1],
      }}
      className={cn(
        "group relative flex w-full min-w-0 flex-col overflow-hidden rounded-2xl bg-card text-left",
        "ring-1 ring-foreground/10 shadow-xs",
        "transform-gpu will-change-transform",
        "transition-all duration-300 ease-out",
        "hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm",
        "focus:outline-none focus-visible:outline-none focus-visible:ring-0",
        isDeleted && "opacity-60 grayscale",
      )}
      role="button"
      tabIndex={0}
      onClick={() => handleAction("variants")}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          handleAction("variants");
        }
      }}
    >
      <div
        className={cn(
          "relative aspect-video overflow-hidden",
          "bg-linear-to-br",
          getRandomGradient(index),
        )}
      >
        <div className="flex h-full w-full items-center justify-center">
          <div className="flex size-18 items-center justify-center rounded-full border border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
            <Tag className="size-7 text-muted-foreground/40" />
          </div>
        </div>

        <div className="absolute inset-0 bg-linear-to-t from-background/90 via-background/35 to-transparent" />

        <div className="absolute left-3 top-3 flex flex-wrap items-center gap-1.5">
          <span
            className={cn(
              "inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-medium backdrop-blur-sm",
              isDeleted
                ? "border-rose-500/20 bg-rose-500/10 text-rose-700"
                : "border-emerald-500/20 bg-emerald-500/10 text-emerald-700",
            )}
          >
            <span
              className={cn(
                "size-1.5 rounded-full",
                isDeleted ? "bg-rose-500" : "bg-emerald-500",
              )}
            />
            <span className="max-w-28 truncate">
              {isDeleted ? "Excluído" : "Ativo"}
            </span>
          </span>
        </div>

        <div className="absolute right-3 top-3">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <button
                  type="button"
                  onClick={(e) => e.stopPropagation()}
                  className={cn(
                    "inline-flex size-9 items-center justify-center rounded-full",
                    "bg-background/85 text-foreground shadow-sm backdrop-blur-sm",
                    "transition-colors hover:bg-background",
                  )}
                  aria-label={`Abrir ações de ${product.vendor_code}`}
                >
                  <MoreVertical className="size-4" />
                </button>
              }
            />
            <DropdownMenuContent align="end" className="w-56">
              <MenuItems
                onEdit={() => handleAction("edit")}
                onDelete={() => handleAction("delete")}
                onManageVariants={() => handleAction("variants")}
              />
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div className="absolute inset-x-0 bottom-0 flex items-end justify-between gap-3 p-4 sm:p-5">
          <div className="min-w-0 space-y-1">
            <h3 className="line-clamp-1 text-balance text-lg font-semibold leading-snug text-foreground transition-colors duration-300 group-hover:text-primary sm:text-xl">
              {product.vendor_code}
            </h3>
          </div>
        </div>
      </div>

      <div className="flex items-center justify-between gap-3 p-4 pt-3 sm:p-5 sm:pt-4">
        <div className="min-w-0 flex-1 space-y-1.5">
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
            <Package className="size-3.5 shrink-0" />
            <span className="truncate font-mono text-xs">
              {product.vendor_code}
            </span>
          </div>

          {product.requires_registration && (
            <div className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
              <ShieldCheck className="size-3.5 shrink-0" />
              <span className="truncate">Exige cadastro</span>
            </div>
          )}
        </div>
      </div>
    </motion.article>
  );

  return (
    <ContextMenu>
      <ContextMenuTrigger render={article} />
      <ContextMenuContent className="w-56">
        <MenuItems
          isContext
          onEdit={() => handleAction("edit")}
          onDelete={() => handleAction("delete")}
          onManageVariants={() => handleAction("variants")}
        />
      </ContextMenuContent>
    </ContextMenu>
  );
}
