import { createEffect, createSignal } from "solid-js";
import { useTheme } from "@/shared/lib/theme";

let counter = 0;

/**
 * A Mermaid diagram, from a fenced ```mermaid block (mdx/remark-mermaid.ts).
 * Mermaid is large, so it loads only on pages that draw one, and redraws when
 * the theme changes.
 */
export function Mermaid(props: { chart: string }) {
  const { theme } = useTheme();
  const id = `mermaid-${++counter}`;
  const [svg, setSvg] = createSignal<string>("", { name: "mermaid-svg" });
  const [failed, setFailed] = createSignal(false, { name: "mermaid-failed" });

  const draw = async (mode: "dark" | "light") => {
    try {
      const { default: mermaid } = await import("mermaid");
      mermaid.initialize({ startOnLoad: false, theme: mode === "dark" ? "dark" : "neutral", fontFamily: "Geist" });
      const { svg } = await mermaid.render(`${id}-${mode}`, props.chart);
      setSvg(svg);
    } catch {
      setFailed(true);
    }
  };

  createEffect(theme, (mode) => {
    void draw(mode);
  });

  return (
    <figure class="not-prose my-6 overflow-x-auto rounded-xl border border-border bg-card p-4">
      {failed() ? (
        <pre class="text-[12px] text-muted-foreground">{props.chart}</pre>
      ) : (
        <div class="flex justify-center [&_svg]:max-w-full" innerHTML={svg()} aria-label="Diagram" role="img" />
      )}
    </figure>
  );
}
