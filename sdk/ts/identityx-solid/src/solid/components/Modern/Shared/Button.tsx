import type { ParentProps } from "solid-js";

export interface ButtonProps extends ParentProps {
  variant?: "default" | "destructive" | "outline" | "secondary" | "ghost" | "link";
  size?: "default" | "sm" | "lg" | "icon";
  class?: string;
  type?: "button" | "submit";
  disabled?: boolean;
  onClick?: (event: MouseEvent) => void;
}

export function Button(props: ButtonProps) {
  const variant = () =>
    props.variant === "destructive"
      ? "bg-destructive text-destructive-foreground hover:opacity-90"
      : props.variant === "outline"
      ? "border border-input bg-background hover:bg-accent hover:text-accent-foreground"
      : props.variant === "secondary"
      ? "bg-secondary text-secondary-foreground hover:opacity-80"
      : props.variant === "ghost"
      ? "hover:bg-accent hover:text-accent-foreground"
      : props.variant === "link"
      ? "text-primary underline-offset-4 hover:underline"
      : "bg-primary text-primary-foreground hover:opacity-90";

  const size = () =>
    props.size === "sm"
      ? "h-9 rounded-md px-3"
      : props.size === "lg"
      ? "h-11 rounded-md px-8"
      : props.size === "icon"
      ? "h-10 w-10"
      : "h-12 px-4 py-2";

  return (
    <button
      type={props.type ?? "button"}
      disabled={props.disabled}
      onClick={props.onClick}
      class={`inline-flex items-center justify-center whitespace-nowrap rounded-sm text-sm font-semibold transition-all focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 active:scale-[0.98] ${variant()} ${size()} ${props.class ?? ""}`}
    >
      {props.children}
    </button>
  );
}
