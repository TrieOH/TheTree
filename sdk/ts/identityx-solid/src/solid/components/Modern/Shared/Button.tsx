import type { ParentProps } from "solid-js";

export interface ButtonProps extends ParentProps {
  variant?: "default" | "outline" | "ghost" | "link";
  size?: "default" | "sm" | "lg" | "icon";
  class?: string;
  type?: "button" | "submit";
  disabled?: boolean;
  onClick?: (event: MouseEvent) => void;
}
export function Button(props: ButtonProps) {
  const variant = () =>
    props.variant === "outline"
      ? "border border-input bg-background hover:bg-accent"
      : props.variant === "ghost"
        ? "hover:bg-accent"
        : props.variant === "link"
          ? "text-primary underline-offset-4 hover:underline"
          : "bg-primary text-primary-foreground hover:opacity-90";
  return (
    <button
      type={props.type ?? "button"}
      disabled={props.disabled}
      onClick={props.onClick}
      class={`inline-flex items-center justify-center whitespace-nowrap rounded-sm text-sm font-semibold transition-all disabled:pointer-events-none disabled:opacity-50 ${variant()} ${props.class ?? ""}`}
    >
      {props.children}
    </button>
  );
}
