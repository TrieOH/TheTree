import { Spinner } from "./spinner";
import WaveText from "./WaveText";

interface WaveTextProps {
  text?: string;
  duration?: number;
  delay?: number;
  lift?: number;
  waveWidth?: number;
}

export default function WaveSpinnerLoading(props: WaveTextProps) {
  return (
    <div class="w-full h-full flex flex-col items-center justify-center gap-5">
      <Spinner
        size={"6rem"}
        activeColor="var(--primary)"
        trackColor="var(--accent)"
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
