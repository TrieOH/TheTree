import { mergeProps } from "@solidjs/web";

export interface SpinnerProps {
  size?: number | string;
  activeColor?: string;
  trackColor?: string;
  duration?: string;
  className?: string;
}

export function Spinner(props: SpinnerProps) {
  const merged = mergeProps(
    {
      size: "3rem" as number | string,
      activeColor: "var(--primary)",
      trackColor: "color-mix(in srgb, var(--primary) 18%, transparent)",
      duration: "8s",
      className: "",
    },
    props,
  );

  const size = () =>
    typeof merged.size === "number" ? `${merged.size}px` : merged.size;

  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 384 384"
      style={`
        --size: ${size()};
        --active: ${merged.activeColor};
        --track: ${merged.trackColor};
        --duration: ${merged.duration};
      `}
      class={`origin-center overflow-visible -rotate-90 motion-safe:animate-[spin_2s_linear_infinite] trieoh-spinner ${merged.className}`}
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle
        class="trieoh-spinner-active"
        pathLength={360}
        fill="transparent"
        stroke-width={32}
        cx={192}
        cy={192}
        r={176}
      />

      <circle
        class="trieoh-spinner-track"
        pathLength={360}
        fill="transparent"
        stroke-width={32}
        cx={192}
        cy={192}
        r={176}
      />
    </svg>
  );
}
