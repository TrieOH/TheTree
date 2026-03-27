import { useRef, useCallback } from "react";

interface WaveTextProps {
  text?: string;
  duration?: number;
  delay?: number;
  lift?: number;
  waveWidth?: number;
}

export default function WaveText({
  text = "Processando pagamento...",
  duration = 1400,
  delay = 500,
  lift = 12,
  waveWidth = 20,
}: WaveTextProps) {
  const spans = useRef<(HTMLSpanElement | null)[]>([]);
  const raf = useRef<number | null>(null);
  const timeout = useRef<ReturnType<typeof setTimeout> | null>(null);
  const startTime = useRef<number | null>(null);
  const mounted = useRef(false);

  const props = useRef({ duration, delay, lift, waveWidth });
  props.current = { duration, delay, lift, waveWidth };

  const frame = useCallback((ts: number) => {
    startTime.current ??= ts

    const p = props.current;

    const halfW = p.waveWidth / 100 / 2;
    const t = Math.min((ts - startTime.current) / p.duration, 1);
    const waveCenter = -halfW + t * (1 + 2 * halfW);
    const els = spans.current;
    const n = els.length - 1;
    const dimBase = 0.18;

    for (let i = 0; i < els.length; i++) {
      const el = els[i];
      if (!el) continue;
      const pos = i / n;
      const dist = (pos - waveCenter) / halfW;
      let sineVal = 0;
      if (dist >= -1 && dist <= 1) {
        sineVal = Math.sin(dist * Math.PI * 0.5 + Math.PI * 0.5);
      }
      el.style.transform = `translateY(${(-lift * sineVal).toFixed(2)}px)`;
      el.style.opacity = (dimBase + (1 - dimBase) * sineVal).toFixed(3);
    }

    if (t < 1) {
      raf.current = requestAnimationFrame(frame);
    } else {
      timeout.current = setTimeout(restart, delay);
    }
  }, []);

  const restart = useCallback(() => {
    if (raf.current) cancelAnimationFrame(raf.current);
    if (timeout.current) clearTimeout(timeout.current);
    startTime.current = null;
    raf.current = requestAnimationFrame(frame);
  }, []);

  const containerRef = useCallback((node: HTMLSpanElement | null) => {
    if (node && !mounted.current) {
      mounted.current = true;
      raf.current = requestAnimationFrame(frame);
    }
    if (!node) {
      if (raf.current) cancelAnimationFrame(raf.current);
      if (timeout.current) clearTimeout(timeout.current);
      mounted.current = false;
    }
  }, []);

  return (
    <span ref={containerRef} className="inline-flex font-bold text-lg text-primary">
      {text.split("").map((char, i) => (
        <span
          key={i}
          ref={(el) => { spans.current[i] = el; }}
          className="inline-block"
          style={{ whiteSpace: "pre" }}
        >
          {char}
        </span>
      ))}
    </span>
  );
}