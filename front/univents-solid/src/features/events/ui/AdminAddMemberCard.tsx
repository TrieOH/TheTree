import type { JSX } from "@solidjs/web";

import UserPlusIcon from "~icons/lucide/user-plus";

import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const UserPlus = UserPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminAddMemberCardProps {
  index?: number;
  animate?: boolean;
  onAdd: () => void;
}

export function AdminAddMemberCard(
  props: AdminAddMemberCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Adicionar membro"
      description="Convidar para a equipe"
      icon={UserPlus}
      onClick={props.onAdd}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[5.75rem]"
    />
  );
}
