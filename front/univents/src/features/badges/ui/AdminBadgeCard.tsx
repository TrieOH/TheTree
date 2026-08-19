import { Copy, MoreVertical, Pencil, Trash2 } from "lucide-react";
import { motion } from "motion/react";
import { useState } from "react";
import { cn } from "@/shared/lib/utils";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/shared/ui/shadcn/alert-dialog";
import { CardContent, CardHeader, CardTitle } from "@/shared/ui/shadcn/card";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/shared/ui/shadcn/context-menu";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/shared/ui/shadcn/dropdown-menu";
import type { BadgePrintItem, BadgeTemplate } from "../model";
import { BadgePreview } from "./badge-preview";

interface AdminBadgeCardProps {
  item: BadgeTemplate | BadgePrintItem;
  kind: "template" | "emission";
  index?: number;
  ticketName?: string;
  onEdit?: () => void;
  onDelete?: () => void;
  onDuplicate?: () => void;
  onView?: () => void;
}

export default function AdminBadgeCard({
  item,
  kind,
  index = 0,
  ticketName,
  onEdit,
  onDelete,
  onDuplicate,
  onView,
}: AdminBadgeCardProps) {
  const title =
    kind === "template"
      ? (item as BadgeTemplate).name
      : (item as BadgePrintItem).event_name;
  const subtitle =
    kind === "template" ? ticketName : (item as BadgePrintItem).edition_name;

  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

  const handleDelete = () => {
    if (!onDelete) return;
    onDelete();
  };

  const card = (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.04, duration: 0.32 }}
      className={cn(
        "group relative flex w-full min-w-60 max-w-full flex-col overflow-hidden rounded-2xl bg-card text-left",
        "ring-1 ring-foreground/5 shadow-xs",
        "transform-gpu will-change-transform",
        "transition-all duration-300 ease-out",
        "hover:-translate-y-0.5 hover:ring-foreground/10 hover:shadow-sm",
      )}
      role="button"
      tabIndex={0}
      onClick={onView}
    >
      <div className="relative aspect-video overflow-hidden bg-muted">
        <BadgePreview
          badge={item}
          className="h-full w-full object-cover"
          contain
          showVariables={kind === "template"}
        />

        {kind === "template" && (
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
                    aria-label={`Abrir ações de ${title}`}
                  >
                    <MoreVertical className="size-4" />
                  </button>
                }
              />
              <DropdownMenuContent align="end" className="w-44">
                {onEdit ? (
                  <DropdownMenuItem
                    onClick={(e) => {
                      e.stopPropagation();
                      onEdit();
                    }}
                  >
                    <Pencil className="size-4" />
                    <span>Editar</span>
                  </DropdownMenuItem>
                ) : null}
                {onDuplicate ? (
                  <DropdownMenuItem
                    onClick={(e) => {
                      e.stopPropagation();
                      onDuplicate();
                    }}
                  >
                    <Copy className="size-4" />
                    <span>Duplicar</span>
                  </DropdownMenuItem>
                ) : null}
                {onDelete ? (
                  <DropdownMenuItem
                    onClick={(e) => {
                      e.stopPropagation();
                      setIsDeleteDialogOpen(true);
                    }}
                    className="text-destructive"
                  >
                    <Trash2 className="size-4" />
                    <span>Excluir</span>
                  </DropdownMenuItem>
                ) : null}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        )}

        <span className="absolute left-3 top-3 inline-flex items-center gap-2 rounded-full border border-border/70 bg-background/85 px-3 py-1 text-[10px] font-semibold text-muted-foreground shadow-sm">
          {kind === "template" ? "Template" : "Crachá"}
        </span>
      </div>

      <CardHeader className="px-4 pt-4">
        <CardTitle className="text-base leading-tight">{title}</CardTitle>
        <p className="mt-1 text-xs text-muted-foreground">{subtitle}</p>
      </CardHeader>
      <CardContent className="px-4 pb-4 text-xs text-muted-foreground">
        {kind === "template"
          ? `${(item as BadgeTemplate).design_data.canvas.width} × ${(item as BadgeTemplate).design_data.canvas.height}px · ${(item as BadgeTemplate).design_data.elements.length} elementos`
          : (item as BadgePrintItem).ticket_name}
      </CardContent>
    </motion.article>
  );

  return (
    <>
      {kind === "template" ? (
        <ContextMenu>
          <ContextMenuTrigger render={card} />
          <ContextMenuContent className="w-44">
            {onEdit ? (
              <ContextMenuItem
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  onEdit();
                }}
              >
                <Pencil className="size-4" />
                <span>Editar</span>
              </ContextMenuItem>
            ) : null}
            {onDelete ? (
              <ContextMenuItem
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  setIsDeleteDialogOpen(true);
                }}
                className="text-destructive"
              >
                <Trash2 className="size-4" />
                <span>Excluir</span>
              </ContextMenuItem>
            ) : null}
          </ContextMenuContent>
        </ContextMenu>
      ) : (
        card
      )}

      <AlertDialog
        open={isDeleteDialogOpen}
        onOpenChange={(open) => setIsDeleteDialogOpen(open)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Excluir template?</AlertDialogTitle>
            <AlertDialogDescription>
              O template “{title}” será removido permanentemente.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              onClick={() => {
                handleDelete();
                setIsDeleteDialogOpen(false);
              }}
            >
              Excluir template
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
