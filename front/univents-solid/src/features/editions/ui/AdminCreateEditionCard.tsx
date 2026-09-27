import type { JSX } from "@solidjs/web";

import CalendarPlusIcon from "~icons/lucide/calendar-plus";

import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const CalendarPlus = CalendarPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateEditionCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateEditionCard(
  props: AdminCreateEditionCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Nova edição"
      description="Criar uma nova edição do evento"
      icon={CalendarPlus}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[5.75rem]"
    />
  );
}
