import { cn } from "./lib/cn";
import { Spinner } from "./spinner";

export interface LoadingProps {
  label?: string;
  showLabel?: boolean;
  size?: "sm" | "md" | "lg";
  tone?: "default" | "accent" | "muted";
  layout?: "inline" | "stacked";
  class?: string;
}

const sizes = { sm: 16, md: 24, lg: 48 };
const tones = {
  default: "text-foreground",
  accent: "text-accent",
  muted: "text-muted-foreground",
};

/** Accessible loading status using the Univents spinner. */
export function Loading(props: LoadingProps) {
  return (
    <div
      role="status"
      class={cn(
        "flex items-center justify-center gap-3 text-sm",
        props.layout === "stacked" && "flex-col gap-4",
        props.class,
      )}
    >
      <Spinner
        size={sizes[props.size ?? "md"]}
        activeColor="currentColor"
        trackColor="color-mix(in srgb, currentColor 18%, transparent)"
        className={cn("shrink-0", tones[props.tone ?? "accent"])}
      />
      <span
        class={props.showLabel === false ? "sr-only" : "text-muted-foreground"}
      >
        {props.label || "Carregando…"}
      </span>
    </div>
  );
}
