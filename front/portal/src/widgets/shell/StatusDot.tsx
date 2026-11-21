import { StatusDot as Dot, type StatusDotTone } from "@trieoh/ui-solid";

export type ServiceStatus = "operational" | "degraded" | "idle" | "unknown";

const LABEL: Record<ServiceStatus, string> = {
  operational: "Operational",
  degraded: "Degraded",
  idle: "Idle",
  unknown: "No data yet",
};

const TONE: Record<ServiceStatus, StatusDotTone> = {
  operational: "success",
  degraded: "warning",
  idle: "idle",
  unknown: "muted",
};

/** Idle (no traffic in the window) is a hollow dot, never "down". */
export function StatusDot(props: { status: ServiceStatus; class?: string }) {
  return <Dot tone={TONE[props.status]} label={LABEL[props.status]} class={props.class} />;
}

export function statusLabel(status: ServiceStatus): string {
  return LABEL[status];
}
