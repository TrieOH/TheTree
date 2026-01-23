import { cn } from "@/lib/utils";

interface PropsI {
  value?: string;
  leftIcon?: React.ReactNode;
  className?: string;
}

export function ShadowButton({ value, leftIcon, className }: PropsI) {
  return (
    <button 
      className={cn(
        "flex cursor-pointer gap-2 font-extralight bg-primary-foreground text-card-foreground",
        "border border-primary rounded-sm p-2",

        // base state
        "shadow-[1px_1px_0_0_var(--color-primary)]",

        // hover
        "hover:shadow-[2px_2px_0_0_var(--color-primary)]",

        // pressed
        "active:translate-x-px active:translate-y-px",
        "active:shadow-none",

        "transition-all duration-300 ease-out",
        className
      )}
    >
      {leftIcon}
      {value}
    </button>
  )
}