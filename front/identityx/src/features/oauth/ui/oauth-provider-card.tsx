import { Code2, Globe, Pencil, Power, Trash2 } from "lucide-react";
import type { OAuthProviderI } from "../model";

interface Props {
  data: OAuthProviderI;
  onEdit: (provider: OAuthProviderI) => void;
  onToggle: (provider: OAuthProviderI) => void;
  onDelete: (provider: OAuthProviderI) => void;
}

const PROVIDER_STYLES = {
  github: {
    icon: Code2,
    accent: "bg-zinc-900 dark:bg-zinc-100",
    iconBg: "bg-zinc-900 dark:bg-zinc-100",
    iconFg: "text-white dark:text-zinc-900",
  },
  default: {
    icon: Globe,
    accent: "bg-blue-500",
    iconBg: "bg-blue-500",
    iconFg: "text-white",
  },
} as const;

export function OAuthProviderCard({ data, onEdit, onToggle, onDelete }: Props) {
  const style =
    data.provider === "github"
      ? PROVIDER_STYLES.github
      : PROVIDER_STYLES.default;
  const Icon = style.icon;

  return (
    <div className="relative flex w-full max-w-full flex-col overflow-hidden rounded-md border border-border bg-card">
      {/* Accent stripe — identidade do provider */}
      <span className={`absolute inset-y-0 left-0 w-1 ${style.accent}`} />

      <div className="flex flex-col gap-3 py-3 pl-4 pr-3">
        {/* Header */}
        <div className="flex items-center gap-2.5 min-w-0">
          <div
            className={`flex size-8 shrink-0 items-center justify-center rounded-md ${style.iconBg}`}
          >
            <Icon className={`size-4 ${style.iconFg}`} />
          </div>

          <div className="min-w-0 flex-1">
            <span className="block truncate text-sm font-semibold capitalize leading-tight">
              {data.provider}
            </span>
            <span className="flex items-center gap-1 text-[11px] text-muted-foreground leading-tight">
              <span
                className={`size-1.5 rounded-full ${
                  data.enabled ? "bg-emerald-500" : "bg-muted-foreground/40"
                }`}
              />
              {data.enabled ? "Active" : "Inactive"}
            </span>
          </div>
        </div>

        {/* Client ID — tratado como dado, não como texto solto */}
        <div
          className="min-w-0 truncate rounded bg-muted/60 px-2 py-1 font-mono text-[11px] text-muted-foreground"
          title={data.client_id}
        >
          {data.client_id}
        </div>
      </div>

      {/* Footer de ações — Edit/Toggle agrupados, Delete isolado */}
      <div className="flex items-stretch border-t border-border">
        <div className="flex flex-1 divide-x divide-border">
          <button
            type="button"
            onClick={() => onToggle(data)}
            className="flex flex-1 items-center justify-center gap-1.5 py-2 text-xs font-medium text-foreground/80 transition-colors hover:bg-muted"
          >
            <Power className="size-3.5" />
            {data.enabled ? "Disable" : "Enable"}
          </button>
          <button
            type="button"
            onClick={() => onEdit(data)}
            className="flex flex-1 items-center justify-center gap-1.5 py-2 text-xs font-medium text-foreground/80 transition-colors hover:bg-muted"
          >
            <Pencil className="size-3.5" />
            Edit
          </button>
        </div>
        <button
          type="button"
          onClick={() => onDelete(data)}
          aria-label="Delete provider"
          title="Delete"
          className="flex w-11 shrink-0 items-center justify-center border-l border-border text-destructive transition-colors hover:bg-destructive/10"
        >
          <Trash2 className="size-3.5" />
        </button>
      </div>
    </div>
  );
}
