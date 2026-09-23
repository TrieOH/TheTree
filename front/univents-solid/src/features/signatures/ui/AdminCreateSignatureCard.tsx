import type { JSX } from "@solidjs/web";
import PenLineIcon from "~icons/lucide/pen-line";
import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateSignatureCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateSignatureCard(
  props: AdminCreateSignatureCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Nova assinatura"
      description="Cadastre uma assinatura desenhando no quadro ou importando uma imagem transparente."
      icon={PenLine}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[13.5rem]"
      class="p-4"
    />
  );
}
