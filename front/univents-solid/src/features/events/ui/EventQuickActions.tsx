import type { JSX } from "@solidjs/web";
import { QuickActions, type QuickActionItem } from "@/widgets/ui/QuickActions";

export type EventQuickAction = QuickActionItem;

export function EventQuickActions(props: {
  actions: EventQuickAction[];
}): JSX.Element {
  return <QuickActions actions={props.actions} ariaLabel="Atalhos do evento" />;
}
