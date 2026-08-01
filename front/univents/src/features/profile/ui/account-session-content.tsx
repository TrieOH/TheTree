import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { CheckCircle2, Fingerprint, Mail } from "lucide-react";
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

  const { email, id } = profile;

  return (
    <div className="space-y-1 pt-1 pb-0">
      <div className="rounded-lg border border-border/40 bg-card p-3 shadow-sm">
        <div className="mb-2 flex items-center gap-3">
          <div className="flex size-5 items-center justify-center rounded-full bg-primary/10 text-primary">
            <Fingerprint className="size-4" />
          </div>
          <div>
            <p className="font-medium m-0!">Sua conta</p>
            <p className="text-xs text-muted-foreground mt-0.5">
              Perfil autenticado
            </p>
          </div>
          <CheckCircle2 className="ml-auto size-5 text-emerald-500" />
        </div>
        <div className="space-y-1.5 border-t border-border/30 pt-2">
          <div className="flex items-center gap-3">
            <Mail className="size-5 shrink-0 text-muted-foreground" />
            <div className="min-w-0 flex-1">
              <p className="text-sm font-medium m-0!">Email</p>
              <p className="truncate text-xs text-muted-foreground mt-0.5">
                {email || "Não disponível"}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3 border-t border-border/30 pt-1.5">
            <Fingerprint className="size-5 shrink-0 text-muted-foreground" />
            <div className="min-w-0">
              <p className="text-sm font-medium m-0!">ID da conta</p>
              <p className="truncate font-mono text-xs text-muted-foreground">
                {id}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
