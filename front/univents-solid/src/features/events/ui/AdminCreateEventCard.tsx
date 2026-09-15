import type { JSX } from "@solidjs/web";

import PlusIcon from "~icons/lucide/plus";

import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateEventCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateEventCard(
  props: AdminCreateEventCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Novo evento"
      description="Criar um novo evento"
      icon={Plus}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      minHeight="min-h-[4rem]"
    />
  );
}
