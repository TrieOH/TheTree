import { cn } from "./lib/cn";

export type StatusDotTone =
  | "success"
  | "warning"
  | "destructive"
  | "muted"
  | "idle";

export interface StatusDotProps {
  tone: StatusDotTone;
  /** Read by screen readers and shown on hover. */
  label: string;
  size?: "sm" | "md";
  class?: string;
}

const tones: Record<StatusDotTone, string> = {
  success: "bg-success",
  warning: "bg-warning",
  destructive: "bg-destructive dark:bg-destructive-foreground",
  muted: "bg-muted-foreground/40",
  idle: "ring-[1.5px] ring-inset ring-muted-foreground",
};

/** A small status light. `idle` is hollow: nothing happening, not broken. */
export function StatusDot(props: StatusDotProps) {
  return (
    <span
      role="img"
      aria-label={props.label}
      title={props.label}
      class={cn(
        "inline-block shrink-0 rounded-full",
        props.size === "md" ? "size-2" : "size-1.5",
        tones[props.tone],
        props.class,
      )}
    />
  );
}
