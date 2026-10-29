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
      class={`origin-center overflow-visible -rotate-90 motion-safe:animate-[spin_2s_linear_infinite] w-[var(--size)] ${merged.className}`}
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle
        class="stroke-[var(--active)] [stroke-linecap:round] [stroke-dashoffset:360] motion-safe:animate-[trieoh-spinner-active_var(--duration)_ease-in-out_infinite] motion-reduce:[stroke-dashoffset:0] motion-reduce:[stroke-dasharray:90_270]"
        pathLength={360}
        fill="transparent"
        stroke-width={32}
        cx={192}
        cy={192}
        r={176}
      />

      <circle
        class="stroke-[var(--track)] [stroke-linecap:round] [stroke-dashoffset:360] motion-safe:animate-[trieoh-spinner-track_var(--duration)_ease-in-out_infinite] motion-reduce:[stroke-dashoffset:0] motion-reduce:[stroke-dasharray:360]"
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
