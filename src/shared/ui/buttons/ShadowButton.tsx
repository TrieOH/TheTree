import { cn } from "@/shared/lib/utils";
import { cva, type VariantProps } from "class-variance-authority";

const buttonVariants = cva(
  "flex cursor-pointer gap-1 font-extralight border rounded-sm p-2 items-center md:text-sm text-xs transition-all duration-300 ease-out",
  {
    variants: {
      variant: {
        default: "bg-primary-foreground text-card-foreground border-primary shadow-[1px_1px_0_0_var(--color-primary)] hover:shadow-[2px_2px_0_0_var(--color-primary)]",
        destructive: "bg-destructive text-destructive-foreground border-destructive shadow-[1px_1px_0_0_var(--color-destructive)] hover:shadow-[2px_2px_0_0_var(--color-destructive)]",
        accent: "text-primary-foreground bg-primary border-accent shadow-[1px_1px_0_0_var(--color-accent)] hover:shadow-[2px_2px_0_0_var(--color-accent)]",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
);

interface PropsI extends VariantProps<typeof buttonVariants> {
  value?: string;
  leftIcon?: React.ReactNode;
  type?: "button" | "submit";
  formId?: string;
  className?: string;
  onClick?: () => void;
  disabled?: boolean;
}

export function ShadowButton({ 
  value, 
  leftIcon, 
  className, 
  onClick, 
  disabled, 
  variant,
  type = "button", 
  formId
}: PropsI) {
  return (
    <button
      type={type}
      form={formId}
      className={cn(
        buttonVariants({ variant, className }),
        // pressed
        "active:translate-x-px active:translate-y-px",
        "active:shadow-none",
        disabled && "opacity-50 cursor-not-allowed shadow-none hover:shadow-none active:translate-x-0 active:translate-y-0",
      )}
      onClick={onClick}
      disabled={disabled}
    >
      {leftIcon}
      {value}
    </button>
  )
}