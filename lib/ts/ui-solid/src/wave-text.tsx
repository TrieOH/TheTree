import { mergeProps } from "@solidjs/web";
import { For, createEffect, onCleanup } from "solid-js";
import { cn } from "./lib/cn";

export interface WaveTextProps {
  text?: string;
  duration?: number;
  delay?: number;
  lift?: number;
  waveWidth?: number;
  class?: string;
}

export function WaveText(props: WaveTextProps) {
  const merged = mergeProps(
    {
      text: "Processando pagamento...",
      duration: 1400,
      delay: 500,
      lift: 12,
      waveWidth: 20,
    },
    props,
  );

  let containerRef: HTMLSpanElement | undefined;
  let raf: number | undefined;
  let timeout: ReturnType<typeof setTimeout> | undefined;
  let startTime: number | undefined;

  function frame(ts: number) {
    if (!containerRef) return;
    const spans = containerRef.children;
    if (spans.length === 0) {
      raf = requestAnimationFrame(frame);
      return;
    }

    startTime ??= ts;

    const halfW = Math.max(merged.waveWidth / 100 / 2, 0.0001);
    const t = Math.min((ts - startTime) / merged.duration, 1);
    const waveCenter = -halfW + t * (1 + 2 * halfW);

    const n = Math.max(spans.length - 1, 1);
    const dimBase = 0.18;

    for (let i = 0; i < spans.length; i++) {
      const el = spans[i] as HTMLElement | undefined;
      if (!el) continue;

      const pos = i / n;
      const dist = (pos - waveCenter) / halfW;

      let sineVal = 0;
      if (dist >= -1 && dist <= 1) {
        sineVal = Math.sin(dist * Math.PI * 0.5 + Math.PI * 0.5);
      }

      el.style.transform = `translateY(${(-merged.lift * sineVal).toFixed(2)}px)`;
      el.style.opacity = (dimBase + (1 - dimBase) * sineVal).toFixed(3);
    }

    if (t < 1) {
      raf = requestAnimationFrame(frame);
    } else {
      timeout = setTimeout(restart, merged.delay);
    }
  }

  function restart() {
    if (raf !== undefined) {
      cancelAnimationFrame(raf);
    }
    if (timeout !== undefined) {
      clearTimeout(timeout);
    }
    startTime = undefined;
    raf = requestAnimationFrame(frame);
  }

  createEffect(
    () => [
      merged.text,
      merged.duration,
      merged.delay,
      merged.lift,
      merged.waveWidth,
    ],
    () => {
      startTime = undefined;
      if (raf !== undefined) cancelAnimationFrame(raf);
      if (timeout !== undefined) clearTimeout(timeout);
      raf = requestAnimationFrame(frame);
    },
  );

  onCleanup(() => {
    if (raf !== undefined) cancelAnimationFrame(raf);
    if (timeout !== undefined) clearTimeout(timeout);
  });

  return (
    <span
      ref={(el) => (containerRef = el)}
      class={cn(
        "inline-flex font-bold text-lg text-primary select-none",
        props.class,
      )}
    >
      <For each={merged.text.split("")}>
        {(char) => (
          <span class="inline-block" style="white-space: pre;">
            {char}
          </span>
        )}
      </For>
    </span>
  );
}

export default WaveText;
