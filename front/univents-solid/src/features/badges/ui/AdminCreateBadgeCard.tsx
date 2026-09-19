import type { JSX } from "@solidjs/web";
import BadgeCheckIcon from "~icons/lucide/badge-check";
import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const BadgeCheck = BadgeCheckIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateBadgeCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateBadgeCard(
  props: AdminCreateBadgeCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Novo template"
      description="Crie um novo modelo visual de crachá para participantes ou equipe desta edição."
      icon={BadgeCheck}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[18rem]"
      class="p-5"
    />
  );
}
