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
  icon: ((props: { class?: string }) => JSX.Element) | JSX.Element;
  onClick: () => void;
  ariaLabel?: string;
  index?: number;
  animate?: boolean;
  class?: string;
  /**
   * Layout style of the card:
   * - "horizontal": Left-to-right icon + text + action indicator (standard for events)
   * - "stacked": Vertical with icon & indicator on top, title & description below (editions, members, tickets)
   * - "centered": Centered vertical layout
   * Default: "horizontal"
   */
  layout?: AdminCreateCardLayout;
  /**
   * Minimum height class (e.g. "min-h-16", "min-h-[5.75rem]", "min-h-[8.5rem]")
   */
  minHeight?: string;
  /**
   * Whether to show the circular action indicator (plus button)
   * Default: true
   */
  showActionIndicator?: boolean;
  /**
   * @deprecated alias for showActionIndicator
   */
  showIndicator?: boolean;
}

export function AdminCreateCard(props: AdminCreateCardProps): JSX.Element {
  const layout = () => props.layout ?? "horizontal";
  const showIndicator = () => props.showActionIndicator ?? props.showIndicator ?? true;

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
          "focus:outline-none focus-visible:ring-2 focus-visible:ring-ring cursor-pointer",
          layout() === "horizontal" && "flex-row items-center gap-3 min-h-16",
          layout() === "stacked" && "flex-col justify-between items-start min-h-23",
          layout() === "centered" && "flex-col justify-center items-center text-center gap-2 min-h-16",
          props.minHeight,
          props.class,
        )}
      >
        <Show
          when={layout() === "stacked"}
          fallback={
            <>
              <div class="flex size-10 shrink-0 items-center justify-center rounded-lg border border-dashed border-border/70 bg-muted/50 text-muted-foreground transition-colors group-hover:border-primary/40 group-hover:bg-primary/10 group-hover:text-primary">
                {renderIcon()}
              </div>

              <div class="min-w-0 flex-1 space-y-0.5">
                <h3 class="truncate text-sm font-medium leading-tight text-foreground transition-colors group-hover:text-primary">
                  {props.title}
                </h3>

                <span class="block truncate text-xs text-muted-foreground transition-colors group-hover:text-muted-foreground/90">
                  {props.description}
                </span>
              </div>

              <Show when={showIndicator()}>
                {indicator}
              </Show>
            </>
          }
        >
          {/* Linha superior: ícone à esquerda e indicador de ação à direita */}
          <div class="flex w-full items-center justify-between">
            <div class="flex size-8.5 shrink-0 items-center justify-center rounded-lg border border-dashed border-border/70 bg-muted/50 text-muted-foreground transition-colors group-hover:border-primary/40 group-hover:bg-primary/10 group-hover:text-primary">
              {renderIcon()}
            </div>

            <Show when={showIndicator()}>
              {indicator}
            </Show>
          </div>

          {/* Linha inferior: título e descrição */}
          <div class="min-w-0 w-full space-y-0.5 pt-1.5">
            <h3 class="truncate text-sm font-medium leading-tight text-foreground transition-colors group-hover:text-primary">
              {props.title}
            </h3>

            <p class="text-xs leading-relaxed text-muted-foreground transition-colors group-hover:text-muted-foreground/90 line-clamp-3">
              {props.description}
            </p>
          </div>
        </Show>
      </button>
    </Reveal>
  );
}
