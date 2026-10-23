import type { JSX } from "@solidjs/web";
import { cn } from "./lib/cn";

export interface TooltipProps {
  label: JSX.Element;
  children: JSX.Element;
  placement?: "top" | "bottom";
  align?: "start" | "center" | "end";
  class?: string;
  tooltipClass?: string;
}

export function Tooltip(props: TooltipProps): JSX.Element {
  return (
    <div class={cn("group relative inline-flex", props.class)}>
      {props.children}
      <span
        role="tooltip"
        class={cn(
          "pointer-events-none absolute z-50 w-max max-w-[calc(100vw-2rem)] rounded-md border border-border/50 bg-popover px-3 py-2 text-xs leading-snug font-medium text-popover-foreground opacity-0 shadow-lg shadow-black/10 transition-opacity duration-150 group-hover:opacity-100 group-focus-within:opacity-100 max-md:hidden",
          props.placement === "bottom"
            ? "top-full mt-1.5"
            : "bottom-full mb-1.5",
          props.align === "start"
            ? "left-0"
            : props.align === "end"
              ? "right-0"
              : "left-1/2 -translate-x-1/2",
          props.tooltipClass,
        )}
      >
        {props.label}
      </span>
    </div>
  );
}
