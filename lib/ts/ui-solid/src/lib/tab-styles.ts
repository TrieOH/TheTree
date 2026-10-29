import { cn } from "./cn";

export type TabsVariant = "underline" | "pills" | "segmented";
export type TabsSize = "sm" | "md";
export type TabsOrientation = "horizontal" | "vertical";

interface ListOptions {
  variant: TabsVariant;
  orientation: TabsOrientation;
  fullWidth: boolean;
}

interface TabOptions extends ListOptions {
  size: TabsSize;
  active: boolean;
}

/** Container around the triggers. */
export function tabListClass(options: ListOptions): string {
  const { variant, orientation, fullWidth } = options;
  const vertical = orientation === "vertical";

  return cn(
    "flex min-w-0 gap-1 scrollbar-none",
    vertical
      ? "flex-col items-stretch overflow-y-auto"
      : "items-center overflow-x-auto",
    fullWidth || vertical ? "w-full" : variant === "segmented" && "w-fit",
    !vertical && variant !== "segmented" && "w-full",
    variant === "underline" &&
      (vertical ? "border-l border-border" : "border-b border-border"),
    variant === "segmented" && "rounded-lg bg-muted p-1",
  );
}

export function tabClass(options: TabOptions): string {
  const { variant, orientation, size, active, fullWidth } = options;
  const vertical = orientation === "vertical";
  const small = size === "sm";

  return cn(
    "group relative inline-flex items-center gap-2 font-medium transition-colors cursor-pointer",
    "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring",
    "disabled:pointer-events-none disabled:opacity-50",
    vertical ? "justify-start" : "shrink-0",
    fullWidth && !vertical && "flex-1 justify-center",
    small ? "text-xs" : "text-sm",
    variant === "underline" && [
      small ? "px-2.5 py-1.5" : "px-3.5 py-2.5",
      vertical ? "-ml-px border-l-2" : "border-b-2",
      active
        ? "border-accent text-foreground font-semibold"
        : "border-transparent text-muted-foreground hover:text-foreground",
    ],
    variant === "pills" && [
      "rounded-full",
      small ? "px-2.5 py-1" : "px-3.5 py-1.5",
      active
        ? "bg-accent text-accent-foreground"
        : "text-muted-foreground hover:bg-muted hover:text-foreground",
    ],
    variant === "segmented" && [
      "rounded-md",
      small ? "px-2.5 py-1" : "px-3 py-1.5",
      active
        ? "bg-background text-foreground shadow-sm"
        : "text-muted-foreground hover:text-foreground",
    ],
  );
}

export function tabIconClass(active: boolean, variant: TabsVariant): string {
  return cn(
    "size-4 shrink-0 transition-colors",
    active
      ? variant === "pills"
        ? "text-accent-foreground"
        : "text-primary"
      : "text-muted-foreground group-hover:text-foreground",
  );
}

export function tabCountClass(active: boolean, variant: TabsVariant): string {
  return cn(
    "ml-1 inline-flex shrink-0 items-center rounded-full px-2 py-0.5 text-[11px] font-mono leading-none transition-colors",
    active
      ? variant === "pills"
        ? "bg-accent-foreground/15 text-accent-foreground font-medium"
        : "bg-accent text-accent-foreground font-medium"
      : "bg-muted text-muted-foreground",
  );
}
