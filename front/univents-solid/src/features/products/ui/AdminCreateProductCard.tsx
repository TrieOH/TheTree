import type { JSX } from "@solidjs/web";

import PackagePlusIcon from "~icons/lucide/package-plus";

import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const PackagePlus = PackagePlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateProductCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateProductCard(
  props: AdminCreateProductCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Novo produto"
      description="Cadastre produtos físicos ou digitais com variações e controle de estoque."
      icon={PackagePlus}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[8.5rem]"
      class="p-4"
    />
  );
}
