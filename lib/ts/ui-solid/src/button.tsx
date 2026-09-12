import type { JSX } from "@solidjs/web";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "./lib/cn";

/**
 * Variant names and tokens mirror `@trieoh/ui-react`'s button so the two
 * frameworks look identical. The moment both files are stable, the class
 * strings move to a shared `ui-core` and this file keeps only the rendering.
 */
export const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        outline: "border border-border bg-background hover:bg-muted hover:text-foreground",
        ghost: "hover:bg-muted hover:text-foreground",
        destructive: "bg-destructive text-destructive-foreground hover:bg-destructive/90",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        sm: "h-8 px-3",
        md: "h-9 px-4",
        lg: "h-10 px-5",
        icon: "size-9",
      },
    },
    defaultVariants: { variant: "default", size: "md" },
  },
);

export type ButtonVariant = NonNullable<VariantProps<typeof buttonVariants>["variant"]>;
export type ButtonSize = NonNullable<VariantProps<typeof buttonVariants>["size"]>;

/**
 * Props are declared explicitly instead of spreading the DOM's: Solid 2 dropped
 * `splitProps`, and destructuring props would break reactivity. Passing a prop
 * that is not listed here is a type error rather than a silently dropped
 * attribute — add to the list when a screen needs more.
 */
export interface ButtonProps {
  variant?: ButtonVariant;
  size?: ButtonSize;
  type?: "button" | "submit" | "reset";
  disabled?: boolean;
  class?: string;
  id?: string;
  name?: string;
  value?: string;
  title?: string;
  form?: string;
  "aria-label"?: string;
  /** Solid 2 types ARIA booleans as the string form. */
  "aria-expanded"?: "true" | "false";
  "aria-controls"?: string;
  onClick?: (event: MouseEvent & { currentTarget: HTMLButtonElement }) => void;
  children?: JSX.Element;
}

export function Button(props: ButtonProps) {
  return (
    <button
      type={props.type ?? "button"}
      id={props.id}
      name={props.name}
      value={props.value}
      title={props.title}
      form={props.form}
      aria-label={props["aria-label"]}
      aria-expanded={props["aria-expanded"]}
      aria-controls={props["aria-controls"]}
      disabled={props.disabled}
      onClick={props.onClick}
      class={cn(
        buttonVariants({ variant: props.variant, size: props.size }),
        props.class,
      )}
    >
      {props.children}
    </button>
  );
}
