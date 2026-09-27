import { createSignal } from "solid-js";
import type { JSX } from "@solidjs/web";

import LogOutIcon from "~icons/lucide/log-out";
import ShieldAlertIcon from "~icons/lucide/shield-alert";

import { useSessionActions } from "@/features/auths/hooks/use-session-actions";

type IconComponent = () => JSX.Element;

const LogOut =
  LogOutIcon as unknown as IconComponent;

const ShieldAlert =
  ShieldAlertIcon as unknown as IconComponent;

export function LogoutCard() {
  const { logout } = useSessionActions();
  const [isLoading, setIsLoading] =
    createSignal(false);

  const onLogout = async () => {
    if (isLoading()) return;

    setIsLoading(true);

    try {
      await logout();
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div class="relative overflow-hidden rounded-lg border border-border bg-background p-3">
      <div class="absolute inset-y-0 left-0 w-0.5 bg-destructive" />

      <div class="flex items-start gap-3">
        <div class="flex size-8 shrink-0 items-center justify-center rounded-lg bg-destructive/10 text-destructive [&>svg]:size-4">
          <ShieldAlert />
        </div>

        <div class="min-w-0 flex-1">
          <p class="text-sm font-semibold">
            Encerrar sessão
          </p>

          <p class="mt-1 text-xs leading-relaxed text-muted-foreground">
            Você será desconectado deste dispositivo.
          </p>
        </div>
      </div>

      <button
        type="button"
        disabled={isLoading()}
        onClick={onLogout}
        class="mt-3 flex h-9 w-full items-center justify-center gap-2 rounded-md border border-destructive/25 bg-destructive/10 px-3 text-xs font-semibold text-destructive transition-colors hover:bg-destructive hover:text-destructive-foreground disabled:cursor-not-allowed disabled:opacity-50 [&>svg]:size-4"
      >
        <LogOut />

        {isLoading()
          ? "Saindo..."
          : "Sair da conta"}
      </button>
    </div>
  );
}