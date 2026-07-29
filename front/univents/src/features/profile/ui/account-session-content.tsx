import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { Globe, Mail, MonitorSmartphone } from "lucide-react";
import { Skeleton } from "@/shared/ui/shadcn/skeleton";

export function AccountSessionContent() {
  const { auth, isInitializing, isAuthenticated } = useAuth();
  const profile = auth.profile();

  if (isInitializing) {
    return (
      <div className="space-y-3 pt-4 pb-2">
        <Skeleton className="h-16 w-full" />
        <Skeleton className="h-16 w-full" />
        <Skeleton className="h-16 w-full" />
      </div>
    );
  }

  if (!isAuthenticated || !profile) {
    return (
      <div className="pt-4 pb-2 text-sm text-muted-foreground">
        Nenhuma sessão ativa encontrada.
      </div>
    );
  }

  const { email, user_ip, user_agent } = profile;

  return (
    <div className="space-y-3 pt-4 pb-2">
      <div className="flex items-center gap-4 rounded-lg border border-border/50 p-4 bg-card">
        <Mail className="size-5 shrink-0 text-muted-foreground" />
        <div className="min-w-0 flex-1">
          <p className="text-sm font-medium">Email</p>
          <p className="truncate text-sm text-muted-foreground">
            {email || "Não disponível"}
          </p>
        </div>
      </div>

      <div className="flex items-center gap-4 rounded-lg border border-border/50 p-4 bg-card">
        <Globe className="size-5 shrink-0 text-muted-foreground" />
        <div className="min-w-0 flex-1">
          <p className="text-sm font-medium">Endereço IP</p>
          <p className="truncate text-sm text-muted-foreground">
            {user_ip || "Não disponível"}
          </p>
        </div>
      </div>

      <div className="flex items-center gap-4 rounded-lg border border-border/50 p-4 bg-card">
        <MonitorSmartphone className="size-5 shrink-0 text-muted-foreground" />
        <div className="min-w-0 flex-1">
          <p className="text-sm font-medium">User Agent</p>
          <p
            className="line-clamp-2 text-sm text-muted-foreground break-all"
            title={user_agent || undefined}
          >
            {user_agent || "Não disponível"}
          </p>
        </div>
      </div>
    </div>
  );
}
