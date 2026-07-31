import {
  AlertTriangle,
  Ban,
  Clock,
  ExternalLink,
  MoreVertical,
  ShieldCheck,
} from "lucide-react";
import { motion } from "motion/react";
import type React from "react";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
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
import type { CertificationEmissionErrorI, CertificationI } from "../model";

const gradients = [
  "from-violet-500/20 via-fuchsia-500/10 to-background",
  "from-emerald-500/20 via-teal-500/10 to-background",
  "from-amber-500/20 via-orange-500/10 to-background",
  "from-blue-500/20 via-cyan-500/10 to-background",
];

function stop(action: () => void) {
  return (event: React.MouseEvent | React.KeyboardEvent) => {
    event.preventDefault();
    event.stopPropagation();
    action();
  };
}

export function CertificationCard({
  certification,
  email,
  scopeName,
  index,
  onView,
  onInvalidate,
}: {
  certification: CertificationI;
  email?: string;
  scopeName: string;
  index: number;
  onView: () => void;
  onInvalidate: () => void;
}) {
  const valid = certification.valid;
  const status = valid
    ? "Válido"
    : certification.invalid_reason
      ? "Invalidado"
      : "Pendente";
  const card = (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05, duration: 0.35 }}
      className="group flex w-full flex-col overflow-hidden rounded-2xl bg-card text-left ring-1 ring-foreground/10 shadow-xs transition-all duration-300 hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm"
    >
      <div
        role="button"
        tabIndex={0}
        onClick={onView}
        onKeyDown={(event) => {
          if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            onView();
          }
        }}
        className={cn(
          "relative aspect-video overflow-hidden bg-linear-to-br text-left",
          gradients[index % gradients.length],
        )}
      >
        <div className="flex h-full items-center justify-center">
          <div className="flex size-18 items-center justify-center rounded-full border border-border/70 bg-background/80 shadow-sm backdrop-blur-sm">
            {status === "Válido" ? (
              <ShieldCheck className="size-7 text-muted-foreground/45" />
            ) : status === "Invalidado" ? (
              <Ban className="size-7 text-muted-foreground/45" />
            ) : (
              <Clock className="size-7 text-muted-foreground/45" />
            )}
          </div>
        </div>
        <div className="absolute inset-0 bg-linear-to-t from-background/95 via-background/25 to-transparent" />
        <div
          className={cn(
            "absolute left-3 top-3 rounded-full border px-2 py-0.5 text-[10px] font-medium",
            valid
              ? "border-emerald-500/20 bg-emerald-500/10 text-emerald-700"
              : "border-amber-500/20 bg-amber-500/10 text-amber-700",
          )}
        >
          <span
            className={cn(
              "mr-1.5 inline-block size-1.5 rounded-full",
              valid ? "bg-emerald-500" : "bg-amber-500",
            )}
          />
          {status}
        </div>
        <div className="absolute right-3 top-3">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <button
                  type="button"
                  onClick={(event) => event.stopPropagation()}
                  className="inline-flex size-9 items-center justify-center rounded-full bg-background/85 text-foreground shadow-sm backdrop-blur-sm hover:bg-background"
                  aria-label="Abrir ações do certificado"
                />
              }
            >
              <MoreVertical className="size-4" />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuItem onClick={stop(onView)}>
                <ExternalLink className="size-4" />
                Verificar certificado
              </DropdownMenuItem>
              <DropdownMenuItem
                disabled={!valid}
                variant="destructive"
                onClick={stop(onInvalidate)}
              >
                <Ban className="size-4" />
                Invalidar
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        <div className="absolute inset-x-0 bottom-0 p-4 sm:p-5">
          <p className="text-[10px] font-medium uppercase tracking-wide text-muted-foreground">
            Certificado
          </p>
          <h3 className="line-clamp-1 text-lg font-semibold leading-snug group-hover:text-primary sm:text-xl">
            {email ?? `Usuário ${certification.user_id.slice(0, 8)}`}
          </h3>
        </div>
      </div>
      <div className="space-y-3 p-4 sm:p-5 sm:pt-4">
        <p className="line-clamp-1 min-h-5 text-sm text-muted-foreground">
          {scopeName}
        </p>
        <span className="inline-flex items-center gap-1.5 text-xs text-muted-foreground">
          Emitido em{" "}
          {new Date(certification.issued_at).toLocaleDateString("pt-BR")}
        </span>
        <Button
          type="button"
          variant="outline"
          className="h-9 w-full gap-2"
          onClick={stop(onView)}
        >
          <ExternalLink className="size-4" />
          Verificar certificado
        </Button>
      </div>
    </motion.article>
  );
  return (
    <ContextMenu>
      <ContextMenuTrigger render={card} />
      <ContextMenuContent className="w-56">
        <ContextMenuItem onClick={stop(onView)}>
          <ExternalLink className="size-4" />
          Verificar certificado
        </ContextMenuItem>
        <ContextMenuItem
          disabled={!valid}
          variant="destructive"
          onClick={stop(onInvalidate)}
        >
          <Ban className="size-4" />
          Invalidar
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  );
}

export function CertificationEmissionErrorCard({
  error,
  email,
  scopeName,
}: {
  error: CertificationEmissionErrorI;
  email?: string;
  scopeName: string;
}) {
  return (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="flex h-full min-h-40 flex-col overflow-hidden rounded-2xl bg-card text-left ring-1 ring-foreground/10 shadow-xs"
    >
      <div className="flex items-start gap-3 border-b border-destructive/10 bg-destructive/5 px-4 py-3">
        <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-destructive/10 text-destructive">
          <AlertTriangle className="size-4" />
        </div>
        <div className="min-w-0">
          <p className="truncate text-sm font-semibold">
            {email ?? `Usuário ${error.user_id.slice(0, 8)}`}
          </p>
          <p className="mt-1 truncate text-xs text-muted-foreground">
            {scopeName}
          </p>
        </div>
      </div>
      <div className="flex flex-1 flex-col justify-between gap-4 p-4">
        <p className="text-sm leading-6 text-foreground/85">
          {error.error_message}
        </p>
        <p className="text-xs text-muted-foreground">
          {new Date(error.created_at).toLocaleString("pt-BR")}
        </p>
      </div>
    </motion.article>
  );
}
