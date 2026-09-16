import type { JSX } from "@solidjs/web";

import PlusIcon from "~icons/lucide/plus";

import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateVariantCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateVariantCard(
  props: AdminCreateVariantCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Nova variação"
      description="Adicione novos tamanhos, cores ou modelos com preços e estoques independentes."
      icon={Plus}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[8.5rem]"
      class="p-4"
    />
  );
}
