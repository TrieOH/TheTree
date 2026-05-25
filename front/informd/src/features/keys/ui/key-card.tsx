import { cn } from "#/shared/lib/utils";
import { Key, KeyRound, ShieldOff } from "lucide-react";
import type { ApiKeyI } from "../model";
import { timeAgo } from "#/shared/lib/helpers/date-utils";
import { Button } from "#/shared/ui/shadcn/button";

interface PropsI {
  data: ApiKeyI;
  onRevoke?: (apiKey: ApiKeyI | null) => void;
}

export function APIKeyCard({ data, onRevoke }: PropsI) {
  const isRevoked = !!data.revoked_at;

  return (
    <div
      className={cn(
        "bg-card rounded-sm w-full cursor-default",
        "ring-1 ring-foreground/10 shadow-xs",
        "flex items-center gap-3 px-4 py-3",
        "hover:ring-foreground/20 duration-150",
        isRevoked && "opacity-60"
      )}
    >
      {/* Icon */}
      <div className="shrink-0 size-9 rounded-full bg-muted ring-1 ring-foreground/10 flex items-center justify-center">
        {isRevoked ? (
          <KeyRound className="size-4 text-muted-foreground" />
        ) : (
          <Key className="size-4 text-muted-foreground" />
        )}
      </div>

      {/* Info */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-sm font-semibold truncate">{data.name}</span>
          <span className="inline-flex items-center gap-1 text-xs font-mono text-muted-foreground bg-muted px-1.5 py-0.5 rounded-sm">
            {data.prefix}…
          </span>
          {isRevoked && (
            <span className="inline-flex items-center gap-1 text-xs font-medium px-1.5 py-0.5 rounded-sm text-destructive bg-destructive/10">
              Revoked
            </span>
          )}
        </div>
        <p className="text-xs text-muted-foreground mt-0.5">
          {isRevoked
            ? `Revoked ${timeAgo(data.revoked_at!)}`
            : `Created ${timeAgo(data.created_at)}`}
        </p>
      </div>

      {/* Revoke button */}
      {!isRevoked && (
        <Button
          variant="ghost"
          size="icon"
          className={cn(
            "shrink-0 text-muted-foreground",
            "hover:text-destructive hover:bg-destructive/10",
            "duration-150 cursor-pointer outline-0"
          )}
          onClick={() => onRevoke?.(data)}
          title="Revoke API key"
        >
          <ShieldOff className="size-4" />
        </Button>
      )}
    </div>
  );
}