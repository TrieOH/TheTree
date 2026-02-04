import { cn } from "@/shared/lib/utils";

interface PropsI {
  value?: string;
  leftIcon?: React.ReactNode;
  type?: "button" | "submit";
  className?: string;
  onClick?: () => void;
  disabled?: boolean;
}

export function ShadowButton({ value, leftIcon, className, onClick, type = "button", disabled }: PropsI) {
  return (
    <button
      type={type}
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
        disabled && "opacity-50 cursor-not-allowed shadow-none hover:shadow-none active:translate-x-0 active:translate-y-0",
        className
      )}
      onClick={onClick}
      disabled={disabled}
    >
      {leftIcon}
      {value}
    </button>
  )
}