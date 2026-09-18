import type { JSX } from "@solidjs/web";

import CalendarPlusIcon from "~icons/lucide/calendar-plus";

import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const CalendarPlus = CalendarPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateOccurrenceCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateOccurrenceCard(
  props: AdminCreateOccurrenceCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Nova ocorrência"
      description="Cadastre data e horário com vagas e controle de presença para esta atividade."
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
