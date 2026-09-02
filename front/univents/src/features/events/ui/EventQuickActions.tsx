import { Link } from "@tanstack/react-router";
import {
  Command,
  Copy,
  ExternalLink,
  Eye,
  Pencil,
  XCircle,
} from "lucide-react";
import { Button } from "@/shared/ui/shadcn/button";

export interface EventQuickAction {
  label: string;
  shortcut: string;
  disabled?: boolean;
  variant: "default" | "destructive";
  onClick?: () => void;
  to?: "/events/$slug";
  params?: { slug: string };
}

function actionIcon(label: string) {
  if (label.includes("Editar")) return Pencil;
  if (label.includes("Copiar")) return Copy;
  if (label.includes("Publicar")) return Eye;
  if (label.includes("Descontinuar")) return XCircle;
  return ExternalLink;
}

function compactLabel(label: string) {
  return label
    .replace(" evento", "")
    .replace(" público", "")
    .replace(" conta", "");
}

function Shortcut({ value }: { value: string }) {
  return (
    <kbd className="hidden rounded border border-border/70 bg-muted/70 px-1 py-0.5 font-mono text-[9px] text-muted-foreground sm:inline-block">
      {value.replace("Mod", "⌘/Ctrl")}
    </kbd>
  );
}

export function EventQuickActions({
  actions,
}: {
  actions: EventQuickAction[];
}) {
  return (
    <div
      className="order-1 space-y-3 rounded-xl border border-border/60 bg-muted/20 p-3"
      role="toolbar"
      aria-label="Atalhos do evento"
    >
      <div className="flex items-center justify-between gap-3 px-1">
        <div className="flex min-w-0 items-center gap-2">
          <div className="flex size-7 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
            <Command className="size-3.5" />
          </div>
          <div className="min-w-0">
            <p className="text-xs font-semibold text-foreground">
              Ações rápidas
            </p>
            <p className="truncate text-[11px] text-muted-foreground">
              Atalhos para as tarefas mais usadas
            </p>
          </div>
        </div>
        <span className="hidden shrink-0 text-[11px] text-muted-foreground sm:inline">
          Ctrl/⌘ + tecla
        </span>
      </div>
      <div className="flex flex-wrap items-center gap-2">
        {actions.map((action) => {
          const Icon = actionIcon(action.label);
          const content = (
            <>
              <span className="flex items-center gap-1.5">
                <Icon className="size-4" />
                <span>{compactLabel(action.label)}</span>
              </span>
              <Shortcut value={action.shortcut} />
            </>
          );

          if (action.to && action.params) {
            return (
              <Link
                key={action.label}
                to={action.to}
                params={action.params}
                aria-disabled={action.disabled}
                title={`${action.label} · ${action.shortcut}`}
                aria-label={action.label}
                className="inline-flex h-9 shrink-0 flex-row items-center justify-center gap-1.5 rounded-md border border-border bg-background px-2 text-xs font-medium text-foreground shadow-xs transition-colors hover:bg-muted aria-disabled:pointer-events-none aria-disabled:opacity-50 sm:h-14! sm:min-w-28! sm:flex-col! sm:gap-1 sm:px-2 sm:py-1.5 sm:text-[11px] sm:leading-tight"
              >
                {content}
              </Link>
            );
          }

          return (
            <Button
              key={action.label}
              size="default"
              variant={
                action.variant === "destructive" ? "destructive" : "outline"
              }
              className={`h-9 shrink-0 flex-row gap-1.5 border px-2 text-xs sm:h-14! sm:min-w-28! sm:flex-col! sm:gap-1 sm:px-2 sm:py-1.5 sm:text-[11px] sm:leading-tight ${action.variant === "destructive" ? "border-destructive/60" : "border-border"}`}
              disabled={action.disabled}
              onClick={action.onClick}
              title={`${action.label} · ${action.shortcut}`}
              aria-label={action.label}
            >
              {content}
            </Button>
          );
        })}
      </div>
    </div>
  );
}
