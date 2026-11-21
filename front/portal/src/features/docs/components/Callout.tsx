import type { JSX } from "@solidjs/web";
import { Callout as RegistryCallout, type CalloutTone } from "@trieoh/ui-solid";

export type CalloutType = "note" | "tip" | "warning" | "danger";

const TONE: Record<CalloutType, CalloutTone> = { note: "info", tip: "success", warning: "warning", danger: "destructive" };

/** `<Callout type="warning" title="Wakes the API">…</Callout>`; type defaults to note. */
export function Callout(props: { type?: CalloutType; title?: string; children?: JSX.Element }) {
  return (
    <RegistryCallout tone={TONE[props.type ?? "note"]} title={props.title} class="not-prose my-5">
      {props.children}
    </RegistryCallout>
  );
}
