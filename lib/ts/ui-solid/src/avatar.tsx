import { Show, createSignal } from "solid-js";
import { cn } from "./lib/cn";
import { shapeClass, type Shape } from "./lib/shape";

export interface AvatarProps {
  name: string;
  src?: string;
  initials?: string;
  size?: "sm" | "md" | "lg";
  shape?: Shape;
  class?: string;
}
const sizes = {
  sm: "size-7 text-[10px]",
  md: "size-9 text-xs",
  lg: "size-12 text-sm",
};
export function Avatar(props: AvatarProps) {
  const [failedSource, setFailedSource] = createSignal<string>();
  const initials = () => {
    if (props.initials) return props.initials;
    const words = props.name.trim().split(/\s+/).filter(Boolean);
    if (!words.length) return "?";
    return (
      words.length === 1
        ? words[0].slice(0, 2)
        : words[0][0] + words[words.length - 1][0]
    ).toLocaleUpperCase();
  };
  return (
    <span
      role="img"
      aria-label={props.name}
      class={cn(
        "inline-flex shrink-0 items-center justify-center overflow-hidden rounded-md border border-border bg-muted font-semibold text-foreground",
        sizes[props.size ?? "md"],
        shapeClass(props.shape),
        props.class,
      )}
    >
      <Show
        when={props.src && props.src !== failedSource()}
        fallback={initials()}
      >
        <img
          src={props.src}
          alt=""
          onError={() => {
            setFailedSource(props.src);
          }}
          class="size-full object-cover"
        />
      </Show>
    </span>
  );
}
