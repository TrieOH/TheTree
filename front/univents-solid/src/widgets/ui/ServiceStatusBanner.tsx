import { For, Show, createMemo } from "solid-js";
import { refetchFailedQueries } from "@/shared/lib/query-client";
import type { ServiceName } from "@/shared/lib/service-error";
import { degradedServices, serviceHealth } from "@/shared/lib/service-health";
import { useServiceHealth } from "@/shared/lib/use-service-health";

const NOTICES: Record<ServiceName, string> = {
  api: "Instabilidade na conexão",
  identityx: "Serviço de conta instável",
};

/**
 * Discreet indicator that warns the user when a service is struggling without
 * getting in the way or shifting layout. Disappears once the service recovers.
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

  return (
    <Show when={services().length > 0}>
      <div class="pointer-events-none fixed top-4 right-4 z-50 flex justify-end print:hidden">
        <div
          role="status"
          aria-live="polite"
          class="pointer-events-auto flex items-center gap-2 rounded-full border border-border/80 bg-background/90 px-3 py-1.5 text-xs text-muted-foreground shadow-md backdrop-blur-md transition-all"
        >
          <span class="inline-block h-2 w-2 shrink-0 rounded-full bg-amber-500 animate-pulse" />
          <div class="flex items-center gap-1.5">
            <For each={services()}>
              {(service) => <span class="font-medium text-foreground">{NOTICES[service]}</span>}
            </For>
          </div>
          <button
            type="button"
            onClick={retry}
            class="ml-1 rounded-full px-2 py-0.5 text-xs font-medium text-foreground hover:bg-muted transition-colors cursor-pointer"
          >
            Tentar agora
          </button>
        </div>
      </div>
    </Show>
  );
}
