import type { JSX } from "@solidjs/web";
import MailIcon from "~icons/lucide/mail";
import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateSignatureRequestCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateSignatureRequestCard(
  props: AdminCreateSignatureRequestCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Novo convite"
      description="Envie um link por e-mail para que um palestrante ou coordenador assine digitalmente."
      icon={Mail}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[11rem]"
      class="p-4"
    />
  );
}
