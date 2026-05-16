import { Spinner } from "./spinner";
import WaveText from "./WaveText";

interface WaveTextProps {
  text?: string;
  duration?: number;
  delay?: number;
  lift?: number;
  waveWidth?: number;
}

export default function WaveSpinnerLoading({
  text = "Processando pagamento...",
  duration = 1400,
  delay = 500,
  lift = 12,
  waveWidth = 20,
}: WaveTextProps) {
  return (
    <div className="w-full h-full flex flex-col items-center justify-center gap-5">
      <Spinner size={"6rem"} activeColor="var(--primary)" trackColor="var(--accent)" />
      <WaveText
        text={text}
        delay={delay}
        duration={duration}
        lift={lift}
        waveWidth={waveWidth}
      />
    </div>
  )
}