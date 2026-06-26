import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import { timeAgo } from "@/shared/lib/date-utils";
import { KeySquare, Copy, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import type { ApiKeyI } from "../model";

interface ApiKeyCardProps {
  data: ApiKeyI;
  onRevoke?: (key: ApiKeyI) => void;
}

export function ApiKeyCard({ data, onRevoke }: ApiKeyCardProps) {
  const isRevoked = !!data.revoked_at;
  const isExpired = data.expires_at ? new Date(data.expires_at) < new Date() : false;

  const handleCopyPrefix = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    navigator.clipboard.writeText(data.key_prefix);
    toast.success("Key prefix copied to clipboard");
  };

  const status = isRevoked
    ? { label: "Revoked", variant: "destructive" as const }
    : isExpired
      ? { label: "Expired", variant: "outline" as const }
      : { label: "Active", variant: "default" as const };

  return (
    <div
      className={cn(
        "bg-card rounded-sm w-full cursor-default",
        "ring-1 ring-foreground/10 shadow-xs",
        "flex items-start gap-3 px-4 py-3",
        "hover:ring-foreground/20 duration-150",
        isRevoked && "opacity-60",
      )}
    >
      {/* Icon */}
      <div className="shrink-0 size-9 rounded-full bg-muted ring-1 ring-foreground/10 flex items-center justify-center mt-0.5">
        <KeySquare className="size-4 text-muted-foreground" />
      </div>

      {/* Info */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-sm font-semibold truncate">
            {data.name}
          </span>
          <Badge variant={status.variant} className="text-[10px] px-1.5 py-0">
            {status.label}
          </Badge>
        </div>
        <div className="flex items-center gap-2 mt-0.5">
          <span className="text-xs text-muted-foreground font-mono">
            {data.key_prefix}...
          </span>
          <ShadowButton
            variant="ghost"
            onClick={handleCopyPrefix}
            className="p-0 h-auto"
            leftIcon={<Copy className="size-3" />}
          />
        </div>
        <p className="text-xs text-muted-foreground mt-0.5">
          Created {timeAgo(data.created_at)}
          {data.last_used_at && ` · Last used ${timeAgo(data.last_used_at)}`}
          {data.expires_at && ` · Expires ${timeAgo(data.expires_at)}`}
        </p>
      </div>

      {/* Revoke button */}
      {!isRevoked && (
        <button
          type="button"
          className={cn(
            "shrink-0 text-muted-foreground p-2 rounded-sm border border-transparent transition-all duration-150 cursor-pointer outline-none self-start mt-0.5",
            "hover:text-destructive hover:bg-destructive/10",
          )}
          onClick={() => onRevoke?.(data)}
          title="Revoke API key"
        >
          <Trash2 className="size-4" />
        </button>
      )}
    </div>
  );
}