import type { JSX } from "@solidjs/web";
import { QuickActions, type QuickActionItem } from "@/widgets/ui/QuickActions";

export type EditionQuickAction = QuickActionItem;

export function EditionQuickActions(props: {
  actions: EditionQuickAction[];
}): JSX.Element {
  return <QuickActions actions={props.actions} ariaLabel="Atalhos da edição" />;
}
