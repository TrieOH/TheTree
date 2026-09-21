import type { JSX } from "@solidjs/web";
import AwardIcon from "~icons/lucide/award";
import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const Award = AwardIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateCertificationCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateCertificationCard(
  props: AdminCreateCertificationCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Novo template"
      description="Crie um novo modelo visual de certificado para participantes ou atividades desta edição."
      icon={Award}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[13.5rem]"
      class="p-4"
    />
  );
}
