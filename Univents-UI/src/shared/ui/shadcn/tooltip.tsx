import { Tooltip as TooltipPrimitive } from "@base-ui/react/tooltip"
import { cn } from "@/shared/lib/utils"

function TooltipProvider({
  delay = 0,
  ...props
}: TooltipPrimitive.Provider.Props) {
  return (
    <TooltipPrimitive.Provider
      data-slot="tooltip-provider"
      delay={delay}
      {...props}
    />
  )
}

function Tooltip({ ...props }: TooltipPrimitive.Root.Props) {
  return <TooltipPrimitive.Root data-slot="tooltip" {...props} />
}

function TooltipTrigger({ ...props }: TooltipPrimitive.Trigger.Props) {
  return <TooltipPrimitive.Trigger data-slot="tooltip-trigger" {...props} />
}

function TooltipContent({
  className,
  side = "top",
  sideOffset = 6,
  align = "center",
  alignOffset = 0,
  children,
  ...props
}: TooltipPrimitive.Popup.Props &
  Pick<
    TooltipPrimitive.Positioner.Props,
    "align" | "alignOffset" | "side" | "sideOffset"
  >) {
  return (
    <TooltipPrimitive.Portal>
      <TooltipPrimitive.Positioner
        align={align}
        alignOffset={alignOffset}
        side={side}
        sideOffset={sideOffset}
        className="isolate z-50"
      >
        <TooltipPrimitive.Popup
          data-slot="tooltip-content"
          className={cn(
            // Layout
            "z-50 inline-flex w-fit max-w-xs origin-(--transform-origin) items-center gap-1.5",
            "rounded-md px-3 py-1.5 text-xs font-medium",

            // Colors
            "bg-popover text-popover-foreground",

            // Border
            "border border-border/50",

            // Shadow
            "shadow-lg shadow-black/10 dark:shadow-black/40",

            // Animations
            "data-[side=bottom]:slide-in-from-top-2",
            "data-[side=left]:slide-in-from-right-2",
            "data-[side=right]:slide-in-from-left-2",
            "data-[side=top]:slide-in-from-bottom-2",
            "data-[state=delayed-open]:animate-in",
            "data-[state=delayed-open]:fade-in-0",
            "data-[state=delayed-open]:zoom-in-95",
            "data-open:animate-in",
            "data-open:fade-in-0",
            "data-open:zoom-in-95",
            "data-closed:animate-out",
            "data-closed:fade-out-0",
            "data-closed:zoom-out-95",
            className
          )}
          {...props}
        >
          {children}
          <TooltipPrimitive.Arrow
            className={cn(
              "z-50 size-2 rotate-45 rounded-[2px]",
              // Arrow
              "bg-popover border-r border-b border-border/50",
              // Positioning
              "data-[side=bottom]:top-1",
              "data-[side=left]:-right-1",
              "data-[side=right]:-left-1",
              "data-[side=top]:-bottom-1"
            )}
          />
        </TooltipPrimitive.Popup>
      </TooltipPrimitive.Positioner>
    </TooltipPrimitive.Portal>
  )
}

export { Tooltip, TooltipTrigger, TooltipContent, TooltipProvider }