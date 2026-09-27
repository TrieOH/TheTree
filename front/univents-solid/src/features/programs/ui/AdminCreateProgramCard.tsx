import type { JSX } from "@solidjs/web";

import CalendarPlusIcon from "~icons/lucide/calendar-plus";

import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const CalendarPlus = CalendarPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateProgramCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateProgramCard(
  props: AdminCreateProgramCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Novo programa"
      description="Crie atividades, palestras ou checkpoints com horários e controle de acesso."
      icon={CalendarPlus}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[8.5rem]"
      class="p-4"
    />
  );
}
