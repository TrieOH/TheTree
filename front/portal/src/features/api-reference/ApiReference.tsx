import { createEffect, createSignal, onCleanup, onSettled } from "solid-js";
import type { ServiceId } from "@/features/catalog/services";
import { useTheme } from "@/shared/lib/theme";
import { loadSpec } from "./specs";

// Scalar's palette points at the portal's tokens, so it follows the portal's
// theme on its own; forceDarkModeState only flips its code highlighting.
const CUSTOM_CSS = `
.light-mode, .dark-mode {
  --scalar-background-1: var(--background);
  --scalar-background-2: var(--card);
  --scalar-background-3: var(--muted);
  --scalar-background-accent: color-mix(in oklch, var(--accent) 14%, transparent);
  --scalar-color-1: var(--foreground);
  --scalar-color-2: var(--muted-foreground);
  --scalar-color-3: color-mix(in oklch, var(--muted-foreground) 75%, transparent);
  --scalar-color-accent: var(--accent);
  --scalar-border-color: var(--border);
  --scalar-font: "Geist", ui-sans-serif, system-ui, sans-serif;
  --scalar-font-code: ui-monospace, SFMono-Regular, Menlo, monospace;
}
.scalar-app .t-doc__sidebar { --scalar-sidebar-background-1: var(--background); }
`;

type Instance = { destroy: () => void; updateConfiguration: (config: Record<string, unknown>) => void };

export function ApiReference(props: { service: ServiceId }) {
  const { theme } = useTheme();
  const [state, setState] = createSignal<"loading" | "ready" | "missing" | "failed">("loading", {
    name: "api-reference",
  });
  let host: HTMLDivElement | undefined;
  let instance: Instance | undefined;

  const config = (content: string) => ({
    content,
    theme: "none",
    customCss: CUSTOM_CSS,
    forceDarkModeState: theme(),
    hideDarkModeToggle: true,
    withDefaultFonts: false,
    // Nothing leaves the portal, and nothing can call the API: a request from
    // here would wake its container.
    hideClientButton: true,
    hideTestRequestButton: true,
    showDeveloperTools: "never",
    telemetry: false,
    agent: { disabled: true },
    mcp: { disabled: true },
    documentDownloadType: "yaml",
  });

  let content = "";
  onSettled(() => {
    const spec = loadSpec(props.service);
    if (!spec) {
      setState("missing");
      return;
    }
    void Promise.all([spec, import("./scalar")])
      .then(([text, { createApiReference }]) => {
        if (!host) return;
        content = text;
        instance = createApiReference(host, config(text) as never) as unknown as Instance;
        setState("ready");
      })
      .catch(() => setState("failed"));
  });

  createEffect(theme, (mode) => {
    instance?.updateConfiguration({ ...config(content), forceDarkModeState: mode });
  });

  onCleanup(() => instance?.destroy());

  return (
    <>
      {state() === "loading" && <p class="px-1 py-10 text-[13px] text-muted-foreground">Loading the API reference…</p>}
      {state() === "missing" && (
        <p class="px-1 py-10 text-[13px] text-muted-foreground">
          This build has no generated spec for the service. Run <code>just gen</code>, or check the release's
          generate step.
        </p>
      )}
      {state() === "failed" && <p class="px-1 py-10 text-[13px] text-muted-foreground">The API reference failed to load.</p>}
      <div ref={(el) => (host = el)} class="min-h-[60vh]" />
    </>
  );
}
