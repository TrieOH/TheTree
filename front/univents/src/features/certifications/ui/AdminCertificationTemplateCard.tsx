import { FileText, MoreVertical, Pencil, Trash2 } from "lucide-react";
import { motion } from "motion/react";
import { cn } from "@/shared/lib/utils";
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
import type { CertificationTemplateI } from "../model";

interface AdminCertificationTemplateCardProps {
  template: CertificationTemplateI;
  index?: number;
  onEdit: () => void;
  onDelete: () => void;
  onView: () => void;
}

function MenuItems({
  isContext = false,
  onEdit,
  onView,
  onDelete,
}: {
  isContext?: boolean;
  onEdit: () => void;
  onView: () => void;
  onDelete: () => void;
}) {
  const Item = isContext ? ContextMenuItem : DropdownMenuItem;
  const stop =
    (action: () => void) => (e: React.MouseEvent | React.KeyboardEvent) => {
      e.preventDefault();
      e.stopPropagation();
      action();
    };

  return (
    <>
      <Item onClick={stop(onView)}>
        <FileText className="size-4" />
        <span>Ver template</span>
      </Item>
      <Item onClick={stop(onEdit)}>
        <Pencil className="size-4" />
        <span>Editar template</span>
      </Item>
      <Item
        onClick={stop(onDelete)}
        className="text-destructive focus:text-destructive"
      >
        <Trash2 className="size-4" />
        <span>Excluir template</span>
      </Item>
    </>
  );
}

export function AdminCertificationTemplateCard({
  template,
  index = 0,
  onEdit,
  onDelete,
  onView,
}: AdminCertificationTemplateCardProps) {
  const handleEdit = () => onEdit();
  const handleView = () => onView();

  return (
    <ContextMenu>
      <ContextMenuTrigger
        render={
          <motion.article
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{
              delay: index * 0.05,
              duration: 0.35,
              ease: [0.25, 0.1, 0.25, 1],
            }}
            className={cn(
              "group relative flex w-full min-w-62.5 max-w-full flex-col overflow-hidden rounded-2xl bg-card text-left",
              "ring-1 ring-foreground/10 shadow-xs",
              "transform-gpu will-change-transform",
              "transition-all duration-300 ease-out",
              "hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm",
              "focus:outline-none focus-visible:outline-none focus-visible:ring-0",
            )}
            role="button"
            tabIndex={0}
            onClick={handleView}
            onKeyDown={(event) => {
              if (event.key === "Enter" || event.key === " ") {
                event.preventDefault();
                handleView();
              }
            }}
          >
            <div className="relative h-24 overflow-hidden bg-muted">
              {template.design_data.background ? (
                <img
                  src={template.design_data.background}
                  alt={template.name}
                  className="h-full w-full object-cover transition-transform duration-700 ease-out group-hover:scale-[1.03]"
                  loading={index < 4 ? "eager" : "lazy"}
                />
              ) : (
                <div className="flex h-full w-full items-center justify-center bg-linear-to-br from-muted via-background to-muted/40">
                  <div className="flex size-12 items-center justify-center rounded-full border border-dashed border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
                    <FileText className="size-5 text-muted-foreground/40" />
                  </div>
                </div>
              )}

              <div className="absolute inset-0 bg-linear-to-t from-background/90 via-background/35 to-transparent" />

              <div className="absolute left-3 top-3 flex flex-wrap items-center gap-1.5">
                <span className="inline-flex items-center gap-1 rounded-full border border-border/60 bg-background/75 px-2 py-0.5 text-[10px] font-medium text-muted-foreground backdrop-blur-sm">
                  <FileText className="size-3.5" />
                  Template
                </span>
              </div>

              <div className="absolute right-2 top-2">
                <DropdownMenu>
                  <DropdownMenuTrigger
                    render={
                      <button
                        type="button"
                        onClick={(e) => e.stopPropagation()}
                        className={cn(
                          "inline-flex size-8 items-center justify-center rounded-full",
                          "bg-background/85 text-foreground shadow-sm backdrop-blur-sm",
                          "transition-colors hover:bg-background",
                        )}
                        aria-label={`Abrir ações de ${template.name}`}
                      >
                        <MoreVertical className="size-4" />
                      </button>
                    }
                  />
                  <DropdownMenuContent align="end" className="w-56">
                    <MenuItems
                      onEdit={handleEdit}
                      onView={handleView}
                      onDelete={onDelete}
                    />
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>

            <div className="flex items-center justify-between gap-3 p-3">
              <div className="min-w-0 flex-1 space-y-2">
                <div className="min-w-0 space-y-1">
                  <h3 className="line-clamp-2 text-sm font-semibold leading-snug text-foreground transition-colors duration-300 group-hover:text-primary">
                    {template.name}
                  </h3>
                  <p className="text-[11px] text-muted-foreground">
                    {template.design_data.background
                      ? "Com fundo configurado"
                      : "Sem fundo"}
                  </p>
                  {template.description ? (
                    <p className="line-clamp-2 text-xs text-muted-foreground">
                      {template.description}
                    </p>
                  ) : null}
                </div>
              </div>
            </div>
          </motion.article>
        }
      />
      <ContextMenuContent className="w-56">
        <MenuItems
          isContext
          onEdit={handleEdit}
          onView={handleView}
          onDelete={onDelete}
        />
      </ContextMenuContent>
    </ContextMenu>
  );
}
