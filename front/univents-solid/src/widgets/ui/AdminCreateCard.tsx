import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import PlusIcon from "~icons/lucide/plus";

import { cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";

const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export type AdminCreateCardLayout = "horizontal" | "stacked" | "centered";

export interface AdminCreateCardProps {
  title: string;
  description: string;
  icon: JSX.Element | ((props: { class?: string }) => JSX.Element);
  onClick: () => void;
  index?: number;
  animate?: boolean;
  class?: string;
  minHeight?: string;
  ariaLabel?: string;
  layout?: AdminCreateCardLayout;
  showIndicator?: boolean;
}

export function AdminCreateCard(props: AdminCreateCardProps): JSX.Element {
  const layout = () => props.layout ?? "horizontal";
  const showIndicator = () => props.showIndicator ?? true;

  const renderIcon = () => {
    if (typeof props.icon === "function") {
      return (props.icon as (p: { class?: string }) => JSX.Element)({
        class: "size-5",
      });
    }
    return props.icon;
  };

  const indicator = (
    <div
      class="flex size-7 shrink-0 items-center justify-center rounded-full border border-dashed border-border/70 bg-background/50 text-muted-foreground/60 transition-all duration-300 group-hover:border-primary/50 group-hover:bg-primary group-hover:text-primary-foreground group-hover:rotate-90"
      aria-hidden="true"
    >
      <Plus class="size-3.5" />
    </div>
  );

  return (
    <Reveal
      delay={(props.index ?? 0) * 0.05}
      animate={props.animate}
      class="h-full"
    >
      <button
        type="button"
        onClick={() => props.onClick()}
        aria-label={props.ariaLabel ?? props.title}
        class={cn(
          "group relative flex h-full w-full min-w-0 rounded-xl",
          "border-2 border-dashed border-border/80 bg-card/40 p-3 text-left",
          "transition-all duration-200 hover:border-primary/50 hover:bg-muted/40",
          "focus:outline-none focus-visible:ring-2 focus-visible:ring-ring",
          layout() === "stacked" && "flex-col justify-between items-start",
          layout() === "centered" && "flex-col justify-center items-center text-center gap-2",
          layout() === "horizontal" && "flex-row items-center gap-3",
          props.minHeight ?? (layout() === "stacked" ? "min-h-23" : "min-h-17"),
          props.class,
        )}
      >
        <Show
          when={layout() === "stacked"}
          fallback={
            <>
              <div
                class={cn(
                  "flex shrink-0 items-center justify-center rounded-lg border border-dashed border-border/70 bg-muted/50 text-muted-foreground transition-colors group-hover:border-primary/40 group-hover:bg-primary/10 group-hover:text-primary",
                  layout() === "centered" ? "size-9" : "size-10",
                )}
              >
                {renderIcon()}
              </div>

              <div
                class={cn(
                  "min-w-0 space-y-0.5",
                  layout() === "centered" ? "w-full" : "flex-1",
                )}
              >
                <h3 class="truncate text-sm font-medium leading-tight text-foreground transition-colors group-hover:text-primary">
                  {props.title}
                </h3>

                <span class="block truncate text-xs text-muted-foreground transition-colors group-hover:text-muted-foreground/90">
                  {props.description}
                </span>
              </div>

              <Show when={showIndicator() && layout() === "horizontal"}>
                {indicator}
              </Show>
            </>
          }
        >
          {/* Top row with icon on the left and indicator on the right */}
          <div class="flex w-full items-center justify-between">
            <div class="flex size-9 shrink-0 items-center justify-center rounded-lg border border-dashed border-border/70 bg-muted/50 text-muted-foreground transition-colors group-hover:border-primary/40 group-hover:bg-primary/10 group-hover:text-primary">
              {renderIcon()}
            </div>

            <Show when={showIndicator()}>
              {indicator}
            </Show>
          </div>

          {/* Bottom row with title and description */}
          <div class="min-w-0 w-full space-y-0.5 pt-2">
            <h3 class="truncate text-sm font-medium leading-tight text-foreground transition-colors group-hover:text-primary">
              {props.title}
            </h3>

            <span class="block truncate text-xs text-muted-foreground transition-colors group-hover:text-muted-foreground/90">
              {props.description}
            </span>
          </div>
        </Show>
      </button>
    </Reveal>
  );
}
