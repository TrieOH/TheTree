import { cn } from "@trieoh/ui-solid";

export type ServiceStatus = "operational" | "degraded" | "idle" | "unknown";

const LABEL: Record<ServiceStatus, string> = {
  operational: "Operational",
  degraded: "Degraded",
  idle: "Idle",
  unknown: "No data yet",
};

/** Idle (no traffic in the window) is a hollow dot, never "down". */
export function StatusDot(props: { status: ServiceStatus; class?: string }) {
  return (
    <span
      role="img"
      aria-label={LABEL[props.status]}
      title={LABEL[props.status]}
      class={cn(
        "inline-block size-1.5 shrink-0 rounded-full",
        props.status === "operational" && "bg-ok",
        props.status === "degraded" && "bg-warn",
        props.status === "idle" && "ring-[1.5px] ring-inset ring-muted-foreground",
        props.status === "unknown" && "bg-muted-foreground/40",
        props.class,
      )}
    />
  );
}

export function statusLabel(status: ServiceStatus): string {
  return LABEL[status];
}
