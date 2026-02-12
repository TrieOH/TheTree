import { cn } from "@/shared/lib/utils";
import { cva, type VariantProps } from "class-variance-authority";
import React from "react";

const buttonVariants = cva(
  "flex cursor-pointer gap-1 font-extralight border rounded-sm p-2 items-center md:text-sm text-xs transition-all duration-300 ease-out",
  {
    variants: {
      variant: {
        default: "bg-primary-foreground text-card-foreground border-primary shadow-[1px_1px_0_0_var(--color-primary)] hover:shadow-[2px_2px_0_0_var(--color-primary)]",
        solid: "bg-primary text-primary-foreground border-primary shadow-[1px_1px_0_0_var(--color-primary)] hover:shadow-[2px_2px_0_0_var(--color-primary)]",
        "secondary-solid": "bg-secondary text-secondary-foreground border-secondary shadow-[1px_1px_0_0_var(--color-secondary)] hover:shadow-[2px_2px_0_0_var(--color-secondary)]",
        outline: "bg-background text-foreground border-input shadow-[1px_1px_0_0_var(--color-input)] hover:bg-muted/50 hover:shadow-[2px_2px_0_0_var(--color-input)]",
        ghost: "bg-transparent text-muted-foreground border-transparent shadow-none hover:bg-muted/50 hover:shadow-none",
        "ghost-primary": "bg-transparent text-primary border-transparent shadow-none hover:bg-primary/10 hover:shadow-none",
        destructive: "bg-destructive text-destructive-foreground border-destructive shadow-[1px_1px_0_0_var(--color-destructive)] hover:shadow-[2px_2px_0_0_var(--color-destructive)]",
        "accent-solid": "text-primary-foreground bg-primary border-accent shadow-[1px_1px_0_0_var(--color-accent)] hover:shadow-[2px_2px_0_0_var(--color-accent)]",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
);

interface PropsI extends VariantProps<typeof buttonVariants> {
  value?: string;
  label?: string;
  leftIcon?: React.ReactNode;
  type?: "button" | "submit";
  formId?: string;
  className?: string;
  onClick?: () => void;
  disabled?: boolean;
}

export const ShadowButton = React.forwardRef<
  HTMLButtonElement,
  PropsI
>(function ShadowButton(
  { 
    value, 
    label = value,
    leftIcon, 
    className, 
    onClick, 
    disabled, 
    variant,
    type = "button", 
    formId,
    ...props
  },
  ref
) {
  return (
    <button
      ref={ref}
      type={type}
      aria-label={label}
      title={label}
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
      {...props}
    >
      {leftIcon}
      {value}
    </button>
  );
});