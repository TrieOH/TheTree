import type { JSX } from "@solidjs/web";
import { EmptyState } from "@trieoh/ui-solid";

/** Placeholder for data a later phase connects. */
export function NotConnected(props: { source: string; children?: JSX.Element }) {
  return <EmptyState variant="inline" title={props.children ?? <>Not connected yet. This fills in from {props.source}.</>} />;
}
