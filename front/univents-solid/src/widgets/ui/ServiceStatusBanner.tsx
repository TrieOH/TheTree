import type { JSX } from "@solidjs/web";
import { Show, createMemo } from "solid-js";
import WifiOffIcon from "~icons/lucide/wifi-off";
import { refetchFailedQueries } from "@/shared/lib/query-client";
import type { ServiceName } from "@/shared/lib/service-error";
import { degradedServices, serviceHealth } from "@/shared/lib/service-health";
import { useServiceHealth } from "@/shared/lib/use-service-health";

const WifiOff = WifiOffIcon as unknown as (props: { class?: string }) => JSX.Element;

/**
 * Minimal, non-intrusive connection indicator.
 * Appears only when services fail, without pulse animations or noisy text.
 * Clicking attempts to reconnect.
 */
export function ServiceStatusBanner() {
  const health = useServiceHealth();
  const key = createMemo(() => degradedServices(health()).join(","), {
    name: "ServiceStatusBanner.key",
  });
  const services = createMemo(
    () => (key() ? (key().split(",") as ServiceName[]) : []),
    { name: "ServiceStatusBanner.services" },
  );

  const retry = () => {
    serviceHealth.halfOpenAll();
    void refetchFailedQueries();
  };

  const title = () =>
    services().includes("identityx")
      ? "Sessão desconectada · Clique para reconectar"
      : "Sem conexão · Clique para reconectar";

  return (
    <Show when={services().length > 0}>
      <div class="pointer-events-none fixed top-4 right-4 z-50 flex justify-end print:hidden">
        <button
          type="button"
          role="status"
          aria-label={title()}
          title={title()}
          onClick={retry}
          class="pointer-events-auto flex size-7 items-center justify-center rounded-full border border-border/80 bg-background/90 text-amber-600 shadow-xs backdrop-blur-sm transition-colors hover:bg-muted hover:text-amber-700 dark:text-amber-400 dark:hover:text-amber-300 cursor-pointer"
        >
          <WifiOff class="size-3.5" />
        </button>
      </div>
    </Show>
  );
}
