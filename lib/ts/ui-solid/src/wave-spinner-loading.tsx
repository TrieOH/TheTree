import { cn } from "./lib/cn";
import { Spinner } from "./spinner";
import { WaveText, type WaveTextProps } from "./wave-text";

export interface WaveSpinnerLoadingProps extends WaveTextProps {
  spinnerSize?: string | number;
  activeColor?: string;
  trackColor?: string;
  class?: string;
}

export function WaveSpinnerLoading(props: WaveSpinnerLoadingProps) {
  return (
    <div
      role="status"
      class={cn(
        "flex h-full w-full flex-col items-center justify-center gap-5",
        props.class,
      )}
    >
      <Spinner
        size={props.spinnerSize ?? "6rem"}
        activeColor={props.activeColor ?? "var(--primary)"}
        trackColor={props.trackColor ?? "var(--accent)"}
      />
      <WaveText
        text={props.text ?? "Processando pagamento..."}
        delay={props.delay ?? 500}
        duration={props.duration ?? 1400}
        lift={props.lift ?? 12}
        waveWidth={props.waveWidth ?? 20}
      />
    </div>
  );
}

export default WaveSpinnerLoading;
