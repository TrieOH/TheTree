import { cn } from "@/shared/lib/utils";

interface PropsI {
  value?: string;
  leftIcon?: React.ReactNode;
  className?: string;
  onClick: () => void;
}

export function ShadowButton({ value, leftIcon, className, onClick }: PropsI) {
  return (
    <button
      type="button"
      className={cn(
        "flex cursor-pointer gap-1 font-extralight bg-primary-foreground text-card-foreground",
        "border border-primary rounded-sm p-2 items-center md:text-sm text-xs ",

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
      onClick={onClick}
    >
      {leftIcon}
      {value}
    </button>
  )
}