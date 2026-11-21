import type { JSX } from "@solidjs/web";
import { cn } from "./lib/cn";
import { shapeClass, type Shape } from "./lib/shape";

export interface BadgeProps {
  children: JSX.Element;
  variant?:
    | "default"
    | "accent"
    | "secondary"
    | "outline"
    | "muted"
    | "success"
    | "warning"
    | "destructive";
  shape?: Shape;
  class?: string;
}
const variants = {
  default: "border-transparent bg-muted text-foreground",
  accent: "border-transparent bg-accent text-accent-foreground",
  secondary: "border-transparent bg-secondary text-secondary-foreground",
  outline: "border-border text-foreground",
  muted: "border-transparent bg-muted text-muted-foreground",
  success: "border-transparent bg-success/15 text-success",
  warning: "border-transparent bg-warning/15 text-warning",
  destructive:
    "border-transparent bg-destructive/20 text-destructive dark:text-destructive-foreground",
};
export function Badge(props: BadgeProps) {
  return (
    <span
      class={cn(
        "inline-flex shrink-0 items-center rounded px-2 py-0.5 text-[11px] font-medium border",
        variants[props.variant ?? "default"],
        shapeClass(props.shape),
        props.class,
      )}
    >
      {props.children}
    </span>
  );
}
