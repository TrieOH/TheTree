import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { Button, cn } from "@trieoh/ui-solid";
import { Match, Show, Switch } from "solid-js";
import CloudOffIcon from "~icons/lucide/cloud-off";
import TriangleAlertIcon from "~icons/lucide/triangle-alert";
import WifiOffIcon from "~icons/lucide/wifi-off";
import { classifyError } from "@/shared/lib/service-error";

const CloudOff = CloudOffIcon as unknown as (props: { class?: string }) => JSX.Element;
const TriangleAlert = TriangleAlertIcon as unknown as (props: { class?: string }) => JSX.Element;
const WifiOff = WifiOffIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface QueryErrorPanelProps {
  error: unknown;
  /** What failed to load, as it reads after "Não foi possível carregar": "os eventos". */
  what: string;
  onRetry: () => void;
  /** Replaces the generic reason for failures that deserve a specific one (e.g. a 404). */
  message?: string;
  /** A retry is in flight: keep the panel up instead of flashing a skeleton. */
  retrying?: boolean;
  class?: string;
}

/**
 * The inline "couldn't load this" state: scoped to the region that failed, so
 * the rest of the page, and the navigation, keep working.
 */
export function QueryErrorPanel(props: QueryErrorPanelProps) {
  const failure = () => classifyError(props.error);
  // A retry clears the previous error before the new one arrives; with nothing
  // known, say nothing, rather than reporting an "unexpected error".
  const reason = () =>
    props.message ?? (props.error === undefined || props.error === null ? "" : failure().message);

  return (
    <div
      role="status"
      class={cn(
        "flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border bg-muted/30 px-6 py-10 text-center",
        props.class,
      )}
    >
      <span class="flex size-10 items-center justify-center rounded-full bg-muted text-muted-foreground">
        <Switch fallback={<TriangleAlert class="size-5" />}>
          <Match when={failure().kind === "offline"}>
            <WifiOff class="size-5" />
          </Match>
          <Match when={failure().kind === "unavailable" || failure().kind === "timeout"}>
            <CloudOff class="size-5" />
          </Match>
        </Switch>
      </span>
      <div class="space-y-1">
        <p class="text-sm font-medium text-foreground">
          Não foi possível carregar {props.what}
        </p>
        <Show when={reason()}>
          <p class="text-sm text-muted-foreground">{reason()}</p>
        </Show>
      </div>
      {/* Asking again can't change a definitive answer (404, 403). */}
      <Switch
        fallback={
          <Button
            variant="outline"
            size="sm"
            disabled={props.retrying}
            onClick={() => props.onRetry()}
          >
            {props.retrying ? "Tentando novamente…" : "Tentar novamente"}
          </Button>
        }
      >
        <Match when={failure().kind === "unauthorized"}>
          <Link
            to="/auth"
            class="inline-flex h-8 items-center rounded-md border border-border bg-background px-3 text-sm font-medium hover:bg-muted"
          >
            Entrar
          </Link>
        </Match>
        <Match when={failure().kind === "not-found" || failure().kind === "forbidden"}>
          {null}
        </Match>
      </Switch>
    </div>
  );
}
