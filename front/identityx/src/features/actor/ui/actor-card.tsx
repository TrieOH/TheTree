import type { ActorI } from "../model";
import { cn } from "@/shared/lib/utils";
import { Copy, Cpu, Fingerprint, Globe, KeyRound, Shield, User2, Github } from "lucide-react";
import { toast } from "sonner";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";

interface ActorCardProps {
  data: ActorI;
}

const typeConfig = {
  human: {
    label: "Human",
    icon: User2,
    className: "text-blue-500 bg-blue-500/10",
  },
  service: {
    label: "Service",
    icon: Shield,
    className: "text-emerald-500 bg-emerald-500/10",
  },
  machine: {
    label: "Machine",
    icon: Cpu,
    className: "text-amber-500 bg-amber-500/10",
  },
} as const;

const authMethodLabel = {
  password: "Password",
  api_key: "API Key",
  google_auth: "Google",
  github_auth: "GitHub",
} as const;

const authMethodIcon = {
  password: KeyRound,
  api_key: Fingerprint,
  google_auth: Globe,
  github_auth: Github,
} as const;

export function ActorCard({ data }: ActorCardProps) {
  const type = typeConfig[data.type] ?? typeConfig.human;
  const TypeIcon = type.icon;
  const AuthIcon = authMethodIcon[data.auth_method] ?? Fingerprint;

  const handleCopyId = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    navigator.clipboard.writeText(data.id);
    toast.success("Actor ID copied to clipboard");
  };

  return (
    <div
      className={cn(
        "bg-card rounded-md w-full min-w-0 cursor-default",
        "ring-1 ring-foreground/10 shadow-xs",
        "flex items-start gap-3 px-3 py-3 sm:px-4",
        "hover:ring-foreground/20 duration-150",
        data.deleted_at && "opacity-60",
      )}
    >
      <div className="shrink-0 size-10 rounded-md bg-muted ring-1 ring-foreground/10 flex items-center justify-center">
        <TypeIcon className="size-4.5 text-muted-foreground" />
      </div>

      <div className="min-w-0 flex-1">
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <span className="block text-sm font-semibold truncate">
              {data.email ?? "No email"}
            </span>
            <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
              <span className="inline-flex items-center gap-1.5">
                <AuthIcon className="size-3.5 shrink-0" />
                <span className="truncate">
                  {authMethodLabel[data.auth_method as keyof typeof authMethodLabel] ?? data.auth_method}
                </span>
              </span>
              {data.project_id && (
                <span className="font-mono truncate">
                  {data.project_id}
                </span>
              )}
            </div>
          </div>

          <ShadowButton
            variant="ghost"
            onClick={handleCopyId}
            className="shrink-0 p-0 h-auto"
            leftIcon={<Copy className="size-3" />}
          />
        </div>

        <div className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-[11px] text-muted-foreground">
          <span>
            Created {new Date(data.created_at).toLocaleDateString()}
          </span>
          {data.verified_at && (
            <span>
              Verified {new Date(data.verified_at).toLocaleDateString()}
            </span>
          )}
          {data.deleted_at && (
            <span>
              Deleted {new Date(data.deleted_at).toLocaleDateString()}
            </span>
          )}
        </div>
      </div>
    </div>
  );
}
